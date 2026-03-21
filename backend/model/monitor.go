package model

import "time"

// Monitor はヘルスチェック対象のサービスを表す構造体
// DBの monitors テーブルと1:1で対応する
type Monitor struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	URL             string    `json:"url"`
	IntervalSeconds int       `json:"interval_seconds"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CreateMonitorRequest はMonitor作成時のリクエストボディ
// IDやタイムスタンプはサーバー側で自動生成するので含まない
type CreateMonitorRequest struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	IntervalSeconds int    `json:"interval_seconds"`
}
