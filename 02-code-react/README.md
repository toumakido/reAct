# 02-code-react: コード解析に特化したReAct実装

## 概要

Goのソースコードを読み込んで関数の実装を解析するReActエージェントです。
`data/`ディレクトリ内の複数のGoファイルを探索し、指定された関数の実装内容を見つけて説明します。

## 特徴

- **コード理解に特化**: 関数の実装内容を調べる質問に対応
- **複数ファイル対応**: 複数のGoファイルから必要な情報を検索
- **2つのアクション**:
  - `ListFiles`: ファイル一覧の表示
  - `ReadFile`: ファイル内容の読み込み

## データ構造

`data/`ディレクトリには以下のサンプルGoファイルが含まれています：

- **math.go**: 数学関数（Add, Subtract, Multiply, Divide, Factorial）
- **string.go**: 文字列関数（Reverse, ToUpperCase, ToLowerCase, IsPalindrome, CountVowels）
- **utils.go**: ユーティリティ関数（Max, Min, Abs, IsEven, IsOdd, Clamp）

## 実行方法

```bash
# プロジェクトルートから
go run ./02-code-react "Reverse関数の実装を教えてください"

# または、ディレクトリ内から
cd 02-code-react
go run . "Factorial関数はどのように実装されていますか？"
```

## 質問例

```bash
# 特定の関数の実装を調べる
go run . "Reverse関数の実装を教えてください"

# 複数の関数を比較
go run . "AddとSubtract関数の実装の違いを説明してください"

# アルゴリズムについて質問
go run . "Factorial関数はどのようなアルゴリズムで実装されていますか？"

# 使用している標準ライブラリを調べる
go run . "string.goではどの標準ライブラリを使っていますか？"
```

## ReActフロー

1. **ListFiles**: まずファイル一覧を取得
2. **ReadFile**: 関連するファイルを読み込み
3. **解析**: 関数の実装を見つけて理解
4. **Final Answer**: 実装コードと説明を返す

## 実装のポイント

### システムプロンプト
- コード解析に特化した指示
- 関数の実装と説明の両方を求める

### アクションハンドリング
```go
case "ListFiles":
    // ファイル一覧を返す
case "ReadFile":
    // ファイルの内容を返す
```

### パース処理
- 正規表現でAction/Action Inputを抽出
- ListFilesはAction Inputが不要

## 01-basic-reactとの違い

| 項目 | 01-basic-react | 02-code-react |
|------|----------------|---------------|
| 目的 | ファイル探索ゲーム | コード解析 |
| データ | テキストファイル | Goソースファイル |
| アクション | ReadFileのみ | ListFiles + ReadFile |
| 質問形式 | 物語的な質問 | 技術的な質問 |

## 拡張アイデア

- **GrepCode**: コード内のパターン検索
- **FindFunction**: 関数名から定義を直接検索
- **AnalyzeImports**: import文の解析
- **CountLines**: ファイルの行数やコメント率の計算
