# 環境構築・起動手順

## 必要ソフト
- Go 1.22+（または 1.23）
- Node.js 18+
- PostgreSQL 14+

## 環境変数
- ルートに `.env`（例は `.env.example`）
- 主要キー: `DATABASE_URL`, `LINE_CHANNEL_SECRET`, `LINE_CHANNEL_TOKEN`, `LINE_LOGIN_*`, `NEXT_PUBLIC_*`, `JWT_SECRET`, `ALLOWED_ORIGIN`, `PORT`

### 環境変数の設定方法（PowerShell）
```powershell
# .env から読み込む（手動で設定する場合）
$env:LINE_CHANNEL_SECRET="YOUR_SECRET"
$env:LINE_CHANNEL_TOKEN="YOUR_TOKEN"
$env:PORT="3001"
$env:ALLOWED_ORIGIN="http://localhost:3000"
```

## DB 準備
- DB 作成: `detour_bot_dev`
- 初期マイグレーション適用（psql）

## 起動
- API: `cd api-server && go run cmd/server/main.go`（ポート8080）
- Web: `cd web-app && npm run dev`（ポート3000）
- Bot: `cd line-bot && go run src/main.go`（ポート3001）

### 起動確認
- API: `http://localhost:8080/healthz` → `{"status":"ok"}` が返ればOK
- Web: `http://localhost:3000` にアクセスしてページが表示されればOK
- Bot: `http://localhost:3001/webhook` はエンドポイント（LINEから呼ばれる）

## LINE Bot 開発フロー（ngrok設定含む）
1. LINE Developers でチャネル作成 → `LINE_CHANNEL_SECRET` と `LINE_CHANNEL_TOKEN` を取得
2. `.env` に設定（環境変数として読み込む）
3. ローカルで起動: `cd line-bot && go run src/main.go`（ポート3001）
4. ngrok で公開:
   ```powershell
   ngrok http 3001
   ```
   → 表示される `Forwarding: https://xxxxx.ngrok.io -> http://localhost:3001` のURLをメモ
5. LINE Developers コンソールで:
   - Webhook URL に `https://xxxxx.ngrok.io/webhook` を設定
   - 「Webhookの利用」を有効化
   - 「応答メッセージ」はOFF推奨（Botで制御するため）
6. 実機のLINEでメッセージ送信 → Flexメッセージが返ってくればOK

## マイグレーション適用例
```
psql "postgres://postgres:YOUR_PASSWORD@localhost:5432/detour_bot_dev?sslmode=disable" -f api-server/migrations/0001_init.sql
```
