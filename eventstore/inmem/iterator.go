package inmem

import (
	"context"

	"gopkg.in/cqrses/eventstore"
	"gopkg.in/cqrses/messages"
)

type (
	// StreamIterator which iterates over events in memory.
	StreamIterator struct {
		Error    error
		Events   []*messages.Event
		position int
	}
)

// Current will return the current event in the stream.
func (s *StreamIterator) Current() *messages.Event {
	return s.Events[s.position]
}

// Next will move the cursor forward.
func (s *StreamIterator) Next(context.Context) error {
	if s.Error != nil {
		return s.Error
	}

	s.position++

	if s.position >= len(s.Events) {
		return eventstore.EOF
	}

	return nil
}

// Rewind will go back to the begining of the stream.
func (s *StreamIterator) Rewind() {
	s.position = -1
}

// Close will clean up resources.
func (s *StreamIterator) Close() {
	s.position = -1
	s.Events = nil
}
