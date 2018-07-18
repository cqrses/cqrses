package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/cqrses/eventstore"
	"gopkg.in/cqrses/messages"

	// MySQL driver for the database/sql package.
	"github.com/go-sql-driver/mysql"
)

const (
	// DefaultBatchSize ...
	DefaultBatchSize uint64 = 1000

	storeTimeFormat = "2006-01-02T15:04:05"
)

type (
	// EventStore will use a MySQL database to manage streams.
	EventStore struct {
		db        *sql.DB
		batchSize uint64
	}
)

// New returns a new MySQL event store, it is best to send a context
// with a deadline so we do not hang.
func New(ctx context.Context, dsn string, batchSize uint64) (*EventStore, error) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	if cfg.DBName == "" {
		cfg.DBName = "eventstore"
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	if err := applyEventStreamsSchema(ctx, db); err != nil {
		db.Close()
		return nil, err
	}

	return &EventStore{
		db:        db,
		batchSize: batchSize,
	}, nil
}

// Ping tests connection to the database is still ok.
func (s *EventStore) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// Load events from the given stream name.
func (s *EventStore) Load(ctx context.Context, streamName string, from, count uint64, matcher eventstore.MetadataMatcher) eventstore.StreamIterator {
	tblName, err := getStreamTableName(ctx, s.db, streamName)
	if err != nil {
		return &ErrorStreamIterator{err}
	}
	return iter(newAggregateBatchHandler(s.db, tblName, true, matcher), s.batchSize, from, count)
}

// LoadReverse Loads events from the given stream name in reverse.
func (s *EventStore) LoadReverse(ctx context.Context, streamName string, from, count uint64, matcher eventstore.MetadataMatcher) eventstore.StreamIterator {
	tblName, err := getStreamTableName(ctx, s.db, streamName)
	if err != nil {
		return &ErrorStreamIterator{err}
	}
	return iter(newAggregateBatchHandler(s.db, tblName, false, matcher), s.batchSize, from, count)
}

// FetchStreamNames gets  stream names that match the filter.
func (s *EventStore) FetchStreamNames(ctx context.Context, filter string, matcher eventstore.MetadataMatcher, limit, offset uint64) ([]string, error) {
	return []string{}, nil
}

// FetchStreamNamesRegex gets stream names that match the regex filter.
func (s *EventStore) FetchStreamNamesRegex(ctx context.Context, filter string, matcher eventstore.MetadataMatcher, limit, offset uint64) ([]string, error) {
	return []string{}, nil
}

// FetchStreamMetadata gets the metadata about a stream.
func (s *EventStore) FetchStreamMetadata(ctx context.Context, streamName string) (eventstore.StreamMetadata, error) {
	return nil, eventstore.ErrStreamDoesNotExist
}

// Create will create the stream with the name and metadata provided.
//
// AFAIK MySQL doesn't have transactions that support schema changes
// along with inserting rows etc so this could leave the database in
// a dodgy state.
func (s *EventStore) Create(ctx context.Context, stream *eventstore.Stream) error {
	if err := createStream(ctx, s.db, stream); err != nil {
		return err
	}

	err := s.AppendTo(ctx, stream.Name, stream.Events)
	if err != nil {
		// Hope and prey we can delete cleanly, if the DB has gone then we are done for.
		s.Delete(ctx, stream.Name)
	}
	return err
}

// AppendTo will append events to the stream.
func (s *EventStore) AppendTo(ctx context.Context, streamName string, events []*messages.Event) error {
	if len(events) == 0 {
		return nil
	}
	tblName, err := getStreamTableName(ctx, s.db, streamName)
	if err != nil {
		return err
	}

	values := "(?, ?, ?, ?, ?) "
	statement := "insert into " + tblName + " (event_id, event_name, payload, metadata, created_at) values " + values
	if l := len(events); l > 1 {
		statement += strings.Repeat(", "+values, l-1)
	}

	bindings := []interface{}{}
	for _, event := range events {
		eJ, err := json.Marshal(event.Data())
		if err != nil {
			return err
		}

		eM, err := json.Marshal(event.Metadata())
		if err != nil {
			return err
		}

		bindings = append(
			bindings,
			event.MessageID(),
			event.MessageName(),
			string(eJ),
			string(eM),
			event.Created().Format(storeTimeFormat),
		)
	}

	res, err := s.db.ExecContext(ctx, statement, bindings...)
	if err != nil {
		return err
	}

	if ra, err := res.RowsAffected(); err != nil {
		return err
	} else if ra != int64(len(events)) {
		return fmt.Errorf(
			"events persisted (%d) did not match events given (%d)",
			ra,
			len(events),
		)
	}

	return nil
}

// Delete will remove the stream.
func (s *EventStore) Delete(ctx context.Context, streamName string) error {
	tblName, err := getStreamTableName(ctx, s.db, streamName)
	if err != nil {
		return err
	}

	if _, err := s.db.ExecContext(ctx, "drop table if exists "+tblName); err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, "delete from event_streams where real_stream_name = ?", streamName)
	return err
}

// UpdateStreamMetadata sets the metadata for the given stream name.
func (s *EventStore) UpdateStreamMetadata(ctx context.Context, streamName string, newMetadata eventstore.StreamMetadata) error {
	return nil
}
