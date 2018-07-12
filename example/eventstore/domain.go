package main

import (
	"time"

	"gopkg.in/cqrses/eventstore"
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
