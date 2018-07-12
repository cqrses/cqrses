package aggregate_test

import (
	"testing"

	"gopkg.in/cqrses"
)

type stubAggregate struct {
	id     string
	stream *cqrses.Stream
}

func (a *stubAggregate) Identifier() string {
	return a.id
}

func (a *stubAggregate) WithStream(h *cqrses.Stream) {
	a.stream = h
}

func (a *stubAggregate) Stream() *cqrses.Stream {
	return a.stream
}

func (a *stubAggregate) Apply(e *cqrses.Event) error {
	switch e.MessageName {

	}
	return nil
}

func TestAggregate(t *testing.T) {
	// a := &stubAggregate{}
	// h := cqrses.NewStream()
	// a.WithStream()
}
