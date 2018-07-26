package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeStreamTableName(t *testing.T) {
	assert.Equal(t, "_5B7DCD14A4FAA2CDD54CF6EB8D4BC35DA31914A1", makeStreamTableName("users"))
}
