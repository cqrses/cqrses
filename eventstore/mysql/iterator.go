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

// Current will return the item we currently have.
func (it *StreamIterator) Current() *messages.Event {
	return nil
}

// Next will get the next result, and if there is an error return it.
// Once next has been called without an error returned you can grab
// the result from Current()
func (it *StreamIterator) Next() error {
	return nil
}

// Rewind will set the position of the stream back to the default
// position and allow you to iterate of the stream again.
func (it *StreamIterator) Rewind() {
	//
}

// Set current will scan the row.
func (it *StreamIterator) setCurrent() {
	e := &messages.Event{}
	err := it.rows.Scan()

	if err != nil {
		it.currentItem = nil
		return
	}

	it.currentItem = e
}

// Close will clean up resources, do not attempt to use stream
// after closing.
func (it *StreamIterator) Close() {
	it.currentItem = nil
	it.rows.Close()
}
