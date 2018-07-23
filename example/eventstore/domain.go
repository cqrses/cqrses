package main

import (
	"context"
	"errors"
	"time"

	"github.com/go-cqrses/cqrses/aggregate"
	"github.com/go-cqrses/cqrses/messages"
)

const (
	userCreated             = "user.created"
	userPasswordChanged     = "user.passwordChanged"
	userEmailAddressUpdated = "user.emailAddressUpdated"
	userRemoved             = "user.removed"
)

type (
	userCreatedPayload struct {
		UserID       string `json:"user_id"`
		EmailAddress string `json:"email_address"`
		Password     string `json:"password"`
	}

	userPasswordChangedPayload struct {
		UserID   string `json:"user_id"`
		Password string `json:"password"`
	}

	user struct {
		id           string
		emailAddress string
		password     string
		created      time.Time
		removed      bool
	}
)

func (u *user) Handle(ctx context.Context, msg messages.Message, er aggregate.EventRecorder) error {
	switch cmd := msg.Data().(type) {
	case *createUserPayload:
		if err := cmd.Validate(); err != nil {
			return err
		}

		if u.id != "" {
			return errors.New("user with that id already exists")
		}

		return er(userCreated, &userCreatedPayload{
			UserID:       cmd.UserID,
			EmailAddress: cmd.EmailAddress,
			Password:     cmd.Password,
		})
	case *changeUserPasswordPayload:
		if err := cmd.Validate(); err != nil {
			return err
		}
		return er(
			userPasswordChanged,
			&userPasswordChangedPayload{
				UserID:   cmd.UserID,
				Password: cmd.Password,
			},
		)
	}
	return nil
}

func (u *user) Apply(msg *messages.Event) error {
	switch event := msg.Data().(type) {
	case *userCreatedPayload:
		u.id = event.UserID
		u.emailAddress = event.EmailAddress
		u.password = event.Password
		u.created = msg.Created()
		u.removed = false
	}

	return nil
}
