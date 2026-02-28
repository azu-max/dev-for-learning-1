# Phase 0: 全体設計 - システムアーキテクチャ設計

## 概要

| 項目 | 内容 |
|---|---|
| 目的 | システム全体の設計を行い、各コンポーネントの役割を理解する |
| 状態 | 🔄 進行中 |
| 成果物 | 構成図（.drawio）、API設計書、DB設計書 |

## ステップ

- [x] Step 1: システムアーキテクチャ設計（構成図作成）
- [x] Step 2: API設計（エンドポイント定義）
- [ ] Step 3: データベース設計（テーブル定義・最小限）

---

## Step 1: システムアーキテクチャ設計

### 構成図

**ファイル:** [../architecture.drawio](../architecture.drawio)

```
User (Browser)
  │ HTTP/HTTPS
  ▼
┌─ AWS Cloud ──────────────────────────────────────┐
│  ┌─ VPC (10.0.0.0/16) ───────────────────────┐   │
│  │  Internet Gateway                          │   │
│  │  ┌─ Public Subnet (10.0.1.0/24) ───────┐  │   │
│  │  │  Security Group                       │  │   │
│  │  │  ┌─ EC2 (t3.micro) ──────────────┐  │  │   │
│  │  │  │  Docker Compose                │  │  │   │
│  │  │  │  ├── Go API Server (:8080)     │  │  │   │
│  │  │  │  │   + Health Check Worker     │  │  │   │
│  │  │  │  │   + Sentry SDK              │  │  │   │
│  │  │  │  ├── Next.js Frontend (:3000)  │  │  │   │
│  │  │  │  │   + Sentry SDK              │  │  │   │
│  │  │  │  └── PostgreSQL (:5432)        │  │  │   │
│  │  │  └────────────────────────────────┘  │  │   │
│  │  └──────────────────────────────────────┘  │   │
│  └────────────────────────────────────────────┘   │
│  CloudWatch (Logs) ← オプション                     │
│  Elastic IP ← オプション                            │
└───────────────────────────────────────────────────┘
     ↕ エラー・パフォーマンスデータ        ↔ Health Check (HTTP)
  Sentry (SaaS)                     Monitored External Services
```

### 設計判断の記録

| 判断 | 選択 | 理由 |
|---|---|---|
| コンピュート | EC2 (t3.micro) | 無料枠あり。中で何が起きているか直接触って学べる |
| DB | PostgreSQL (Docker on EC2) | ローカルと本番で同じDBエンジンに統一 |
| コンテナ管理 | Docker Compose | ローカルと同じ環境をAWSでも再現できる |
| ネットワーク | Public Subnet 1つ | 学習用の最小構成 |
| CloudWatch | オプション（後から追加可能） | 学習のメインはSentry。CloudWatchは優先度低 |
| Elastic IP | オプション（あると便利） | EC2再起動時にIPが変わらない。EC2紐付け時は無料 |
| ALB | 不使用 | コスト削減（月$16節約）。学習に不要 |
| ECS Fargate | 不使用 | 無料枠なし。EC2の方が基礎を学べる |
| RDS | 不使用 | EC2内のDockerでPostgreSQLを動かす（無料） |

---

## 座学ノート

### 1. AWSの全体像（たとえ話）

```
AWS Cloud      = ビル全体（Amazonが管理するデータセンター）
VPC            = あなたが借りたフロア（専用ネットワーク空間）
Public Subnet  = そのフロアの中の「外から入れる部屋」
Security Group = 部屋のドアの鍵（誰を入れるかのルール）
EC2            = 部屋の中に置いた「パソコン1台」
Internet GW    = ビルの正面玄関
Elastic IP     = そのビルの「固定の住所」
CloudWatch     = ビルの監視カメラ
```

### 2. VPC と Subnet

- **VPC (Virtual Private Cloud):** AWS上にあなた専用の仮想ネットワークを作ること
- AWSアカウント作成時にデフォルトVPCが自動作成される
- **Public Subnet:** インターネットからアクセスできる区画
- **Private Subnet:** 外からアクセスできない区画（DBを守りたい時などに使用）
- VPCの中にSubnetは**複数作れる**（可用性・セキュリティ・役割分離のため）
- 今回は学習用なのでPublic Subnet 1つで十分

### 3. Internet Gateway

- VPCはデフォルトでインターネットに繋がっていない
- Internet Gatewayを付けることで外の世界と通信可能になる

### 4. Elastic IP

- EC2にはデフォルトでIPが付くが、**再起動するたびにIPが変わる**
- Elastic IP = 固定IPアドレス。再起動してもIPが変わらない
- 料金: EC2に紐付けて使用中は**無料**。紐付けずに放置すると**約$0.005/時間**

### 5. CloudWatch vs Sentry

| 観点 | CloudWatch | Sentry |
|---|---|---|
| 例えるなら | 防犯カメラの生映像 | 警備レポート |
| 見るもの | サーバーの健康状態（CPU, メモリ, ログ） | アプリのバグ・エラーの詳細 |
| ログの見方 | テキストをそのまま表示 | エラーを自動分類・集計・通知 |
| バックエンド | ✅ ログが見れる | ✅ エラーの文脈まで見れる |
| フロントエンド | ❌ 見れない | ✅ 見れる |

**フロントエンドのエラーがCloudWatchで取れない理由:**
フロントエンドのエラーはユーザーのブラウザ内で発生する。サーバー（EC2）はブラウザの中を覗けないため、CloudWatchでは取得できない。SentryのSDKはブラウザ内で動作するためキャッチ可能。

**CloudWatchがあればSentryは不要か？ → No。役割が違う:**
- CloudWatch = サーバーの健康管理（インフラ層）
- Sentry = アプリのバグ管理（アプリケーション層）

### 6. EC2 vs Lambda vs Fargate

| | EC2 | Lambda | Fargate |
|---|---|---|---|
| 例えるなら | 部屋を借りて自分で管理 | 電話1本で用事を頼む | 部屋付きの執事を雇う |
| 管理するもの | OS, ミドルウェア, アプリ全部 | コードだけ | コンテナだけ |
| 課金 | 起動時間で課金 | 実行回数+時間で課金 | CPU+メモリで課金 |
| 無料枠 | ✅ 750h/月 | ✅ 100万リクエスト/月 | ❌ なし |
| OS管理 | 必要 | 不要 | 不要 |

**今回EC2を選んだ理由:** 無料枠があり、中で何が起きているかを直接触って学べるから。

### 7. EC2 1台 vs 複数台

| | 1台構成 | 複数台構成 |
|---|---|---|
| メリット | シンプル、安い、管理が楽、無料枠で収まる | 1台落ちても他が動く（高可用性）、負荷分散 |
| デメリット | その1台が落ちたら全部止まる（SPOF） | 料金が台数分かかる、構成が複雑（ALB必要） |

**SPOF = Single Point of Failure（単一障害点）**

### 8. Docker Compose on AWS

- Docker ComposeはローカルでもAWSでも同じように使える
- メリット: ローカルで動いたものがそのままAWSでも動く
- EC2内でDocker Composeを使い、Go API / Next.js / PostgreSQLの3コンテナを管理

### 9. Monitored External Services

- AWSのサービスではない
- 今回作るヘルスチェックサービスが**監視する対象**の外部サイトやAPI
- 例: example.com, api.xxx.com など

---

## 次のステップ

→ [Step 2: API設計](./api-design.md)（未作成）
→ [Step 3: データベース設計](./db-design.md)（未作成）
