# API 仕様

## エンドポイント一覧

| HTTPメソッド | エンドポイント | パスパラメータ | リクエストボディ主要パラメータ | 説明 |
| --- | --- | --- | --- | --- |
| GET | `/healthz` | - | - | サーバーとDBの状態確認 |
| POST | `/api/users` | - | `line_id`, `name` | ユーザー作成/取得 |
| GET | `/api/users/:line_id` | `line_id` | - | ユーザー情報取得 |
| POST | `/api/routes` | - | `line_id`, `start_lat`, `start_lng`, `end_lat`, `end_lng` | ルート作成 |
| GET | `/api/users/:line_id/routes` | `line_id` | - | ユーザーのルート一覧取得 |
| POST | `/api/routes/:route_id/spots` | `route_id` | `name`, `category`, `lat`, `lng`, `url`, `rating` | スポット追加 |
| GET | `/api/routes/:route_id/spots` | `route_id` | - | ルートのスポット一覧取得 |

---

## ヘルスチェック

### GET /healthz
サーバーとデータベースの状態を確認します。

**レスポンス**
- 200 OK: `{ "status": "ok", "db": "ok" }` または `{ "status": "ok", "db": "disabled" }`
- 503 Service Unavailable: `{ "status": "degraded", "db": "エラーメッセージ" }`（DB接続失敗時）

---

## Users API

### POST /api/users
新しいユーザーを作成、または既存ユーザーを取得します（LINE IDが既に存在する場合は既存のIDを返します）。

**リクエストボディ**
```json
{
  "line_id": "U1234567890abcdef",
  "name": "ユーザー名"  // オプション
}
```

**レスポンス**
- 201 Created: `{ "id": "uuid" }`
- 400 Bad Request: `{ "error": "invalid json" }`
- 503 Service Unavailable: `{ "error": "db not configured" }`

---

### GET /api/users/:line_id
LINE IDでユーザー情報を取得します。

**パスパラメータ**
- `line_id`: LINEユーザー識別子

**レスポンス**
- 200 OK: `{ "id": "uuid", "line_id": "U1234567890abcdef", "name": "ユーザー名", "created_at": "2024-01-01T00:00:00Z" }`
- 404 Not Found: `{ "error": "user not found" }`
- 503 Service Unavailable: `{ "error": "db not configured" }`

---

## Routes API

### POST /api/routes
新しいルート（経路）を作成します。

**リクエストボディ**
```json
{
  "line_id": "U1234567890abcdef",
  "start_lat": 35.681236,
  "start_lng": 139.767125,
  "end_lat": 35.658034,
  "end_lng": 139.701636
}
```

**レスポンス**
- 201 Created: `{ "id": "uuid" }`
- 400 Bad Request: `{ "error": "invalid json" }`
- 404 Not Found: `{ "error": "user not found" }`（指定されたLINE IDが存在しない場合）
- 503 Service Unavailable: `{ "error": "db not configured" }`

---

### GET /api/users/:line_id/routes
指定したユーザーのルート一覧を取得します（作成日時の降順）。

**パスパラメータ**
- `line_id`: LINEユーザー識別子

**レスポンス**
- 200 OK: `{ "routes": [{ "id": "uuid", "start_lat": 35.681236, "start_lng": 139.767125, "end_lat": 35.658034, "end_lng": 139.701636, "created_at": "2024-01-01T00:00:00Z" }, ...] }`
- 503 Service Unavailable: `{ "error": "db not configured" }`

---

## Spots API

### POST /api/routes/:route_id/spots
ルートに新しいスポット（地点）を追加します。

**パスパラメータ**
- `route_id`: ルートID（UUID）

**リクエストボディ**
```json
{
  "route_id": "uuid",
  "name": "カフェ名",
  "category": "cafe",
  "lat": 35.681236,
  "lng": 139.767125,
  "url": "https://example.com"  // オプション
  "rating": 4.5  // オプション
}
```

**レスポンス**
- 201 Created: `{ "id": "uuid" }`
- 400 Bad Request: `{ "error": "invalid json" }`
- 503 Service Unavailable: `{ "error": "db not configured" }`

---

### GET /api/routes/:route_id/spots
指定したルートのスポット一覧を取得します（作成日時の降順）。

**パスパラメータ**
- `route_id`: ルートID（UUID）

**レスポンス**
- 200 OK: `{ "spots": [{ "id": "uuid", "name": "カフェ名", "category": "cafe", "lat": 35.681236, "lng": 139.767125, "url": "https://example.com", "rating": 4.5, "created_at": "2024-01-01T00:00:00Z" }, ...] }`
- 503 Service Unavailable: `{ "error": "db not configured" }`

---

## セキュリティ

### CORS設定
- `ALLOWED_ORIGIN` 環境変数で指定したドメインのみ許可（デフォルト: `http://localhost:3000`）
- 許可メソッド: GET, POST, PUT, DELETE, OPTIONS
- 許可ヘッダー: Authorization, Content-Type

**設定方法**: `.env.local` に `ALLOWED_ORIGIN` を設定（例: `ALLOWED_ORIGIN=http://localhost:3000`）

### 認証
- 現在は未実装（将来、JWT/LINE Loginを追加予定）
