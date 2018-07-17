package mysql

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"strings"

	"gopkg.in/cqrses/eventstore"
)

const (
	eventStreamsTable = "" +
		"CREATE TABLE IF NOT EXISTS `event_streams` (" +
		"  `no` BIGINT(20) NOT NULL AUTO_INCREMENT," +
		"  `real_stream_name` VARCHAR(150) NOT NULL," +
		"  `stream_name` CHAR(41) NOT NULL," +
		"  `metadata` JSON" +
		" PRIMARY KEY (`no`)," +
		" UNIQUE KEY `ix_rsn` (`real_stream_name`)," +
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

func createStreamTable(ctx context.Context, db *sql.DB, stream *eventstore.Stream) error {
	statement := strings.Replace(eventStreamTable, "{tableName}", makeStreamTableName(stream.Name), 1)
	_, err := db.ExecContext(ctx, statement)
	return err
}

// Create a table name for a stream with the given name.
// We do this to avoid conflicts.
func makeStreamTableName(streamName string) string {
	return "_" + strings.ToUpper(fmt.Sprintf("%x", sha256.Sum256([]byte(streamName))))
}

// Get the table name for a stream.
func getStreamTableName(ctx context.Context, db *sql.DB, streamName string) (out string, err error) {
	row := db.QueryRowContext(ctx, "select stream_name from `event_streams` where real_stream_name = ?", streamName)
	err = row.Scan(&out)
	return
}
