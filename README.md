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
   - チャネルシークレット → `.env.local` の `LINE_CHANNEL_SECRET` に設定
   - チャネルアクセストークン → `.env.local` の `LINE_CHANNEL_TOKEN` に設定
   - **注意**: `.env.local` は `line-bot/src/main.go` が自動で読み込みます（godotenv使用）

### その他の外部API（後で設定可能）
- OpenRouteService API: https://openrouteservice.org/ （経路探索）
- ぐるなびAPI / HotPepper API: 必要に応じて申請

詳細は [外部APIと鍵の扱い](docs/external-apis.md) を参照

### 環境変数について

`.env.local` に以下の環境変数を設定します（`.env.example` をコピーして実際の値を設定）。

**必須（LINE Bot動作に必要）：**
- `LINE_CHANNEL_SECRET`: LINE Developersで取得
- `LINE_CHANNEL_TOKEN`: LINE Developersで取得

**その他の主要環境変数：**
- `DATABASE_URL`: PostgreSQL接続文字列
- `ALLOWED_ORIGIN`: CORS許可ドメイン（デフォルト: `http://localhost:3000`）
- `PORT`: Botサーバーポート（デフォルト: `3001`）
- `JWT_SECRET`: JWT署名用
- `LINE_LOGIN_*`: LINE Login設定（将来実装）
- `NEXT_PUBLIC_*`: フロントエンド公開変数
- `ORS_API_KEY`, `GURUNAVI_API_KEY`, `HOTPEPPER_API_KEY`: 外部APIキー（オプション）

**詳細は [環境構築・起動手順](docs/setup.md#環境変数) を参照**

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
├── line-bot\                     # LINE Bot（Go）
│   ├── src\
│   │   ├── main.go               # Webhook実装済み（署名検証/Flex返信）
│   │   ├── handlers\             # メッセージハンドラー（将来実装）
│   │   ├── services\             # ビジネスロジック（将来実装）
│   │   ├── models\               # データモデル（将来実装）
│   │   └── utils\                 # ユーティリティ（将来実装）
│   ├── go.mod
│   └── go.sum
│
├── api-server\                   # APIサーバー（Go + Gin）
│   ├── cmd\server\main.go       # /healthz, CORS（厳格設定）
│   ├── internal\
│   │   ├── handlers\            # HTTP ハンドラー（将来実装）
│   │   ├── services\            # ビジネスロジック（将来実装）
│   │   ├── repositories\        # データアクセス層（将来実装）
│   │   ├── models\              # データモデル（将来実装）
│   │   └── middleware\          # ミドルウェア（将来実装）
│   ├── pkg\
│   │   ├── database\
│   │   ├── config\
│   │   └── utils\
│   ├── migrations\              # DBマイグレーション（将来実装）
│   ├── go.mod
│   └── go.sum
│
├── web-app\                      # 管理UI（Next.js 14）
│   ├── src\app\
│   │   ├── layout.tsx
│   │   ├── page.tsx
│   │   └── globals.css
│   ├── package.json
│   └── next.config.ts
│
├── shared\                       # 共通定義（将来実装）
│   ├── types\
│   └── constants\
│
├── scripts\                      # 開発用スクリプト（将来実装）
│
├── docs\                         # ドキュメント
│   ├── architecture.md
│   ├── setup.md
│   ├── api.md
│   ├── data-model.md
│   ├── external-apis.md
│   ├── operations.md
│   └── line-bot-explanation.md
│
├── .env.example                  # 環境変数サンプル
├── .env.local                    # 実際の設定（Git管理外）
├── .gitignore
└── README.md
```

