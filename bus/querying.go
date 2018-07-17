package bus

// import (
// 	"context"
// 	"errors"

// 	"gopkg.in/cqrses/messages"
// )

// var (
// 	ErrQueryAlreadyRegistered = errors.New("query already registered")
// )

// type (
// 	QueryHandler func(context.Context, messages.Message) (interface{}, error)

// 	QueryBus struct {
// 		handlers map[string]QueryHandler
// 	}
// )

// func (b *QueryBus) Register(n string, h QueryHandler) error {
// 	if _, ok := b.handlers[n]; ok {
// 		return ErrCommandAlreadyRegistered
// 	}

// 	b.handlers[n] = h

// 	return nil
// }

// func (c *QueryBus) Handle(ctx context.Context, m messages.Message) (interface{}, error) {
// 	h, ok := c.handlers[m.MessageName()]
// 	if !ok {
// 		return nil, ErrNoHandlerFound
// 	}

// 	res, err := h(ctx, m)
// 	if err != nil {
// 		return nil, &Error{original: err}
// 	}

// 	return res, nil
// }
