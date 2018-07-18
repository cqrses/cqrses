package bus_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/go-cqrses/cqrses/bus"
	"github.com/go-cqrses/cqrses/messages"
)

func TestEventBus(t *testing.T) {
	sut := bus.NewEventBus()
	called := map[string]int{
		"hello-world":   0,
		"goodbye-world": 0,
	}
	h := func(_ context.Context, m messages.Message) error {
		called[m.MessageName()]++
		return nil
	}

	sut.Register(bus.MatchAny(), h)
	sut.Register(bus.MatchMessageNameRaw("hello-world"), h)

	e1 := messages.NewEvent("1", "hello-world", map[string]interface{}{}, map[string]interface{}{}, 0, time.Now())
	e2 := messages.NewEvent("2", "goodbye-world", map[string]interface{}{}, map[string]interface{}{}, 0, time.Now())

	sut.Handle(context.Background(), e1)
	sut.Handle(context.Background(), e2)

	assert.Equal(t, 2, called["hello-world"])
	assert.Equal(t, 1, called["goodbye-world"])
}
