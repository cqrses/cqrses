package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"sync"

	"github.com/go-cqrses/cqrses/projection"
)

type (
	// ProjectionManager ...
	ProjectionManager struct {
		es *EventStore
	}
)

// NewProjectionManager will get a projection manager that uses the MySQL backend
// to store projection states.
func NewProjectionManager(es *EventStore) projection.Manager {
	return &ProjectionManager{
		es: es,
	}
}

// Create ...
func (m *ProjectionManager) Create(ctx context.Context, name string, opts []projection.ProjectorOpt) (projection.Projector, error) {
	options, err := projection.BuildOptionsFrom(opts)
	if err != nil {
		return nil, err
	}
	return &StreamProjection{
		name:        name,
		es:          m.es,
		streamNames: []string{},
		opts:        options,
		handlers:    map[string][]projection.Handler{},
		any:         []projection.Handler{},
		close:       make(chan struct{}, 1),
		modLock:     &sync.Mutex{},
	}, nil
}

// Delete ...
func (m *ProjectionManager) Delete(ctx context.Context, projectionName string) error {
	_, err := m.es.db.ExecContext(ctx, "delete from projections where name = ?", projectionName)
	return err
}

// Reset ...
func (m *ProjectionManager) Reset(ctx context.Context, projectionName string) error {
	_, err := m.es.db.ExecContext(ctx, "updates projections set position = 0 where name = ?", projectionName)
	return err
}

// Stop ...
func (m *ProjectionManager) Stop(ctx context.Context, projectionName string) error {
	_, err := m.es.db.ExecContext(ctx, "updates projections set status = ? where name = ?", projection.StatusStopping, projectionName)
	return err
}

// FetchProjectionNames ...
func (m *ProjectionManager) FetchProjectionNames(ctx context.Context, filter string, start, limit uint64) ([]string, error) {
	var res *sql.Rows
	var err error
	out := make([]string, 0, limit)

	if len(filter) == 0 {
		res, err = m.es.db.QueryContext(ctx, "select `name` from projections limit ?,?", start, limit)
	} else {
		res, err = m.es.db.QueryContext(ctx, "select `name` from where name like ? limit ?,?", filter, start, limit)
	}

	if err == nil {
		for res.Next() {
			var name string
			if err = res.Scan(&name); err != nil {
				return out, err
			}
			out = append(out, name)
		}
	}

	return out, err
}

// FetchPojectionStatus ...
func (m *ProjectionManager) FetchPojectionStatus(ctx context.Context, projectionName string) (projection.Status, error) {
	var status projection.Status
	row := m.es.db.QueryRowContext(ctx, "select `status` from projections where name = ?", projectionName)
	err := row.Scan(&status)
	return status, err
}

// FetchPojectionStreamPositions ...
func (m *ProjectionManager) FetchPojectionStreamPositions(ctx context.Context, projectionName string) (projection.StreamPositions, error) {
	return fetchPojectionStreamPositions(ctx, m.es.db, projectionName)
}

func fetchPojectionStreamPositions(ctx context.Context, db *sql.DB, projectionName string) (projection.StreamPositions, error) {
	row := db.QueryRowContext(ctx, "select `position` from projections where name = ?", projectionName)

	var rawPositions string
	var positions projection.StreamPositions

	if err := row.Scan(&rawPositions); err != nil {
		return positions, err
	}

	return positions, json.Unmarshal([]byte(rawPositions), &positions)
}
