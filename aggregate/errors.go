package aggregate

import (
	"strings"
)

// ErrPayloadValidationFailed ...
type ErrPayloadValidationFailed map[string][]string

// Push a new validation error message for a key.
func (e ErrPayloadValidationFailed) Push(key, msg string) {
	s, ok := e[key]
	if !ok {
		s = []string{}
	}
	e[key] = append(s, msg)
}

// FailedKeys returns a slice of keys that have error messages
// associated with them.
func (e ErrPayloadValidationFailed) FailedKeys() []string {
	res := []string{}
	for k := range e {
		res = append(res, k)
	}
	return res
}

func (e ErrPayloadValidationFailed) Error() string {
	res := "invalid payload:"
	for key, msgs := range e {
		res += " " + key + ": " + strings.Join(msgs, ", ")
	}
	return res
}
