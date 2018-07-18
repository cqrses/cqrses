package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/go-cqrses/cqrses/eventstore"
	"github.com/go-cqrses/cqrses/messages"
)

type (
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
	}
)

func iter(bh batchHandler, batchSize, fromNumber, count uint64) *StreamIterator {
	return &StreamIterator{
		currentItem:       nil,
		currentKey:        -1,
		batchHandler:      bh,
		batchSize:         batchSize,
		batchPosition:     0,
		fromNumber:        fromNumber,
		currentFromNumber: fromNumber,
		count:             count,
	}
}

// Current will return the item we currently have.
func (it *StreamIterator) Current() *messages.Event {
	return it.currentItem
}

// Next will get the next result, and if there is an error return it.
// Once next has been called without an error returned you can grab
// the result from Current()
func (it *StreamIterator) Next(ctx context.Context) error {
	if it.rows == nil {
		limit := it.batchSize
		if it.count > 0 {
			limit = it.count
		}
		rows, err := it.batchHandler(ctx, it.fromNumber, limit)
		if err != nil {
			return err
		}
		it.rows = rows
	}

	if !it.rows.Next() {
		return eventstore.EOF
	}

	var no, eventID, eventName, payload, metadata, createdAt, aggregateID string
	var aggregateVersion uint64

	err := it.rows.Scan(&no, &eventID, &eventName, &payload, &metadata, &createdAt, &aggregateVersion, &aggregateID)
	if err != nil {
		return err
	}

	var jp, jm map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &jp); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(metadata), &jm); err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return err
	}

	it.currentItem = messages.NewEvent(eventID, eventName, jp, jm, aggregateVersion, t)

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
	if it.rows != nil {
		it.rows.Close()
	}
}

// ErrorStreamIterator is returned when an error occured getting
// stream data, maybe it didn't exist.
type ErrorStreamIterator struct {
	err error
}

// Current always returns nil.
func (it *ErrorStreamIterator) Current() *messages.Event {
	return nil
}

// Next return the error provided.
func (it *ErrorStreamIterator) Next(ctx context.Context) error {
	return it.err
}

// Rewind is empty.
func (it *ErrorStreamIterator) Rewind() {}

// Close is empty.
func (it *ErrorStreamIterator) Close() {}

// Error will return the inner error's Error method result.
func (it *ErrorStreamIterator) Error() string {
	return it.err.Error()
}
