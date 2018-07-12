package messages

import (
	"encoding/json"
	"time"
)

type (
	MessageFactory interface {
		Serialize(Message) ([]byte, error)
		Unserialize([]byte) (Message, error)
	}

	JSONMessageFactory struct{}

	JSONMessage struct {
		MessageID   string                 `json:"message_id"`
		MessageName string                 `json:"message_name"`
		Data        map[string]interface{} `json:"data"`
		Metadata    map[string]interface{} `json:"metadata"`
		Version     uint64                 `json:"version"`
		Created     string                 `json:"created"`
	}

	JSONMessageWrapper struct {
		values *JSONMessage
	}
)

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

func (f *JSONMessageFactory) Unserialize(m []byte) (Message, error) {
	out := new(JSONMessage)
	return &JSONMessageWrapper{
		values: out,
	}, json.Unmarshal(m, out)
}

func (m *JSONMessageWrapper) MessageID() string {
	return m.values.MessageID
}

func (m *JSONMessageWrapper) MessageName() string {
	return m.values.MessageName
}

func (m *JSONMessageWrapper) Data() map[string]interface{} {
	return m.values.Data
}

func (m *JSONMessageWrapper) Metadata() map[string]interface{} {
	return m.values.Metadata
}

func (m *JSONMessageWrapper) Version() uint64 {
	return m.values.Version
}

func (m *JSONMessageWrapper) Created() time.Time {
	t, _ := time.Parse(time.RFC3339Nano, m.values.Created)
	return t
}
