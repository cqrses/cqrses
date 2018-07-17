package main

import (
	"context"

	"gopkg.in/cqrses/messages"
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

func (p *userProjector) Handle(ctx context.Context, event messages.Message) error {
	switch event.MessageName() {
	case userCreated:
		data := event.Data()
		id, _ := data["user_id"].(string)
		emailAddress, _ := data["email_address"].(string)
		p.all[id] = &userDTO{id, emailAddress}
	}
	return nil
}
