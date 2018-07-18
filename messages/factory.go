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

	keyFormatter func(interface{}, bool) interface{}

	JSONMessageFactory struct {
		keyFormatters map[string]keyFormatter
	}

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

func NewJSONMessageFactory() *JSONMessageFactory {
	return &JSONMessageFactory{
		keyFormatters: map[string]keyFormatter{},
	}
}

func (f *JSONMessageFactory) AddPayloadKeyFormatter(name string, formatter keyFormatter) {
	f.keyFormatters[name] = formatter
}

func (f *JSONMessageFactory) Serialize(m Message) ([]byte, error) {
	data := m.Data()

	for key, val := range data {
		if formatter, ok := f.keyFormatters[key]; ok {
			data[key] = formatter(val, true)
		}
	}

	return json.Marshal(JSONMessage{
		MessageID:   m.MessageID(),
		MessageName: m.MessageName(),
		Data:        data,
		Metadata:    m.Metadata(),
		Version:     m.Version(),
		Created:     m.Created().Format(time.RFC3339Nano),
	})
}

func (f *JSONMessageFactory) Unserialize(m []byte) (Message, error) {
	out := new(JSONMessage)
	if err := json.Unmarshal(m, out); err != nil {
		return nil, err
	}

	for key, val := range out.Data {
		if formatter, ok := f.keyFormatters[key]; ok {
			out.Data[key] = formatter(val, false)
		}
	}

	return &JSONMessageWrapper{
		values: out,
	}, nil
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
