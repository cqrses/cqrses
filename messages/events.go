package messages

import (
	"context"
	"time"
)

const (
	// EventCausationID should be set on a context where a command
	// is called, this will allow us to pull the command id out and
	// set it as the causation ID.
	EventCausationID eventMetaKey = "causation_id"

	// EventCorrelationID if another event caused this new event we should
	// set that as the correlation id on the event's metadata.
	EventCorrelationID eventMetaKey = "correlation_id"
)

type (
	eventMetaKey string

	// Event describes something that has happened.
	Event struct {
		messageID   string
		messageName string
		data        map[string]interface{}
		metadata    map[string]interface{}
		version     uint64
		created     time.Time
	}
)

// NewEvent will return an immutable event.
func NewEvent(id, name string, data, metadata map[string]interface{}, version uint64, created time.Time) *Event {
	return &Event{
		messageID:   id,
		messageName: name,
		data:        data,
		metadata:    metadata,
		version:     version,
		created:     created,
	}
}

// NewEventFromContext will return an immutable
// event, filling metadata with any information
// that we know can be added to events such as
// causation id.
func NewEventFromContext(ctx context.Context, id, name string, data, metadata map[string]interface{}, version uint64, created time.Time) *Event {
	return &Event{
		messageID:   id,
		messageName: name,
		data:        data,
		metadata:    metadata,
		version:     version,
		created:     created,
	}
}

// MessageID returns the id of the message.
func (e *Event) MessageID() string {
	return e.messageID
}

// MessageName returns the name of the message.
func (e *Event) MessageName() string {
	return e.messageName
}

// Data will return information related to the event.
func (e *Event) Data() map[string]interface{} {
	return e.data
}

// Metadata will return metadata about the event.
func (e *Event) Metadata() map[string]interface{} {
	return e.metadata
}

// Version returns the version of the event.
func (e *Event) Version() uint64 {
	return e.version
}

// Created returns the time the event was created.
func (e *Event) Created() time.Time {
	return e.created
}

// EventWithMetadata returns an event copied from the previous
// event with the updated metadata.
func EventWithMetadata(e *Event, m map[string]interface{}) *Event {
	em := &Event{}
	*em = *e
	em.metadata = m
	return em
}

// EventWithVersion returns an event copied from the previous
// event with the updated version.
func EventWithVersion(e *Event, v uint64) *Event {
	em := &Event{}
	*em = *e
	em.version = v
	return em
}
