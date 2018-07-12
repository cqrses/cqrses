package bus

import (
	"context"
	"errors"

	"gopkg.in/cqrses/messages"
)

var (
	ErrCommandAlreadyRegistered = errors.New("command already registered")
)

type (
	CommandBus struct {
		handlers map[string]Handler
	}
)

func (b *CommandBus) Register(n string, h Handler) error {
	if _, ok := b.handlers[n]; ok {
		return ErrCommandAlreadyRegistered
	}

	b.handlers[n] = h

	return nil
}

func (c *CommandBus) Handle(ctx context.Context, m messages.Message) error {
	h, ok := c.handlers[m.MessageName()]
	if !ok {
		return ErrNoHandlerFound
	}

	if err := h(ctx, m); err != nil {
		return &Error{original: err}
	}

	return nil
}
