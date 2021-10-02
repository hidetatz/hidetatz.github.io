Goで使いやすいPackageを作るために知っておきたいこと---2018-11-09 22:57:53

Goのコードをいろいろと読んできて感じた、Goらしいパッケージのデザインを書いていく。

### パッケージ

Goのパッケージは、ディレクトリ構造と一致する。
ひとつのディレクトリには同じパッケージのファイルだけがある必要がある。
また、パッケージ内には複数のGoファイルがあり、それらは `適切に` 分けられているべきだ。

例えば、[strings](https://github.com/golang/go/tree/master/src/strings)パッケージを見てみよう。
stringsを操作するドメインごとにファイルが分かれている。
とはいえ、Goではパッケージ = 名前空間なので、そんなに気にしすぎる必要はない。

```
builder.go
builder_test.go
compare.go
compare_test.go
example_test.go
export_test.go
reader.go
reader_test.go
replace.go
replace_test.go
search.go
search_test.go
strings.go
strings.s
strings_test.go
```

### doc.go

doc.goという文化がある。コメントでパッケージの説明を書き、最後に `package xxx` だけ書いておく、というものだ。

[go/src/cmd/asm/doc.go](https://github.com/golang/go/blob/master/src/cmd/asm/doc.go)では、そのパッケージによって提供されるコマンドのusageを書いている。

[go/src/runtime/race/doc.go](https://github.com/golang/go/blob/master/src/runtime/race/doc.go)では、単にパッケージの説明が書いてある。

[go/src/net/http/doc.go](https://github.com/golang/go/blob/master/src/net/http/doc.go)では、パッケージのサンプルコードを書いている。

書く内容は上記のように多岐にわたるのだが、ある程度規模の大きいパッケージを作るときには覚えておきたいテクニックだ。

### 型とファイルはなるべく一致させる

goでは、例えば[http.Header](https://github.com/golang/go/blob/master/src/net/http/header.go) は `net/http/header.go` にあることが一般的だ。

```go
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"io"
	"net/http/httptrace"
	"net/textproto"
	"sort"
	"strings"
	"sync"
	"time"
)

// A Header represents the key-value pairs in an HTTP header.
type Header map[string][]string

// Add adds the key, value pair to the header.
// It appends to any existing values associated with key.
// The key is case insensitive; it is canonicalized by
// textproto.CanonicalMIMEHeaderKey.
func (h Header) Add(key, value string) {
	textproto.MIMEHeader(h).Add(key, value)
}
```

このように、ファイル名と一致するstructがあり、そのメソッドが定義されている、というのは一般的だ。

### パッケージは責務ごとに切る

Ruby on Railsや他のMVCフレームワークのように、 `model` や `controller` のようなパッケージはGoでは切らない。
そうではなく、機能に紐付いた責務でパッケージを切る。

例えば、userテーブルに紐づくuser構造体のあるべき場所は、 `models` パッケージではなく、 `user` 構造体を使用する
ドメインのパッケージである。
これがうまく決められない場合は、大きなひとつのmainパッケージにしても構わない。

### godocを書く

すべてのGoプログラマは[godoc](https://godoc.org/)についてよく知っているべきだ。
godocはわざわざそのために書いているわけではなく、コードの中のコメントを勝手に読んでくれる。
godocフレンドリーなコメントを常に書くよう心がける。

### Exampleを書く

パッケージの関数の使い方を説明するために、godocに加え `Example` という手法がある。
[strings.Replace](https://golang.org/pkg/strings/#Replace)を見てみよう。

![](/images/go-package-style/1.png)

このように、 `Example` が表示されている。
これはどのように生成されているかというと、 [example_test.go](https://golang.org/src/strings/example_test.go#L200)にある。

```go
func ExampleReplace() {
	fmt.Println(strings.Replace("oink oink oink", "k", "ky", 2))
	fmt.Println(strings.Replace("oink oink oink", "oink", "moo", -1))
	// Output:
	// oinky oinky oink
	// moo moo moo
}
```

godocのExampleは、

* xxx_test.go に定義している(つまりgoからはテストとみなされ、普通にgo buildしてもビルド対象にはならない)
* ExampleXxx() という名前の関数である

ときに、生成される。
strings.Replaceのケースでは `example_test.go` にまとめてstringsパッケージのexampleが定義してあるが、
`example_test` にしなくても動作する。

また、上記のように `Output: xxx` と書くことで、これ自体がテストの役割を果たす。
goから見るとExampleであってもtestとみなされ、 `go test` 時には `ExampleXxx` 関数を実行し、結果が `Output: xxx` と一致しているかを
勝手にチェックしてくれるのだ。(すごい)

また、関数名を工夫することで、単なる関数以外のテストもできるようになっている。
詳しくは[Testable Examples in Go](https://blog.golang.org/examples) を参照する。
