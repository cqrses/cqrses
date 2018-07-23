package messages

type (
	// MessageFactory ...
	MessageFactory interface {
		Serialize(Message) ([]byte, error)
		Unserialize([]byte) (Message, error)
	}

	// When you have a struct that represents a message's data you can use payload
	// struct factories that should return a struct which the data will be unmarshalled
	// into.
	dataTypeFactory func() interface{}

	// PayloadBuilder ...
	PayloadBuilder interface {
		Build(msgName string, pl []byte) (interface{}, bool)
		Builds(msgName string, with dataTypeFactory)
	}
)
