# 01-basic-react: 基本的なReAct実装

正規表現でLLM出力をパースする最もシンプルなReAct実装です。

## 概要

この実装では、以下を学べます：
- ReActパターンの基本構造（Thought → Action → Observation）
- LLM出力の手動パース
- フレームワークなしでのエージェントループ実装

## 特徴

### ReActフォーマット

System Promptで以下のフォーマットを指示：

```
Thought: [次に何をすべきか考える]
Action: ReadFile
Action Input: [filename]
```

LLMからの応答を正規表現で解析し、`Action`と`Action Input`を抽出します。

### 宝探しシナリオ

5つの日本語ファイルが連鎖した宝探しデータ：

- **start.txt**: 冒険の開始点
- **library.txt**: 黄金の鍵が3パーツに分かれていることを説明
- **garden.txt**: パート1の場所 + 数字7
- **tower.txt**: パート2の場所 + 数字5
- **treasure.txt**: パート3の場所（7 × 5 + 3 = 部屋38）

最終回答には複数ファイルの情報統合と計算が必要です。

## 実行方法

```bash
# このディレクトリから実行
go run . "黄金の鍵はどこにありますか？"

# プロジェクトルートから実行
go run ./01-basic-react "質問"

# プロファイルを指定して実行
AWS_PROFILE=your-profile go run . "質問"
```

## 試すべき質問

### 基本的な質問
```bash
go run . "黄金の鍵はどこにありますか？"
```

### 複数ファイルの統合が必要
```bash
go run . "黄金の鍵の3つのパーツの場所と、それぞれの管理人や関連情報を教えてください"
```

### 計算要素を含む
```bash
go run . "treasure.txtのどの部屋にパート3がありますか？部屋番号の計算方法も説明してください"
```

### 歴史的情報
```bash
go run . "黄金の鍵は誰が、いつ作りましたか？"
```

## コードの構造

### main.go

```go
func main() {
    // 1. Bedrockクライアントを初期化
    client, err := bedrock.NewClient(ctx)

    // 2. ReActループを実行
    runReActLoop(ctx, client, question)
}

func runReActLoop() {
    for i := 0; i < maxIterations; i++ {
        // 3. LLMに送信
        response := client.InvokeModel(ctx, systemPrompt, messages)

        // 4. "Final Answer:"をチェック
        if strings.Contains(response, "Final Answer:") {
            return // 完了
        }

        // 5. Actionをパース
        action, input := parseAction(response)

        // 6. ツールを実行
        content := tools.ReadFile(input)

        // 7. Observationとして履歴に追加
        messages = append(messages, Message{
            Role: "user",
            Content: "Observation: " + content,
        })
    }
}
```

### アクションパース

```go
func parseAction(response string) (action, input string, found bool) {
    actionRegex := regexp.MustCompile(`(?i)Action:\s*(\w+)`)
    inputRegex := regexp.MustCompile(`(?i)Action Input:\s*(.+?)(?:\n|$)`)

    // 正規表現でAction と Action Inputを抽出
}
```

## 出力例

```
=== Starting ReAct Agent ===
Question: 黄金の鍵はどこにありますか？

--- Iteration 1 ---
Thought: まずstart.txtを読んで冒険を始める必要があります。
Action: ReadFile
Action Input: start.txt

Observation: 宝探しへようこそ！
あなたは古代の城の入り口に立っています...

--- Iteration 2 ---
Thought: 次にlibrary.txtを読んで、鍵の詳細情報を確認します。
Action: ReadFile
Action Input: library.txt

Observation: 王立図書館
...

[続く...]

Final Answer: 黄金の鍵は3つのパーツに分かれています：
1. パート1: 秘密の庭園の不死鳥の像の下
2. パート2: 北の塔の宝箱の中
3. パート3: 宝物庫の部屋38

=== Agent Complete ===
```

## 学習ポイント

1. **ReActパターンの理解**
   - Thought: AIの推論プロセス
   - Action: 実行するアクション
   - Observation: アクション結果

2. **LLM出力のパース**
   - 正規表現を使った構造化データの抽出
   - エラーハンドリング（Actionが見つからない場合）

3. **会話履歴の管理**
   - `[]Message`での状態管理
   - AssistantとUserロールの切り替え

4. **終了条件**
   - `Final Answer:`の検出
   - 最大反復回数の設定

## 制限事項

- ツールは`ReadFile`のみ
- 正規表現パースのため、フォーマット崩れに弱い
- エラーリカバリーは基本的なもののみ

## 次のステップ

この実装を理解したら、以下を試してください：

- **02-converse-api**: Converse APIでの実装
- **03-structured-output**: JSON出力による堅牢なパース
- **04-multi-tool**: 複数ツールの実装
