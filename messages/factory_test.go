package messages_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/cqrses/messages"
)

func TestJSONMessageFactory(t *testing.T) {
	sut := &messages.JSONMessageFactory{}

	event := messages.NewEvent(
		"hello-world",
		"test.event",
		map[string]interface{}{
			"hello": "world",
		},
		map[string]interface{}{
			"1+1": "2",
		},
		1,
		time.Now(),
	)

	in, err := sut.Serialize(event)
	assert.Nil(t, err)
	assert.NotEmpty(t, in)

	out, err := sut.Unserialize(in)
	assert.Nil(t, err)
	assert.Equal(t, event.MessageID(), out.MessageID())
	assert.Equal(t, event.MessageName(), out.MessageName())
	assert.Equal(t, event.Data(), out.Data())
	assert.Equal(t, event.Metadata(), out.Metadata())
	assert.Equal(t, event.Version(), out.Version())
	assert.Equal(t, event.Created().Format(time.RFC3339Nano), out.Created().Format(time.RFC3339Nano))
}
