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

## 必要な外部サービスアカウント登録（開発開始前）

### ngrok（ローカル開発でLINE Webhook受信に必要）
1. **アカウント登録**: https://dashboard.ngrok.com/signup （無料）
2. **authtoken取得**: https://dashboard.ngrok.com/get-started/your-authtoken
3. **authtoken設定**:
   ```powershell
   npx ngrok config add-authtoken YOUR_AUTHTOKEN
   ```
4. **インストール**（未導入の場合）:
   ```powershell
   npm install -g ngrok
   ```

### LINE Developers（LINE Bot作成に必要）
1. **アカウント登録**: https://developers.line.biz/
2. **チャネル作成**: 新規プロバイダー作成 → Messaging APIチャネル作成
3. **必要情報の取得**:
   - チャネルシークレット → `.env.local` の `LINE_CHANNEL_SECRET`
   - チャネルアクセストークン → `.env.local` の `LINE_CHANNEL_TOKEN`

### その他の外部API（後で設定可能）
- OpenRouteService API: https://openrouteservice.org/ （経路探索）
- ぐるなびAPI / HotPepper API: 必要に応じて申請

詳細は [外部APIと鍵の扱い](docs/external-apis.md) を参照

### チーム開発時の注意点

#### Webhook URLの設定方法（選択肢）

**パターンA: ngrokを個別に起動（推奨）**
- 各開発者が各自のngrokを起動
- 各自のWebhook URLをLINE Developersに設定（使う前に共有）
- **設定手順：**
  1. 各自で `npx ngrok http 3001` を起動
  2. 表示されたURL（例: `https://xxxxx.ngrok.io`）をチームで共有
  3. 使う前に「開発中です。Webhook URL: https://xxxxx.ngrok.io/webhook」と通知
  4. LINE DevelopersでWebhook URLを設定・有効化
  5. 終了時に「終了しました」と通知
  6. 次の開発者が同じ手順で自分のURLを設定

**パターンB: ngrokを共有で使う**
- 1人がngrokを起動し、そのURLを全員で共有
- 全員が同じWebhook URLを使用
- **設定手順：**
  1. 1人が `npx ngrok http 3001` を起動してURLを共有
  2. 全員が同じWebhook URL（例: `https://xxxxx.ngrok.io/webhook`）をLINE Developersに設定
  3. 各開発者は自分のローカルBotサーバー（localhost:3001）を起動
  4. **注意**: ngrokは起動した人のローカルサーバー（localhost:3001）に転送するため、**起動した人のBotサーバーしか動作しない**

**共通設定：**
- `.env.local` の `LINE_CHANNEL_SECRET` / `LINE_CHANNEL_TOKEN` は同じチャンネルの値を全員で共有
- 「応答メッセージ」はOFF推奨
- ngrokのWeb UI（`http://127.0.0.1:4040`）でリクエスト履歴を確認可能


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

