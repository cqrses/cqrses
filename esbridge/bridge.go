package esbridge

import (
	"context"

	"gopkg.in/cqrses/bus"
	"gopkg.in/cqrses/eventstore"
	"gopkg.in/cqrses/messages"
)

type (
	eventStoreCtxKey struct{}
)

// AttachEventStoreToBus is a command bus middleware which will attach a
func AttachEventStoreToBus(es eventstore.EventStore) bus.CommandBusMiddleware {
	return func(ctx context.Context, msg messages.Message, next func(context.Context, messages.Message) error) error {
		if _, ok := ctx.Value(eventStoreCtxKey{}).(eventstore.EventStore); ok {
			return next(ctx, msg)
		}

		ctx = context.WithValue(ctx, eventStoreCtxKey{}, es)
		ctx = context.WithValue(ctx, messages.MetaCausationID, msg.MessageID())
		ctx = context.WithValue(ctx, messages.MetaCorrelationID, msg.MessageID())

		return next(ctx, msg)
	}
}

// GetEventStoreFromContext will return the event bus from the context, if it does not
// exist the second return value will be false.
func GetEventStoreFromContext(ctx context.Context) (eventstore.EventStore, bool) {
	es, ok := ctx.Value(eventStoreCtxKey{}).(eventstore.EventStore)
	return es, ok
}

// MustGetEventStoreFromContext calls GetEventStoreFromContext but if the bus is not present
// we will panic instead of returning gracefully.
func MustGetEventStoreFromContext(ctx context.Context) eventstore.EventStore {
	es, ok := GetEventStoreFromContext(ctx)
	if !ok {
		panic("expecting event bus on context but none found")
	}
	return es
}
