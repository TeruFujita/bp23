# 運用・本番移行

## 認証/認可
- 管理UIは管理者認証（LINE Login + adminロール推奨）

## ネットワーク/セキュリティ
- CORSは本番ドメインに限定（ALLOWED_ORIGIN）
- HTTPS 終端（CDN/リバプロ）

## デプロイ/監視
- ログ出力（構造化ログ）
- 監視（ヘルスチェック, 重要メトリクス）

## 秘匿情報
- .envは本番では使わず、Secret Manager/KMS等を利用
