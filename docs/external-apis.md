# 外部APIと鍵の扱い

## 候補（無料重視）
- 経路探索: OpenRouteService API（無料 2,500/日, 要キー）
- 店舗検索: ぐるなびAPI / HotPepper API（無料枠・商用可, 要申請）
- 地図表示: OpenStreetMap（無料・キー不要）/ Mapbox / Google Maps
- チャット連携: LINE Messaging API（無料でBot開発可）

## 環境変数（.env.local に設定）

### LINE Messaging API（必須）
- `LINE_CHANNEL_SECRET`: LINE Developersで取得
- `LINE_CHANNEL_TOKEN`: LINE Developersで取得

### 経路探索
- `ORS_API_KEY`: OpenRouteService APIキー（無料枠: 2,500リクエスト/日）

### 店舗検索
- `GURUNAVI_API_KEY`: ぐるなびAPIキー
- `HOTPEPPER_API_KEY`: HotPepper APIキー

### 地図表示
- `NEXT_PUBLIC_MAPS_API_KEY`: Google Maps等のAPIキー
- `NEXT_PUBLIC_MAPBOX_TOKEN`: Mapboxトークン（代替）

**注意**: `NEXT_PUBLIC_*` はNext.jsでクライアント側に公開されるため、機密情報を設定しないこと

## セキュリティ
- 秘密鍵は `.env.local` で管理し、Gitに含めない（`.gitignore` に追加済み）
- レート制限/利用規約を順守
