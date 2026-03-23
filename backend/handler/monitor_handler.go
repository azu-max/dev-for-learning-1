package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/azu-max/dev-for-learning-1/backend/model"
	"github.com/azu-max/dev-for-learning-1/backend/service"
)

// MonitorHandler はMonitor関連のHTTPハンドラを担当する
type MonitorHandler struct {
	svc *service.MonitorService
}

// NewMonitorHandler はMonitorHandlerを生成する
func NewMonitorHandler(svc *service.MonitorService) *MonitorHandler {
	return &MonitorHandler{svc: svc}
}

// HandleMonitors は /api/monitors へのリクエストを処理する
// GET → 一覧取得、POST → 新規作成
func (h *MonitorHandler) HandleMonitors(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAll(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleMonitorByID は /api/monitors/{id} へのリクエストを処理する
// GET → 詳細取得、DELETE → 削除
func (h *MonitorHandler) HandleMonitorByID(w http.ResponseWriter, r *http.Request) {
	// URLから ID を取り出す: /api/monitors/xxx → xxx
	id := strings.TrimPrefix(r.URL.Path, "/api/monitors/")
	if id == "" {
		http.Error(w, "monitor id is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getByID(w, r, id)
	case http.MethodDelete:
		h.delete(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getAll はすべてのMonitorを返す
// include=latest_result クエリパラメータで最新チェック結果も含められる
func (h *MonitorHandler) getAll(w http.ResponseWriter, r *http.Request) {
	include := r.URL.Query().Get("include")

	if include == "latest_result" {
		results, err := h.svc.GetAllMonitorsWithLatestResult(r.Context())
		if err != nil {
			log.Printf("ERROR: failed to get monitors with results: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if results == nil {
			results = []model.MonitorWithLatestResult{}
		}
		respondJSON(w, http.StatusOK, results)
		return
	}

	monitors, err := h.svc.GetAllMonitors(r.Context())
	if err != nil {
		log.Printf("ERROR: failed to get monitors: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if monitors == nil {
		monitors = []model.Monitor{}
	}

	respondJSON(w, http.StatusOK, monitors)
}

// create は新しいMonitorを作成する
func (h *MonitorHandler) create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateMonitorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	monitor, err := h.svc.CreateMonitor(r.Context(), req)
	if err != nil {
		// Service層のビジネスエラーは 400 で返す
		if errors.Is(err, service.ErrNameRequired) ||
			errors.Is(err, service.ErrURLRequired) ||
			errors.Is(err, service.ErrURLInvalid) ||
			errors.Is(err, service.ErrIntervalTooShort) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("ERROR: failed to create monitor: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusCreated, monitor)
}

// getByID は指定IDのMonitorを返す
func (h *MonitorHandler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	monitor, err := h.svc.GetMonitor(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Monitor not found", http.StatusNotFound)
			return
		}
		log.Printf("ERROR: failed to get monitor: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, monitor)
}

// delete は指定IDのMonitorを削除する
func (h *MonitorHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	err := h.svc.DeleteMonitor(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Monitor not found", http.StatusNotFound)
			return
		}
		log.Printf("ERROR: failed to delete monitor: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204: 成功、返すボディなし
}

// respondJSON はJSONレスポンスを返すヘルパー
func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
