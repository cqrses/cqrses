package bus

import (
	"context"

	"github.com/go-cqrses/cqrses/eventstore"
	"github.com/go-cqrses/cqrses/messages"
)

type (
	// publishingEventStore reads events going via appendTo
	// and will publish them on a success result.
	publishingEventStore struct {
		bus   *EventBus
		store eventstore.EventStore
	}
)

func EventStoreWithBus(bus *EventBus, store eventstore.EventStore) eventstore.EventStore {
	return &publishingEventStore{
		bus:   bus,
		store: store,
	}
}

// Load proxies to underlying store.
func (s *publishingEventStore) Load(ctx context.Context, streamName string, from, count uint64, matcher eventstore.MetadataMatcher) eventstore.StreamIterator {
	return s.store.Load(ctx, streamName, from, count, matcher)
}

// LoadReverse proxies to underlying store.
func (s *publishingEventStore) LoadReverse(ctx context.Context, streamName string, from, count uint64, matcher eventstore.MetadataMatcher) eventstore.StreamIterator {
	return s.store.LoadReverse(ctx, streamName, from, count, matcher)
}

// FetchStreamNames proxies to underlying store.
func (s *publishingEventStore) FetchStreamNames(ctx context.Context, filter string, matcher eventstore.MetadataMatcher, limit, offset uint64) ([]string, error) {
	return s.store.FetchStreamNames(ctx, filter, matcher, limit, offset)
}

// FetchStreamNamesRegex proxies to underlying store.
func (s *publishingEventStore) FetchStreamNamesRegex(ctx context.Context, filter string, matcher eventstore.MetadataMatcher, limit, offset uint64) ([]string, error) {
	return s.store.FetchStreamNamesRegex(ctx, filter, matcher, limit, offset)
}

// FetchStreamMetadata proxies to underlying store.
func (s *publishingEventStore) FetchStreamMetadata(ctx context.Context, streamName string) (eventstore.StreamMetadata, error) {
	return s.store.FetchStreamMetadata(ctx, streamName)
}

// Create proxies to underlying store.
func (s *publishingEventStore) Create(ctx context.Context, stream *eventstore.Stream) error {
	return s.store.Create(ctx, stream)
}

// AppendTo proxies to underlying store.
func (s *publishingEventStore) AppendTo(ctx context.Context, streamName string, events []*messages.Event) error {
	err := s.store.AppendTo(ctx, streamName, events)

	if err == nil {
		for _, e := range events {
			s.bus.Handle(ctx, e)
		}
	}

	return err
}

// Delete proxies to underlying store.
func (s *publishingEventStore) Delete(ctx context.Context, streamName string) error {
	return s.store.Delete(ctx, streamName)
}

// UpdateStreamMetadata proxies to underlying store.
func (s *publishingEventStore) UpdateStreamMetadata(ctx context.Context, streamName string, newMetadata eventstore.StreamMetadata) error {
	return s.store.UpdateStreamMetadata(ctx, streamName, newMetadata)
}
