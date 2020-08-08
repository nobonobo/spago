# Getting Started Guide

## コマンドラインツールのインストール

```shell
go get -u github.com/nobonobo/spago/cmd/spago
```

## プロジェクトの作成

```shell
mkdir sample1
cd sample1
git mod init sample1
```

"main.go"ファイルの内容を作成します。
まずは以下の内容にしてください。

```go
package main

import (
	"github.com/nobonobo/spago"
)

func main() {
	spago.RenderBody(&Top{})
	select {}
}
```

## コンポーネントの作成

spago コマンドを使って、`Top`コンポーネントの雛形を作成します。

```shell
spago new -p main Top
```

フォルダ内は以下のようなファイル一覧になります。

- sample1/
  - main.go
  - top.go
  - top.html
  - go.mod

## HTML の記述

body タグ以下の記述を
top.html に記述します。

```html
<body>
  <h5>Hello World!</h5>
</body>
```

## 開発サーバーの起動

```shell
spago server
```

## 最初の表示を確認

以下の URL をブラウザで開きます。

open http://localhost:8080

`top_gen.go`が生成され、画面に top.html 相当の表示が出ます。

top_gen.go の内容

```go
package main

import (
	"github.com/nobonobo/spago"
)

// Render ...
func (c *Top) Render() spago.HTML {
	return spago.Tag("body",
		spago.Tag("h5",
			spago.T("Hello World!"),
		),
	)
}
```

## HTML のアップデート

このまま、top.html の内容を編集し、リロードすれば反映されるのを確認してみましょう。
h5 タグを button タグに書き換えます。

```html
<body>
  <button>Hello World!</button>
</body>
```

リロードすると top_gen.go の内容は以下のように変更され、ブラウザの表示もボタンになっています。

```go
package main

import (
	"github.com/nobonobo/spago"
)

// Render ...
func (c *Top) Render() spago.HTML {
	return spago.Tag("body",
		spago.Tag("button",
			spago.T("Hello World!"),
		),
	)
}
```

## クリックイベントのハンドリング

次にボタンクリックイベントを実装してみます。
top.go の後半に`func (c *Top) OnClick(ev js.Value) interface{}`メソッドを追記します。

```go
package main

import (
	"syscall/js"

	"github.com/nobonobo/spago"
)

//go:generate spago generate -c Top -p main top.html

// Top  ...
type Top struct {
	spago.Core
}

func (c *Top) OnClick(ev js.Value) interface{} {
    js.Global().Call("alert", "button clicked!")
    return nil
}
```

## HTML でのイベントマッピング

top.html を以下のように修正します。

```html
<body>
  <button @click="{{c.OnClick}}">Hello World!</button>
</body>
```

リロード後、ボタンクリックでアラートが表示されることが確認できるはずです。

top_gen.go の内容をみてみましょう。

```go
package main

import (
	"github.com/nobonobo/spago"
)

// Render ...
func (c *Top) Render() spago.HTML {
	return spago.Tag("body",
		spago.Tag("button",
			spago.Event("click", c.OnClick),
			spago.T("Hello World!"),
		),
	)
}
```

spago.Event の記述がイベントマッピングの記述です。

## プロパティと再描画

次は DOM アップデートを試します。

- Top 構造体に「Counter int」プロパティを追加。
- クリックイベントは以下のように変更します。

```go
func (c *Top) OnClick(ev js.Value) interface{} {
	c.Count++
	spago.Rerender(c)
	return nil
}
```

top.html は以下のように変更します

```html
<body>
  <button @click="{{c.OnClick}}">{{spago.T(c.Count)}}</button>
</body>
```

`spago.T(...)`は引数を文字列化しテキストノードにするマークアップです。

- クリック->Count プロパティインクリメント
- `spago.Rerender(c)`により Top コンポーネントが再描画
- 結果、ボタンの表記の数値が変更される

DOM 更新には差分抽出してから適用されますので、ボタン表記だけの差分であればボタン表記だけが変更されます。
（この場合、ボタンに乗ったフォーカスが失われることはありません。）

## HTML ライク DSL とリロード時の出来事

html ファイルの内容は簡易な DSL になっていて、
「@イベント名=」属性でイベントリスナーの指定ができます。
多くの属性値や子ノードにて`{{...}}`を記入する時、`...`部分が Go のコードとして展開されます。
`c`は慣習的にコンポーネントのプレースフォルダとして扱っています。
この DSL はあまり親切には作っていません。基本はコードによる記述を理解した上で変動部分を記述し、
固定部分は HTML そのままで記述するだけです。

また、WASM ファイルのリクエストが来た時、`go generate ./...`が実行されます。
この際、Go ソースコード中のコメントに以下のような記述があれば、そのソースのあるフォルダで
コメントにあるコマンドが実行される仕組みがあります。

```go
//go:generate spago generate -c Top -p main top.html
```

## spago generate について

```shell
spago generate -c Top -p main top.html
```

このコマンドの意味は以下の内容です。

- top.html を Go の spago を使った記述のコードを生成します。
- コンポーネント名は「Top」で
- パッケージ名は「main」で
- top_gen.go を出力

パッケージ名を省略すると「main」扱いになります。コンポーネント名は必須パラメータです。

## spago server の仕組み

### 基本の挙動

- main.wasm リソースを要求されたら？
  - go generate ./...
  - go build -o main.wasm .
  - main.wasm を gzip に圧縮してサーブ
- index.html リソースを要求されたら？
  - index.html があればそれをサーブ
  - なければ最小限の index.html をサーブ
- wasm_exec.js リソースを要求されたら？
  - wasm_exec.js があればそれをサーブ
  - 「go env GOROOT」配下にある wasm_exec.js をサーブ
- 上記以外のリソースを要求されたら？
  - 該当ファイルがあればそれをサーブ
  - なければ 404 ステータスをサーブ

### 上級者向けの機能

#### リバースプロキシ機能

バックエンドとの連携を仮に作り込むときに使います。

- 指定パスに対して別の Web サーバーへプロキシします
- 複数ルール記述できます
- 基本の挙動よりもこのプロキシルール指定のほうが優先されます

#### TinyGo バージョン

`spago server -tinygo`にて wasm ビルドに TinyGo を使います

## spago deploy の仕組み

`spago deploy -tinygo dist`とすると

- TinyGo でビルドした main.wasm
- index.html
- wasm_exec.js

以上を dist フォルダに出力します。
その他の参照リソースファイルはユーザーの手で dist フォルダにコピーしてください。

## TinyGo を使う時のメリット・デメリット

メリット

- 出力 WASM ファイルのサイズが本家版の 20〜40%以下
- LLVM 最適化の恩恵もある（末尾再帰最適化などが効く？）
- TinyGo でビルド可能なコードならおおむね本家でビルド可能（例外は TinyGo 専用ライブラリを利用など）

デメリット

- 別途 TinyGo のインストールが必要
- ビルドに時間がかかる
- reflect が部分的に対応していないためライブラリの一部が使えない
- goroutine は LLVM のコルーチンに置き換えられており並列での効率化はサポートされない