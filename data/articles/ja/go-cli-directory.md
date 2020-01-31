go getすれば即コマンドとして使えるようにCLIツールを作る---2018-01-30 20:08:16

### はじめに
Goはソースコードをクロスコンパイルして、各プラットフォーム向けに配布することができる。
本エントリでは、Goの環境があるユーザに対して、goで作ったCLIツールを配布する方法を書いていく。

### 事前準備

ツールを使うユーザは、以下の環境変数を持っている必要がある。

```shell
PATH=$PATH:$GOPATH/bin
```

[公式](https://golang.org/doc/code.html#GOPATH)にもある通り。
作りたいツールのREADMEにでも書いておくのがいいと思っている。

### ツールの作り方

ディレクトリ構成を、こんな感じにしている。ここでは、 `hello` というコマンドを作るとする。

```go
.
├── cmd
│   └── hello
│       └── main.go
├── hello_lib.go
├── hello_lib_test.go
├── LICENSE
└── README.md
```

## main.go

main.goはコマンドの起点となる。これは、 `./cmd/コマンド名/main.go` という名前にする。
こんな感じで書いていく。

```go
package main

import (
  "fmt"
  "github.com/yagi5/hello_lib"
)

func main() {
  fmt.Println(hello_lib.World())
}
```

ポイントとしては、

* 依存ライブラリ(`hello_lib`)はgithub上のパスを書く

ことです。
次に、 `hello_lib` を書いていく。

## hello_lib

hello_libはmainから呼ばれるパッケージである。

```go
package hello_lib

func World() string {
  return "hello world"
}
```

### ダウンロード方法
この状態で、

```shell
$ hello
```

と叩いたら

```shell
hello world
```

と表示されるようにするには、ユーザーに以下のコマンドを叩いてもらえば良い。

```shell
$ go get github.com/ygnmhdtt/hello/cmd/hello
```
これで

```shell
$ hello
```

が叩けるようになります。

### なぜこれだけでいいのか？

`go get` は渡されたパスのソースをダウンロードして、ユーザのプラットフォームに合わせてビルドし、それを `$GOPATH/bin` に配置してくれる。
この時、main.goがあるディレクトリ名(上記の例では `hello` )のバイナリになる。
ユーザが `$GOPATH/bin` にPATHを通してくれていれば、 `hello` というバイナリを叩けるようになる。
また、 `main` で `hello_lib` というパッケージをimportしていますが、このような依存しているパッケージも勝手に落としてくれる。

### コマンドを増やすことができる

```shell
.
├── cmd
│   └── hello
│       └── main.go
├── hello_lib.go
├── hello_lib_test.go
├── LICENSE
└── README.md
```

当初の状態から

```shell
.
├── cmd
│   └── hello
│       └── main.go
│   └── dog
│       └── main.go
│   └── cat
│       └── main.go
├── hello_lib.go
├── hello_lib_test.go
├── LICENSE
└── README.md
```

こんな風に `cmd` 配下を増やして、

```shell
$ go get github.com/ygnmhdtt/hello/cmd/dog
```

```shell
$ go get github.com/ygnmhdtt/hello/cmd/cat
```

とすれば、hello_libを共通で使うようなコマンドを簡単に増やすこともできる。
そのため、 `hello_lib` にユーティリティ的な関数を定義しておいて、それらを小さなコマンドに分割して、パイプでつないで使うような、使い方ができる。

### ライブラリとしても使える

このやり方だと、 `hello_lib` をライブラリとしても提供できるようになる。
ユーザは `go get github.com/ygnmhdtt/hello_lib`するだけでいい。

### まとめ
このやり方を採用することで、

* `go get`すれば即コマンドとして使える
* 小さいコマンドを増やしていける
* ライブラリとしてGoのソースからも使える

のようなメリットが有る。
