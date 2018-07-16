package messages

type (
	// SupportsPayloadManipulation is a helper interface that allows the following
	// method BuildPayload and ReadPayload to interace with a message payload object.
	SupportsPayloadManipulation interface {
		Payload() map[string]interface{}
		FromPayload(map[string]interface{})
	}
)

// BuildPayload will take a map[string]interface{} representation from a struct.
func BuildPayload(in SupportsPayloadManipulation) map[string]interface{} {
	return in.Payload()
}

// ReadPayload will provide the structs FromPayload method with the data from
// a message.
func ReadPayload(m Message, in SupportsPayloadManipulation) {
	in.FromPayload(m.Data())
}
