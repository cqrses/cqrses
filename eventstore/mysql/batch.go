package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"gopkg.in/cqrses/eventstore"
)

type (
	batchHandler func(ctx context.Context, offset, limit uint64) (*sql.Rows, error)

	aggregateBatchHandler struct {
		db              *sql.DB
		tblName         string
		whereConditions string
		whereBindings   []interface{}
		orderBy         string
		matcher         eventstore.MetadataMatcher
	}
)

func newAggregateBatchHandler(db *sql.DB, tblName string, forward bool, matcher eventstore.MetadataMatcher) batchHandler {
	wc, wb := metadataMatcherConditionsToSQL(matcher)

	if wc == "" {
		wc = "1"
	}

	ah := &aggregateBatchHandler{
		db:              db,
		tblName:         tblName,
		whereConditions: wc,
		whereBindings:   wb,
		orderBy:         "`no`",
	}

	if !forward {
		ah.orderBy = "`no` DESC"
	}

	return ah.next
}

func (b *aggregateBatchHandler) next(ctx context.Context, offset, limit uint64) (*sql.Rows, error) {
	statement := fmt.Sprintf(
		"select * from `%s` where %s order by %s limit %d,%d",
		b.tblName,
		b.whereConditions,
		b.orderBy,
		offset,
		limit,
	)

	return b.db.QueryContext(ctx, statement, b.whereBindings...)
}
