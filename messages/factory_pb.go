package messages

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

type (
	ProtoMessageFactory struct{}

	ProtoMessageWrapper struct {
		values *DomainMessage
	}
)

// NewProtoMessageFactory will return a new message factory
// that serializes and unserialises Proto payloads.
func NewProtoMessageFactory() *ProtoMessageFactory {
	return &ProtoMessageFactory{}
}

// Builds is obsolete for protobuff as the generated Go library registers the messages itself.
func (f *ProtoMessageFactory) Builds(name string, factory dataTypeFactory) {
	panic("ProtoMessageFactory does not need registering, just import the Go library")
}

// Build ...
func (f *ProtoMessageFactory) Build(msgName string, pl []byte) (interface{}, bool) {
	mt := proto.MessageType(msgName)
	if mt == nil {
		return nil, false
	}

	out, ok := reflect.New(mt.Elem()).Interface().(proto.Message)

	if !ok {
		return nil, false
	}

	return out, proto.Unmarshal(pl, out) == nil
}

// Serialize ...
func (f *ProtoMessageFactory) Serialize(m Message) ([]byte, error) {
	d, _ := proto.Marshal(m.Data().(proto.Message))
	meta, _ := json.Marshal(m.Metadata())
	created, _ := ptypes.TimestampProto(m.Created())

	return proto.Marshal(&DomainMessage{
		MessageId:   m.MessageID(),
		MessageName: m.MessageName(),
		Data:        d,
		Metadata:    map[string][]byte{"__json": meta},
		Version:     m.Version(),
		Created:     created,
	})
}

// Unserialize ...
func (f *ProtoMessageFactory) Unserialize(m []byte) (Message, error) {
	out := &DomainMessage{}

	if err := proto.Unmarshal(m, out); err != nil {
		return nil, err
	}

	return &ProtoMessageWrapper{values: out}, nil
}

// MessageID ...
func (m *ProtoMessageWrapper) MessageID() string {
	return m.values.MessageId
}

// MessageName ...
func (m *ProtoMessageWrapper) MessageName() string {
	return m.values.MessageName
}

// Data ...
func (m *ProtoMessageWrapper) Data() interface{} {
	mt := proto.MessageType(m.MessageName())
	if mt == nil {
		panic("cannot unmarshal DomainMessage data, proto message not registered: " + m.MessageName())
	}

	out, ok := reflect.New(mt.Elem()).Interface().(proto.Message)
	if !ok {
		panic("cannot unmarshal DomainMessage destination not a proto.Message for type: " + m.MessageName())
	}

	if err := proto.Unmarshal(m.values.Data, out); err != nil {
		panic("cannot unmarshal DomainMessage data is invalid: " + err.Error())
	}

	return out
}

// Metadata ...
func (m *ProtoMessageWrapper) Metadata() map[string]interface{} {
	out := map[string]interface{}{}
	if rb, ok := m.values.Metadata["__json"]; ok && len(rb) > 0 {
		_ = json.Unmarshal(rb, &out)
	}
	return out
}

// Version ...
func (m *ProtoMessageWrapper) Version() uint64 {
	return m.values.Version
}

// Created ...
func (m *ProtoMessageWrapper) Created() time.Time {
	return time.Unix(m.values.Created.Seconds, int64(m.values.Created.Nanos))
}
