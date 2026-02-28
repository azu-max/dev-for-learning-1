# Phase 0 - Step 3: データベース設計

## 概要

| 項目 | 内容 |
|---|---|
| 目的 | ヘルスチェックサービスのテーブル定義を行う（最小限） |
| 状態 | ✅ 完了 |
| 成果物 | このドキュメント（テーブル定義書） |

---

## テーブル一覧

| テーブル | 説明 | API リソース |
|---|---|---|
| `monitors` | 監視対象のURL設定 | Monitor |
| `check_results` | ヘルスチェックの実行結果 | Check Result |

---

## ER図

```
┌─────────────────────┐       ┌─────────────────────────┐
│     monitors         │       │    check_results         │
├─────────────────────┤       ├─────────────────────────┤
│ id (PK, UUID)        │──┐   │ id (PK, UUID)            │
│ name (VARCHAR)       │  │   │ monitor_id (FK, UUID)    │
│ url (TEXT)           │  └──→│   REFERENCES monitors(id) │
│ interval_seconds     │      │ status_code (INT)        │
│ timeout_seconds      │      │ response_time_ms (INT)   │
│ is_active (BOOL)     │      │ is_healthy (BOOL)        │
│ current_status       │      │ error_message (TEXT)      │
│ last_checked_at      │      │ checked_at (TIMESTAMPTZ) │
│ created_at           │      └─────────────────────────┘
│ updated_at           │
└─────────────────────┘        1 Monitor : N Check Results
                                    (1対多)
```

---

## テーブル定義

### monitors

| カラム | 型 | 制約 | デフォルト | 説明 |
|---|---|---|---|---|
| id | UUID | PK | gen_random_uuid() | 一意識別子 |
| name | VARCHAR(255) | NOT NULL | - | 監視対象の名前（例: 「本番サーバー」） |
| url | TEXT | NOT NULL | - | 監視対象のURL |
| interval_seconds | INTEGER | NOT NULL | 60 | チェック間隔（秒）。何秒おきにチェックするか |
| timeout_seconds | INTEGER | NOT NULL | 10 | タイムアウト（秒）。何秒待って返事がなければ異常とするか |
| is_active | BOOLEAN | NOT NULL | true | 監視が有効かどうか |
| current_status | VARCHAR(20) | NOT NULL | 'unknown' | 現在のステータス（healthy / unhealthy / unknown） |
| last_checked_at | TIMESTAMPTZ | - | NULL | 最後にチェックした日時 |
| created_at | TIMESTAMPTZ | NOT NULL | NOW() | 作成日時 |
| updated_at | TIMESTAMPTZ | NOT NULL | NOW() | 更新日時 |

```sql
CREATE TABLE monitors (
    id               UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name             VARCHAR(255) NOT NULL,
    url              TEXT         NOT NULL,
    interval_seconds INTEGER      NOT NULL DEFAULT 60,
    timeout_seconds  INTEGER      NOT NULL DEFAULT 10,
    is_active        BOOLEAN      NOT NULL DEFAULT true,
    current_status   VARCHAR(20)  NOT NULL DEFAULT 'unknown',
    last_checked_at  TIMESTAMPTZ,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
```

### check_results

| カラム | 型 | 制約 | デフォルト | 説明 |
|---|---|---|---|---|
| id | UUID | PK | gen_random_uuid() | 一意識別子 |
| monitor_id | UUID | FK, NOT NULL | - | 紐づくMonitorのID |
| status_code | INTEGER | - | NULL | HTTPレスポンスのステータスコード |
| response_time_ms | INTEGER | - | NULL | レスポンス時間（ミリ秒） |
| is_healthy | BOOLEAN | NOT NULL | - | 正常かどうかの判定結果 |
| error_message | TEXT | - | NULL | エラー時のメッセージ（正常時はNULL） |
| checked_at | TIMESTAMPTZ | NOT NULL | NOW() | チェック実行日時 |

```sql
CREATE TABLE check_results (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    monitor_id       UUID        NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    status_code      INTEGER,
    response_time_ms INTEGER,
    is_healthy       BOOLEAN     NOT NULL,
    error_message    TEXT,
    checked_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_check_results_monitor_id ON check_results(monitor_id);
CREATE INDEX idx_check_results_checked_at ON check_results(checked_at DESC);
```

---

## インデックス

| インデックス | カラム | 理由 |
|---|---|---|
| `idx_check_results_monitor_id` | monitor_id | 「このMonitorの結果一覧」を高速取得（GET /api/monitors/:id/results） |
| `idx_check_results_checked_at` | checked_at DESC | 「最新の結果から取得」を高速化 |

---

## 設計判断の記録

| 判断 | 選択 | 理由 |
|---|---|---|
| 主キーの型 | UUID | 連番（INT）だとURLから他のリソースを推測可能。UUIDなら推測不可能 |
| 外部キー制約 | REFERENCES monitors(id) | 存在しないMonitorのCheck Resultが作られるのを防ぐ |
| 削除時の振る舞い | ON DELETE CASCADE | Monitor削除時にCheck Resultも自動削除。孤児データを防ぐ |
| タイムスタンプ型 | TIMESTAMPTZ | タイムゾーン情報を保持。どこからアクセスしても正しい時刻 |
| current_status の管理 | monitors テーブルに持つ | 一覧取得時に毎回check_resultsを集計するのを避ける |

---

## 座学ノート

### 1. APIリソースからテーブルへの対応

REST APIをリソース指向で設計しておくと、テーブル設計がほぼ自動的に決まる。APIの2つのリソース（Monitor, Check Result）がそのまま2つのテーブルになった。

### 2. UUIDを主キーに使う理由

INT（連番）は推測可能（id=1の次はid=2）。UUIDは推測不可能でセキュア。PostgreSQLの `gen_random_uuid()` で自動生成できる。

### 3. interval_seconds と timeout_seconds

- `interval_seconds`: 何秒おきにヘルスチェックを実行するか
- `timeout_seconds`: 何秒待っても返事がなければ異常とみなすか。これがないとサーバーが応答しない場合に永遠に待ち続けてしまう

### 4. 外部キー制約（REFERENCES）

`monitor_id REFERENCES monitors(id)` = monitor_idの値はmonitorsテーブルに実際に存在するIDでなければならない。存在しないMonitorのCheck Resultが作られるのを防ぐ。

### 5. ON DELETE CASCADE

親（Monitor）が削除されたら子（Check Result）も自動削除される。API設計の「DELETE /api/monitors/:id で紐づくCheck Resultも削除」をDBレベルで保証する仕組み。

### 6. インデックス（INDEX）

頻繁に検索に使うカラムにインデックスを貼ると検索が高速化される。インデックスなしだとテーブル全件を走査（フルスキャン）する必要がある。
