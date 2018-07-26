package mysql

import (
	"fmt"
	"strings"

	"github.com/go-cqrses/cqrses/eventstore"
)

func metadataMatcherConditionsToSQL(conditions eventstore.MetadataMatcher) (string, []interface{}) {
	sql := []string{}
	bindings := []interface{}{}

	for field, condition := range conditions {
		var op, val string
		switch condition.Operation {
		case eventstore.MatchOpIn:
			op = "IN"
			val = "(?" + strings.Repeat(",?", len(condition.Values)-1) + ")"

			vb := make([]interface{}, len(condition.Values))
			for vi, vv := range condition.Values {
				vb[vi] = vv
			}
			bindings = append(bindings, vb...)
		case eventstore.MatchOpNotIn:
			op = "NOT IN"
			val = "(?" + strings.Repeat(",?", len(condition.Values)-1) + ")"

			vb := make([]interface{}, len(condition.Values))
			for vi, vv := range condition.Values {
				vb[vi] = vv
			}
			bindings = append(bindings, vb...)
		case eventstore.MatchOpRegex:
			op = "REGEX"
			// todo (bweston92) add support for MySQL REGEX match.
			panic("todo")
		case eventstore.MatchOpEq:
			op = "="
			val = "?"
			bindings = append(bindings, condition.Values[0])
		}
		sql = append(sql, fmt.Sprintf("(`%s` %s %s)", field, op, val))
	}

	return strings.Join(sql, " AND "), bindings
}
