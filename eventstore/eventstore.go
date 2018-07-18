package eventstore

import (
	"context"
	"errors"

	"github.com/go-cqrses/cqrses/messages"
)

var (
	// ErrStreamDoesNotExist is returned when attempting to read from
	// a stream that does not exist.
	ErrStreamDoesNotExist = errors.New("stream does not exist")

	// ErrStreamAlreadyExists is returned when attempting to create a
	// stream that does already exists.
	ErrStreamAlreadyExists = errors.New("stream already exists")
)

type (
	// ReadOnlyEventStore contains the methods to read from an event store.
	ReadOnlyEventStore interface {
		// Load events from the given stream name.
		Load(ctx context.Context, streamName string, from, count uint64, matcher MetadataMatcher) StreamIterator

		// LoadReverse Loads events from the given stream name in reverse.
		LoadReverse(ctx context.Context, streamName string, from, count uint64, matcher MetadataMatcher) StreamIterator

		// FetchStreamNames gets  stream names that match the filter.
		FetchStreamNames(ctx context.Context, filter string, matcher MetadataMatcher, limit, offset uint64) ([]string, error)

		// FetchStreamNamesRegex gets stream names that match the regex filter.
		FetchStreamNamesRegex(ctx context.Context, filter string, matcher MetadataMatcher, limit, offset uint64) ([]string, error)

		// FetchStreamMetadata gets the metadata about a stream.
		FetchStreamMetadata(ctx context.Context, streamName string) (StreamMetadata, error)
	}

	// EventStore contains the methods to read and write to an event store.
	EventStore interface {
		ReadOnlyEventStore

		// Create will create the stream with the name and metadata provided.
		Create(ctx context.Context, stream *Stream) error

		// AppendTo will append events to the stream.
		AppendTo(ctx context.Context, streamName string, events []*messages.Event) error

		// Delete will remove the stream.
		Delete(ctx context.Context, streamName string) error

		// UpdateStreamMetadata sets the metadata for the given stream name.
		UpdateStreamMetadata(ctx context.Context, streamName string, newMetadata StreamMetadata) error
	}
)
