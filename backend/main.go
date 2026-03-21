package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/azu-max/dev-for-learning-1/backend/handler"
	"github.com/azu-max/dev-for-learning-1/backend/repository"
	"github.com/azu-max/dev-for-learning-1/backend/service"

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
		log.Fatalf("Failed to create table: %v", err)
	}
	log.Println("Database table ready")

	// --- レイヤーの組み立て ---
	// Repository → Service → Handler の順に生成
	monitorRepo := repository.NewMonitorRepository(db)
	monitorSvc := service.NewMonitorService(monitorRepo)
	monitorHandler := handler.NewMonitorHandler(monitorSvc)

	// --- ルーティング ---
	http.HandleFunc("/api/health", healthHandler)
	http.HandleFunc("/api/monitors", monitorHandler.HandleMonitors)
	http.HandleFunc("/api/monitors/", monitorHandler.HandleMonitorByID)

	log.Println("API server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
