package bus

import (
	"context"
	"errors"
	"sync"

	"github.com/go-cqrses/cqrses/messages"
)

var (
	ErrCommandAlreadyRegistered = errors.New("command already registered")
)

type (
	// CommandBus can handle the dispatching of commands.
	CommandBus struct {
		handlers   map[string]Handler
		middleware []CommandBusMiddleware
		lock       *sync.Mutex
	}

	// CommandBusMiddleware can guard and alter the context or message that is about to be handled.
	CommandBusMiddleware func(ctx context.Context, msg messages.Message, next func(context.Context, messages.Message) error) error
)

// NewCommandBus returns a new initialised command bus.
func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers:   map[string]Handler{},
		middleware: []CommandBusMiddleware{},
		lock:       &sync.Mutex{},
	}
}

// Register a handler for the message name provided.
func (c *CommandBus) Register(n string, h Handler) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.handlers[n]; ok {
		return ErrCommandAlreadyRegistered
	}

	c.handlers[n] = h

	return nil
}

// PushMiddleware to the middleware slice.
func (c *CommandBus) PushMiddleware(in CommandBusMiddleware) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.middleware = append(c.middleware, in)
}

// Handle a command and return an error if there is any.
func (c *CommandBus) Handle(ctx context.Context, msg messages.Message) error {
	n := func(finalCtx context.Context, finalMsg messages.Message) error {
		h, ok := c.handlers[finalMsg.MessageName()]
		if !ok {
			return ErrNoHandlerFound
		}
		return h(finalCtx, finalMsg)
	}

	c.lock.Lock()
	m := c.middleware[:]
	c.lock.Unlock()

	var err error
	if len(m) > 0 {
		ml := len(m)
		mc := 0
		var mn func(mCtx context.Context, mMsg messages.Message) error
		mn = func(mCtx context.Context, mMsg messages.Message) error {
			mc++
			if mc == ml {
				return n(mCtx, mMsg)
			}

			return m[mc](mCtx, mMsg, mn)
		}

		err = m[mc](ctx, msg, mn)
	} else {
		err = n(ctx, msg)
	}

	if err != nil {
		return &Error{original: err}
	}

	return nil
}
