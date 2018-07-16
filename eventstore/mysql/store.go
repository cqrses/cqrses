package mysql

import (
	"context"
	"database/sql"

	"gopkg.in/cqrses/eventstore"
	"gopkg.in/cqrses/messages"

	// MySQL driver for the database/sql package.
	"github.com/go-sql-driver/mysql"
)

type (
	// EventStore will use a MySQL database to manage streams.
	EventStore struct {
		db *sql.DB
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

	return &EventStore{
		db: db,
	}, nil
}

// Ping tests connection to the database is still ok.
func (s *EventStore) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// Load events from the given stream name.
func (s *EventStore) Load(ctx context.Context, streamName string, from, count uint64, matcher eventstore.MetadataMatcher) eventstore.StreamIterator {
	return &StreamIterator{}
}

// LoadReverse Loads events from the given stream name in reverse.
func (s *EventStore) LoadReverse(ctx context.Context, streamName string, from, count uint64, matcher eventstore.MetadataMatcher) eventstore.StreamIterator {
	return &StreamIterator{}
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
func (s *EventStore) Create(ctx context.Context, stream *eventstore.Stream) error {
	return nil
}

// AppendTo will append events to the stream.
func (s *EventStore) AppendTo(ctx context.Context, streamName string, events []*messages.Event) error {
	return nil
}

// Delete will remove the stream.
func (s *EventStore) Delete(ctx context.Context, streamName string) error {
	return nil
}

// UpdateStreamMetadata sets the metadata for the given stream name.
func (s *EventStore) UpdateStreamMetadata(ctx context.Context, streamName string, newMetadata eventstore.StreamMetadata) error {
	return nil
}
