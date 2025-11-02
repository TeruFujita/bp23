# アーキテクチャ概要

## システム構成

### メイン導線
- **LINE Bot**: 会話形式で経路・寄り道検索

### 補助導線
- **Web管理UI**: 履歴・お気に入り・統計・設定の閲覧・管理

### バックエンド
- **API Server**: Go (Gin) + PostgreSQL
- **通信方式**: 
  - Web → API: REST API
  - Bot → API: 内部呼び出し（将来拡張）

---

## 構成要素

| コンポーネント | 技術 | ポート | 役割 |
|---|---|---|---|
| **line-bot** | Go | 3001 | LINE Messaging API Webhook受信・返信 |
| **api-server** | Go + Gin | 8080 | REST API（認証/履歴/お気に入り/統計） |
| **web-app** | Next.js | 3000 | 管理UI（ブラウザ） |
| **shared** | - | - | 共通型/定数（将来実装） |

---

## 技術スタック

- **LINE Bot**: Go
- **Webアプリ**: Next.js (App Router), TypeScript, Tailwind CSS
- **APIサーバー**: Go, Gin, PostgreSQL
- **認証**: LINE Login（管理UIは管理者向け、将来実装）

---

## データフロー

### LINE Botの流れ

```
[LINEアプリ] 
  → LINE Messaging API
  → ngrok（開発時のみ）
  → line-bot (localhost:3001)
  → 処理・返信
  → LINE Messaging API
  → [LINEアプリ]
```

### Web管理UIの流れ

```
[ブラウザ] 
  → web-app (localhost:3000)
  → api-server (localhost:8080)
  → PostgreSQL
  → データ取得
  → JSONレスポンス
  → [ブラウザ]
```

---

## 管理UIの方針

- **役割**: 管理者が履歴/お気に入り/統計の閲覧と設定を行う
- **公開範囲**: 管理者専用（開発中はローカル、本番は認証必須）

---

## セキュリティ

- **CORS**: 許可ドメイン限定（環境変数 `ALLOWED_ORIGIN`、デフォルト: `http://localhost:3000`）
- **秘密情報**: `.env.local` で管理（`.gitignore` に含まれ、Git管理外）
- **LINE Webhook**: 署名検証を必須（`LINE_CHANNEL_SECRET` を使用）
- **環境変数**: LINE Botは `godotenv` で `.env.local` を自動読み込み

---

## プロジェクト構造

```
line_bot/
├── line-bot/              # LINE Botサーバー
│   ├── src/main.go        # Webhook処理
│   └── ...
│
├── api-server/            # REST APIサーバー
│   ├── cmd/server/        # エントリーポイント
│   ├── internal/          # 内部パッケージ
│   │   ├── handlers/      # HTTPハンドラー
│   │   ├── services/      # ビジネスロジック
│   │   ├── repositories/  # データアクセス層
│   │   └── models/        # データモデル
│   ├── migrations/        # データベースマイグレーション
│   └── ...
│
├── web-app/               # 管理UI
│   ├── src/app/           # Next.js App Router
│   └── ...
│
└── docs/                  # ドキュメント
```

---

## 将来の拡張

- Bot → API: 内部呼び出しによるデータ保存
- 認証機能: LINE Login統合
- 外部API連携: 経路探索、店舗検索
