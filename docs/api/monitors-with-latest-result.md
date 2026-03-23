# API設計: Monitor一覧 + 最新チェック結果

## 背景

現在の `GET /api/monitors` は Monitor の情報のみ返す。
フロントエンドのダッシュボードでは各 Monitor の「正常/異常」「レスポンスタイム」を表示する必要があるため、最新のチェック結果も一緒に返す API が必要。

## 設計

### エンドポイント

```
GET /api/monitors?include=latest_result
```

### 方針

- `include` クエリパラメータがない場合 → 従来通り `Monitor[]` を返す（後方互換）
- `include=latest_result` がある場合 → `MonitorWithLatestResult[]` を返す

### レスポンス（`include=latest_result`）

```json
[
  {
    "id": "abc-123",
    "name": "Google",
    "url": "https://google.com",
    "interval_seconds": 30,
    "is_active": true,
    "created_at": "2026-03-22T10:00:00Z",
    "updated_at": "2026-03-22T10:00:00Z",
    "latest_result": {
      "id": "res-456",
      "monitor_id": "abc-123",
      "status_code": 200,
      "response_time": 312,
      "is_healthy": true,
      "error_message": "",
      "checked_at": "2026-03-23T03:42:46Z"
    }
  },
  {
    "id": "def-789",
    "name": "New Service",
    "url": "https://example.com",
    "interval_seconds": 60,
    "is_active": true,
    "created_at": "2026-03-23T01:00:00Z",
    "updated_at": "2026-03-23T01:00:00Z",
    "latest_result": null
  }
]
```

- `latest_result` はチェック結果がまだない場合は `null`

### 実装方針

```
Handler → Service → Repository
                       │
                       ▼
          1回のSQLで monitors + 最新の check_results を JOIN して取得
          （N+1 問題を避ける）
```

### SQL（LEFT JOIN + サブクエリ）

```sql
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
ORDER BY m.created_at DESC
```

- `LEFT JOIN LATERAL` で各 Monitor に対して最新1件だけ取得
- チェック結果がない Monitor も `LEFT JOIN` なので含まれる（`cr.*` が NULL になる）

### 変更箇所

| レイヤー | ファイル | 変更内容 |
|---------|---------|---------|
| Model | `model/monitor.go` | `MonitorWithLatestResult` 構造体を追加 |
| Repository | `repository/monitor_repository.go` | `GetAllWithLatestResult()` メソッドを追加 |
| Service | `service/monitor_service.go` | `GetAllMonitorsWithLatestResult()` メソッドを追加 |
| Handler | `handler/monitor_handler.go` | `include` クエリパラメータの分岐を追加 |
