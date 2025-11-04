# 環境構築・起動手順

## 📋 目次

- [必要ソフト](#必要ソフト)
- [外部サービスアカウント準備](#外部サービスアカウント準備)
- [環境変数設定](#環境変数設定)
- [データベース準備](#データベース準備)
- [サーバー起動](#サーバー起動)
- [チーム開発](#チーム開発)

---

## 必要ソフト

以下をインストールしてください：

- **Go** 1.22以上（または 1.23）
- **Node.js** 18以上
- **PostgreSQL** 14以上

---

## 外部サービスアカウント準備

### LINE Developers（必須）

LINE Botを動作させるために必要です。

1. **アカウント登録**
   - https://developers.line.biz/ にアクセス
   - LINEアカウントでログイン

2. **チャネル作成**
   - 「新規プロバイダー作成」をクリック
   - 「Messaging APIチャネル」を作成

3. **必要情報の取得**
   - **チャネルシークレット**: チャネル設定 → チャネルシークレット
   - **チャネルアクセストークン**: Messaging API設定 → チャネルアクセストークン
   - これらを `.env.local` に設定します（後述）

### ngrok（LINE Bot開発時のみ必要）

ローカル開発でLINE Webhookを受信するために必要です。

1. **アカウント登録**
   - https://dashboard.ngrok.com/signup （無料）
   
2. **authtoken取得**
   - https://dashboard.ngrok.com/get-started/your-authtoken

3. **authtoken設定**
   ```powershell
   npx ngrok config add-authtoken YOUR_AUTHTOKEN
   ```

4. **インストール**（未導入の場合）
   ```powershell
   npm install -g ngrok
   ```

詳細は [LINE Botの仕組み](line-bot-explanation.md) を参照

### その他の外部API（オプション）

必要に応じて後で設定できます：

- **OpenRouteService API**: 経路探索（無料 2,500リクエスト/日）
- **ぐるなびAPI / HotPepper API**: 店舗検索

詳細は [外部API連携](external-apis.md) を参照

---

## 環境変数設定

### ファイル作成

プロジェクトルートで `.env.example` を `.env.local` にコピー：

```powershell
# プロジェクトルートから
copy .env.example .env.local
```

### 必須項目（LINE Bot動作に必要）

```env
LINE_CHANNEL_SECRET=取得したチャネルシークレット
LINE_CHANNEL_TOKEN=取得したチャネルアクセストークン
```

### 主要な環境変数

| 変数名 | 説明 | デフォルト値 |
|---|---|---|
| `LINE_CHANNEL_SECRET` | LINEチャネルシークレット（必須） | - |
| `LINE_CHANNEL_TOKEN` | LINEチャネルアクセストークン（必須） | - |
| `DATABASE_URL` | PostgreSQL接続文字列 | - |
| `PORT` | Botサーバーのポート | `3001` |
| `ALLOWED_ORIGIN` | CORS許可ドメイン | `http://localhost:3000` |
| `JWT_SECRET` | JWT署名用（将来実装） | - |

**完全な一覧**: [環境変数詳細](#環境変数詳細)

### 環境変数の読み込み

- **LINE Bot**: `.env.local` を自動読み込み（`godotenv`使用）
- **API Server**: 手動で環境変数を設定するか、`.env.local`読み込みを追加が必要

**注意**: `.env.local` は `.gitignore` に含まれているため、Gitにはコミットされません

---

## データベース準備

### 1. データベース作成

```powershell
psql -U postgres -c "CREATE DATABASE detour_bot_dev;"
```

### 2. マイグレーション適用

テーブルを作成するため、マイグレーションファイルを実行します。

**方法1: psqlコマンドで直接適用（推奨）**
```powershell
# プロジェクトルートから実行
psql "postgres://postgres:YOUR_PASSWORD@localhost:5432/detour_bot_dev?sslmode=disable" -f api-server/migrations/0001_init.sql
```

**方法2: psql対話モードで適用**
```powershell
# データベースに接続
psql "postgres://postgres:YOUR_PASSWORD@localhost:5432/detour_bot_dev?sslmode=disable"

# psql内で実行
\i api-server/migrations/0001_init.sql
```

### 3. 確認

```sql
-- テーブル一覧を確認
\dt

-- 各テーブルの構造を確認
\d users
\d routes
\d spots
```

**作成されるテーブル:**

- `users`: LINEユーザー情報
- `routes`: 経路情報
- `spots`: 寄り道スポット情報

詳細は [データモデル](data-model.md) を参照

---

## サーバー起動

### LINE Botを動作させる場合

**同時に起動が必要:**
1. LINE Botサーバー（`localhost:3001`）
2. ngrok（Botサーバーを公開）

**起動手順:**

**ターミナル1: Botサーバー起動**
```powershell
cd line-bot
go run src/main.go
```

**ターミナル2: ngrok起動**
```powershell
# プロジェクトルートから実行
npx ngrok http 3001
```

**LINE Developersでの設定:**
1. ngrokで表示されたURL（例: `https://xxxxx.ngrok.io`）をコピー
2. LINE Developers → Messaging API設定 → Webhook URL に設定
   - 設定値: `https://xxxxx.ngrok.io/webhook`
3. 「Webhookの利用」を有効化
4. 「応答メッセージ」はOFF推奨（Botで制御するため）

**動作確認:**
- LINEでメッセージ送信 → Flexメッセージが返ってくればOK

### Web管理UIを動作させる場合

**同時に起動が必要:**
1. APIサーバー（`localhost:8080`）
2. Webアプリ（`localhost:3000`）

**起動手順:**

**ターミナル1: APIサーバー起動**
```powershell
cd api-server
go run cmd/server/main.go
```

**ターミナル2: Webアプリ起動**
```powershell
cd web-app
npm install  # 初回のみ
npm run dev
```

**動作確認:**
- `http://localhost:3000` にアクセス → ページが表示される
- `http://localhost:8080/healthz` → `{"status":"ok"}` が返る

### すべてを同時に動作させる場合

**3つのターミナルが必要:**
1. LINE Botサーバー + ngrok
2. APIサーバー
3. Webアプリ

**最小構成:**
- LINE Botだけ使う場合: Botサーバー + ngrok のみで動作可

---

## チーム開発

### Webhook URLの設定方法

**パターンA: ngrokを個別に起動（推奨）**

各開発者が各自のngrokを起動し、使う前にURLを共有します。

1. 各自で `npx ngrok http 3001` を起動
2. 表示されたURL（例: `https://xxxxx.ngrok.io`）をチームで共有
3. 使う前に「開発中です。Webhook URL: https://xxxxx.ngrok.io/webhook」と通知
4. LINE DevelopersでWebhook URLを設定・有効化
5. 終了時に「終了しました」と通知
6. 次の開発者が同じ手順で自分のURLを設定

**パターンB: ngrokを共有で使う**

1人がngrokを起動し、そのURLを全員で共有します。

⚠️ **注意**: ngrokは起動した人のローカルサーバー（localhost:3001）に転送するため、**起動した人のBotサーバーしか動作しません**

### 共通設定

- `.env.local` の `LINE_CHANNEL_SECRET` / `LINE_CHANNEL_TOKEN` は同じチャンネルの値を全員で共有
- 「応答メッセージ」はOFF推奨

---

## 環境変数詳細

### Database
```env
DATABASE_URL=postgres://postgres:YOUR_PASSWORD@localhost:5432/detour_bot_dev?sslmode=disable
```

### LINE Messaging API（必須）
```env
LINE_CHANNEL_SECRET=取得したチャネルシークレット
LINE_CHANNEL_TOKEN=取得したチャネルアクセストークン
```

### LINE Login（将来実装）
```env
LINE_LOGIN_CLIENT_ID=LINE Loginチャネルで取得
LINE_LOGIN_CLIENT_SECRET=LINE Loginチャネルで取得
LINE_LOGIN_CALLBACK_URL=http://localhost:3000/(auth)/callback
```

### Frontend（公開変数）
```env
NEXT_PUBLIC_MAPS_API_KEY=Google Maps等のAPIキー
```

### API Server
```env
JWT_SECRET=JWT署名用（変更推奨）
ALLOWED_ORIGIN=http://localhost:3000
```

### LINE Bot Server
```env
PORT=3001
```

### External APIs（オプション）
```env
ORS_API_KEY=OpenRouteService APIキー
GURUNAVI_API_KEY=ぐるなびAPIキー
HOTPEPPER_API_KEY=HotPepper APIキー
NEXT_PUBLIC_MAPBOX_TOKEN=Mapboxトークン
```

---

## トラブルシューティング

### PostgreSQLサービスが起動しない

PowerShellを**管理者として実行**してから以下を実行してください：

```powershell
Start-Service "postgresql-x64-18"
```

### データベース接続エラー

- PostgreSQLサーバーが起動しているか確認
- `DATABASE_URL` の接続文字列が正しいか確認
- データベース `detour_bot_dev` が作成されているか確認

### LINE Botが応答しない

- ngrokが起動しているか確認
- LINE DevelopersでWebhook URLが正しく設定されているか確認
- Botサーバーが起動しているか確認（`localhost:3001`）

### 環境変数が読み込まれない

- `.env.local` がプロジェクトルートにあるか確認
- LINE Bot: 自動読み込み（`godotenv`）
- API Server: 手動で環境変数を設定する必要がある場合あり
