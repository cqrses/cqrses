package mysql

import (
	"testing"

	"github.com/go-cqrses/cqrses/eventstore"

	"github.com/stretchr/testify/assert"
)

func TestMetadataMatcherConditionsToSQLMatchOpIn(t *testing.T) {
	{
		m := eventstore.MetadataMatcher{
			"field1": eventstore.MetadataMatcherCondition{
				Operation: eventstore.MatchOpIn,
				Values:    []string{"abcd", "def"},
			},
		}
		sql, bindings := metadataMatcherConditionsToSQL(m)

		assert.Contains(t, sql, "field1")
		assert.Contains(t, sql, "?,?")
		assert.Len(t, bindings, 2)
		assert.Equal(t, "abcd", bindings[0])
		assert.Equal(t, "def", bindings[1])
	}
	{
		m := eventstore.MetadataMatcher{
			"field1": eventstore.MetadataMatcherCondition{
				Operation: eventstore.MatchOpIn,
				Values:    []string{"abcd"},
			},
		}
		sql, bindings := metadataMatcherConditionsToSQL(m)

		assert.Contains(t, sql, "field1")
		assert.Contains(t, sql, "?")
		assert.Len(t, bindings, 1)
		assert.Equal(t, "abcd", bindings[0])
	}
}

func TestMetadataMatcherConditionsToSQLMatchOpNotIn(t *testing.T) {
	{
		m := eventstore.MetadataMatcher{
			"field1": eventstore.MetadataMatcherCondition{
				Operation: eventstore.MatchOpNotIn,
				Values:    []string{"abcd", "def"},
			},
		}
		sql, bindings := metadataMatcherConditionsToSQL(m)

		assert.Contains(t, sql, "field1")
		assert.Contains(t, sql, "?,?")
		assert.Len(t, bindings, 2)
		assert.Equal(t, "abcd", bindings[0])
		assert.Equal(t, "def", bindings[1])
	}
	{
		m := eventstore.MetadataMatcher{
			"field1": eventstore.MetadataMatcherCondition{
				Operation: eventstore.MatchOpNotIn,
				Values:    []string{"abcd"},
			},
		}
		sql, bindings := metadataMatcherConditionsToSQL(m)

		assert.Contains(t, sql, "field1")
		assert.Contains(t, sql, "?")
		assert.Len(t, bindings, 1)
		assert.Equal(t, "abcd", bindings[0])
	}
}

func TestMetadataMatcherConditionsToSQLMatchOpEq(t *testing.T) {
	m := eventstore.MetadataMatcher{
		"field1": eventstore.MetadataMatcherCondition{
			Operation: eventstore.MatchOpEq,
			Values:    []string{"abcd"},
		},
	}
	sql, bindings := metadataMatcherConditionsToSQL(m)

	assert.Contains(t, sql, "field1")
	assert.Len(t, bindings, 1)
	assert.Equal(t, "abcd", bindings[0])
}
