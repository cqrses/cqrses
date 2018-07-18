package aggregate_test

import (
	"context"
	"testing"

	"github.com/go-cqrses/cqrses/aggregate"
	"github.com/go-cqrses/cqrses/eventstore"
	"github.com/go-cqrses/cqrses/eventstore/inmem"
	"github.com/go-cqrses/cqrses/messages"

	"github.com/stretchr/testify/assert"
)

type state struct {
	appliedCount int
}

func (s *state) Handle(_ context.Context, _ messages.Message, _ aggregate.EventRecorder) error {
	return nil
}

func (s *state) Apply(*messages.Event) error {
	s.appliedCount++
	return nil
}

func TestHistory(t *testing.T) {
	ctx := context.Background()
	es := inmem.New()
	es.Create(ctx, eventstore.EmptyStreamWithName("users"))
	aID := "1df0d42f-596c-4fbb-8d8b-363524d50195"

	{ // Write events to the store directly.
		as := &state{}
		h := aggregate.New(aID, es, "users", as)

		_ = h.RecordThat(ctx, "itHappened", map[string]interface{}{})
		_ = h.RecordThat(ctx, "itHappened", map[string]interface{}{})
		_ = h.RecordThat(ctx, "itHappened", map[string]interface{}{})
		_ = h.RecordThat(ctx, "itHappened", map[string]interface{}{})

		if err := h.Close(ctx); err != nil {
			t.Fatalf("unable to persist recorded itHappened events: %s", err)
		}
	}

	{ // Read events from the store.
		as := &state{}
		_, err := aggregate.Load(ctx, aID, es, "users", as)
		assert.Nil(t, err)
		assert.Equal(t, 4, as.appliedCount)
	}
}
