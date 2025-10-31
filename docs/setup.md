# 環境構築・起動手順

## 必要ソフト
- Go 1.22+（または 1.23）
- Node.js 18+
- PostgreSQL 14+

## 環境変数
- ルートに `.env`（例は `.env.example`）
- 主要キー: `DATABASE_URL`, `LINE_CHANNEL_SECRET`, `LINE_CHANNEL_TOKEN`, `LINE_LOGIN_*`, `NEXT_PUBLIC_*`, `JWT_SECRET`, `ALLOWED_ORIGIN`

## DB 準備
- DB 作成: `detour_bot_dev`
- 初期マイグレーション適用（psql）

## 起動
- API: `cd api-server && go run cmd/server/main.go`
- Web: `cd web-app && npm run dev`
- Bot: `cd line-bot && go run src/main.go`

## LINE Bot 開発フロー（最小）
1. LINEチャネル作成 → `LINE_CHANNEL_SECRET` と `LINE_CHANNEL_TOKEN` を `.env` に設定
2. `/webhook` を実装（署名検証 + エコー返信から開始）
3. ngrok で公開: `ngrok http <bot-port>`
4. 発行URLを Webhook に設定・有効化
5. 実機のLINEで送受信を確認

## マイグレーション適用例
```
psql "postgres://postgres:YOUR_PASSWORD@localhost:5432/detour_bot_dev?sslmode=disable" -f api-server/migrations/0001_init.sql
```
