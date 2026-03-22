package service

import (
	"context"
	"errors"
	"net/url"

	"github.com/azu-max/dev-for-learning-1/backend/model"
	"github.com/azu-max/dev-for-learning-1/backend/repository"
)

// ビジネスルールのエラー定義
var (
	ErrNameRequired     = errors.New("name is required")
	ErrURLRequired      = errors.New("url is required")
	ErrURLInvalid       = errors.New("url must be a valid HTTP/HTTPS URL")
	ErrIntervalTooShort = errors.New("interval_seconds must be at least 10")
)

// MonitorService はMonitorのビジネスロジックを担当する
type MonitorService struct {
	repo *repository.MonitorRepository
}

// NewMonitorService はMonitorServiceを生成する
func NewMonitorService(repo *repository.MonitorRepository) *MonitorService {
	return &MonitorService{repo: repo}
}

// CreateMonitor はバリデーション後にMonitorを作成する
func (s *MonitorService) CreateMonitor(ctx context.Context, req model.CreateMonitorRequest) (*model.Monitor, error) {
	// ビジネスルールの適用
	if req.Name == "" {
		return nil, ErrNameRequired
	}
	if req.URL == "" {
		return nil, ErrURLRequired
	}
	// URLが有効な形式かチェック
	parsedURL, err := url.Parse(req.URL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return nil, ErrURLInvalid
	}
	if req.IntervalSeconds < 10 {
		return nil, ErrIntervalTooShort
	}

	return s.repo.Create(ctx, req)
}

// GetAllMonitors はすべてのMonitorを取得する
func (s *MonitorService) GetAllMonitors(ctx context.Context) ([]model.Monitor, error) {
	return s.repo.GetAll(ctx)
}

// GetMonitor は指定IDのMonitorを取得する
func (s *MonitorService) GetMonitor(ctx context.Context, id string) (*model.Monitor, error) {
	return s.repo.GetByID(ctx, id)
}

// GetActiveMonitors はアクティブなMonitorをすべて取得する（Worker用）
func (s *MonitorService) GetActiveMonitors(ctx context.Context) ([]model.Monitor, error) {
	return s.repo.GetAllActive(ctx)
}

// DeleteMonitor は指定IDのMonitorを削除する
func (s *MonitorService) DeleteMonitor(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
