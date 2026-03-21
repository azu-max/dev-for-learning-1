FROM golang:1.24-alpine

WORKDIR /app

# 依存関係のキャッシュ（go.mod/go.sumが変わらなければ再ダウンロード不要）
COPY go.mod go.sum* ./
RUN go mod download

# ソースコードをコピー
COPY . .

# ビルド
RUN go build -o server .

# 起動
CMD ["./server"]
