# データモデル

## users
- id: UUID（PK）
- line_id: TEXT
- name: TEXT
- created_at: TIMESTAMPTZ

## routes
- id: UUID（PK）
- user_id: UUID（FK → users.id）
- start_lat, start_lng: DOUBLE PRECISION
- end_lat, end_lng: DOUBLE PRECISION
- created_at: TIMESTAMPTZ

## spots
- id: UUID（PK）
- route_id: UUID（FK → routes.id）
- name: TEXT
- category: TEXT（例: cafe/restaurant/...）
- lat, lng: DOUBLE PRECISION
- url: TEXT
- rating: NUMERIC（将来拡張用, NULL許容）
