# ============================================
# Stage 1: Builder（ビルド用の環境）
# ============================================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 依存関係ファイルを先にコピー（キャッシュ活用）
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# ビルド（CGO無効 → 静的バイナリになる）
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# ============================================
# Stage 2: Runner（実行用の最小環境）
# ============================================
FROM alpine:3.21

# CA証明書を追加（HTTPS通信に必要）
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Builder から実行ファイルだけコピー
COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
