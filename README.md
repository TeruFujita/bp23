# Detour Bot - 寄り道検索システム

LINE Bot + Web管理UI + Go APIサーバーによる、経路検索と寄り道スポット提案システム

## 📚 目次

- [🚀 クイックスタート](#-クイックスタート)
- [📖 ドキュメント](#-ドキュメント)
- [🛠️ 必要なもの](#️-必要なもの)
- [🏗️ プロジェクト構成](#️-プロジェクト構成)

---

## 🚀 クイックスタート

### 1. 必要なソフトをインストール

- **Go** 1.22以上
- **Node.js** 18以上  
- **PostgreSQL** 14以上

### 2. 外部サービスを準備

#### LINE Developers（必須）
1. https://developers.line.biz/ に登録
2. Messaging APIチャネルを作成
3. `LINE_CHANNEL_SECRET` と `LINE_CHANNEL_TOKEN` を取得

#### ngrok（LINE Bot開発時のみ必要）
1. https://dashboard.ngrok.com/signup でアカウント作成
2. authtokenを取得して設定:
   ```powershell
   npx ngrok config add-authtoken YOUR_AUTHTOKEN
   ```

### 3. 環境変数を設定

プロジェクトルートで `.env.example` を `.env.local` にコピーし、以下の値を設定:

**必須項目:**
- `LINE_CHANNEL_SECRET`: LINE Developersで取得
- `LINE_CHANNEL_TOKEN`: LINE Developersで取得

詳細は [環境構築ガイド](docs/setup.md) を参照

### 4. データベースを準備

PostgreSQL サービスが起動していることを確認してください（Windows: サービス / macOS: brew services / Linux: systemctl）。

```powershell
# データベース作成
psql -U postgres -c "CREATE DATABASE detour_bot_dev;"

# マイグレーション実行（プロジェクトルートから）
psql "postgres://postgres:YOUR_PASSWORD@localhost:5432/detour_bot_dev?sslmode=disable" -f api-server/migrations/0001_init.sql
```

#### 4-1. PostgreSQL が起動していない場合（手動起動）

- Windows（PowerShell 管理者）
  ```powershell
  Start-Service "postgresql-x64-18"
  ```
  サービスが無い/使えない場合は、`pg_ctl` でも起動できます（パスは環境に合わせて変更）。
  ```powershell
  pg_ctl -D "D:\develop\Installers\PostgreSQL\data" start
  ```

- macOS（Homebrew）
  ```bash
  brew services start postgresql@14
  brew services list
  ```

- Linux（systemd）
  ```bash
  sudo systemctl start postgresql
  sudo systemctl status postgresql
  ```

### 5. サーバーを起動

**LINE Botを使う場合:**
```powershell
# ターミナル1: Botサーバー
cd line-bot
go run src/main.go

# ターミナル2: ngrok（別ウィンドウ）
npx ngrok http 3001
# → 表示されたURLをLINE DevelopersのWebhook URLに設定
```

**Web管理UIを使う場合:**
```powershell
# ターミナル1: APIサーバー
cd api-server
go run cmd/server/main.go

# ターミナル2: Webアプリ
cd web-app
npm install  # 初回のみ
npm run dev
```

詳細な手順は [環境構築ガイド](docs/setup.md) を参照してください。

---

## 📖 ドキュメント

### 📘 開発者向け

| ドキュメント | 説明 |
|---|---|
| [環境構築ガイド](docs/setup.md) | セットアップ手順、環境変数、起動方法 |
| [アーキテクチャ概要](docs/architecture.md) | システム構成、技術スタック、設計方針 |
| [データモデル](docs/data-model.md) | データベーステーブル設計 |


### 🔧 技術詳細

| ドキュメント | 説明 |
|---|---|
| [API仕様](docs/api.md) | REST APIエンドポイント一覧 |
| [外部API連携](docs/external-apis.md) | 使用する外部API（経路探索、店舗検索など） |
| [LINE Botの仕組み](docs/line-bot-explanation.md) | Webhook、ngrok、署名検証の解説 |

### 🚀 運用・本番

| ドキュメント | 説明 |
|---|---|
| [運用・本番移行](docs/operations.md) | 認証設定、セキュリティ、デプロイ方針 |

---

## 🛠️ 必要なもの

### 開発環境

- **Go** 1.22以上（または 1.23）
- **Node.js** 18以上
- **PostgreSQL** 14以上

### 外部サービスアカウント

- **LINE Developers**（必須）
  - Messaging APIチャネル作成
  - チャネルシークレットとアクセストークンを取得

- **ngrok**（LINE Bot開発時のみ）
  - 無料アカウントでOK
  - ローカルサーバーを公開するため

- **その他**（オプション、後で設定可能）
  - OpenRouteService API（経路探索）
  - ぐるなびAPI / HotPepper API（店舗検索）

---

## 🏗️ プロジェクト構成

```
line_bot/
├── line-bot/                     # LINE Botサーバー（Go）
│   ├── src/
│   │   ├── main.go               # Webhook実装済み（署名検証/Flex返信）
│   │   ├── handlers/             # メッセージハンドラー（将来実装）
│   │   ├── services/             # ビジネスロジック（将来実装）
│   │   ├── models/               # データモデル（将来実装）
│   │   └── utils/                 # ユーティリティ（将来実装）
│   ├── go.mod
│   └── go.sum
│
├── api-server/                   # APIサーバー（Go + Gin）
│   ├── cmd/
│   │   └── server/
│   │       └── main.go           # /healthz, CORS（厳格設定）
│   ├── internal/
│   │   ├── handlers/             # HTTP ハンドラー（将来実装）
│   │   ├── services/             # ビジネスロジック（将来実装）
│   │   ├── repositories/         # データアクセス層（将来実装）
│   │   ├── models/               # データモデル（将来実装）
│   │   └── middleware/           # ミドルウェア（将来実装）
│   ├── pkg/
│   │   ├── database/             # データベース接続（将来実装）
│   │   ├── config/               # 設定管理（将来実装）
│   │   └── utils/                 # ユーティリティ（将来実装）
│   ├── migrations/               # DBマイグレーション
│   │   └── 0001_init.sql         # 初期スキーマ（users, routes, spots）
│   ├── go.mod
│   └── go.sum
│
├── web-app/                      # 管理UI（Next.js 14）
│   ├── src/
│   │   └── app/
│   │       ├── layout.tsx
│   │       ├── page.tsx
│   │       └── globals.css
│   ├── package.json
│   ├── next.config.ts
│   └── tsconfig.json
│
├── shared/                       # 共通定義（将来実装）
│   ├── types/                    # 共通型定義
│   └── constants/                # 定数定義
│
├── scripts/                      # 開発用スクリプト（将来実装）
│
├── docs/                         # ドキュメント
│   ├── architecture.md           # アーキテクチャ概要
│   ├── setup.md                  # 環境構築・起動手順
│   ├── api.md                    # API仕様
│   ├── data-model.md             # データモデル
│   ├── external-apis.md          # 外部APIと鍵の扱い
│   ├── operations.md             # 運用・本番移行
│   └── line-bot-explanation.md   # LINE Botの仕組みと設定
│
├── .env.example                  # 環境変数サンプル
├── .env.local                    # 実際の設定（Git管理外）
├── .gitignore
└── README.md
```

### 各コンポーネントの役割

| コンポーネント | ポート | 役割 |
|---|---|---|
| **LINE Bot** | 3001 | LINEからのWebhookを受信して返信 |
| **API Server** | 8080 | Web管理UIから呼ばれるREST API |
| **Web App** | 3000 | 管理UI（ブラウザでアクセス） |

---

## 💡 チーム開発時の注意点

### 環境変数の共有

- `LINE_CHANNEL_SECRET` / `LINE_CHANNEL_TOKEN` は全員で同じ値を共有（同じLINEチャネルを使用）
- `.env.local` はGit管理外なので、各自で設定が必要

### ngrokの使い方

**推奨: 各自でngrokを起動**
- 各開発者が `npx ngrok http 3001` を実行
- 取得したURLをチームで共有してから、LINE Developersで設定
- 使う前に「開発中です」と通知、終了時に「終了しました」と通知

**注意**: ngrokを共有する場合、起動した人のBotサーバーしか動作しません

詳細は [環境構築ガイド - チーム開発](docs/setup.md#チーム開発) を参照

---

## 📝 現在の実装状況

### ✅ 実装済み

- LINE Bot: Webhook受信、署名検証、Flexメッセージ返信
- API Server: CORS設定、以下のRESTエンドポイント

  - Health
    - `GET /healthz`

  - Users（LINE ID 連携）
    - `POST /api/users` 例: `{ "line_id": "Uxxx", "name": "Taro" }`（既存line_idなら既存を返す）
    - `GET /api/users/:line_id`

  - Routes（ユーザーのルート）
    - `POST /api/routes` 例: `{ "line_id": "Uxxx", "start_lat": 35.6, "start_lng": 139.7, "end_lat": 35.7, "end_lng": 139.8 }`
    - `GET /api/users/:line_id/routes`

  - Spots（ルート内スポット）
    - `POST /api/routes/:route_id/spots` 例: `{ "name": "Cafe", "category": "cafe", "lat": 35.6, "lng": 139.7, "url": "https://...", "rating": 4.5 }`
    - `GET /api/routes/:route_id/spots`

- データベース: マイグレーションファイル（users, routes, spotsテーブル）

### 🔄 将来実装予定

- API Server: 認証・認可、更新・削除API、入力バリデーション強化
- Web App: 履歴・お気に入り・統計機能
- 認証: LINE Login統合
- 外部API連携: 経路探索、店舗検索

---

## 🤝 質問・問題があれば

ドキュメントを確認してください：
- 環境構築で困った → [環境構築ガイド](docs/setup.md)
- LINE Botの仕組みが知りたい → [LINE Botの仕組み](docs/line-bot-explanation.md)
- アーキテクチャを知りたい → [アーキテクチャ概要](docs/architecture.md)
