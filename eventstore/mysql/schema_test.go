package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeStreamTableName(t *testing.T) {
	assert.Equal(t, "_7DFB4CF67742CB0660305E56EF816C53FCEC892CAE7F6EE39B75F34E659D672C", makeStreamTableName("users"))
}
