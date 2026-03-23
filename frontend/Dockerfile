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

EXPOSE 3000

CMD ["npm", "run", "start"]
