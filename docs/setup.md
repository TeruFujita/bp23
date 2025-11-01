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

## 起動が必要なサービスと組み合わせ

### LINE Botを動作させる場合
**同時に起動が必要：**
1. LINE Botサーバー（`localhost:3001`）
2. ngrok（Botサーバーを公開）

**起動手順：**
```powershell
# 1. Botサーバー起動（PowerShell タブ1）
cd line-bot
go run src/main.go

# 2. ngrok起動（PowerShell タブ2）
cd ..
npx ngrok http 3001
```

**動作確認：**
- LINEでメッセージ送信 → Flexメッセージが返ってくればOK

### Web管理UIを動作させる場合
**同時に起動が必要：**
1. APIサーバー（`localhost:8080`）
2. Webアプリ（`localhost:3000`）

**起動手順：**
```powershell
# 1. APIサーバー起動（PowerShell タブ1）
cd api-server
go run cmd/server/main.go

# 2. Webアプリ起動（PowerShell タブ2）
cd ../web-app
npm run dev
```

**動作確認：**
- `http://localhost:3000` にアクセス → ページ表示
- `http://localhost:8080/healthz` → `{"status":"ok"}` が返る

### すべてを同時に動作させる場合
**3つのターミナルが必要：**
1. LINE Botサーバー + ngrok（LINE Bot用）
2. APIサーバー（Web管理UI用）
3. Webアプリ（管理UI用）

**最小構成（LINE Botだけ使う場合）：**
- Botサーバー + ngrok のみで動作可（API/Web不要）

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
