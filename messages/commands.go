package messages

import (
	"context"
	"time"
)

type (
	// Command describes something that has happened.
	Command struct {
		messageID   string
		messageName string
		data        interface{}
		metadata    map[string]interface{}
		version     uint64
		created     time.Time
	}
)

// NewCommand will return an immutable command.
func NewCommand(id, name string, data interface{}, metadata map[string]interface{}, version uint64, created time.Time) *Command {
	return &Command{
		messageID:   id,
		messageName: name,
		data:        data,
		metadata:    metadata,
		version:     version,
		created:     created,
	}
}

// NewCommandFromContext will return an immutable
// command, filling metadata with any information
// that we know can be added to commands such as
// causation id.
func NewCommandFromContext(ctx context.Context, id, name string, data, metadata map[string]interface{}, version uint64, created time.Time) *Command {
	return &Command{
		messageID:   id,
		messageName: name,
		data:        data,
		metadata:    metadata,
		version:     version,
		created:     created,
	}
}

// MessageID returns the id of the message.
func (c *Command) MessageID() string {
	return c.messageID
}

// MessageName returns the name of the message.
func (c *Command) MessageName() string {
	return c.messageName
}

// Data will return information related to the command.
func (c *Command) Data() interface{} {
	return c.data
}

// Metadata will return metadata about the command.
func (c *Command) Metadata() map[string]interface{} {
	return c.metadata
}

// Version returns the version of the command.
func (c *Command) Version() uint64 {
	return c.version
}

// Created returns the time the Command was created.
func (c *Command) Created() time.Time {
	return c.created
}

// CommandWithMetadata returns an command copied from the previous
// command with the updated metadata.
func CommandWithMetadata(c *Command, m map[string]interface{}) *Command {
	cm := &Command{}
	*cm = *c
	cm.metadata = m
	return cm
}

// CommandWithVersion returns an command copied from the previous
// command with the updated version.
func CommandWithVersion(c *Command, v uint64) *Command {
	cm := &Command{}
	*cm = *c
	cm.version = v
	return cm
}
