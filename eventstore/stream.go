package eventstore

import (
	"context"
	"io"
	"regexp"

	"github.com/go-cqrses/cqrses/messages"
)

const (
	// MatchOpEq will make sure the value is exact.
	MatchOpEq matchOp = "eq"
	// MatchOpIn will check the value is in the given condition values.
	MatchOpIn matchOp = "in"
	// MatchOpNotIn will check the value is not in the given condition values.
	MatchOpNotIn matchOp = "not_in"
	// MatchOpRegex will check the value against the regular expression condition.
	MatchOpRegex matchOp = "regex"
)

var (
	// EOF is returned from a stream when you're at the end of the stream.
	EOF = io.EOF
)

type (
	matchOp string

	// StreamMetadata should have string keys and string values.
	StreamMetadata map[string]string

	// MetadataMatcherCondition ...
	MetadataMatcherCondition struct {
		Operation matchOp
		Values    []string
	}

	// MetadataMatcher will match a string ID
	// with the value using the operation provided.
	MetadataMatcher map[string]MetadataMatcherCondition

	// Stream is a struct used to create an initial stream, also
	// used by the in memory store.
	Stream struct {
		Name     string
		Metadata StreamMetadata
		Events   []*messages.Event
	}

	// StreamIterator iterates over events in a stream.
	StreamIterator interface {
		// Current will return the current event in the stream.
		Current() *messages.Event

		// Next will move the cursor forward.
		Next(context.Context) error

		// Rewind will go back to the begining of the stream.
		Rewind()

		// Close will clean up resources.
		Close()
	}
)

// EmptyStreamWithName returns a new empty stream with the name
// provided.
func EmptyStreamWithName(name string) *Stream {
	return &Stream{
		Name:     name,
		Metadata: StreamMetadata{},
		Events:   []*messages.Event{},
	}
}

// NewStreamWithName returns a new stream with the data provided.
func NewStreamWithName(name string, metadata StreamMetadata, events []*messages.Event) *Stream {
	return &Stream{
		Name:     name,
		Metadata: metadata,
		Events:   events,
	}
}

// MatchStreamMetadata will test all keys that require their
// values testing exist inside the metadata and then check
// the value is valid.
func (m MetadataMatcher) MatchStreamMetadata(in StreamMetadata) bool {
	for k, matcher := range m {
		v, ok := in[k]

		if !ok || !matcher.Match(v) {
			return false
		}
	}
	return true
}

// MatchEventMetadata will test all keys that require their
// values testing exist inside the metadata and then check
// the value is valid. At this moment in time it will only
// test string values.
func (m MetadataMatcher) MatchEventMetadata(in map[string]interface{}) bool {
	for k, matcher := range m {
		v, ok := in[k].(string)

		if !ok || !matcher.Match(v) {
			return false
		}
	}
	return true
}

// Match will test the provided value again the conditions
// given in the metadata matcher condition.
func (m MetadataMatcherCondition) Match(v string) bool {
	switch m.Operation {
	case MatchOpIn:
		for _, mv := range m.Values {
			if mv == v {
				return true
			}
		}
		return false
	case MatchOpNotIn:
		for _, mv := range m.Values {
			if mv == v {
				return false
			}
		}
		return true
	case MatchOpRegex:
		// We expect only 1 possible value to match against.
		if len(m.Values) != 1 {
			return false
		}

		ok, _ := regexp.Match(m.Values[0], []byte(v))
		return ok
	default: // MatchOpEq
		// We expect only 1 possible value to match against.
		if len(m.Values) != 1 {
			return false
		}

		return m.Values[0] == v
	}
}
