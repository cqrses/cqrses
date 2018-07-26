package projection

import (
	"context"
	"time"

	"github.com/go-cqrses/cqrses/messages"
)

const (
	// StatusRunning is set when the projector is running.
	StatusRunning Status = "running"
	// StatusStopping is set when the projector is stopping.
	StatusStopping Status = "stopping"
	// StatusDeleting is set when the projector is deleting.
	StatusDeleting Status = "deleting"
	// StatusResetting is set when the projector is resetting.
	StatusResetting Status = "resetting"
	// StatusIdle is set when the projector is idle.
	StatusIdle Status = "idle"
)

type (
	// Status of the projection.
	Status string

	// Handler should handle the event provided.
	Handler func(context.Context, messages.Message) error

	// Projector describes what should be possible by an implementor.
	Projector interface {
		// FromStream will limit the Projector to events from 1 stream.
		FromStream(streamName string) Projector
		// FromStreams will limit the Projector to events from many streams.
		FromStreams(streamNames []string) Projector
		// When the event with the event name is given the Handler will be called.
		When(eventName string, cb Handler) Projector
		// WhenAny event is given the Handler will be called.
		WhenAny(cb Handler) Projector
		// Stop will stop the processing of the events.
		Stop(ctx context.Context) error
		// Run will start the processing of the events.
		Run(ctx context.Context) error
	}

	// ProjectorOpt applies configuration to projector options.
	ProjectorOpt func(*ProjectorOpts) error

	// ProjectorOpts contains options for the projector.
	ProjectorOpts struct {
		// Sleep how long after reading all the events should we sleep before reading more.
		Sleep time.Duration
	}

	// Manager manages projections.
	Manager interface {
		// Create a new stream.
		Create(ctx context.Context, name string, options []ProjectorOpt) (Projector, error)

		// FetchProjections ...
		FetchProjectionNames(ctx context.Context, filter string, start, limit uint64) ([]string, error)
	}
)

// BuildOptionsFrom ...
func BuildOptionsFrom(opts []ProjectorOpt) (*ProjectorOpts, error) {
	out := &ProjectorOpts{
		Sleep: 5 * time.Second,
	}
	for _, opt := range opts {
		if err := opt(out); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// WithSleep will set the Sleep on the projector options.
func WithSleep(v time.Duration) ProjectorOpt {
	return func(o *ProjectorOpts) error {
		o.Sleep = v
		return nil
	}
}
