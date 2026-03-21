FROM node:20-alpine

WORKDIR /app

# 依存関係のキャッシュ
COPY package.json package-lock.json* ./
RUN npm install

# ソースコードをコピー
COPY . .

# 開発サーバー起動
CMD ["npm", "run", "dev"]
