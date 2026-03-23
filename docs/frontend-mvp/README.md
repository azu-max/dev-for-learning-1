# Frontend MVP 設計書

## 概要

| 項目 | 内容 |
|---|---|
| 目的 | Monitor の CRUD と ヘルスチェック結果をブラウザから操作・確認できるようにする |
| 位置づけ | MVP（仮UI）。Phase 3 の AWS デプロイ前に最低限の操作性を確保する |
| 技術 | Next.js 15 (App Router) + React 19 + TypeScript |
| スタイリング | CSS Modules（追加ライブラリなし） |

## 画面構成

```
┌──────────────────────────────────────────────────┐
│  🟢 Health Check Monitor              [API: ok]  │  ← ヘッダー（APIステータス表示）
├──────────────────────────────────────────────────┤
│                                                  │
│  ┌─ サマリーカード ───────────────────────────┐  │
│  │  ✅ 3 Healthy    ❌ 1 Unhealthy    📊 4 Total │  │
│  └───────────────────────────────────────────┘  │
│                                                  │
│  ┌─ Monitor 追加フォーム ────────────────────┐  │
│  │  Name: [________]  URL: [________]        │  │
│  │  Interval: [30] sec        [+ 追加]       │  │
│  └───────────────────────────────────────────┘  │
│                                                  │
│  ┌─ Monitor リスト ──────────────────────────┐  │
│  │                                           │  │
│  │  🟢 Google          https://google.com    │  │
│  │     200 OK | 312ms | 30秒前              │  │
│  │                               [🗑 削除]   │  │
│  │  ─────────────────────────────────────    │  │
│  │  🔴 Test Error      http://invalid.test   │  │
│  │     接続エラー | 10024ms | 30秒前         │  │
│  │                               [🗑 削除]   │  │
│  │  ─────────────────────────────────────    │  │
│  │  ⚫ New Service     https://example.com   │  │
│  │     未チェック                            │  │
│  │                               [🗑 削除]   │  │
│  └───────────────────────────────────────────┘  │
│                                                  │
│  ┌─ フッター ────────────────────────────────┐  │
│  │  30秒ごとに自動更新 | 最終更新: 12:34:56  │  │
│  └───────────────────────────────────────────┘  │
└──────────────────────────────────────────────────┘
```

### ステータス表示ルール

| アイコン | 条件 | 色 |
|---------|------|-----|
| 🟢 | 最新チェックが `is_healthy: true` | 緑 (`#10b981`) |
| 🔴 | 最新チェックが `is_healthy: false` | 赤 (`#ef4444`) |
| ⚫ | チェック結果がまだない | グレー (`#6b7280`) |

## コンポーネント構成

```
app/
├── layout.tsx          ... 全体レイアウト（ヘッダー含む）
├── page.tsx            ... ダッシュボードページ（メイン）
├── components/
│   ├── SummaryCards.tsx ... サマリーカード（Healthy / Unhealthy / Total）
│   ├── MonitorForm.tsx ... Monitor 追加フォーム
│   ├── MonitorList.tsx ... Monitor 一覧
│   └── MonitorCard.tsx ... 個別 Monitor カード（ステータス + 最新結果）
├── lib/
│   └── api.ts          ... API クライアント（fetch ラッパー）
└── types/
    └── index.ts        ... 型定義（Monitor, CheckResult）
```

## 型定義

```typescript
// Monitor
type Monitor = {
  id: string;
  name: string;
  url: string;
  interval_seconds: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

// チェック結果
type CheckResult = {
  id: string;
  monitor_id: string;
  status_code: number;
  response_time: number;
  is_healthy: boolean;
  error_message: string;
  checked_at: string;
};

// Monitor + 最新チェック結果（表示用）
type MonitorWithLatestResult = Monitor & {
  latest_result: CheckResult | null;
};
```

## 自動更新

- 30秒間隔でポーリング（Worker のチェック間隔と同じ）
- `setInterval` + `fetch` でシンプルに実装
- 最終更新時刻をフッターに表示

## スタイリング方針

- **CSS Modules** を使用（追加パッケージ不要）
- ダークテーマベース（背景: `#0f172a`、カード: `#1e293b`）
- ステータスカラーで直感的に状態を把握できるように
- レスポンシブ対応は不要（MVP のため）

## カラーパレット

| 用途 | カラー | コード |
|------|--------|--------|
| 背景 | ダークネイビー | `#0f172a` |
| カード背景 | スレートグレー | `#1e293b` |
| テキスト | ホワイト | `#f1f5f9` |
| サブテキスト | グレー | `#94a3b8` |
| 正常 | エメラルド | `#10b981` |
| 異常 | レッド | `#ef4444` |
| 未チェック | グレー | `#6b7280` |
| アクセント | ブルー | `#3b82f6` |
| 削除 | ローズ | `#f43f5e` |

## 実装ステップ

1. **型定義・API クライアント**: `types/index.ts`, `lib/api.ts`
2. **コンポーネント実装**: SummaryCards → MonitorForm → MonitorCard → MonitorList
3. **ページ組み立て**: `page.tsx` でコンポーネントを配置
4. **自動更新**: 30秒ポーリング追加
5. **動作確認**: Docker Compose で確認
