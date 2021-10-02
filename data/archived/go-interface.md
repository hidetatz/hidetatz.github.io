Goにおけるインタフェースの実装について---2018-10-28 22:01:45

Goではinterfaceに定義されているメソッドリストのメソッドを実装することで
interfaceをimplementしているとみなされる。
[以下のコード](https://play.golang.org/p/hVDGz3HJO_A)を見てみよう。

```go
package main

type Shower interface {
	Show()
}

type S1 struct {
}

func (this *S1) Show() {
	println("Show from S1")
}

func InvokeShow(s Shower) {
	if s1, ok := s.(S1); ok {
		s1.Show()
	}
}

func main() {
	var s1 S1
	InvokeShow(s1)
}
```

このコードを実行すると、 `Show from S1` と出力されるだろうか？
実行してみるとわかるが、実際はコンパイルエラーになる。
どこがエラーになるだろうか。
実際は、以下のようなエラーメッセージが表示されるだろう。

```
prog.go:15:16: impossible type assertion:
	S1 does not implement Shower (Show method has pointer receiver)
prog.go:22:12: cannot use s1 (type S1) as type Shower in argument to InvokeShow:
	S1 does not implement Shower (Show method has pointer receiver)
```

2行のエラーが出るが、問題点はひとつだ。
`Showメソッドのレシーバはポインタである` ということである。

該当の部分は以下である。

```go
func (this *S1) Show() {
	println("Show from S1")
}
```

Showメソッドのレシーバは `*S1` である。 `S1` ではない。
Showerインタフェースを実装しているのは `*S1` であり、 `S1` ではないのである。

しかしながら、このコードはこれを満たしていない。
エラーメッセージをひとつずつ見ていこう。

```go
func InvokeShow(s Shower) {
	if s1, ok := s.(S1); ok {
		s1.Show()
	}
}
```

仮引数の `s` は、Showerインタフェース型である。
Showerインタフェースを実装しているのは `*S1` 型なので、
Showerインタフェース型の変数を `S1` にキャストすることはできない。
正しいコードは

```go
func InvokeShow(s Shower) {
	if s1, ok := s.(*S1); ok {
		s1.Show()
	}
}
```

になる。

```
func main() {
	var s1 S1
	InvokeShow(s1)
}
```

こちらは、InvokeShowメソッドにS1型の変数s1を渡している。
前述の通り、Showerインタフェースを実装しているのは `*S1` なため、
s1は `*S1` 型でないと渡せない。
正しいコードは以下のようになる。

```go
func main() {
	var s1 *S1
	InvokeShow(s1)
}
```

これらを踏まえた、動作するコードは以下のようになる。

```go
package main

type Shower interface {
	Show()
}

type S1 struct {
}

func (this *S1) Show() {
	println("Show from S1")
}

func InvokeShow(s Shower) {
	if s1, ok := s.(*S1); ok {
		s1.Show()
	}
}

func main() {
	var s1 *S1
	InvokeShow(s1)
}
```

さて、ここまで書けばわかるが、Showメソッドを実装するのを
ポインタでなくポインタのデリファレンスにすることでも動作させることができる。
以下のようなコードだ。

```go
package main

type Shower interface {
	Show()
}

type S1 struct {
}

func (this S1) Show() {
	println("Show from S1")
}

func InvokeShow(s Shower) {
	if s1, ok := s.(S1); ok {
		s1.Show()
	}
}

func main() {
	var s1 S1
	InvokeShow(s1)
}
```

これでも当然動く。
しかしながら、特別理由のない限り、ポインタレシーバにすべきだ。
理由はいくつかあるが、また別の機会に書こうと思う。
