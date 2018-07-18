package inmem_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/go-cqrses/cqrses/eventstore"
	"github.com/go-cqrses/cqrses/eventstore/inmem"
	"github.com/go-cqrses/cqrses/messages"
)

func TestStore(t *testing.T) {
	ctx := context.Background()
	store := inmem.New()

	if err := store.Create(ctx, eventstore.EmptyStreamWithName("todo")); err != nil {
		t.Fatalf("unable to create stream: %s", err)
	}

	{ // Check to make sure we get an error returned when a stream does not exist and we attempt to iterate over it.
		stream := store.Load(ctx, "na", 0, 1000, eventstore.MetadataMatcher{})
		assert.Equal(t, eventstore.ErrStreamDoesNotExist, stream.Next())
	}

	{ // Check fetching stream names.
		{ // Make sure our stream is in the streams available.
			names, err := store.FetchStreamNames(ctx, "", eventstore.MetadataMatcher{}, 10, 0)
			assert.Nil(t, err)
			assert.Len(t, names, 1)
			assert.Equal(t, "todo", names[0])
		}

		{ // Attempt to get a stream with filter containing partial name of a stream.
			names, err := store.FetchStreamNames(ctx, "to", eventstore.MetadataMatcher{}, 10, 0)
			assert.Nil(t, err)
			assert.Len(t, names, 1)
			assert.Equal(t, "todo", names[0])
		}

		{ // Attempt to get a stream with a non existing partial of a stream.
			names, err := store.FetchStreamNames(ctx, "idonotexists", eventstore.MetadataMatcher{}, 10, 0)
			assert.Nil(t, err)
			assert.Len(t, names, 0)
		}
	}

	err := store.AppendTo(
		ctx,
		"todo",
		[]*messages.Event{
			messages.NewEvent("ev1", "TodoAdded", map[string]interface{}{}, map[string]interface{}{}, 0, time.Now()),
			messages.NewEvent("ev2", "TodoAdded", map[string]interface{}{}, map[string]interface{}{}, 0, time.Now()),
			messages.NewEvent("ev3", "TodoAdded", map[string]interface{}{}, map[string]interface{}{}, 0, time.Now()),
			messages.NewEvent("ev4", "TodoAdded", map[string]interface{}{}, map[string]interface{}{}, 0, time.Now()),
			messages.NewEvent("ev5", "TodoAdded", map[string]interface{}{}, map[string]interface{}{}, 0, time.Now()),
			messages.NewEvent("ev6", "TodoAdded", map[string]interface{}{}, map[string]interface{}{}, 0, time.Now()),
			messages.NewEvent("ev7", "TodoAdded", map[string]interface{}{}, map[string]interface{}{}, 0, time.Now()),
			messages.NewEvent("ev8", "TodoAdded", map[string]interface{}{}, map[string]interface{}{}, 0, time.Now()),
		},
	)
	assert.Nil(t, err)

	{ // Load a stream and check the contents is what we are expecting.
		stream := store.Load(ctx, "todo", 0, 1000, eventstore.MetadataMatcher{})
		expecteds := []string{"ev1", "ev2", "ev3", "ev4", "ev5", "ev6", "ev7", "ev8"}
		for i := 0; i < len(expecteds); i++ {
			if err := stream.Next(); err != nil {
				t.Fatalf("unable to get next item in stream: %s", err)
			}
			event := stream.Current()
			assert.Equal(t, expecteds[i], event.MessageID())
		}

		if err := stream.Next(); err != eventstore.EOF {
			t.Errorf("expected eventstore.EOF but got: %+v", err)
		}

		// Test rewinding and doing the same again.

		stream.Rewind()
		for i := 0; i < len(expecteds); i++ {
			if err := stream.Next(); err != nil {
				t.Fatalf("unable to get next item in stream: %s", err)
			}
			event := stream.Current()
			assert.Equal(t, expecteds[i], event.MessageID())
		}

		if err := stream.Next(); err != eventstore.EOF {
			t.Errorf("expected eventstore.EOF but got: %+v", err)
		}

		stream.Close()
	}

	{ // Load a stream with a from number and a limit.
		stream := store.Load(ctx, "todo", 2, 5, eventstore.MetadataMatcher{})
		expecteds := []string{"ev3", "ev4", "ev5", "ev6", "ev7"}
		for i := 0; i < len(expecteds); i++ {
			if err := stream.Next(); err != nil {
				t.Fatalf("unable to get next item in stream: %s", err)
			}
			event := stream.Current()
			assert.Equal(t, expecteds[i], event.MessageID())
		}

		if err := stream.Next(); err != eventstore.EOF {
			t.Errorf("expected eventstore.EOF but got: %+v", err)
		}
		stream.Close()
	}

	{ // Load a stream in reverse.
		stream := store.LoadReverse(ctx, "todo", 2, 5, eventstore.MetadataMatcher{})
		expecteds := []string{"ev5", "ev4", "ev3", "ev2", "ev1"}
		for i := 0; i < len(expecteds); i++ {
			if err := stream.Next(); err != nil {
				t.Fatalf("unable to get next item in stream: %s", err)
			}
			event := stream.Current()
			assert.Equal(t, expecteds[i], event.MessageID())
		}

		if err := stream.Next(); err != eventstore.EOF {
			t.Errorf("expected eventstore.EOF but got: %+v", err)
		}
		stream.Close()
	}

}
