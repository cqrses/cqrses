package mysql

import (
	"database/sql"

	"gopkg.in/cqrses/messages"
)

type (
	batchHandler func(offset uint64, limit uint64) *sql.Rows

	// StreamIterator iterates over events from a MySQL database.
	StreamIterator struct {
		rows              *sql.Rows
		currentItem       *messages.Event
		currentKey        int64
		batchHandler      batchHandler
		batchSize         uint64
		batchPosition     uint64
		fromNumber        uint64
		currentFromNumber uint64
		count             uint64
		forward           bool
	}
)

func iter(bh batchHandler, batchSize, fromNumber, count uint64, forward bool) *StreamIterator {
	return &StreamIterator{
		currentItem:       nil,
		currentKey:        -1,
		batchHandler:      bh,
		batchSize:         batchSize,
		batchPosition:     0,
		fromNumber:        fromNumber,
		currentFromNumber: fromNumber,
		count:             count,
		forward:           forward,
	}
}

func (it *StreamIterator) Next() error {
	return nil
}

func (it *StreamIterator) Reset() {
	//
}

func (it *StreamIterator) setCurrent() {
	e := &messages.Event{}
	err := it.rows.Scan()

	if err != nil {
		it.currentItem = nil
		return
	}

	it.currentItem = e
}
