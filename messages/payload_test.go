package messages_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/go-cqrses/cqrses/messages"
)

type myPayload struct {
	todo string
	done bool
}

func (p *myPayload) Payload() map[string]interface{} {
	return map[string]interface{}{
		"todo": p.todo,
		"done": p.done,
	}
}

func (p *myPayload) FromPayload(in map[string]interface{}) {
	p.todo, _ = in["todo"].(string)
	p.done, _ = in["done"].(bool)
}

func Test(t *testing.T) {
	todoMsg := "Get tests passing"
	done := true
	msg := messages.NewEvent(
		"test",
		"AddTodo",
		messages.BuildPayload(&myPayload{
			todo: todoMsg,
			done: done,
		}),
		map[string]interface{}{},
		0,
		time.Now(),
	)

	out := &myPayload{}
	messages.ReadPayload(msg, out)

	assert.Equal(t, todoMsg, out.todo)
	assert.Equal(t, done, out.done)
}
