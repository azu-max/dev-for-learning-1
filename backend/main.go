package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/azu-max/dev-for-learning-1/backend/handler"
	"github.com/azu-max/dev-for-learning-1/backend/repository"
	"github.com/azu-max/dev-for-learning-1/backend/service"
	"github.com/azu-max/dev-for-learning-1/backend/worker"

	_ "github.com/lib/pq" // PostgreSQLドライバ（init()だけ使う）
)

// HealthResponse はヘルスチェックのレスポンス
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func main() {
	// --- DB接続 ---
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// DB接続テスト
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database")

	// --- テーブル作成（開発用の自動マイグレーション）---
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS monitors (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		url TEXT NOT NULL,
		interval_seconds INTEGER NOT NULL DEFAULT 60,
		is_active BOOLEAN NOT NULL DEFAULT true,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatalf("Failed to create monitors table: %v", err)
	}

	createCheckResultsSQL := `
	CREATE TABLE IF NOT EXISTS check_results (
		id TEXT PRIMARY KEY,
		monitor_id TEXT NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
		status_code INTEGER NOT NULL DEFAULT 0,
		response_time INTEGER NOT NULL DEFAULT 0,
		is_healthy BOOLEAN NOT NULL DEFAULT false,
		error_message TEXT NOT NULL DEFAULT '',
		checked_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`
	if _, err := db.Exec(createCheckResultsSQL); err != nil {
		log.Fatalf("Failed to create check_results table: %v", err)
	}
	log.Println("Database tables ready")

	// --- レイヤーの組み立て ---
	// Repository → Service → Handler の順に生成
	monitorRepo := repository.NewMonitorRepository(db)
	monitorSvc := service.NewMonitorService(monitorRepo)
	monitorHandler := handler.NewMonitorHandler(monitorSvc)

	// CheckResult用のRepository
	checkResultRepo := repository.NewCheckResultRepository(db)

	// --- Worker起動 ---
	checker := worker.NewHealthChecker(monitorSvc, checkResultRepo, 30*time.Second)
	go checker.Start() // goroutineでバックグラウンド実行

	// --- ルーティング ---
	http.HandleFunc("/api/health", healthHandler)
	http.HandleFunc("/api/monitors", monitorHandler.HandleMonitors)
	http.HandleFunc("/api/monitors/", monitorHandler.HandleMonitorByID)

	// --- HTTPサーバー起動 ---
	server := &http.Server{Addr: ":8080"}

	// goroutineでサーバー起動（ブロッキングなので別スレッドで）
	go func() {
		log.Println("API server starting on :8080")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// --- グレースフルシャットダウン ---
	// Ctrl+C や docker stop のシグナルを待つ
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // シグナルが来るまでここでブロック

	log.Println("Shutting down...")
	checker.Stop() // Workerを停止

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}
