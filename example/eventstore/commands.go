package main

import (
	"context"
	"errors"
	"time"

	"github.com/satori/go.uuid"

	"github.com/go-cqrses/cqrses/aggregate"
	"github.com/go-cqrses/cqrses/bus"
	"github.com/go-cqrses/cqrses/esbridge"
	"github.com/go-cqrses/cqrses/messages"
)

const (
	createUserCommand         = "user.create"
	changeUserPasswordCommand = "user.changePassword"
)

type (
	// createUserPayload is used to provide data for the createUserCommand.
	createUserPayload struct {
		UserID       string `json:"user_id"`
		EmailAddress string `json:"email_address"`
		Password     string `json:"password"`
	}
	// changeUserPasswordPayload is used to provide data for the changeUserPasswordCommand.
	changeUserPasswordPayload struct {
		UserID   string `json:"user_id"`
		Password string `json:"password"`
	}
)

func registerCommandBusHandlers(cmdBus *bus.CommandBus) {
	cmdBus.Register(createUserCommand, func(ctx context.Context, msg messages.Message) error {
		es := esbridge.MustGetEventStoreFromContext(ctx)
		pl, ok := msg.Data().(*createUserPayload)
		if !ok {
			return errors.New("invalid command payload provided")
		}

		u := &user{}
		a := aggregate.New(pl.UserID, es, "users", u)
		a.Handle(ctx, msg)

		return a.Close(ctx)
	})
	cmdBus.Register(changeUserPasswordCommand, func(ctx context.Context, msg messages.Message) error {
		es := esbridge.MustGetEventStoreFromContext(ctx)
		pl, ok := msg.Data().(*changeUserPasswordPayload)
		if !ok {
			return errors.New("invalid command payload provided")
		}

		u := &user{}
		a, err := aggregate.Load(ctx, pl.UserID, es, "users", u)
		if err != nil {
			return err
		}
		a.Handle(ctx, msg)

		return a.Close(ctx)
	})
}

func createUserWith(ctx context.Context, id, emailAddress, hashedPassword string) (*messages.Command, error) {
	pl := &createUserPayload{
		UserID:       id,
		EmailAddress: emailAddress,
		Password:     hashedPassword,
	}

	if err := pl.Validate(); err != nil {
		return nil, err
	}

	return messages.NewCommandFromContext(
		ctx,
		uuid.Must(uuid.NewV4()).String(),
		createUserCommand,
		pl,
		map[string]interface{}{},
		0,
		time.Now(),
	), nil
}

// Validate will ensure all related information for creating the user
// is availabile within the payload payload.
func (c *createUserPayload) Validate() error {
	if len(c.UserID) < 34 {
		return errors.New("user id not valid")
	}

	if len(c.EmailAddress) < 4 {
		return errors.New("email address not valid")
	}

	if len(c.Password) == 0 {
		return errors.New("password hash is required")
	}

	return nil
}

func changeUserPassword(ctx context.Context, id, newHashedPassword string) (*messages.Command, error) {
	pl := &changeUserPasswordPayload{
		UserID:   id,
		Password: newHashedPassword,
	}

	if err := pl.Validate(); err != nil {
		return nil, err
	}

	return messages.NewCommandFromContext(
		ctx,
		uuid.Must(uuid.NewV4()).String(),
		changeUserPasswordCommand,
		pl,
		map[string]interface{}{},
		0,
		time.Now(),
	), nil
}

// Validate will ensure all related information for creating the user
// is availabile within the payload payload.
func (c *changeUserPasswordPayload) Validate() error {
	if len(c.UserID) < 34 {
		return errors.New("user id not valid")
	}

	if len(c.Password) == 0 {
		return errors.New("password hash for new password is required")
	}

	return nil
}
