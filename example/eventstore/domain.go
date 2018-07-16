package main

import (
	"context"
	"time"

	"gopkg.in/cqrses/eventstore"
	"gopkg.in/cqrses/messages"
)

const (
	userCreated             = "user.created"
	userPasswordChanged     = "user.passwordChanged"
	userEmailAddressUpdated = "user.emailAddressUpdated"
	userRemoved             = "user.removed"
)

type (
	user struct {
		id           string
		emailAddress string
		password     string
		created      time.Time
		removed      bool
	}

	userHandlers struct {
		store eventstore.EventStore
	}
)

func (h *userHandlers) Handle(ctx context.Context, m messages.Message) error {
	switch m.MessageName() {
	case createUserCommand:
	}
	return nil
}

func (*user) create(ctx context.Context, m messages.Message) {
	data := &createUserPayload{}
	data.Payload()
}
