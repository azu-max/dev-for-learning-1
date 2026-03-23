# Phase 5: Sentry深掘り - ダッシュボード活用 & エラーキャプチャ設計

## 概要

| 項目 | 内容 |
|---|---|
| 目的 | Phase 2 で後回しにしたSentryの深掘りトピックを学習する |
| 状態 | ⬜ 未着手 |
| 前提 | Phase 4 完了（本番環境でSentryが動いている状態） |
| 学習ポイント | ダッシュボードの実践的な使い方、エラーキャプチャの設計パターン |

## 残タスク

### 1. Sentryダッシュボードの見方を深掘り

- [ ] タグでのフィルタリング（monitor_name / monitor_url の活用）
- [ ] Issues ページの見方（グルーピング、発生頻度、影響ユーザー数）
- [ ] Breadcrumbs（エラーに至るまでの経過）の読み方
- [ ] Discover（カスタムクエリ）でエラー傾向を分析
- [ ] アラートルールの設定（Slack / メール通知）

### 2. エラーキャプチャの設計

- [ ] HTTP 4xx/5xx レスポンスもSentryに送る（現在は接続エラーのみ）
- [ ] エラーレベルの使い分け（Warning vs Error vs Fatal）
  - 例: 4xx → Warning、5xx → Error、接続失敗 → Error
- [ ] Sentry に送るべきエラー / 送らないべきエラーの判断基準
- [ ] Context の充実（リクエストID、チェック回数、前回の状態など）
- [ ] カスタムエラー型の設計（Go の errors パッケージ活用）
- [ ] パフォーマンスモニタリング設定（トランザクショントレーシング）

### 3. フロントエンドMVP実装の振り返り

以下の変更について座学 → クイズで理解を深める。

#### Next.js / React
- [ ] Next.js App Router の仕組み（`app/` ディレクトリ、`layout.tsx` と `page.tsx` の役割）
- [ ] `"use client"` ディレクティブの意味（Server Components vs Client Components）
- [ ] `useEffect` + `setInterval` によるポーリングの仕組みと注意点
- [ ] `useCallback` を使う理由（不要な再レンダリングの防止）
- [ ] CSS Modules の仕組み（クラス名のスコープ化、`.module.css` の命名規則）

#### API 連携
- [ ] Next.js `rewrites` によるプロキシの仕組み（なぜ CORS を回避できるのか）
- [ ] `fetch` API の基本（GET / POST / DELETE、ヘッダー、レスポンス処理）
- [ ] フロントエンドからのエラーハンドリングパターン

#### バックエンド追加API
- [ ] `LEFT JOIN LATERAL` の仕組み（通常の LEFT JOIN との違い、N+1 問題の回避）
- [ ] `sql.NullString` / `sql.NullInt32` 等の Nullable 型の扱い方
- [ ] クエリパラメータ `?include=latest_result` による後方互換の設計判断

#### 実装したファイル一覧（参照用）

| ファイル | 役割 |
|---------|------|
| `frontend/app/types/index.ts` | Monitor, CheckResult 等の型定義 |
| `frontend/app/lib/api.ts` | API クライアント（fetch ラッパー） |
| `frontend/app/components/SummaryCards.tsx` | Healthy / Unhealthy / Total のサマリー表示 |
| `frontend/app/components/MonitorForm.tsx` | Monitor 追加フォーム |
| `frontend/app/components/MonitorCard.tsx` | 個別 Monitor のステータスカード |
| `frontend/app/components/MonitorList.tsx` | Monitor 一覧表示 |
| `frontend/app/page.tsx` | ダッシュボードページ（30秒ポーリング） |
| `frontend/next.config.ts` | rewrites 設定（API プロキシ） |
| `backend/repository/monitor_repository.go` | `GetAllWithLatestResult()`（LEFT JOIN LATERAL） |
| `backend/model/monitor.go` | `MonitorWithLatestResult` 構造体 |

### 4. Docker / docker-compose 設定の振り返り

以下の変更について座学 → クイズで理解を深める。

#### Docker 基礎
- [ ] マルチステージビルドの仕組み（なぜ Stage を分けるのか、イメージサイズへの影響）
- [ ] `CGO_ENABLED=0` の意味（静的バイナリと動的バイナリの違い）
- [ ] `npm ci` vs `npm install` の違い（なぜ Docker / CI では `npm ci` を使うのか）
- [ ] `--omit=dev` で devDependencies を除外する理由
- [ ] `.dockerignore` の役割（なぜ `node_modules` をコピーしないのか）

#### docker-compose
- [ ] `ports` vs `expose` の違い（ホスト公開 vs コンテナ間通信のみ）
- [ ] `depends_on` + `condition: service_healthy` による起動順序制御
- [ ] `restart: unless-stopped` の動作（コンテナが落ちたら自動再起動）
- [ ] DB を `expose` のみにするセキュリティ上の理由

## 背景

Phase 2 で Sentry SDK の導入と基本的なエラーキャプチャ（ヘルスチェック接続エラー → Sentry送信）まで完了。
本番環境（AWS）での動作確認を優先するため、これらの深掘りトピックは Phase 5 に移動した。

## 座学ノート

_Phase 5 開始時に記録_

## 理解度クイズ結果

_Phase 5 完了時に記録_
