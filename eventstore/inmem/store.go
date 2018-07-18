package inmem

import (
	"context"
	"regexp"
	"strings"
	"sync"

	"github.com/go-cqrses/cqrses/eventstore"
	"github.com/go-cqrses/cqrses/messages"
)

type (
	// EventStore that stores stream in memory.
	EventStore struct {
		streams map[string]*eventstore.Stream
		lock    *sync.Mutex
	}
)

// New returns a new in memory event store.
func New() *EventStore {
	return &EventStore{
		streams: map[string]*eventstore.Stream{},
		lock:    &sync.Mutex{},
	}
}

// Load events from the given stream name.
func (s EventStore) Load(ctx context.Context, streamName string, from, count uint64, matcher eventstore.MetadataMatcher) eventstore.StreamIterator {
	stream, ok := s.streams[streamName]

	if !ok {
		return &StreamIterator{Error: eventstore.ErrStreamDoesNotExist}
	}

	total := uint64(len(stream.Events))
	events := make([]*messages.Event, 0, count)
	taken := uint64(0)

	if from < total {
		for _, e := range stream.Events[from:] {
			if count > 0 && taken == count {
				break
			}

			if matcher.MatchEventMetadata(e.Metadata()) {
				events = append(events, e)
				taken++
			}
		}
	}

	return &StreamIterator{
		Events:   events,
		position: -1,
	}
}

// LoadReverse Loads events from the given stream name in reverse.
func (s EventStore) LoadReverse(ctx context.Context, streamName string, from, count uint64, matcher eventstore.MetadataMatcher) eventstore.StreamIterator {
	stream, ok := s.streams[streamName]

	if !ok {
		return &StreamIterator{Error: eventstore.ErrStreamDoesNotExist}
	}

	total := uint64(len(stream.Events))
	events := make([]*messages.Event, 0, count)
	skipped := uint64(0)
	taken := uint64(0)

	if from < total {

		for i := total - 1; i >= 0; i-- {
			if count > 0 && taken == count {
				break
			}

			if matcher.MatchEventMetadata(stream.Events[i].Metadata()) {
				if skipped > from {
					events = append(events, stream.Events[i])
					taken++
				} else {
					skipped++
				}
			}
		}
	}

	return &StreamIterator{
		Events:   events,
		position: -1,
	}
}

// FetchStreamNames gets  stream names that match the filter.
func (s EventStore) FetchStreamNames(ctx context.Context, filter string, matcher eventstore.MetadataMatcher, limit, offset uint64) ([]string, error) {
	sn := make([]string, 0, limit)
	i := uint64(0)
	for k, stream := range s.streams {
		if !strings.Contains(k, filter) {
			continue
		}

		if !matcher.MatchStreamMetadata(stream.Metadata) {
			continue
		}

		i++

		if offset >= i {
			continue
		}

		sn = append(sn, k)

		if limit == i {
			break
		}
	}
	return sn, nil
}

// FetchStreamNamesRegex gets stream names that match the regex filter.
func (s EventStore) FetchStreamNamesRegex(ctx context.Context, filter string, matcher eventstore.MetadataMatcher, limit, offset uint64) ([]string, error) {
	rx := regexp.MustCompile(filter)
	sn := make([]string, 0, limit)
	i := uint64(0)
	for k, stream := range s.streams {
		if !rx.MatchString(k) {
			continue
		}

		if !matcher.MatchStreamMetadata(stream.Metadata) {
			continue
		}

		i++

		if offset >= i {
			continue
		}

		sn = append(sn, k)

		if limit == i {
			break
		}
	}
	return sn, nil
}

// FetchStreamMetadata gets the metadata about a stream.
func (s EventStore) FetchStreamMetadata(ctx context.Context, streamName string) (eventstore.StreamMetadata, error) {
	stream, ok := s.streams[streamName]
	if !ok {
		return eventstore.StreamMetadata{}, eventstore.ErrStreamDoesNotExist
	}

	return stream.Metadata, nil
}

// Create will create the stream with the name and metadata provided.
func (s EventStore) Create(ctx context.Context, stream *eventstore.Stream) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, ok := s.streams[stream.Name]
	if !ok {
		return eventstore.ErrStreamAlreadyExists
	}

	s.streams[stream.Name] = stream
	return nil
}

// AppendTo will append events to the stream.
func (s EventStore) AppendTo(ctx context.Context, streamName string, events []*messages.Event) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	stream, ok := s.streams[streamName]
	if !ok {
		return eventstore.ErrStreamDoesNotExist
	}

	stream.Events = append(stream.Events, events...)

	return nil
}

// Delete will remove the stream.
func (s EventStore) Delete(ctx context.Context, streamName string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.streams, streamName)
	return nil
}

// UpdateStreamMetadata sets the metadata for the given stream name.
func (s EventStore) UpdateStreamMetadata(ctx context.Context, streamName string, newMetadata eventstore.StreamMetadata) error {
	stream, ok := s.streams[streamName]
	if !ok {
		return eventstore.ErrStreamDoesNotExist
	}

	stream.Metadata = newMetadata
	return nil
}
