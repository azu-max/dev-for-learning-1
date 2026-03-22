package model

import "time"

// CheckResult はヘルスチェックの実行結果を表す構造体
// 「いつ、どのMonitorに対して、どんな結果だったか」を記録する
type CheckResult struct {
	ID           string    `json:"id"`
	MonitorID    string    `json:"monitor_id"`    // どのMonitorのチェックか
	StatusCode   int       `json:"status_code"`   // HTTPステータスコード（200, 500など）
	ResponseTime int       `json:"response_time"` // 応答時間（ミリ秒）
	IsHealthy    bool      `json:"is_healthy"`    // 正常かどうか
	ErrorMessage string    `json:"error_message"` // エラーがあればその内容
	CheckedAt    time.Time `json:"checked_at"`    // チェック実行日時
}
