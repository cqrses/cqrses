package bus_test

import (
	"context"
	"testing"
	"time"

	"gopkg.in/cqrses/bus"
	"gopkg.in/cqrses/messages"

	"github.com/stretchr/testify/assert"
)

func TestCommandBus(t *testing.T) {
	sut := bus.NewCommandBus()
	called := false

	sut.Register("test", func(ctx context.Context, msg messages.Message) error {
		called = true
		return nil
	})

	sut.Handle(context.Background(), messages.NewCommand("123", "test", map[string]interface{}{}, map[string]interface{}{}, 0, time.Now()))

	assert.True(t, called)
}

func TestCommandBusMiddleware(t *testing.T) {
	sut := bus.NewCommandBus()
	called := false

	sut.Register("test", func(ctx context.Context, msg messages.Message) error {
		called = true

		if v, ok := ctx.Value("from_test_1").(string); !ok || v != "correct" {
			t.Errorf("did not get valid value from_test_1: %+v", ctx)
		}

		if v, ok := ctx.Value("from_test_2").(string); !ok || v != "correct" {
			t.Errorf("did not get valid value from_test_2: %+v", ctx)
		}

		return nil
	})

	sut.PushMiddleware(func(ctx context.Context, msg messages.Message, next func(context.Context, messages.Message) error) error {
		return next(context.WithValue(ctx, "from_test_1", "correct"), msg)
	})

	sut.PushMiddleware(func(ctx context.Context, msg messages.Message, next func(context.Context, messages.Message) error) error {
		return next(context.WithValue(ctx, "from_test_2", "correct"), msg)
	})

	sut.Handle(context.Background(), messages.NewCommand("123", "test", map[string]interface{}{}, map[string]interface{}{}, 0, time.Now()))

	assert.True(t, called)
}
