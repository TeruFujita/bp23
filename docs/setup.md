# 環境構築・起動手順

## 必要ソフト
- Go 1.22+（または 1.23）
- Node.js 18+
- PostgreSQL 14+

## 環境変数
- ルートに `.env.local` を作成（`.env.example` をコピーして実際の値を設定）
- **LINE Botは自動で `.env.local` を読み込む**（godotenv使用）

### 必要な環境変数一覧

#### Database
- `DATABASE_URL`: PostgreSQL接続文字列（例: `postgres://postgres:YOUR_PASSWORD@localhost:5432/detour_bot_dev?sslmode=disable`）

#### LINE Messaging API（必須）
- `LINE_CHANNEL_SECRET`: LINE Developersで取得
- `LINE_CHANNEL_TOKEN`: LINE Developersで取得

#### LINE Login（Web管理UI用、将来実装）
- `LINE_LOGIN_CLIENT_ID`: LINE Loginチャネルで取得
- `LINE_LOGIN_CLIENT_SECRET`: LINE Loginチャネルで取得
- `LINE_LOGIN_CALLBACK_URL`: `http://localhost:3000/(auth)/callback`

#### Frontend（公開変数）
- `NEXT_PUBLIC_MAPS_API_KEY`: 地図APIキー（Google Maps等）

#### API Server
- `JWT_SECRET`: JWT署名用（変更推奨）
- `ALLOWED_ORIGIN`: CORS許可ドメイン（デフォルト: `http://localhost:3000`）

#### LINE Bot Server
- `PORT`: Botサーバーのポート（デフォルト: `3001`）

#### External APIs（オプション）
- `ORS_API_KEY`: OpenRouteService APIキー
- `GURUNAVI_API_KEY`: ぐるなびAPIキー
- `HOTPEPPER_API_KEY`: HotPepper APIキー
- `NEXT_PUBLIC_MAPBOX_TOKEN`: Mapboxトークン（地図の代替）

### 設定方法
1. `.env.example` を `.env.local` にコピー
2. 実際の値を設定（LINE Developersで取得した値など）
3. 必須項目: `LINE_CHANNEL_SECRET`, `LINE_CHANNEL_TOKEN`（LINE Bot動作に必要）

**注意**: `.env.local` は `.gitignore` に含まれているため、Gitにはコミットされません

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
2. `.env.local` に設定（`line-bot/src/main.go` が自動で読み込む）
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
