package messages_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-cqrses/cqrses/messages"
)

func TestJSONMessageFactory(t *testing.T) {
	sut := &messages.JSONMessageFactory{}

	event := messages.NewEvent(
		"hello-world",
		"test.event",
		map[string]interface{}{
			"hello": "world",
		},
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
	assert.Equal(t, event.Data(), out.Data())
	assert.Equal(t, event.Metadata(), out.Metadata())
	assert.Equal(t, event.Version(), out.Version())
	assert.Equal(t, event.Created().Format(time.RFC3339Nano), out.Created().Format(time.RFC3339Nano))
}

type myCustonType struct {
	Sentence string `json:"sentence"`
	Words    int64  `json:"words"`
}

func TestJSONMessageFactoryCustomTypes(t *testing.T) {
	sut := messages.NewJSONMessageFactory()
	sut.AddPayloadKeyFormatter("customSentence", func(in interface{}, serialise bool) interface{} {
		if serialise {
			if st, ok := in.(*myCustonType); ok {
				return map[string]interface{}{
					"sentence": st.Sentence,
					"words":    st.Words,
				}
			}
		} else {
			st := in.(map[string]interface{})
			words, _ := st["words"].(float64)
			return &myCustonType{
				Sentence: st["sentence"].(string),
				Words:    int64(words),
			}
		}

		return in
	})

	ctIn := &myCustonType{
		Sentence: "my name is test",
		Words:    4,
	}

	event := messages.NewEvent(
		"hello-world",
		"test.event",
		map[string]interface{}{
			"hello":          "world",
			"customSentence": ctIn,
		},
		map[string]interface{}{
			"1+1": "2",
		},
		1,
		time.Now(),
	)

	in, err := sut.Serialize(event)
	assert.Nil(t, err)

	out, err := sut.Unserialize(in)
	assert.Nil(t, err)

	ctOut, ok := out.Data()["customSentence"].(*myCustonType)
	if !ok {
		t.Fatalf("did not get *myCustomType back out")
	}

	assert.Equal(t, ctIn, ctOut)
}
