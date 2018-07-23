package main

import (
	"context"

	"github.com/go-cqrses/cqrses/messages"
)

type (
	userDTO struct {
		id           string
		emailAddress string
	}

	userCollection map[string]*userDTO

	userProjector struct {
		all userCollection
	}
)

func (p *userProjector) Handle(ctx context.Context, msg messages.Message) error {
	switch event := msg.Data().(type) {
	case *userCreatedPayload:
		p.all[event.UserID] = &userDTO{event.UserID, event.EmailAddress}
	}
	return nil
}
