Golangにおけるnilの扱い方---2018-10-28 18:32:50

goにはnilが存在する。
他の多くの言語と違って、goのnilは型情報を保持している。
[以下のような](https://play.golang.org/p/mG06lmbu-Ud)コードを見てみよう。

```go
package main

type s1 struct{}

type s2 struct{}

func main() {
	var x *s1 = nil
	var y *s2 = nil

	compare(x, y)
}

func compare(x, y interface{}) {
	println(x == y)
}
```

このコードを実行すると、 `false` が出力される。
これは、nilがそれぞれ型情報を保持しているからである。

[以下のように](https://play.golang.org/p/PKaxqlYZYRJ)修正してやると、trueを出力する。

```go
package main

type s1 struct{}

type s2 struct{}

func main() {
	var x *s1 = nil
	var y *s1 = nil // 修正

	compare(x, y)
	return
}

func compare(x, y interface{}) {
	println(x == y)
}
```

[試しに](https://play.golang.org/p/rb_aEcwibhM)、compare関数を外してみよう。

```go
package main

type s1 struct{}

type s2 struct{}

func main() {
	var x *s1 = nil
	var y *s2 = nil
	println(x == y)
}
```

前述の通り、falseが出力される。と思いきや、そうではない。
これは実は、コンパイルエラーになる。

```
prog.go:10:12: invalid operation: x == y (mismatched types *s1 and *s2)
```

では、nilチェックをどのように行えばよいか？
実は、[nilはキャストできる。](https://play.golang.org/p/UGXSzHFenjz)

```go
package main

type s1 struct{}

func main() {
	var x *s1 = nil

	isnil(x)
	return
}

func isnil(x interface{}) {
	println(x == nil)
	println(x == (*s1)(nil))
}
```

これを実行すると、

```
false
true
```

と出力される。

これらは、関数の仮引数の型をinterface{}にしていることによる。
[以下のように](https://play.golang.org/p/ZsNHwWkbuWC)修正すると、意図通りに動くだろう。

```go
package main

type s1 struct{}

func main() {
	var x *s1 = nil

	isnil(x)
	return
}

func isnil(x *s1) { // 仮引数の型を修正
	println(x == nil)
	println(x == (*s1)(nil))
}
```

現実的に、nilチェックを行うべきときに常にnilの型情報やキャストを意識しなければいけないのは非常に煩雑である。コードの独立性を下げることにもつながるだろう。
これを避けるには、明示的に型情報を持たないnilリテラルを返すのが望ましい。

つまり、

```go
package main

import "fmt"

type s1 struct{}

func newS1() interface{} {
	s := new(s1)
	s = nil
	return s
}

func main() {
	s := newS1()
	fmt.Println(s == nil)
}
```

こうではなく、

```go
package main

import "fmt"

type s1 struct{}

func newS1() interface{} {
	_ = new(s1)
	return nil
}

func main() {
	s := newS1()
	fmt.Println(s == nil)
}
```

このように、常に `return nil` という書き方をするようにすると良い。
