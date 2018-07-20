package messages

type (
	// MessageFactory ...
	MessageFactory interface {
		Serialize(Message) ([]byte, error)
		Unserialize([]byte) (Message, error)
	}

	// PayloadBuilder ...
	PayloadBuilder interface {
		Build(msgName string, pl []byte) (interface{}, bool)
	}
)
