package messages

import (
	"time"
)

type (
	// Message provides base functionality to commands,
	// events and queries.
	Message interface {
		// MessageID returns the id of the message.
		MessageID() string

		// MessageName returns the name of the message.
		MessageName() string

		// Data will return information related to the event.
		Data() map[string]interface{}

		// Metadata will return metadata about the event.
		Metadata() map[string]interface{}

		// Version returns the version of the event.
		Version() uint64

		// Created returns the time the event was created.
		Created() time.Time
	}
)
