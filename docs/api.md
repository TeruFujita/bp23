# API 仕様（現状）

## ヘルスチェック
- GET /healthz
  - 200 OK: `{ "status": "ok" }`

## セキュリティ
- CORS: `ALLOWED_ORIGIN` 環境変数で指定したドメインのみ許可（デフォルト: `http://localhost:3000`）
- 許可メソッド: GET, POST, PUT, DELETE, OPTIONS
- 許可ヘッダー: Authorization, Content-Type
- 認証: 今後追加（JWT/LINE Login想定）

**設定方法**: `.env.local` に `ALLOWED_ORIGIN` を設定（例: `ALLOWED_ORIGIN=http://localhost:3000`）
