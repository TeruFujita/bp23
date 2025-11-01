# LINE Bot の仕組みと設定

## .env.example で設定すべき項目（優先順位順）

### 🔴 必須（LINE Bot動作に必要）
1. **LINE_CHANNEL_SECRET**
   - LINE Developers コンソール → チャネル設定 → チャネルシークレット
   - Webhook署名検証に使用
2. **LINE_CHANNEL_TOKEN**
   - LINE Developers コンソール → Messaging API → チャネルアクセストークン
   - BotからLINE APIを呼び出すために使用

### 🟡 開発時に設定（推奨）
3. **PORT**
   - Botサーバーのポート番号（既定: 3001）
4. **ALLOWED_ORIGIN**
   - APIサーバーのCORS許可オリジン（既定: http://localhost:3000）

### 🟢 後で設定（今は空でOK）
- DATABASE_URL（DB接続時）
- LINE_LOGIN_*（管理UI認証時）
- NEXT_PUBLIC_MAPS_API_KEY（地図表示時）
- ORS_API_KEY等（外部API利用時）

## LINE Bot の仕組み

### 全体の流れ

```
[あなたのスマホのLINEアプリ]
        |
        | メッセージ送信
        v
[LINE Messaging API サーバー（クラウド）]
        |
        | Webhook（HTTPS POST）
        | 署名付きで送信
        v
[ngrok 公開URL]
        |
        | HTTPS → ローカルへ転送
        v
[あなたのPC: localhost:3001]
        |
        | /webhook エンドポイントで受信
        v
[line-bot/src/main.go]
        |
        | 1. 署名検証（LINE_CHANNEL_SECRET使用）
        | 2. メッセージ解析
        | 3. Flexメッセージ作成
        | 4. LINE APIへ返信（LINE_CHANNEL_TOKEN使用）
        v
[LINE Messaging API サーバー]
        |
        | ユーザーにメッセージ配信
        v
[あなたのスマホのLINEアプリ]
```

## ngrok とは？

### 概要
- **ローカルサーバーをインターネット経由でアクセス可能にするツール（中継/プロキシ）**
- 開発中、LINEのWebhookはHTTPSの公開URLが必要だが、ローカルは `http://localhost:3001` しかない
- ngrokが一時的な公開URL（例: `https://xxxxx.ngrok.io`）を提供し、ローカルサーバーに転送

### なぜ必要？
- LINE Messaging APIは**インターネット上の公開URLにしかWebhookを送れない**
- ローカルの `localhost` には直接アクセスできない
- ngrokが「トンネル」を作り、公開URL → ローカルサーバーを接続

### ngrokの役割（正確な理解）

ngrokは「仮想的な公開サーバーのように見せて、ローカルサーバーとLINEの橋渡しをする中継役」です。

**正確な流れ：**
```
[LINE Messaging API]
        |
        | POST /webhook (HTTPS) ← 公開URLが必要
        v
[ngrok 公開URL] ← 仮想的な公開サーバーとして振る舞う
https://xxxxx.ngrok.io
        |
        | 中継（トンネル/プロキシ）
        | リクエストをローカルの localhost:3001 に転送
        v
[あなたのPC: line-bot/src/main.go]
        |
        | ★実際の処理はここで行われる★
        | - 署名検証
        | - メッセージ解析
        | - Flexメッセージ作成
        | - レスポンス生成
        v
[ngrok] ← レスポンスを受け取りLINEに中継
        |
        | レスポンスをLINEに返す
        v
[LINE Messaging API]
        |
        v
[あなたのLINEアプリにメッセージ表示]
```

**重要なポイント：**
- ngrok自体はサーバーではなく「中継役（トンネル/プロキシ）」
- 公開URLとして見せかけて、リクエスト/レスポンスを中継する
- 実際の処理（署名検証、Flex作成など）は**ローカルの line-bot サーバーで実行される**
- 「ngrokが仮想的な公開サーバーのように見せて、ローカルサーバーとLINEの橋渡しをしている」という理解でOK

### 使い方
```powershell
# 1. ngrokを起動（Botサーバーが3001で起動している状態で）
ngrok http 3001

# 2. 表示されるURLをコピー
# Forwarding: https://xxxxx.ngrok.io -> http://localhost:3001

# 3. LINE Developers の Webhook URL に設定
# https://xxxxx.ngrok.io/webhook
```

## サーバーの設定

### 現在の構成

| サーバー | ポート | 役割 | 起動コマンド |
| --- | --- | --- | --- |
| API | 8080 | REST API（Webから呼ばれる） | `cd api-server && go run cmd/server/main.go` |
| Web | 3000 | 管理UI（ブラウザで見る） | `cd web-app && npm run dev` |
| Bot | 3001 | LINE Webhook受信（LINEから呼ばれる） | `cd line-bot && go run src/main.go` |

### LINE Botサーバー（line-bot/src/main.go）の処理

1. **起動時**
   - 環境変数から `LINE_CHANNEL_SECRET` と `LINE_CHANNEL_TOKEN` を読み込む
   - `/webhook` エンドポイントを開く（ポート3001）

2. **Webhook受信時**
   - LINEが送ってくるリクエストの署名を検証（`LINE_CHANNEL_SECRET`使用）
   - メッセージイベントを解析
   - Flexメッセージを作成して返信（`LINE_CHANNEL_TOKEN`使用）

## LINE画面からBotが見えるようにするには

### 必要な設定

1. **LINE Developers でチャネル作成**
   - https://developers.line.biz/ でログイン
   - 「新規プロバイダー作成」→「Messaging APIチャネル作成」

2. **チャネル設定で取得**
   - チャネルシークレット → `.env` の `LINE_CHANNEL_SECRET`
   - チャネルアクセストークン → `.env` の `LINE_CHANNEL_TOKEN`

3. **Webhook設定**
   - 「Messaging API設定」タブ
   - Webhook URL: `https://xxxxx.ngrok.io/webhook`（ngrokで取得したURL）
   - 「Webhookの利用」をON

4. **友だち追加**
   - 「チャネル基本設定」→ QRコードまたは友だち追加URL
   - 自分のLINEで友だち追加

5. **動作確認**
   - メッセージ送信 → BotからFlexメッセージが返ってくればOK

## よくある質問

### Q: ngrokは常に起動しておく必要がある？
A: 開発中は必要。本番環境（クラウドにデプロイ）では不要。

### Q: ngrokのURLは毎回変わる？
A: 無料版は起動するたびに変わる。固定したい場合は有料版を使う。

### Q: 署名検証って何？
A: LINEが本当に送ってきたリクエストか確認する仕組み。悪意のあるリクエストを防ぐ。

