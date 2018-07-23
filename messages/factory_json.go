package messages

import (
	"encoding/json"
	"time"

	"github.com/buger/jsonparser"
)

type (
	// A key formatter is called with a JSON value, for example if the payload looks
	// like {"address": {"building": "abc"}} interface would be:
	// `map[string]interface{}{"building": "abc"}``
	// You can return any type as a value and it will be assigned to the key.
	//
	// This is run against both metadata and data where data was not handled by another
	// message building stratergy.
	keyFormatter func(interface{}, bool) interface{}

	JSONMessageFactory struct {
		dataTypeFactories map[string]dataTypeFactory
	}

	JSONMessage struct {
		MessageID   string                 `json:"message_id"`
		MessageName string                 `json:"message_name"`
		Data        interface{}            `json:"data"`
		Metadata    map[string]interface{} `json:"metadata"`
		Version     uint64                 `json:"version"`
		Created     string                 `json:"created_at"`
	}

	JSONMessageWrapper struct {
		values *JSONMessage
	}
)

// NewJSONMessageFactory will return a new message factory
// that serializes and unserialises JSON payloads.
func NewJSONMessageFactory() *JSONMessageFactory {
	return &JSONMessageFactory{
		dataTypeFactories: map[string]dataTypeFactory{},
	}
}

// Builds ...
func (f *JSONMessageFactory) Builds(name string, factory dataTypeFactory) {
	f.dataTypeFactories[name] = factory
}

// Build ...
func (f *JSONMessageFactory) Build(msgName string, pl []byte) (interface{}, bool) {
	b, ok := f.dataTypeFactories[msgName]
	if !ok {
		return nil, false
	}

	out := b()
	return out, json.Unmarshal(pl, &out) == nil
}

// Serialize ...
func (f *JSONMessageFactory) Serialize(m Message) ([]byte, error) {
	return json.Marshal(JSONMessage{
		MessageID:   m.MessageID(),
		MessageName: m.MessageName(),
		Data:        m.Data(),
		Metadata:    m.Metadata(),
		Version:     m.Version(),
		Created:     m.Created().Format(time.RFC3339Nano),
	})
}

// Unserialize ...
func (f *JSONMessageFactory) Unserialize(m []byte) (Message, error) {
	msgName, err := jsonparser.GetString(m, "message_name")
	if err != nil {
		return nil, err
	}

	out := new(JSONMessage)

	// Set known object.
	out.MessageID, _ = jsonparser.GetString(m, "message_id")
	out.MessageName = msgName
	version, _ := jsonparser.GetInt(m, "version")
	out.Version = uint64(version)
	out.Created, _ = jsonparser.GetString(m, "created_at")

	// Get the data and build correct type.
	data, _, _, err := jsonparser.Get(m, "data")
	if err != nil {
		return nil, err
	}
	dtf, ok := f.Build(msgName, data)
	if ok {
		out.Data = dtf
	} else {
		var dtm map[string]interface{}
		if err := json.Unmarshal(data, &dtm); err != nil {
			return nil, err
		}
		out.Data = dtm
	}

	// Get the payload.
	if md, _, _, err := jsonparser.Get(m, "metadata"); err != nil {
		return nil, err
	} else if mErr := json.Unmarshal(md, &out.Metadata); mErr != nil {
		return nil, mErr
	}

	return &JSONMessageWrapper{
		values: out,
	}, nil
}

// MessageID ...
func (m *JSONMessageWrapper) MessageID() string {
	return m.values.MessageID
}

// MessageName ...
func (m *JSONMessageWrapper) MessageName() string {
	return m.values.MessageName
}

// Data ...
func (m *JSONMessageWrapper) Data() interface{} {
	return m.values.Data
}

// Metadata ...
func (m *JSONMessageWrapper) Metadata() map[string]interface{} {
	return m.values.Metadata
}

// Version ...
func (m *JSONMessageWrapper) Version() uint64 {
	return m.values.Version
}

// Created ...
func (m *JSONMessageWrapper) Created() time.Time {
	t, _ := time.Parse(time.RFC3339Nano, m.values.Created)
	return t
}
