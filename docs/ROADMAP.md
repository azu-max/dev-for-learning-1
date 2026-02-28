# Health Check Monitoring Service - 学習ロードマップ

## プロジェクト概要

| 項目 | 内容 |
|---|---|
| プロジェクト名 | Health Check Monitoring Service（ヘルスチェックモニタリングサービス） |
| 目的 | Sentry + AWSの学習（学習特化） |
| 技術スタック | Go + Next.js + PostgreSQL + Sentry + AWS |
| 学習時間 | 週に数時間 |

## 学習者プロフィール

| 項目 | 内容 |
|---|---|
| 経験 |  フルスタック |
| 業務技術 | Go, AWS, React/Next.js |
| 課題感 | システム全体の設計力、DB/データ設計 |
| Sentry経験 | 基本導入のみ |
| 学習の深さ | 仕組みまで理解したい |

## 全体フェーズ

```
Phase 0: 全体設計                    ← 今ここ
  ↓
Phase 1: アプリ基盤構築
  ↓
Phase 2: Sentry統合                  ← メイン学習1
  ↓
Phase 3: AWSデプロイ                  ← メイン学習2
  ↓
Phase 4: Sentry + AWS連携
```

## 各Phase詳細

### Phase 0: 全体設計 - システムアーキテクチャ設計（Sentry + AWS特化）

| 項目 | 内容 |
|---|---|
| 目的 | システム全体の設計を行い、各コンポーネントの役割を理解する |
| 成果物 | 構成図（.drawio）、API設計書、DB設計書 |
| 学習ポイント | AWS各サービスの役割、REST API設計、設計してからコードを書く習慣 |
| 詳細 | [Phase 0 詳細](./phase-0/README.md) |

**ステップ:**
- [x] システムアーキテクチャ設計（構成図作成）
- [ ] API設計（エンドポイント定義）
- [ ] データベース設計（テーブル定義・最小限）

---

### Phase 1: アプリ基盤構築（Go API + Next.js + PostgreSQL）

| 項目 | 内容 |
|---|---|
| 目的 | ローカル開発環境を構築し、基本的なCRUD APIを実装する |
| 成果物 | 動くGo API + Next.jsフロントエンド + Docker Compose環境 |
| 学習ポイント | レイヤー分離、責務分割、Docker Composeによる開発環境構築 |
| 詳細 | [Phase 1 詳細](./phase-1/README.md) |

**ステップ（予定）:**
- [ ] Go APIサーバーのセットアップ
- [ ] Next.jsフロントエンドのセットアップ
- [ ] Docker Compose（Go + Next.js + PostgreSQL）
- [ ] 基本CRUD API実装
- [ ] ヘルスチェックWorker実装

---

### Phase 2: Sentry統合 - エラー追跡 & パフォーマンス監視

| 項目 | 内容 |
|---|---|
| 目的 | Sentryを深く理解し、バックエンド・フロントエンド両方に統合する |
| 成果物 | Sentry統合済みアプリケーション |
| 学習ポイント | エラーキャプチャ設計、パフォーマンスモニタリング、アラート設定、可観測性 |
| 詳細 | [Phase 2 詳細](./phase-2/README.md) |

**ステップ（予定）:**
- [ ] Sentry SDKの導入（Go / Next.js）
- [ ] エラーキャプチャの設計
- [ ] パフォーマンスモニタリング設定
- [ ] アラート設定
- [ ] Sentryダッシュボードの活用

---

### Phase 3: AWSデプロイ - インフラ構築 & デプロイ

| 項目 | 内容 |
|---|---|
| 目的 | AWSの基本的なインフラ構築を学び、アプリをデプロイする |
| 成果物 | AWS上で動くアプリケーション |
| 学習ポイント | VPC/Subnet/Security Group、EC2セットアップ、SSH接続、Docker on EC2 |
| 詳細 | [Phase 3 詳細](./phase-3/README.md) |

**ステップ（予定）:**
- [ ] VPC / Subnet / Security Group の構築
- [ ] EC2 (t3.micro) のセットアップ
- [ ] Docker / Docker Compose のインストール
- [ ] アプリケーションのデプロイ
- [ ] 動作確認

---

### Phase 4: Sentry + AWS連携 - 本番環境での可観測性

| 項目 | 内容 |
|---|---|
| 目的 | 本番環境でのSentryの動きを理解し、実際の障害対応を体験する |
| 成果物 | 本番環境でSentry統合された監視サービス |
| 学習ポイント | Source Maps、リリース管理、障害シナリオのシミュレーション |
| 詳細 | [Phase 4 詳細](./phase-4/README.md) |

**ステップ（予定）:**
- [ ] 本番環境でのSentry設定
- [ ] Source Maps の設定
- [ ] リリース管理
- [ ] 障害シナリオのシミュレーション
- [ ] 振り返り・まとめ

---

## アーキテクチャ概要

```
User (Browser)
  │ HTTP/HTTPS
  ▼
┌─ AWS Cloud ──────────────────────────────┐
│  ┌─ VPC ──────────────────────────────┐  │
│  │  ┌─ Public Subnet ─────────────┐   │  │
│  │  │  ┌─ EC2 (t3.micro) ──────┐  │   │  │
│  │  │  │  Go API + Worker      │  │   │  │
│  │  │  │  Next.js Frontend     │  │   │  │
│  │  │  │  PostgreSQL (Docker)  │  │   │  │
│  │  │  └───────────────────────┘  │   │  │
│  │  │  Security Group             │   │  │
│  │  └─────────────────────────────┘   │  │
│  └────────────────────────────────────┘  │
│  CloudWatch (Logs) ← オプション           │
└──────────────────────────────────────────┘
    ↕ Sentry SDK
Sentry (SaaS)

構成図: docs/architecture.drawio
```

## コスト見積もり

| サービス | 月額目安 |
|---|---|
| EC2 (t3.micro) | 無料枠（12ヶ月間 750h/月） |
| Elastic IP | EC2紐付け時は無料 |
| CloudWatch | 無料枠内（学習レベル） |
| Sentry | 無料プラン（Developer） |
| **合計** | **約 $0/月**（無料枠内） |
