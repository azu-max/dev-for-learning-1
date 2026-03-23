# ============================================
# Stage 1: Builder（ビルド用）
# ============================================
FROM node:20-alpine AS builder

WORKDIR /app

# 依存関係ファイルを先にコピー（キャッシュ活用）
COPY package.json package-lock.json ./
RUN npm ci

# ソースコードをコピーしてビルド
COPY . .
RUN mkdir -p public
# Docker Compose内ではAPIコンテナのサービス名で通信する
ENV API_URL=http://api:8080
RUN npm run build

# ============================================
# Stage 2: Runner（実行用）
# ============================================
FROM node:20-alpine

WORKDIR /app

# 本番用の依存関係だけインストール
COPY package.json package-lock.json ./
RUN npm ci --omit=dev

# ビルド成果物をコピー
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/next.config.ts ./next.config.ts

EXPOSE 3000

CMD ["npm", "run", "start"]
