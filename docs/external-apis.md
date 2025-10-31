# 外部APIと鍵の扱い

## 候補（無料重視）
- 経路探索: OpenRouteService API（無料 2,500/日, 要キー）
- 店舗検索: ぐるなびAPI / HotPepper API（無料枠・商用可, 要申請）
- 地図表示: OpenStreetMap（無料・キー不要）/ Mapbox / Google Maps
- チャット連携: LINE Messaging API（無料でBot開発可）

## 環境変数
- `ORS_API_KEY`
- `GURUNAVI_API_KEY` / `HOTPEPPER_API_KEY`
- `NEXT_PUBLIC_MAPBOX_TOKEN` / `NEXT_PUBLIC_MAPS_API_KEY`
- `LINE_CHANNEL_SECRET`, `LINE_CHANNEL_TOKEN`

## セキュリティ
- 秘密鍵は `.env`/Secret Manager で管理し、Gitに含めない
- レート/利用規約を順守
