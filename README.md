# ReAct Agent in Go

ReActアーキテクチャを学ぶためのGo実装サンプル集です。
AWS Bedrockを使用し、さまざまなReAct実装パターンを試せる構成になっています。

## プロジェクト構造

```
reAct/
├── 01-basic-react/          # 基本的なReAct実装（正規表現パース）
│   ├── main.go
│   ├── data/                # この実験専用データ
│   └── README.md
│
├── lib/                     # 共通ライブラリ
│   ├── bedrock/             # Bedrock API クライアント
│   │   └── client.go
│   ├── tools/               # 共通ツール
│   │   └── readfile.go
│   └── types/               # 共通型定義
│       └── types.go
│
├── go.mod
└── README.md                # このファイル
```

## セットアップ

### 前提条件

- Go 1.21以上
- AWS認証情報の設定（Bedrockへのアクセス権限が必要）

### AWS認証情報の設定方法

**デフォルトプロファイルを使う場合:**
- `~/.aws/credentials`の`[default]`が自動的に使用される
- 特に指定不要

**特定のプロファイルを使う場合:**
```bash
# 環境変数で指定
export AWS_PROFILE=your-profile-name

# または実行時に指定
AWS_PROFILE=your-profile-name go run ./01-basic-react "質問"
```

**環境変数で直接指定する場合:**
```bash
export AWS_ACCESS_KEY_ID=your-access-key
export AWS_SECRET_ACCESS_KEY=your-secret-key
export AWS_REGION=us-east-1
```

### インストール

```bash
# 依存関係のインストール
go mod download
```

## 各実装の実行方法

### 01-basic-react: 基本的なReAct実装

```bash
# 実行
go run ./01-basic-react "黄金の鍵はどこにありますか？"

# または、ディレクトリに移動して実行
cd 01-basic-react
go run . "黄金の鍵の3つのパーツの場所を教えてください"
```

詳細は各ディレクトリの`README.md`を参照してください。

## 実装パターン

### 01-basic-react
- **特徴**: 正規表現でLLM出力をパースする基本実装
- **目的**: ReActの仕組みを理解する
- **パターン**: Thought → Action → Observation のループ
- **学習ポイント**: フレームワークなしでの生のLLM呼び出し

### 今後追加予定の実装例

- **02-converse-api**: Converse API + Tool Use実装
- **03-structured-output**: 構造化出力を使った実装
- **04-multi-tool**: 複数ツールを持つエージェント
- **05-streaming**: ストリーミングレスポンス対応版

## 共通ライブラリ

### `lib/bedrock`
- AWS Bedrock RuntimeのクライアントWrapper
- `NewClient()`: Bedrockクライアントの初期化
- `InvokeModel()`: Claude APIの呼び出し

### `lib/tools`
- エージェントが使用するツール群
- `ReadFile()`: ファイル読み込みツール

### `lib/types`
- 共通型定義
- `Message`: LLMとのメッセージ型

## 新しい実装の追加方法

1. 新しいディレクトリを作成
```bash
mkdir 02-your-experiment
cd 02-your-experiment
```

2. `main.go`を作成し、`lib/`の共通ライブラリをimport
```go
import (
    "github.com/toumakido/reAct/lib/bedrock"
    "github.com/toumakido/reAct/lib/tools"
    "github.com/toumakido/reAct/lib/types"
)
```

3. 必要に応じて`data/`ディレクトリを作成

4. `README.md`でその実装の特徴を説明

## トラブルシューティング

### AWS認証エラー

```
Failed to create Bedrock client: failed to load AWS config
```

**原因と対処法:**
- AWS認証情報が設定されていない
  - `~/.aws/credentials`を確認
  - または環境変数`AWS_PROFILE`を設定
- プロファイルが存在しない
  - `aws configure list-profiles`でプロファイル一覧を確認
  - 正しいプロファイル名を指定

### Bedrockアクセスエラー

```
failed to invoke model: operation error
```

**原因と対処法:**
- Bedrockが利用可能なリージョンではない
  - Bedrockは特定リージョンのみで利用可能（us-east-1, us-west-2など）
  - `~/.aws/config`または環境変数`AWS_REGION`でリージョンを設定
  ```bash
  export AWS_REGION=us-east-1
  ```
- Claudeモデルへのアクセス権限がない
  - AWS Bedrockコンソールでモデルアクセスを有効化
  - IAMポリシーでbedrockの権限を確認

## ライセンス

MIT
