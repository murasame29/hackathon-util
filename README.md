# hackathon-util

ハッカソンでDiscordのロール・チャンネル・カテゴリを自動管理するツール集

[サンプルのスプレッドシート](https://docs.google.com/spreadsheets/d/1kOFmbrdYd4gsF3i0bo5PuteUYWqq5R-g0i65jdRZMy0/edit?usp=sharing)

![](./image/img1.png)

## 提供ツール

### cmd/sheet-to-discord

Googleスプレッドシートからチーム情報を読み取り、Discordにロール・カテゴリ・チャンネルを自動生成するスクリプト

**機能:**
- 全参加者用の共通ロール `@参加者_{EVENT_NAME}` の作成・付与
- メンター用ロール `@メンター_{EVENT_NAME}` の作成（色: #3498db）
- チームごとのロール作成
- チームごとのカテゴリ作成（テキストチャンネル「やりとり」とボイスチャンネル「会話」を含む）
  - カテゴリ・チャンネルがすでに存在する場合は権限のみ更新
  - `PRIVATE_VC=true` でボイスチャンネル「会話」を参加者ロール・メンターロール保持者のみに表示
  - `PRIVATE_CATEGORY=true` でカテゴリ全体をチームロール・メンターロール保持者のみに表示
  - `VORTEX_MUTEROLE_ID` を設定すると、そのロールに対してメッセージ送信・リアクション・VC接続などを禁止
- スプレッドシートの各行 B〜F列のユーザー名（最大5名）にチームロールと参加者ロールを付与
- Discord上に存在しないユーザーの一覧を実行後に表示

**環境変数による権限設定:**

| `PRIVATE_CATEGORY` | `PRIVATE_VC` | `#やりとり`                   | `#会話`                       |
| ------------------ | ------------ | ----------------------------- | ----------------------------- |
| `false`            | `false`      | `@everyone`                   | `@everyone`                   |
| `false`            | `true`       | `@everyone`                   | 参加者ロール + メンターロール |
| `true`             | true/false   | チームロール + メンターロール | チームロール + メンターロール |

※ `PRIVATE_CATEGORY=true` の場合、`PRIVATE_VC` の設定は上書きされます

**実行方法:**

```bash
# 通常実行
go run cmd/sheet-to-discord/main.go

# ボイスチャンネル「会話」のみをプライベートにする
PRIVATE_VC=true go run cmd/sheet-to-discord/main.go

# カテゴリをプライベートにする
PRIVATE_CATEGORY=true go run cmd/sheet-to-discord/main.go
```

---

### cmd/sheet-to-discord-delete

ハッカソン終了後のクリーンアップ用スクリプト

**機能:**
- スプレッドシートに記載されたチームのカテゴリ・配下チャンネルを削除
- チームロールをメンバーから剥奪し、ロール自体を削除
- `REMOVE_ALL_MEMBERS=true` で全参加者ロール `@参加者_{EVENT_NAME}` からメンバーを剥奪（ロール自体は削除しない）
- `DRY_RUN=true`（デフォルト）で実際の削除を行わず対象を確認できる

**実行方法:**

```bash
# ドライラン（デフォルト）: 削除対象を確認するだけで実際には削除しない
go run cmd/sheet-to-discord-delete/main.go

# 実際に削除
DRY_RUN=false go run cmd/sheet-to-discord-delete/main.go

# 参加者ロールからもメンバーを剥奪する場合
DRY_RUN=false REMOVE_ALL_MEMBERS=true go run cmd/sheet-to-discord-delete/main.go
```

---

## セットアップ

### 環境変数を設定

```bash
cp .env.example .env
```

### 環境変数一覧

| 変数名 | 必須 | 説明 |
|---|---|---|
| `GOOGLE_SPREADSHEET_ID` | ✅ | 対象のスプレッドシートID |
| `TEAM_RANGE` | ✅ | チーム情報の範囲（例: `チームシート!A2:F15`） |
| `GOOGLE_CREDENTIALS_FILE` | ✅ | Google認証情報ファイルのパス |
| `EVENT_NAME` | ✅ | イベント名。参加者・メンターロール名のサフィックスに使用 |
| `DISCORD_BOT_TOKEN` | ✅ | DiscordのBotトークン |
| `DISCORD_GUILD_ID` | ✅ | 対象のDiscordサーバーID |
| `PRIVATE_VC` | - | `true` でボイスチャンネルをプライベートにする（デフォルト: `false`） |
| `PRIVATE_CATEGORY` | - | `true` でカテゴリをプライベートにする（デフォルト: `false`） |
| `VORTEX_MUTEROLE_ID` | - | ミュートロールのID。設定するとそのロールの発言・VC接続を禁止 |
| `DRY_RUN` | - | `sheet-to-discord-delete` 用。`false` で実際に削除（デフォルト: `true`） |
| `REMOVE_ALL_MEMBERS` | - | `sheet-to-discord-delete` 用。`true` で参加者ロールからもメンバーを剥奪（デフォルト: `false`） |

### スプレッドシートのフォーマット

| A列（チーム名） | B列（メンバー1） | C列（メンバー2） | D列（メンバー3） | E列（メンバー4） | F列（メンバー5） |
|---|---|---|---|---|---|
| チームA | username1 | username2 | username3 | | |
| チームB | username4 | username5 | | | |

- ユーザー名はDiscordの `@` なしのユーザー名（小文字）で入力してください
- 1チーム1行、メンバーは最大5名まで

### DISCORD_BOT_TOKEN の取得

1. [Discord Developer Portal](https://discord.com/developers/applications) にアクセス
2. New Application からアプリケーションを作成
3. サイドバーの Bot タブからトークンを作成してコピー

### DISCORD_GUILD_ID の取得

1. Discordでサーバー名を右クリック
2. メニュー最下部の「サーバーIDをコピー」をクリック

### Google認証情報ファイルの生成

1. Google Cloudで [スプレッドシートAPI](https://console.cloud.google.com/apis/library/sheets.googleapis.com?hl=ja) を有効化
2. 認証情報 → 認証情報の作成 → サービスアカウントを選択
3. 作成したサービスアカウントを選択 → キー → 鍵を追加 → 新しい鍵 → JSON
4. ダウンロードしたJSONを `hackathon-util/` 直下に `credential.json` として保存
