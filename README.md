# bp23

LINE Bot + Web アプリケーション + Go API サーバーによる寄り道検索システム

## プロジェクト構成

```
d:\develop\
├── line-bot/                    # LINE Bot アプリケーション
│   ├── src/
│   │   ├── handlers/           # メッセージハンドラー
│   │   ├── services/           # ビジネスロジック
│   │   ├── models/             # データモデル
│   │   └── utils/              # ユーティリティ
│   ├── go.mod
│   └── go.sum
│
├── web-app/                     # Next.js Webアプリケーション
│   ├── src/
│   │   ├── app/                # App Router
│   │   │   ├── (auth)/         # 認証関連ページ
│   │   │   │   ├── login/
│   │   │   │   └── callback/
│   │   │   ├── (dashboard)/    # ダッシュボード関連
│   │   │   │   ├── home/       # ホーム（履歴+マップ）
│   │   │   │   ├── favorites/  # お気に入り一覧
│   │   │   │   └── statistics/ # 統計ページ
│   │   │   ├── api/            # API Routes（必要に応じて）
│   │   │   ├── globals.css
│   │   │   ├── layout.tsx
│   │   │   └── page.tsx
│   │   ├── components/         # 共通コンポーネント
│   │   │   ├── ui/             # 基本UIコンポーネント
│   │   │   ├── map/            # 地図関連コンポーネント
│   │   │   └── layout/         # レイアウトコンポーネント
│   │   ├── features/           # Feature-based構造
│   │   │   ├── auth/           # 認証機能
│   │   │   │   ├── components/
│   │   │   │   ├── hooks/
│   │   │   │   ├── services/
│   │   │   │   └── types/
│   │   │   ├── favorites/      # お気に入り機能
│   │   │   │   ├── components/
│   │   │   │   ├── hooks/
│   │   │   │   ├── services/
│   │   │   │   └── types/
│   │   │   ├── history/        # 履歴機能
│   │   │   │   ├── components/
│   │   │   │   ├── hooks/
│   │   │   │   ├── services/
│   │   │   │   └── types/
│   │   │   └── map/            # 地図機能
│   │   │       ├── components/
│   │   │       ├── hooks/
│   │   │       └── services/
│   │   ├── lib/                # 共通ライブラリ
│   │   │   ├── api/            # API クライアント
│   │   │   ├── auth/           # 認証設定
│   │   │   ├── utils/          # ユーティリティ
│   │   │   └── constants/      # 定数
│   │   ├── types/              # グローバル型定義
│   │   └── styles/             # スタイル
│   ├── package.json
│   ├── next.config.js
│   └── tailwind.config.js
│
├── api-server/                  # Go API サーバー
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── handlers/           # HTTP ハンドラー
│   │   │   ├── auth/
│   │   │   ├── favorites/
│   │   │   ├── history/
│   │   │   └── places/
│   │   ├── services/           # ビジネスロジック
│   │   │   ├── auth/
│   │   │   ├── favorites/
│   │   │   ├── history/
│   │   │   └── places/
│   │   ├── repositories/       # データアクセス層
│   │   │   ├── auth/
│   │   │   ├── favorites/
│   │   │   ├── history/
│   │   │   └── places/
│   │   ├── models/             # データモデル
│   │   └── middleware/         # ミドルウェア
│   ├── pkg/                    # 共通パッケージ
│   │   ├── database/
│   │   ├── config/
│   │   └── utils/
│   ├── migrations/             # DBマイグレーション
│   ├── go.mod
│   └── go.sum
│
├── shared/                      # 共通定義
│   ├── types/                   # 共通型定義
│   └── constants/               # 共通定数
│
├── scripts/                     # 開発用スクリプト
│   ├── setup-dev.sh            # 開発環境セットアップ
│   ├── start-dev.sh            # 開発環境起動
│   └── migrate.sh              # DBマイグレーション実行
│
├── docs/                        # ドキュメント
│   ├── api/                     # API仕様書
│   ├── deployment/              # デプロイ手順
│   └── development/             # 開発ガイド
│
├── .env.example                 # 環境変数サンプル
├── .gitignore
└── README.md
```

## アーキテクチャ概要

### メイン導線
- **LINE Bot**: チャットで経路・寄り道検索を行う

### 補助導線
- **Web (Next.js)**: 履歴やお気に入りを閲覧する

### データ管理
- **Go サーバー + PostgreSQL**: バックエンドAPIとデータベース

### API通信
- **Next.js → Go**: REST API経由でデータ取得

## 開発環境セットアップ

### 前提条件
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+

### 起動手順

1. **PostgreSQL セットアップ**
```bash
createdb detour_bot_dev
```

2. **API サーバー起動**
```bash
cd api-server
go mod tidy
go run cmd/server/main.go
```

3. **Web アプリ起動**
```bash
cd web-app
npm install
npm run dev
```

4. **LINE Bot 起動**
```bash
cd line-bot
go mod tidy
go run src/main.go
```

## 技術スタック

- **LINE Bot**: Go
- **Web アプリ**: Next.js 14 (App Router), TypeScript, Tailwind CSS
- **API サーバー**: Go, Gin, PostgreSQL
- **認証**: LINE Login
- **地図**: Google Maps API / Mapbox

