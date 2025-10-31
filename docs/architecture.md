# アーキテクチャ概要

- メイン導線: LINE Bot（会話で経路・寄り道検索）
- 補助導線: Web(Next.js) 管理UI（履歴/お気に入り/統計/設定）
- バックエンド: Go(Gin) + PostgreSQL
- 通信: Web → API は REST、Bot → API は内部呼び出し（将来拡張）

## 構成要素
- line-bot: LINE Messaging API Webhook 受信・返信
- api-server: REST API（認証/履歴/お気に入り/統計）
- web-app: Next.js 管理UI
- shared: 共通型/定数

## 技術スタック
- LINE Bot: Go
- Web アプリ: Next.js (App Router), TypeScript, Tailwind CSS
- API サーバー: Go, Gin, PostgreSQL
- 認証: LINE Login（管理UIは管理者向け）

## 管理UIの方針
- 役割: 管理者が履歴/お気に入り/統計の閲覧と設定を行う
- 公開範囲: 管理者専用（開発中はローカル、本番は認証必須）

## セキュリティ
- CORS は許可ドメイン限定（環境変数 ALLOWED_ORIGIN）
- 秘密情報は .env/.env.local で管理（Git管理外）
- LINE Webhook は署名検証を必須
