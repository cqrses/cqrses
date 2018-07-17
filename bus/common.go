package bus

import (
	"context"
	"errors"
	"fmt"

	"gopkg.in/cqrses/messages"
)

var (
	// ErrNoHandlerFound is returned from a command or query bus where
	// there is no handler found.
	ErrNoHandlerFound = errors.New("no command handler found")
)

type (
	// Handler handles messages of different kinds.
	Handler func(ctx context.Context, msg messages.Message) error

	// Error is returned from dispatch functions.
	Error struct {
		messageID   string
		messageName string
		original    error
	}
)

// Error returns an error description.
func (e *Error) Error() string {
	return fmt.Sprintf(
		"Error processing %s (id:%s): %s",
		e.messageName,
		e.messageID,
		e.original,
	)
}
