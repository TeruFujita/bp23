# bp23

LINE Bot + Web アプリケーション + Go API サーバーによる寄り道検索システム

## ドキュメント

- [アーキテクチャ概要](docs/architecture.md)
- [環境構築・起動手順](docs/setup.md)
- [API 仕様](docs/api.md)
- [データモデル](docs/data-model.md)
- [外部APIと鍵の扱い](docs/external-apis.md)
- [運用・本番移行](docs/operations.md)

### 学習・理解用

- [LINE Bot の仕組みと設定（ngrok/Webhookの理解）](docs/line-bot-explanation.md)

## プロジェクト構成

```
d:\develop\line_bot\
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

