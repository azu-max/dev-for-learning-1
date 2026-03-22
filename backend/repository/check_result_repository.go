package repository

import (
	"context"
	"database/sql"

	"github.com/azu-max/dev-for-learning-1/backend/model"
	"github.com/google/uuid"
)

// CheckResultRepository はチェック結果のDB操作を担当する
type CheckResultRepository struct {
	db *sql.DB
}

// NewCheckResultRepository はCheckResultRepositoryを生成する
func NewCheckResultRepository(db *sql.DB) *CheckResultRepository {
	return &CheckResultRepository{db: db}
}

// Create はチェック結果をDBに保存する
func (r *CheckResultRepository) Create(ctx context.Context, result *model.CheckResult) error {
	result.ID = uuid.New().String()

	query := `
		INSERT INTO check_results (id, monitor_id, status_code, response_time, is_healthy, error_message, checked_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		result.ID,
		result.MonitorID,
		result.StatusCode,
		result.ResponseTime,
		result.IsHealthy,
		result.ErrorMessage,
		result.CheckedAt,
	)
	return err
}
