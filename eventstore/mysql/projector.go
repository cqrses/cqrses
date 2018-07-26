package mysql

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/go-cqrses/cqrses/eventstore"
	"github.com/go-cqrses/cqrses/messages"
	"github.com/go-cqrses/cqrses/projection"
	"github.com/pkg/errors"
)

type (
	// StreamProjection ...
	StreamProjection struct {
		name        string
		es          *EventStore
		streamNames []string
		opts        *projection.ProjectorOpts
		handlers    map[string][]projection.Handler
		any         []projection.Handler
		modLock     *sync.Mutex
		close       chan struct{}
	}

	// ProjectionManager ...
	ProjectionManager struct {
		es *EventStore
	}
)

// NewProjectionManager will get a projection manager that uses the MySQL backend
// to store projection states.
func NewProjectionManager(es *EventStore) projection.Manager {
	return &ProjectionManager{
		es: es,
	}
}

// Create ...
func (m *ProjectionManager) Create(ctx context.Context, name string, opts []projection.ProjectorOpt) (projection.Projector, error) {
	options, err := projection.BuildOptionsFrom(opts)
	if err != nil {
		return nil, err
	}
	return &StreamProjection{
		name:        name,
		es:          m.es,
		streamNames: []string{},
		opts:        options,
		handlers:    map[string][]projection.Handler{},
		any:         []projection.Handler{},
		close:       make(chan struct{}, 1),
		modLock:     &sync.Mutex{},
	}, nil
}

// FetchProjectionNames ...
func (m *ProjectionManager) FetchProjectionNames(ctx context.Context, filter string, start, limit uint64) ([]string, error) {
	var res *sql.Rows
	var err error
	out := make([]string, 0, limit)

	if len(filter) == 0 {
		res, err = m.es.db.QueryContext(ctx, "select `name` from projections limit ?,?", start, limit)
	} else {
		res, err = m.es.db.QueryContext(ctx, "select `name` from where name like ? limit ?,?", filter, start, limit)
	}

	if err == nil {
		for res.Next() {
			var name string
			if err = res.Scan(&name); err != nil {
				return out, err
			}
			out = append(out, name)
		}
	}

	return out, err
}

// FromStream will limit the Projector to events from 1 stream.
func (p *StreamProjection) FromStream(streamName string) projection.Projector {
	p.streamNames = []string{streamName}
	return p
}

// FromStreams will limit the Projector to events from many streams.
func (p *StreamProjection) FromStreams(streamNames []string) projection.Projector {
	if true {
		panic("MySQL projector currently only supports 1 stream.")
	}
	p.streamNames = streamNames
	return p
}

// When the event with the event name is given the callback will be called.
func (p *StreamProjection) When(eventName string, cb projection.Handler) projection.Projector {
	p.modLock.Lock()
	defer p.modLock.Unlock()

	h, ok := p.handlers[eventName]
	if ok {
		h = append(h, cb)
	} else {
		h = []projection.Handler{cb}
	}

	p.handlers[eventName] = h

	return p
}

// WhenAny event is given the callback will be called.
func (p *StreamProjection) WhenAny(cb projection.Handler) projection.Projector {
	p.modLock.Lock()
	defer p.modLock.Unlock()

	p.any = append(p.any, cb)

	return p
}

// Stop will stop the processing of the events.
func (p *StreamProjection) Stop(ctx context.Context) error {
	p.close <- struct{}{}
	return nil
}

// Run will start the processing of the events.
func (p *StreamProjection) Run(ctx context.Context) error {
	if err := p.ensureProjectionExists(ctx); err != nil {
		return err
	}

	cEvents := make(chan *messages.Event, 1)
	cErr := make(chan error, 1)

	go p.handleEvents(ctx, cEvents)

	go func() {
		select {
		case <-p.close:
			cErr <- errors.New("projection was closed")
			return
		case <-time.After(p.opts.Sleep):
			if err := p.retreiveEventsFromStream(ctx, cEvents); err != nil {
				cErr <- err
				return
			}
		}
	}()

	err := <-cErr
	close(cEvents)
	return err
}

func (p *StreamProjection) ensureProjectionExists(ctx context.Context) error {
	var total uint32
	res := p.es.db.QueryRowContext(ctx, "select count(*) from projections where name = ?", p.name)
	if err := res.Scan(&total); err != nil {
		return errors.Wrap(err, "unable to determine if projection exists")
	}

	if total == 1 {
		return nil
	}

	_, err := p.es.db.ExecContext(ctx, "insert projections (name, position, status) values (?, ?, ?)", p.name, 0, projection.StatusIdle)
	return errors.Wrap(err, "unable to store projection in projections store")
}

func (p *StreamProjection) retreiveEventsFromStream(ctx context.Context, cEvents chan *messages.Event) error {
	if snl := len(p.streamNames); snl != 1 {
		return errors.Errorf("expected 1 stream to query got %d", snl)
	}

	var position uint64
	res := p.es.db.QueryRowContext(ctx, "select position from projections where name = ?", p.name)
	if err := res.Scan(&position); err != nil {
		return err
	}

	it := p.es.Load(ctx, p.streamNames[0], position, 1000, eventstore.MetadataMatcher{})
	for {
		if err := it.Next(ctx); err != nil {
			if err == eventstore.EOF {
				return nil
			}

			return err
		}

		cEvents <- it.Current()

		if _, err := p.es.db.ExecContext(ctx, "update projections set position = position+1 where name = ?", p.name); err != nil {
			return errors.Wrap(err, "unable to set projection position, this will lead to duplicate events when next run")
		}
	}
}

func (p *StreamProjection) handleEvents(ctx context.Context, cEvents chan *messages.Event) {
	for {
		select {
		case event := <-cEvents:
			for _, h := range p.any {
				h(ctx, event)
			}

			if hs, ok := p.handlers[event.MessageName()]; ok {
				for _, h := range hs {
					h(ctx, event)
				}
			}
		}
	}
}
