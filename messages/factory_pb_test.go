package messages

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/golang/protobuf/proto"
)

func TestProtoMessageFactory(t *testing.T) {
	sut := &ProtoMessageFactory{}
	exp := &TestPayload{
		AString: "boo",
		AInt:    1,
		ABool:   true,
	}

	event := NewEvent(
		"hello-world",
		"com.github.go_cqrses.cqrses.messages.TestPayload",
		exp,
		map[string]interface{}{
			"1+1": "2",
		},
		1,
		time.Now(),
	)

	in, err := sut.Serialize(event)
	assert.Nil(t, err)
	assert.NotEmpty(t, in)

	out, err := sut.Unserialize(in)
	assert.Nil(t, err)
	assert.Equal(t, event.MessageID(), out.MessageID())
	assert.Equal(t, event.MessageName(), out.MessageName())
	assert.Equal(t, event.Metadata(), out.Metadata())
	assert.Equal(t, event.Version(), out.Version())
	assert.Equal(t, event.Created().Format(time.RFC3339Nano), out.Created().Format(time.RFC3339Nano))

	// We have to test the message payload because of the `XXX_*` keys that we cannot control.
	eventData := event.Data().(*TestPayload)
	outData := out.Data().(*TestPayload)
	assert.Equal(t, eventData.AString, outData.AString)
	assert.Equal(t, eventData.AInt, outData.AInt)
	assert.Equal(t, eventData.ABool, outData.ABool)

}

func TestProtoMessageFactoryBuild(t *testing.T) {
	fac := &ProtoMessageFactory{}
	in := &TestPayload{
		AString: "a test string",
		AInt:    4321,
		ABool:   true,
	}

	rb, err := proto.Marshal(in)
	if err != nil {
		t.Fatal("unable to seralise *TestPayload:", err)
	}

	out, valid := fac.Build("com.github.go_cqrses.cqrses.messages.TestPayload", rb)
	if !valid {
		t.Error("expected valid to be true but returned false")
		return
	}

	evnt := out.(*TestPayload)

	assert.Equal(t, in.AString, evnt.AString)
	assert.Equal(t, in.AInt, evnt.AInt)
	assert.Equal(t, in.ABool, evnt.ABool)
}
