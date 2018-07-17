package main

import (
	"context"
	"time"

	"gopkg.in/cqrses/aggregate"
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
)

func (u *user) Handle(ctx context.Context, msg messages.Message, er aggregate.EventRecorder) error {
	switch msg.MessageName() {
	case createUserCommand:
		data := &createUserPayload{}
		data.FromPayload(msg.Data())
		if err := data.Validate(); err != nil {
			return err
		}
		return er(userCreated, data.Payload())
	}
	return nil
}

func (u *user) Apply(msg *messages.Event) error {
	switch msg.MessageName() {
	case userCreated:
		data := &createUserPayload{}
		data.FromPayload(msg.Data())
		u.id = data.userID
		u.emailAddress = data.emailAddress
		u.password = data.password
		u.created = msg.Created()
		u.removed = false
	}

	return nil
}
