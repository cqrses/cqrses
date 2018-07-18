package mysql

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/cqrses/eventstore"

	"github.com/go-sql-driver/mysql"
)

const (
	eventStreamsTable = "" +
		"CREATE TABLE IF NOT EXISTS `event_streams` (" +
		"  `no` BIGINT(20) NOT NULL AUTO_INCREMENT," +
		"  `real_stream_name` VARCHAR(150) NOT NULL," +
		"  `stream_name` CHAR(41) NOT NULL," +
		"  `metadata` JSON," +
		" PRIMARY KEY (`no`)," +
		" UNIQUE KEY `ix_rsn` (`real_stream_name`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;"

	eventStreamTable = "" +
		"CREATE TABLE `{tableName}` (" +
		"    `no` BIGINT(20) NOT NULL AUTO_INCREMENT," +
		"    `event_id` CHAR(36) COLLATE utf8mb4_bin NOT NULL," +
		"    `event_name` VARCHAR(100) COLLATE utf8mb4_bin NOT NULL," +
		"    `payload` JSON NOT NULL," +
		"    `metadata` JSON NOT NULL," +
		"    `created_at` DATETIME(6) NOT NULL," +
		"    `aggregate_version` INT(11) UNSIGNED GENERATED ALWAYS AS (JSON_EXTRACT(metadata, '$.aggregate_version')) STORED NOT NULL," +
		"    `aggregate_id` CHAR(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin GENERATED ALWAYS AS (JSON_UNQUOTE(JSON_EXTRACT(metadata, '$.aggregate_id'))) STORED NOT NULL," +
		"    PRIMARY KEY (`no`)," +
		"    UNIQUE KEY `ix_event_id` (`event_id`)," +
		"    UNIQUE KEY `ix_unique_event` (`aggregate_id`, `aggregate_version`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin"
)

func applyEventStreamsSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, eventStreamsTable)
	return err
}

func createStream(ctx context.Context, db *sql.DB, stream *eventstore.Stream) error {
	tblName := makeStreamTableName(stream.Name)

	meta, err := json.Marshal(stream.Metadata)
	if err != nil {
		return err
	}

	if _, err := db.ExecContext(
		ctx,
		"insert into event_streams (real_stream_name, stream_name, metadata) values (?, ?, ?)",
		stream.Name,
		tblName,
		string(meta),
	); err != nil {
		if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1062 {
			return eventstore.ErrStreamAlreadyExists
		}
		return err
	}

	return createStreamTable(ctx, db, tblName)
}

func createStreamTable(ctx context.Context, db *sql.DB, name string) error {
	statement := strings.Replace(eventStreamTable, "{tableName}", name, 1)
	_, err := db.ExecContext(ctx, statement)
	return err
}

// Create a table name for a stream with the given name.
// We do this to avoid conflicts.
func makeStreamTableName(streamName string) string {
	return "_" + strings.ToUpper(fmt.Sprintf("%x", sha1.Sum([]byte(streamName))))
}

// Get the table name for a stream.
func getStreamTableName(ctx context.Context, db *sql.DB, streamName string) (out string, err error) {
	row := db.QueryRowContext(ctx, "select stream_name from `event_streams` where real_stream_name = ?", streamName)
	err = row.Scan(&out)
	if err == sql.ErrNoRows {
		return "", eventstore.ErrStreamDoesNotExist
	}
	return
}
