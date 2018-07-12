package bus

import (
	"context"
	"errors"
	"fmt"

	"gopkg.in/cqrses/messages"
)

var (
	ErrNoHandlerFound = errors.New("no command handler found")
)

type (
	Handler func(context.Context, messages.Message) error

	Error struct {
		messageID   string
		messageName string
		original    error
	}
)

func (e *Error) Error() string {
	return fmt.Sprintf(
		"Error processing %s (id:%s): %s",
		e.messageName,
		e.messageID,
		e.original,
	)
}
