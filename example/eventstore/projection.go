package main

import (
	"context"

	"gopkg.in/cqrses/messages"
)

type (
	userDTO struct {
	}

	userCollection map[string]*userDTO

	userProjector struct {
		all userCollection
	}
)

func (p *userProjector) Handle(ctx context.Context, event *messages.Event) error {
	return nil
}
