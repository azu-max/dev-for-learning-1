# Phase 1 振り返りクイズ（全15問）

Phase 1（アプリ基盤構築）で学んだ内容の確認クイズ。

---

## Docker 基礎

### Q1: Docker Compose の主な役割はどれ？

- A) イメージビルド高速化
- B) 複数コンテナの一括管理
- C) 本番オートスケール
- D) CI/CDパイプライン

### Q3: Dockerfile で COPY go.mod go.sum を先に行い、RUN go mod download してからソースをCOPYする理由は？

- A) セキュリティ対策
- B) Dockerキャッシュ活用
- C) ファイルサイズ削減
- D) 並列ビルド

### Q4: Docker Compose の volumes でホストのソースコードをコンテナにマウントする主な目的は？

- A) データ永続化
- B) ホットリロード
- C) バックアップ
- D) セキュリティ

---

## Docker Compose

### Q2: docker-compose.yml で depends_on を指定すると何が保証される？

- A) 起動完了を待つ
- B) コンテナ開始順のみ（中身の準備完了は保証しない）
- C) ネットワーク分離
- D) 連鎖停止

### Q14: docker compose up でコンテナ名で通信できる（例: db:5432）のはなぜ？

- A) /etc/hostsに自動登録
- B) Docker内部DNSが名前解決
- C) ポートマッピング
- D) 環境変数でIP指定

---

## レイヤー分離

### Q5: Go API の3層構造で、Handler 層の責務はどれ？

- A) HTTPリクエストの受付・レスポンス返却
- B) ビジネスロジックの実行
- C) SQLクエリの組み立て
- D) Workerのスケジュール管理

### Q6: Service 層が直接 SQL を書かず、Repository 層を経由するメリットは？

- A) パフォーマンス向上
- B) DB変更の影響を限定
- C) セキュリティ対策
- D) Goの言語仕様

### Q8: Repository 層を interface で定義する主なメリットは？

- A) コード量削減
- B) テスト時にモックに差し替え可能
- C) 実行速度向上
- D) Goの必須構文

### Q15: 「Handler / Service / Repository」の3層分離で得られる最大のメリットは？

- A) 実行速度の向上
- B) 変更の影響範囲を限定
- C) コード行数の削減
- D) デプロイが簡単になる

---

## REST API

### Q7: RESTful API で「リソースの新規作成」に使う HTTP メソッドは？

- A) GET
- B) POST
- C) PUT
- D) DELETE

---

## Worker・並行処理

### Q9: ヘルスチェック Worker はレイヤー構造のどこに位置する？

- A) Handler層
- B) Service層
- C) Repository層
- D) Handlerと同じ「エントリーポイント」層

### Q10: Go の goroutine の説明として正しいのは？

- A) OSスレッドと同じ
- B) 軽量な並行処理単位
- C) 別プロセス
- D) async/awaitと同じ

### Q11: Worker で time.NewTicker を使う目的は？

- A) 一定間隔で繰り返し実行
- B) タイムアウト設定
- C) CPU使用率制限
- D) ログタイムスタンプ

### Q12: Worker が単一 goroutine で動く場合の弱点は？

- A) メモリ使用量が大きい
- B) 1つのチェックが遅いと全体が詰まる
- C) DB接続が切れる
- D) ログが出せない

---

## Graceful Shutdown

### Q13: Graceful Shutdown で signal.Notify(quit, SIGINT, SIGTERM) を使う目的は？

- A) アイドル時の自動終了
- B) 終了シグナルを受け取り安全に停止
- C) メモリ監視
- D) DBタイムアウト設定
