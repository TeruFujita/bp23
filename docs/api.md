# API 仕様（現状）

## ヘルスチェック
- GET /healthz
  - 200 OK: `{ "status": "ok" }`

## セキュリティ
- CORS: `ALLOWED_ORIGIN` で指定したドメインのみ許可
- 認証: 今後追加（JWT/LINE Login想定）
