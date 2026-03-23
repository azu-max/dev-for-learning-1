package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/azu-max/dev-for-learning-1/backend/model"
	"github.com/google/uuid"
)

// MonitorRepository はMonitorのDB操作を担当する
type MonitorRepository struct {
	db *sql.DB
}

// NewMonitorRepository はMonitorRepositoryを生成する
func NewMonitorRepository(db *sql.DB) *MonitorRepository {
	return &MonitorRepository{db: db}
}

// Create は新しいMonitorをDBに保存する
func (r *MonitorRepository) Create(ctx context.Context, req model.CreateMonitorRequest) (*model.Monitor, error) {
	monitor := &model.Monitor{
		ID:              uuid.New().String(),
		Name:            req.Name,
		URL:             req.URL,
		IntervalSeconds: req.IntervalSeconds,
		IsActive:        true,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	query := `
		INSERT INTO monitors (id, name, url, interval_seconds, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		monitor.ID,
		monitor.Name,
		monitor.URL,
		monitor.IntervalSeconds,
		monitor.IsActive,
		monitor.CreatedAt,
		monitor.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return monitor, nil
}

// GetAll はすべてのMonitorをDBから取得する
func (r *MonitorRepository) GetAll(ctx context.Context) ([]model.Monitor, error) {
	query := `SELECT id, name, url, interval_seconds, is_active, created_at, updated_at FROM monitors ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monitors []model.Monitor
	for rows.Next() {
		var m model.Monitor
		err := rows.Scan(&m.ID, &m.Name, &m.URL, &m.IntervalSeconds, &m.IsActive, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		monitors = append(monitors, m)
	}

	return monitors, nil
}

// GetByID は指定IDのMonitorをDBから取得する
func (r *MonitorRepository) GetByID(ctx context.Context, id string) (*model.Monitor, error) {
	query := `SELECT id, name, url, interval_seconds, is_active, created_at, updated_at FROM monitors WHERE id = $1`

	var m model.Monitor
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.Name, &m.URL, &m.IntervalSeconds, &m.IsActive, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// GetAllActive はアクティブなMonitorをすべて取得する（Worker用）
func (r *MonitorRepository) GetAllActive(ctx context.Context) ([]model.Monitor, error) {
	query := `SELECT id, name, url, interval_seconds, is_active, created_at, updated_at
		FROM monitors WHERE is_active = true ORDER BY created_at`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monitors []model.Monitor
	for rows.Next() {
		var m model.Monitor
		err := rows.Scan(&m.ID, &m.Name, &m.URL, &m.IntervalSeconds, &m.IsActive, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		monitors = append(monitors, m)
	}
	return monitors, nil
}

// GetAllWithLatestResult はすべてのMonitorと最新チェック結果を取得する
func (r *MonitorRepository) GetAllWithLatestResult(ctx context.Context) ([]model.MonitorWithLatestResult, error) {
	query := `
		SELECT
			m.id, m.name, m.url, m.interval_seconds, m.is_active, m.created_at, m.updated_at,
			cr.id, cr.status_code, cr.response_time, cr.is_healthy, cr.error_message, cr.checked_at
		FROM monitors m
		LEFT JOIN LATERAL (
			SELECT * FROM check_results
			WHERE monitor_id = m.id
			ORDER BY checked_at DESC
			LIMIT 1
		) cr ON true
		ORDER BY m.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.MonitorWithLatestResult
	for rows.Next() {
		var mwr model.MonitorWithLatestResult
		var crID, crErrorMessage sql.NullString
		var crStatusCode, crResponseTime sql.NullInt32
		var crIsHealthy sql.NullBool
		var crCheckedAt sql.NullTime

		err := rows.Scan(
			&mwr.ID, &mwr.Name, &mwr.URL, &mwr.IntervalSeconds, &mwr.IsActive, &mwr.CreatedAt, &mwr.UpdatedAt,
			&crID, &crStatusCode, &crResponseTime, &crIsHealthy, &crErrorMessage, &crCheckedAt,
		)
		if err != nil {
			return nil, err
		}

		if crID.Valid {
			mwr.LatestResult = &model.CheckResult{
				ID:           crID.String,
				MonitorID:    mwr.ID,
				StatusCode:   int(crStatusCode.Int32),
				ResponseTime: int(crResponseTime.Int32),
				IsHealthy:    crIsHealthy.Bool,
				ErrorMessage: crErrorMessage.String,
				CheckedAt:    crCheckedAt.Time,
			}
		}

		results = append(results, mwr)
	}

	return results, nil
}

// Delete は指定IDのMonitorをDBから削除する
func (r *MonitorRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM monitors WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
