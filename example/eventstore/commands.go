package main

import (
	"context"
	"errors"
	"time"

	"github.com/satori/go.uuid"

	"gopkg.in/cqrses/messages"
)

const (
	createUserCommand = "user.create"
)

type (
	// CreateUserPayload is used to provide data for the CreateUser command.
	CreateUserPayload struct {
		userID       string
		emailAddress string
		password     string
	}
)

// CreateUserWith ...
func CreateUserWith(ctx context.Context, id, emailAddress, hashedPassword string) (*messages.Command, error) {
	pl := &CreateUserPayload{
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
func (c *CreateUserPayload) Validate() error {
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

func (c *CreateUserPayload) Payload() map[string]interface{} {
	return map[string]interface{}{
		"user_id":       c.userID,
		"email_address": c.emailAddress,
		"password":      c.password,
	}
}

func (c *CreateUserPayload) FromPayload(data map[string]interface{}) {
	c.userID, _ = data["user_id"].(string)
	c.emailAddress, _ = data["email_address"].(string)
	c.password, _ = data["password"].(string)
}
