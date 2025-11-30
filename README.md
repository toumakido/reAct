# ReAct Agent in Go

ReActアーキテクチャを学ぶためのシンプルなGo実装サンプルです。
AWS Bedrock (Claude 3.5 Sonnet) を使用し、ローカルファイルを読み込みながら推論を重ねて質問に答えるエージェントです。

## 特徴

- **フレームワーク不使用**: LangChainなどのフレームワークを使わず、生のLLM呼び出しとループ処理で実装
- **ReActパターン**: Thought → Action → Observation のサイクルを可視化
- **宝探しシナリオ**: 5つの連鎖したファイルから情報を収集し、統合して回答

## 前提条件

- Go 1.21以上
- AWS認証情報の設定（Bedrockへのアクセス権限が必要）
  - `~/.aws/credentials` または環境変数 `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`
  - Bedrock Claudeモデルへのアクセス権限

## セットアップ

```bash
# 依存関係のインストール
go mod download

# ビルド（オプション）
go build -o react-agent
```

## 使い方

```bash
# 直接実行
go run . "あなたの質問をここに入力"

# またはビルド後に実行
./react-agent "あなたの質問をここに入力"
```

## 試すべき質問の例

### 例1: 基本的な情報収集
```bash
go run . "Golden Keyはどこにありますか？"
```

### 例2: 複数ファイルの情報統合が必要
```bash
go run . "Golden Keyの3つのパーツはそれぞれどこにありますか？各パーツの場所と守護者や関連情報を教えてください。"
```

### 例3: 計算が必要
```bash
go run . "treasure.txtのどの部屋にPart 3がありますか？その部屋番号はどうやって計算しますか？"
```

### 例4: 歴史的情報
```bash
go run . "Golden Keyは誰が、いつ作りましたか？"
```

### 例5: 包括的な質問
```bash
go run . "宝探しの完全なストーリーを教えてください。誰が関わっていて、どの場所を訪れる必要がありますか？"
```

## プロジェクト構造

```
.
├── main.go           # ReActメインループ
├── bedrock.go        # AWS Bedrock クライアント
├── tools.go          # ReadFile ツール実装
├── go.mod
├── go.sum
├── README.md
└── data/             # 探索対象データ
    ├── start.txt     # 開始地点
    ├── library.txt   # 図書館（歴史情報）
    ├── garden.txt    # 秘密の庭園（Part 1）
    ├── tower.txt     # 北の塔（Part 2）
    └── treasure.txt  # 宝物庫（Part 3）
```

## コードの仕組み

### 1. ReActループ (`main.go`)

```
ユーザーの質問
    ↓
┌─────────────────────────┐
│  LLMに送信              │
│  (System Prompt含む)     │
└─────────────────────────┘
    ↓
┌─────────────────────────┐
│  応答を解析              │
│  - Thought を表示       │
│  - Action を検出        │
└─────────────────────────┘
    ↓
┌─────────────────────────┐
│  Action実行             │
│  (ReadFile)             │
└─────────────────────────┘
    ↓
┌─────────────────────────┐
│  Observation として     │
│  履歴に追加              │
└─────────────────────────┘
    ↓
  Final Answer? ─No→ ループ継続
    ↓ Yes
  完了
```

### 2. System Prompt

LLMに以下のフォーマットを厳密に守るよう指示：

```
Thought: [次にすべきことの推論]
Action: ReadFile
Action Input: [filename]
```

### 3. アクションの解析

正規表現を使って応答から `Action:` と `Action Input:` を抽出：

```go
actionRegex := regexp.MustCompile(`(?i)Action:\s*(\w+)`)
actionInputRegex := regexp.MustCompile(`(?i)Action Input:\s*(.+?)(?:\n|$)`)
```

## 実行例の出力

```
=== Starting ReAct Agent ===
Question: Golden Keyの3つのパーツの場所を教えてください

--- Iteration 1 ---
Thought: まずstart.txtから開始して、宝探しの情報を収集する必要があります。
Action: ReadFile
Action Input: start.txt

Observation: Welcome to the Treasure Hunt!
...

--- Iteration 2 ---
Thought: 次にlibrary.txtを読んで、Golden Keyの詳細情報を確認します。
Action: ReadFile
Action Input: library.txt

Observation: The Royal Library
...

[続く...]

Final Answer: Golden Keyは3つのパーツに分かれています：
1. Part 1: 秘密の庭園のPhoenix像の下（守護者：Isabella）
2. Part 2: 北の塔の鍵付き箱の中（建築家の名前で開錠）
3. Part 3: 宝物庫のRoom 38（7×5+3=38で計算）

=== Agent Complete ===
```

## カスタマイズ

### データファイルの変更

`data/` ディレクトリ内のファイルを編集することで、独自のシナリオを作成できます。

### ツールの追加

`tools.go` に新しい関数を追加し、`main.go` の `parseAction` 関数で処理を追加することで、新しいアクションを実装できます。

### モデルの変更

`bedrock.go` の `modelID` を変更することで、異なるClaudeモデルを使用できます：

```go
modelID: "anthropic.claude-3-sonnet-20240229-v1:0"  // Claude 3 Sonnet
```

## トラブルシューティング

### AWS認証エラー

```
Failed to create Bedrock client: failed to load AWS config
```

→ AWS認証情報が正しく設定されているか確認してください。

### Bedrockアクセスエラー

```
failed to invoke model: operation error
```

→ AWSアカウントでBedrockのClaudeモデルへのアクセスが有効になっているか確認してください。

### 最大反復回数到達

```
max iterations (15) reached without final answer
```

→ `main.go` の `maxIterations` を増やすか、質問をより明確にしてください。

## 学習ポイント

このプロジェクトから学べること：

1. **ReActパターンの理解**: Thought → Action → Observation のループ
2. **LLM統合**: フレームワークなしでのLLM呼び出し
3. **プロンプトエンジニアリング**: System Promptでの振る舞い制御
4. **パース処理**: LLM出力からの構造化データ抽出
5. **エラーハンドリング**: ファイルI/OとAPI呼び出しの堅牢な処理

## ライセンス

MIT
