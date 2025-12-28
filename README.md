# hackathon-util

ハッカソンでDiscordのロール・チャンネル・カテゴリを自動管理するツール集

[サンプルのスプレッドシート](https://docs.google.com/spreadsheets/d/1kOFmbrdYd4gsF3i0bo5PuteUYWqq5R-g0i65jdRZMy0/edit?usp=sharing)

![](./image/img1.png)

## 提供ツール

### cmd/sheet-to-discord
Googleスプレッドシートからチーム情報を読み取り、Discordにロール・カテゴリ・チャンネルを自動生成するスクリプト

**機能:**
- 全参加者用の共通ロール`@参加者_{EVENT_NAME}`の付与
- メンター用ロール`@メンター_{EVENT_NAME}`の作成
- チームごとのロール作成
- チームごとのカテゴリ作成（テキストチャンネル「やりとり」とボイスチャンネル「会話」を含む）
  - オプション`PRIVATE_VC`でボイスチャンネル「会話」は、全参加者用ロールまたはメンター用ロールを持っている人のみに表示する
  - オプション`PRIVATE_CATEGORY`でチームカテゴリは、該当チームロールまたはメンター用ロールを持っている人のみに表示する（他チームのカテゴリが非表示になる）
- メンバーへのロール自動付与
- Discord上に存在しないユーザーの一覧表示

**実行方法:**

環境変数を埋めてから下記を実行

```bash
go run cmd/sheet-to-discord/main.go

# ボイスチャンネル「会話」のみをプライベートチャンネルにする
PRIVATE_VC=true go run cmd/sheet-to-discord/main.go

# カテゴリをプライベートにする
PRIVATE_CATEGORY=true go run cmd/sheet-to-discord/main.go
```

### cmd/sheet-to-discord-delete
ハッカソン終了後のクリーンアップ用スクリプト

**機能:**
- スプレッドシートに記載されたチームのカテゴリ・チャンネル削除
- チームロールの削除とメンバーからの削除
- オプションで全参加者ロール`@参加者_{EVENT_NAME}`からメンバーを削除

**実行方法:**
```bash
# ドライラン（実際には削除しない）
DRY_RUN=true go run cmd/sheet-to-discord-delete/main.go

# 実際に削除
DRY_RUN=false go run cmd/sheet-to-discord-delete/main.go

# 参加者ロールからもメンバーを削除する場合
DRY_RUN=false REMOVE_ALL_MEMBERS=true go run cmd/sheet-to-discord-delete/main.go
```

## 注意点

1. 欠席というチーム名を付けた場合無視されます。
2. メンバーは1チーム1行にしてください。
3. ユーザ名はdiscordの @から始まるIDを @ なしで入力してください。

## 開発

### 環境変数を設定

```bash
# .env.example からコピー
cp .env.example .env
```

### .env の設定

```env
# Application Env
ENV=dev                                                                # 開発環境: dev, 本番環境: prod

# Google Env
GOOGLE_SPREADSHEET_ID=    # 対象のスプレッドシートID
TEAM_RANGE=チームシート!A2:F15                                         # チーム情報の範囲
GOOGLE_CREDENTIALS_FILE=./credential.json                             # Google認証情報ファイルのパス
EVENT_NAME=Progateハッカソン_25.6                                       # 参加者・メンターロール用

# Discord Env
DISCORD_BOT_TOKEN=                                                     # DiscordのBotトークン
DISCORD_GUILD_ID=                                                      # 対象のDiscordサーバーID
```

#### DISCORD_BOT_TOKENの設定方法

1. [Discord Developer Portal](https://discord.com/developers/applications) にアクセス
2. New Application からアプリケーションを作成
3. サイドバーのBotタブ から トークン を作成し、コピーして DISCORD_BOT_TOKEN に設定

#### DISCORD_GUILD_ID の設定方法（スクリプトで実行したい場合のみ）

1. サーバー名で右クリックしてメニューの一番下にある「サーバーIDをコピー」をクリック

#### credential ファイルの生成（任意）

1. Google Cloudから [スプレッドシートAPI](https://console.cloud.google.com/apis/library/sheets.googleapis.com?hl=ja)を有効にする 
2. スプレッドシートAPIの管理から認証情報 -> 認証情報の作成 -> サービスアカウントを選択
3. 適切な権限のSAを作ったら保存
4. サービス アカウントから保存されたSAを選択 -> キー -> 鍵を追加 -> 新しい鍵 -> json を選択
5. 生成されたjsonを落として、hackathon-utilの直下に"credential.json"として保存

### Botのローカルでの実行

```bash
# パッケージのインストール
go install github.com/air-verse/air@latest
```

`.docker/app/sheetless.air.toml` のローカル用のコメントアウトを外す

```bash
# ローカルで実行
air -c .docker/app/sheetless.air.toml
```
