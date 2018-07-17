package main

import (
	"context"
	"errors"
	"time"

	"github.com/satori/go.uuid"

	"gopkg.in/cqrses/aggregate"
	"gopkg.in/cqrses/bus"
	"gopkg.in/cqrses/esbridge"
	"gopkg.in/cqrses/messages"
)

const (
	createUserCommand         = "user.create"
	changeUserPasswordCommand = "user.changePassword"
)

type (
	// createUserPayload is used to provide data for the createUserCommand.
	createUserPayload struct {
		userID       string
		emailAddress string
		password     string
	}
	// changeUserPasswordPayload is used to provide data for the changeUserPasswordCommand.
	changeUserPasswordPayload struct {
		userID   string
		password string
	}
)

func registerCommandBusHandlers(cmdBus *bus.CommandBus) {
	cmdBus.Register(createUserCommand, func(ctx context.Context, msg messages.Message) error {
		es := esbridge.MustGetEventStoreFromContext(ctx)
		userID, ok := msg.Data()["user_id"].(string)
		if !ok {
			return errors.New("expected user id and did not get one")
		}

		u := &user{}
		a := aggregate.New(userID, es, "users", u)
		a.Handle(ctx, msg)

		return a.Close(ctx)
	})
	cmdBus.Register(changeUserPasswordCommand, func(ctx context.Context, msg messages.Message) error {
		es := esbridge.MustGetEventStoreFromContext(ctx)
		userID, ok := msg.Data()["user_id"].(string)
		if !ok {
			return errors.New("expected user id and did not get one")
		}

		u := &user{}
		a, err := aggregate.Load(ctx, userID, es, "users", u)
		if err != nil {
			return err
		}
		a.Handle(ctx, msg)

		return a.Close(ctx)
	})
}

func createUserWith(ctx context.Context, id, emailAddress, hashedPassword string) (*messages.Command, error) {
	pl := &createUserPayload{
		userID:       id,
		emailAddress: emailAddress,
		password:     hashedPassword,
	}

	if err := pl.Validate(); err != nil {
		return nil, err
	}

	return messages.NewCommandFromContext(
		ctx,
		uuid.Must(uuid.NewV4()).String(),
		createUserCommand,
		pl.Payload(),
		map[string]interface{}{},
		0,
		time.Now(),
	), nil
}

// Validate will ensure all related information for creating the user
// is availabile within the payload payload.
func (c *createUserPayload) Validate() error {
	if len(c.userID) < 34 {
		return errors.New("user id not valid")
	}

	if len(c.emailAddress) < 4 {
		return errors.New("email address not valid")
	}

	if len(c.password) == 0 {
		return errors.New("password hash is required")
	}

	return nil
}

func (c *createUserPayload) Payload() map[string]interface{} {
	return map[string]interface{}{
		"user_id":       c.userID,
		"email_address": c.emailAddress,
		"password":      c.password,
	}
}

func (c *createUserPayload) FromPayload(data map[string]interface{}) {
	c.userID, _ = data["user_id"].(string)
	c.emailAddress, _ = data["email_address"].(string)
	c.password, _ = data["password"].(string)
}

func changeUserPassword(ctx context.Context, id, newHashedPassword string) (*messages.Command, error) {
	pl := &changeUserPasswordPayload{
		userID:   id,
		password: newHashedPassword,
	}

	if err := pl.Validate(); err != nil {
		return nil, err
	}

	return messages.NewCommandFromContext(
		ctx,
		uuid.Must(uuid.NewV4()).String(),
		changeUserPasswordCommand,
		pl.Payload(),
		map[string]interface{}{},
		0,
		time.Now(),
	), nil
}

// Validate will ensure all related information for creating the user
// is availabile within the payload payload.
func (c *changeUserPasswordPayload) Validate() error {
	if len(c.userID) < 34 {
		return errors.New("user id not valid")
	}

	if len(c.password) == 0 {
		return errors.New("password hash for new password is required")
	}

	return nil
}

func (c *changeUserPasswordPayload) Payload() map[string]interface{} {
	return map[string]interface{}{
		"user_id":      c.userID,
		"new_password": c.password,
	}
}

func (c *changeUserPasswordPayload) FromPayload(data map[string]interface{}) {
	c.userID, _ = data["user_id"].(string)
	c.password, _ = data["new_password"].(string)
}
