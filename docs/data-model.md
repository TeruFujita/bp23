# データモデル

## users

| カラム名 | 型 | 制約・備考 |
| --- | --- | --- |
| id | UUID | PK |
| line_id | TEXT | LINEユーザー識別子、将来 UNIQUE 検討 |
| name | TEXT | 表示名、NULL可 |
| created_at | TIMESTAMPTZ | 既定: now() |

## routes

| カラム名 | 型 | 制約・備考 |
| --- | --- | --- |
| id | UUID | PK |
| user_id | UUID | FK → users.id（ON DELETE CASCADE 推奨） |
| start_lat | DOUBLE PRECISION | 緯度（出発） |
| start_lng | DOUBLE PRECISION | 経度（出発） |
| end_lat | DOUBLE PRECISION | 緯度（到着） |
| end_lng | DOUBLE PRECISION | 経度（到着） |
| created_at | TIMESTAMPTZ | 既定: now() |

## spots

| カラム名 | 型 | 制約・備考 |
| --- | --- | --- |
| id | UUID | PK |
| route_id | UUID | FK → routes.id（ON DELETE CASCADE 推奨） |
| name | TEXT | 店舗・地点名 |
| category | TEXT | 例: cafe / restaurant / ... |
| lat | DOUBLE PRECISION | 緯度 |
| lng | DOUBLE PRECISION | 経度 |
| url | TEXT | 詳細URL、NULL可 |
| rating | NUMERIC | 将来拡張用、NULL可 |
