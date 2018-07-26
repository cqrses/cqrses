package projection

import (
	"context"
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

	// StreamPositions provides the stream name and current position in the projection.
	StreamPositions map[string]uint64

	// Manager manages projections.
	Manager interface {
		// Create a new stream.
		Create(ctx context.Context, name string, options []ProjectorOpt) (Projector, error)

		// Delete will remove the projection from the projections store.
		Delete(ctx context.Context, projectionName string) error

		// Reset will reset the position of the stream to 0.
		Reset(ctx context.Context, projectionName string) error

		// Stop will set the status to stop, the running projection should handle this.
		Stop(ctx context.Context, projectionName string) error

		// FetchProjections grabs projections matching the filter.
		FetchProjectionNames(ctx context.Context, filter string, start, limit uint64) ([]string, error)

		// FetchPojectionStatus will return the status of a projection.
		FetchPojectionStatus(ctx context.Context, projectionName string) (Status, error)

		// FetchPojectionStreamPositions will return the status of a projection.
		FetchPojectionStreamPositions(ctx context.Context, projectionName string) (StreamPositions, error)
	}

	// HasProjectionManager should be implemented by an event store that can return a projection
	// manager. We have decided to not put this in the event store package to keep is separated,
	// however your event store implementation should have one.
	HasProjectionManager interface {
		GetProjectionManager([]ProjectorOpt) Manager
	}
)
