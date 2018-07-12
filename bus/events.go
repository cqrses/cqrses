package bus

import (
	"context"

	"gopkg.in/cqrses/messages"
)

type (
	MessageMatcher func(messages.Message) bool

	eventBusHandler struct {
		matches MessageMatcher
		handler Handler
	}

	EventBus struct {
		handlers []*eventBusHandler
	}
)

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: []*eventBusHandler{},
	}
}

func (b *EventBus) Register(m MessageMatcher, h Handler) {
	b.handlers = append(b.handlers, &eventBusHandler{m, h})
}

func (c *EventBus) Handle(ctx context.Context, m messages.Message) error {
	for _, h := range c.handlers {
		if h.matches(m) {
			_ = h.handler(ctx, m) // we ignore error as this is an event rather then a command
		}
	}
	return nil
}

func MatchAny() MessageMatcher {
	return func(messages.Message) bool {
		return true
	}
}

func MatchMessageName(m messages.Message) MessageMatcher {
	return MatchMessageNameRaw(m.MessageName())
}

func MatchMessageNameRaw(n string) MessageMatcher {
	return func(m messages.Message) bool {
		return m.MessageName() == n
	}
}
