package worker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/azu-max/dev-for-learning-1/backend/model"
	"github.com/azu-max/dev-for-learning-1/backend/repository"
	"github.com/azu-max/dev-for-learning-1/backend/service"
	"github.com/getsentry/sentry-go"
)

// HealthChecker は定期的にMonitorのURLをチェックするWorker
type HealthChecker struct {
	monitorService *service.MonitorService
	resultRepo     *repository.CheckResultRepository
	httpClient     *http.Client
	interval       time.Duration
	stop           chan struct{} // 停止シグナルを受け取るチャネル
}

// NewHealthChecker はHealthCheckerを生成する
func NewHealthChecker(
	monitorService *service.MonitorService,
	resultRepo *repository.CheckResultRepository,
	interval time.Duration,
) *HealthChecker {
	return &HealthChecker{
		monitorService: monitorService,
		resultRepo:     resultRepo,
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // 1つのチェックに最大10秒
		},
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Start はWorkerを開始する（goroutineで実行される想定）
func (h *HealthChecker) Start() {
	log.Printf("[Worker] ヘルスチェッカー開始（間隔: %s）", h.interval)

	// 起動直後に1回実行
	h.runChecks()

	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// ticker が発火するたびにチェック実行
			h.runChecks()
		case <-h.stop:
			// 停止シグナルを受けたら終了
			log.Println("[Worker] ヘルスチェッカー停止")
			return
		}
	}
}

// Stop はWorkerを停止する
func (h *HealthChecker) Stop() {
	close(h.stop)
}

// runChecks はアクティブなMonitor全件をチェックする
func (h *HealthChecker) runChecks() {
	ctx := context.Background()

	monitors, err := h.monitorService.GetActiveMonitors(ctx)
	if err != nil {
		log.Printf("[Worker] Monitor取得エラー: %v", err)
		return
	}

	if len(monitors) == 0 {
		log.Println("[Worker] チェック対象のMonitorがありません")
		return
	}

	log.Printf("[Worker] %d件のMonitorをチェック開始", len(monitors))

	for _, m := range monitors {
		result := h.checkOne(m)

		if err := h.resultRepo.Create(ctx, result); err != nil {
			log.Printf("[Worker] 結果保存エラー（%s）: %v", m.Name, err)
			continue
		}

		status := "✅ 正常"
		if !result.IsHealthy {
			status = "❌ 異常"
		}
		log.Printf("[Worker] %s | %s | %dms | HTTP %d",
			m.Name, status, result.ResponseTime, result.StatusCode)
	}
}

// checkOne は1つのMonitorに対してHTTP GETを実行し結果を返す
func (h *HealthChecker) checkOne(m model.Monitor) *model.CheckResult {
	result := &model.CheckResult{
		MonitorID: m.ID,
		CheckedAt: time.Now().UTC(),
	}

	start := time.Now()
	resp, err := h.httpClient.Get(m.URL)
	elapsed := time.Since(start)

	result.ResponseTime = int(elapsed.Milliseconds())

	if err != nil {
		// 接続エラー（タイムアウト、DNS解決失敗など）
		result.StatusCode = 0
		result.IsHealthy = false
		result.ErrorMessage = err.Error()

		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetTag("monitor_name", m.Name)
			scope.SetTag("monitor_url", m.URL)
			scope.SetLevel(sentry.LevelError)
			sentry.CaptureException(fmt.Errorf("health check failed for %s: %w", m.Name, err))
		})
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.IsHealthy = resp.StatusCode >= 200 && resp.StatusCode < 300

	return result
}
