package bus

import (
	"context"

	"github.com/go-cqrses/cqrses/eventstore"
	"github.com/go-cqrses/cqrses/messages"
)

type (
	// MessageMatcher decides whether or not a messages should be
	// handled by the handler provided.
	MessageMatcher func(messages.Message) bool

	eventBusHandler struct {
		matches MessageMatcher
		handler Handler
	}

	// EventBus is used to dispatch messages to various handlers.
	EventBus struct {
		handlers []*eventBusHandler
	}
)

// NewEventBus returns a new initialised event bus.
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: []*eventBusHandler{},
	}
}

// Register a handler that will be called if a dispatched event returns
// a positive result from the message matcher.
func (c *EventBus) Register(m MessageMatcher, h Handler) {
	c.handlers = append(c.handlers, &eventBusHandler{m, h})
}

// Handle disptaches the event to matched handlers.
func (c *EventBus) Handle(ctx context.Context, m messages.Message) error {
	for _, h := range c.handlers {
		if h.matches(m) {
			_ = h.handler(ctx, m) // we ignore error as this is an event rather then a command
		}
	}
	return nil
}

// WrapStore will return an event store where persisted events will be dispatched
// on this event bus.
func (c *EventBus) WrapStore(store eventstore.EventStore) eventstore.EventStore {
	return EventStoreWithBus(c, store)
}

// MatchAny will also return a positive match to process a message.
func MatchAny() MessageMatcher {
	return func(messages.Message) bool {
		return true
	}
}

// MatchMessageName will return a positive match if the message name
// matched the message name on the message provided.
func MatchMessageName(m messages.Message) MessageMatcher {
	return MatchMessageNameRaw(m.MessageName())
}

// MatchMessageNameRaw will return a positive match if the message
// name matched the message name provided.
func MatchMessageNameRaw(n string) MessageMatcher {
	return func(m messages.Message) bool {
		return m.MessageName() == n
	}
}