## 管理UI（Web）の役割と方針

- 目的: 管理者が履歴/お気に入り/統計の閲覧と設定を行うためのUI。
- 公開範囲: 管理者専用（開発中はローカル、将来は認証を前提に限定公開）。
- LINE Bot はユーザー向けの会話UIであり、体験確認は実機のLINE上で行う。

## 必要な環境変数

- DB: `DATABASE_URL`
- LINE Messaging: `LINE_CHANNEL_SECRET`, `LINE_CHANNEL_TOKEN`
- LINE Login: `LINE_LOGIN_CLIENT_ID`, `LINE_LOGIN_CLIENT_SECRET`, `LINE_LOGIN_CALLBACK_URL`
- Frontend 公開用: `NEXT_PUBLIC_MAPS_API_KEY`
- API: `JWT_SECRET`, `ALLOWED_ORIGIN`（CORS許可オリジン。既定: `http://localhost:3000`）

## セキュリティポリシー（開発/本番の基本線）

- CORS は許可ドメインを限定（既定はローカル開発用）。
- 秘密情報は `.env`/`.env.local` で管理し、Git には含めない（`.env.example` のみ共有）。
- LINE Webhook は署名検証を必須とする。
- DB は最小権限ユーザーを利用（本番）。

## 開発手順（最短）

1. DB作成: `detour_bot_dev`（PostgreSQL）
2. マイグレーション適用: `api-server/migrations/0001_init.sql` を `psql -f` で適用
3. API起動: `go run cmd/server/main.go`（`ALLOWED_ORIGIN` に `http://localhost:3000` を設定推奨）
4. Web起動: `npm run dev`（管理UI）
5. LINE Bot: `/webhook` 実装（署名検証+エコー）→ ngrok で公開 → Webhook設定

## LINE Bot 開発フロー（最小）

- チャネル作成 → トークン/シークレットを `.env` に設定
- `/webhook` に署名検証+簡易応答（エコー）を実装
- `ngrok http <bot-port>` でトンネル作成 → 発行URLをLINEのWebhookに設定/有効化
- 実機のLINEで送受信を確認

## API 概要（現状）

- `GET /healthz`: 稼働確認
- CORS: `ALLOWED_ORIGIN` で指定したオリジンのみ許可

## マイグレーション

- 初期テーブル: `users`, `favorites`, `history`
- 適用例:
  - `psql "postgres://postgres:YOUR_PASSWORD@localhost:5432/detour_bot_dev?sslmode=disable" -f api-server/migrations/0001_init.sql`

## 起動コマンド（概要）

- API: `cd api-server && go run cmd/server/main.go`
- Web: `cd web-app && npm run dev`
- Bot: `cd line-bot && go run src/main.go`（後で `/webhook` を実装）

## 本番移行時の検討

- 管理UIの認証方式（LINE Login＋ロール付与推奨）
- HTTPS/Reverse Proxy（例: nginx）
- ログ/監視の整備

## データ設計（PostgreSQL）

- users
  - id: UUID（主キー）
  - line_id: TEXT（LINEのユーザー識別子）
  - name: TEXT（表示名）
  - created_at: TIMESTAMPTZ（作成日時）

- routes
  - id: UUID（主キー）
  - user_id: UUID（FK → users.id）
  - start_lat: DOUBLE PRECISION
  - start_lng: DOUBLE PRECISION
  - end_lat: DOUBLE PRECISION
  - end_lng: DOUBLE PRECISION
  - created_at: TIMESTAMPTZ（作成日時）

- spots
  - id: UUID（主キー）
  - route_id: UUID（FK → routes.id）
  - name: TEXT
  - category: TEXT（例: cafe / restaurant / ...）
  - lat: DOUBLE PRECISION
  - lng: DOUBLE PRECISION
  - url: TEXT
  - rating: NUMERIC（将来拡張用、NULL 許容）

## 外部API候補（無料重視）

| 目的 | 候補API | 無料枠・備考 |
| --- | --- | --- |
| 経路探索 | OpenRouteService API | 無料枠 2,500 リクエスト/日（キー要）。ルート・距離・時間取得可 |
| 店舗検索 | ぐるなびAPI / HotPepper API | 無料枠あり・商用可（要申請）。ジャンル/エリア検索 |
| 地図表示 | OpenStreetMap（+ Leaflet など） | 無料・キー不要（自前タイル/制限配慮）。もしくは Mapbox/GoogleMaps（キー要） |
| チャット連携 | LINE Messaging API | 無料でBot開発可（チャネル作成・トークン/シークレット要） |

使用予定の主な外部APIと環境変数例:

- OpenRouteService: `ORS_API_KEY`
- ぐるなび/HotPepper: `GURUNAVI_API_KEY` / `HOTPEPPER_API_KEY`
- 地図表示（どちらか）:
  - OSM+Leaflet: キー不要（ただし利用規約/レート配慮）
  - Mapbox: `NEXT_PUBLIC_MAPBOX_TOKEN`
  - Google Maps: `NEXT_PUBLIC_MAPS_API_KEY`
- LINE Messaging: `LINE_CHANNEL_SECRET`, `LINE_CHANNEL_TOKEN`