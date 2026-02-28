# Phase 0 - Step 2: API設計

## 概要

| 項目 | 内容 |
|---|---|
| 目的 | ヘルスチェックサービスのREST APIエンドポイントを定義する |
| 状態 | ✅ 完了 |
| 成果物 | このドキュメント（API仕様書） |

---

## リソース定義

本サービスで扱うリソースは2つ:

| リソース | 説明 | 例 |
|---|---|---|
| **Monitor** | 監視対象のURL設定 | 「本番サーバー」https://example.com/health を60秒ごとに監視 |
| **Check Result** | Monitorに対するヘルスチェックの実行結果 | 2026-03-01 10:00 → 200 OK, 150ms |

**親子関係:**
```
Monitor（親）
  └── Check Result（子）← このMonitorのチェック結果（1対多）
```

---

## エンドポイント一覧

| メソッド | エンドポイント | 説明 | 成功 | エラー |
|---|---|---|---|---|
| `POST` | `/api/monitors` | 監視対象の登録 | 201 | 422 |
| `GET` | `/api/monitors` | 監視対象の一覧取得 | 200 | - |
| `GET` | `/api/monitors/:id` | 監視対象の詳細取得 | 200 | 404 |
| `PUT` | `/api/monitors/:id` | 監視対象の更新 | 200 | 404, 422 |
| `DELETE` | `/api/monitors/:id` | 監視対象の削除 | 204 | 404 |
| `GET` | `/api/monitors/:id/results` | チェック結果の一覧取得 | 200 | 404 |
| `GET` | `/api/monitors/:id/results/:result_id` | チェック結果の詳細取得 | 200 | 404 |
| `GET` | `/api/health` | APIサーバー自体の死活確認 | 200 | - |

---

## 詳細仕様

### POST /api/monitors

監視対象を新規登録する。

**リクエスト:**
```json
{
  "name": "本番サーバー",
  "url": "https://example.com/health",
  "interval_seconds": 60,
  "timeout_seconds": 10
}
```

**レスポンス: 201 Created**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "本番サーバー",
  "url": "https://example.com/health",
  "interval_seconds": 60,
  "timeout_seconds": 10,
  "is_active": true,
  "created_at": "2026-03-01T10:00:00Z",
  "updated_at": "2026-03-01T10:00:00Z"
}
```

**エラー: 422 Unprocessable Entity**
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "入力内容に問題があります",
    "details": [
      { "field": "url", "message": "有効なURLを入力してください" },
      { "field": "interval_seconds", "message": "10以上の値を指定してください" }
    ]
  }
}
```

---

### GET /api/monitors

監視対象の一覧を取得する。

**クエリパラメータ:**

| パラメータ | 型 | デフォルト | 説明 |
|---|---|---|---|
| `status` | string | - | `healthy`, `unhealthy` で絞り込み |
| `limit` | int | 20 | 取得件数（最大100） |
| `offset` | int | 0 | 取得開始位置 |

**レスポンス: 200 OK**
```json
{
  "monitors": [
    {
      "id": "550e8400-...",
      "name": "本番サーバー",
      "url": "https://example.com/health",
      "interval_seconds": 60,
      "timeout_seconds": 10,
      "is_active": true,
      "current_status": "healthy",
      "last_checked_at": "2026-03-01T10:05:00Z",
      "created_at": "2026-03-01T10:00:00Z",
      "updated_at": "2026-03-01T10:00:00Z"
    }
  ],
  "pagination": {
    "total": 45,
    "limit": 20,
    "offset": 0
  }
}
```

---

### GET /api/monitors/:id

監視対象の詳細を取得する。

**レスポンス: 200 OK**
```json
{
  "id": "550e8400-...",
  "name": "本番サーバー",
  "url": "https://example.com/health",
  "interval_seconds": 60,
  "timeout_seconds": 10,
  "is_active": true,
  "current_status": "healthy",
  "last_checked_at": "2026-03-01T10:05:00Z",
  "created_at": "2026-03-01T10:00:00Z",
  "updated_at": "2026-03-01T10:00:00Z"
}
```

