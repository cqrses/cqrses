package aggregate

import (
	"sync"

	"gopkg.in/cqrses/messages"
)

type (
	Stream struct {
		aggregate Aggregate
		events    []*messages.Event
		pending   []*messages.Event
		version   uint64
		lock      *sync.Mutex
	}

	Aggregate interface {
		Identifier() string
		WithStream(*Stream)
		Stream() *Stream
		Apply(*messages.Event) error
	}
)

// InitialiseAggregate will create a fresh stream on the aggregate.
func InitialiseAggregate(a Aggregate) {
	a.WithStream(&Stream{
		aggregate: a,
		events:    []*messages.Event{},
		pending:   []*messages.Event{},
		version:   0,
		lock:      &sync.Mutex{},
	})
}

// ReconstituteAggregate will apply all the events in the stream to
// the aggregate restoring it's state.
func ReconstituteAggregate(a Aggregate, e []*messages.Event) error {
	InitialiseAggregate(a)
	h := a.Stream()

	for _, e := range e {
		h.version = e.Version()
		if err := a.Apply(e); err != nil {
			return err
		}
	}
	return nil
}

// RecordThat records the event that has happened increasing
// the version of the aggregate.
func (r *Stream) RecordThat(e *messages.Event) error {
	e = e.WithVersion(r.version + 1)
	r.events = append(r.events, e)
	return r.aggregate.Apply(e)
}

func (r *Stream) popRecordedEvents() []*messages.Event {
	r.lock.Lock()
	defer func() {
		r.pending = []*messages.Event{}
		r.lock.Unlock()
	}()
	return r.pending
}
