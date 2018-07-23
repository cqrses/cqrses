package aggregate

import (
	"context"
	"errors"
	"sync"

	"github.com/go-cqrses/cqrses/esbridge"
	"github.com/go-cqrses/cqrses/eventstore"
	"github.com/go-cqrses/cqrses/messages"
)

type (
	Command interface {
		AggregateID() string
	}

	StateFactory func() State

	// Aggregate manages the history.
	Aggregate struct {
		// The ID of the aggregate we are dealing with.
		aggregateID string
		// The event store to store events in.
		store eventstore.EventStore
		// The name of the stream to store events in.
		streamName string
		// A slice of events pending to go into the event store.
		pending []*messages.Event
		// The current version.
		version uint64
		// Aggregate state.
		state State
		// A lock used to ensure no race conditions within the instance.
		lock *sync.Mutex
	}

	// EventRecorder will record aggregate events.
	EventRecorder func(eventName string, data interface{}) error

	// State is the user land implementation of the aggregate.
	State interface {
		Handle(context.Context, messages.Message, EventRecorder) error
		Apply(*messages.Event) error
	}
)

func Make(af StateFactory, streamName string) func(ctx context.Context, msg messages.Message) error {
	return func(ctx context.Context, msg messages.Message) error {
		cmd, ok := msg.Data().(Command)
		if !ok {
			return errors.New("this command payload cannot be handled by cqrses.aggregate.Make")
		}

		es := esbridge.MustGetEventStoreFromContext(ctx)
		ag := af()
		history, err := Load(ctx, cmd.AggregateID(), es, streamName, ag)

		if err != nil {
			return err
		}

		if err := history.Handle(ctx, msg); err != nil {
			return err
		}

		return history.Close(ctx)
	}
}

// New should be used when intiailising an aggregate.
func New(aID string, store eventstore.EventStore, streamName string, state State) *Aggregate {
	return &Aggregate{
		aggregateID: aID,
		store:       store,
		streamName:  streamName,
		pending:     []*messages.Event{},
		version:     0,
		state:       state,
		lock:        &sync.Mutex{},
	}
}

// Load will get an aggregates Aggregate, reconstitue the aggregate using the
// event handler provided and then returning the Aggregate to allow adding more events.
func Load(ctx context.Context, aID string, store eventstore.EventStore, streamName string, state State) (*Aggregate, error) {
	a := &Aggregate{
		aggregateID: aID,
		store:       store,
		streamName:  streamName,
		pending:     []*messages.Event{},
		version:     0,
		state:       state,
		lock:        &sync.Mutex{},
	}

	events := store.Load(ctx, streamName, 0, 0, eventstore.MetadataMatcher{
		"aggregate_id": eventstore.MetadataMatcherCondition{
			Operation: eventstore.MatchOpEq,
			Values:    []string{aID},
		},
	})
	defer events.Close()

	for {
		if err := events.Next(ctx); err != nil {
			if err == eventstore.EOF {
				break
			}

			return nil, err
		}

		event := events.Current()
		if err := a.state.Apply(event); err != nil {
			return nil, err
		}
		a.version = event.Version()
	}

	return a, nil
}

// Handle will execute the callback with the message and persist any events.
func (h *Aggregate) Handle(ctx context.Context, msg messages.Message) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	return h.state.Handle(ctx, msg, func(eventName string, data interface{}) error {
		return h.record(ctx, eventName, data)
	})
}

// RecordThat an event that has happened increasing the version of the aggregate.
func (h *Aggregate) RecordThat(ctx context.Context, eventName string, data interface{}) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	return h.record(ctx, eventName, data)
}

func (h *Aggregate) record(ctx context.Context, eventName string, data interface{}) error {
	h.version++
	event := messages.NewAggregateEvent(ctx, h.aggregateID, h.version, eventName, data)
	h.pending = append(h.pending, event)
	return h.state.Apply(event)
}

// Close will persist any pending events, returning an error if anything failed,
// if an error is returned all pending events will be missing still.
func (h *Aggregate) Close(ctx context.Context) error {
	h.lock.Lock()
	defer func() {
		h.pending = []*messages.Event{}
		h.lock.Unlock()
	}()

	return h.store.AppendTo(ctx, h.streamName, h.pending)
}