**エラー: 404 Not Found**
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "指定された監視対象が見つかりません"
  }
}
```

---

### PUT /api/monitors/:id

監視対象を更新する（全フィールドを送信）。

**リクエスト:**
```json
{
  "name": "本番サーバー（更新後）",
  "url": "https://example.com/health",
  "interval_seconds": 120,
  "timeout_seconds": 10,
  "is_active": false
}
```

**レスポンス: 200 OK**
```json
{
  "id": "550e8400-...",
  "name": "本番サーバー（更新後）",
  "url": "https://example.com/health",
  "interval_seconds": 120,
  "timeout_seconds": 10,
  "is_active": false,
  "current_status": "healthy",
  "last_checked_at": "2026-03-01T10:05:00Z",
  "created_at": "2026-03-01T10:00:00Z",
  "updated_at": "2026-03-01T10:06:00Z"
}
```

---

### DELETE /api/monitors/:id

監視対象を削除する。紐づくCheck Resultも一括削除される。

**レスポンス: 204 No Content**

（レスポンスボディなし）

---

### GET /api/monitors/:id/results

特定の監視対象のチェック結果一覧を取得する（ネストURL）。

**クエリパラメータ:**

| パラメータ | 型 | デフォルト | 説明 |
|---|---|---|---|
| `limit` | int | 50 | 取得件数（最大200） |
| `offset` | int | 0 | 取得開始位置 |

**レスポンス: 200 OK**
```json
{
  "results": [
    {
      "id": "result-001",
      "monitor_id": "550e8400-...",
      "status_code": 200,
      "response_time_ms": 150,
      "is_healthy": true,
      "checked_at": "2026-03-01T10:05:00Z"
    },
    {
      "id": "result-002",
      "monitor_id": "550e8400-...",
      "status_code": 503,
      "response_time_ms": 5000,
      "is_healthy": false,
      "error_message": "Service Unavailable",
      "checked_at": "2026-03-01T10:06:00Z"
    }
  ],
  "pagination": {
    "total": 1440,
    "limit": 50,
    "offset": 0
  }
}
```

---

### GET /api/monitors/:id/results/:result_id

チェック結果の詳細を取得する。

**レスポンス: 200 OK**
```json
{
  "id": "result-002",
  "monitor_id": "550e8400-...",
  "status_code": 503,
  "response_time_ms": 5000,
  "is_healthy": false,
  "error_message": "Service Unavailable",
  "checked_at": "2026-03-01T10:06:00Z"
}
```

---

### GET /api/health

**このAPIサーバー自体** の死活確認用エンドポイント。

> 注意: 監視対象のヘルスチェックではなく、我々のサービスが正常に動作しているかを確認するもの。
> 外部の監視ツール（UptimeRobotなど）がこのエンドポイントを叩いて、我々のサービス自体を監視する。

**レスポンス: 200 OK**
```json
{
  "status": "ok",
  "timestamp": "2026-03-01T10:05:00Z"
}
```

---

## 設計判断の記録

| 判断 | 選択 | 理由 |
|---|---|---|
| URL構造 | リソースベース（名詞） | REST原則に従い、予測しやすいAPI |
| Check ResultのURL | ネストURL `/monitors/:id/results` | Check Resultは必ずMonitorに属する親子関係 |
| Check ResultのCRUD | GETのみ（POST/PUT/DELETE なし） | チェック結果はサーバーが自動生成。手動作成・改ざん・削除不要 |
| レスポンス形式 | オブジェクトで包む `{ "monitors": [...] }` | paginationなどメタ情報を後から追加可能 |
| エラー形式 | 構造化エラー `{ "error": { "code", "message", "details" } }` | フロントエンドのエラー表示とSentryのデバッグに有用 |
| ページネーション | offset/limit方式 | シンプルで理解しやすい。学習用として十分 |
| IDの形式 | UUID v4 | 推測不可能でセキュア。連番IDだと他人のデータが推測される |

---

## 座学ノート

### 1. REST APIとは

クライアント（フロントエンド）がサーバーに「やってほしいこと」を伝える窓口。図書館の窓口のようなもの:
- 本を借りたい → 窓口で申請 → 書庫から取り出す
- APIも同じ構造: リクエスト → 処理 → レスポンス

### 2. HTTPメソッドとCRUD

| HTTPメソッド | 意味 | DB操作 | べき等性 |
|---|---|---|---|
| `GET` | 取得 | SELECT | ✅ べき等 |
| `POST` | 作成 | INSERT | ❌ べき等でない |
| `PUT` | 更新 | UPDATE | ✅ べき等 |
| `DELETE` | 削除 | DELETE | ✅ べき等 |

**べき等性**: 同じリクエストを何回送っても結果が同じになる性質。GET/PUT/DELETEはリトライ安全。POSTは重複作成の可能性あり。

### 3. リソース指向設計

- URLは **リソース（名詞）** で構成する（動詞ベースはNG）
- HTTPメソッドで「何をするか」を表現する
- 厳格なルールだからこそ命名がブレない → チーム開発で強い

### 4. ネストURL vs フラットURL

- **ネストURL** `/monitors/:id/results`: 親子関係が強い場合に使用
- **フラットURL** `/results?monitor_id=xxx`: 独立性が高い場合に使用
- 今回はCheck Result → Monitor の親子関係が強いためネストを採用

### 5. クエリパラメータ

- URLパス = リソースの特定
- クエリパラメータ = 絞り込み・ページネーション・並び替え
- ページネーション: 大量データを分割取得する仕組み（`limit` + `offset`）

### 6. ステータスコード

- `2xx` 成功（200 OK / 201 Created / 204 No Content）
- `4xx` クライアントエラー（400 / 404 / 422）
- `5xx` サーバーエラー（500）→ Sentryが自動キャプチャする対象

### 7. エラーレスポンスの構造化

フロントエンドが適切なエラーメッセージを表示するために、エラーレスポンスは `code` + `message` + `details` で構造化する。

### 8. /api/health の役割

監視対象のヘルスチェックではなく、**このAPIサーバー自体**の死活確認用。外部の監視ツールが叩く。監視する側も誰かに監視されている入れ子構造。

