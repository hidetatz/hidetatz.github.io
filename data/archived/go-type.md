Goの型システムとGoらしい書き方---2018-10-07 11:35:53

Goにはクラスがない。しかし、データ構造をまとめるためのstructがあり、structにはメソッドを紐付けることができる。
そのため、プログラマによってはオブジェクト指向言語のようにGoを書くこともある。
しかし、Goの型に継承関係という概念は無く、RubyやPythonのような考え方でコードを書く訳にはいかないだろう。

Goの作者[Rob Pike](https://twitter.com/rob_pike/status/942528032887029760)は以下のように述べている。

> … the more important idea is the separation of concept: data and behavior are two distinct concepts in Go, not conflated into a single notion of “class”.

Goにはclassは存在しない。classとはデータとふるまいをまとめあげるものである。
Goでは、データと振る舞いは分離している。
structはデータを提供するための軽量な手段であり、それ以上のことはしない。
structは型の階層を表現することは決して無い。

### interfaceの埋め込み

interfaceの埋め込みは継承ではない。
Goでは、コードの再利用は継承ではなくコンポジションで実現される。
継承はしばしば、プログラムの構造を複雑にさせ、メンテナンス性の低下をもたらす。
Goは継承の代わりに、**コンポジション**と、 **interfaceを使ったメソッドのdispatch**を提供している。

[以下のような](https://play.golang.org/p/Pj9VZTUmbE5)コードを見てみよう。

```go
package main

import (
  "io"
  "sync"
)

type File struct {
  sync.Mutex
  rw io.ReadWriter
}

func main() {
  f := File{}
  defer f.Unlock()
  f.Lock()
}

```

File構造体は[sync.Mutex](https://golang.org/src/sync/mutex.go)構造体を埋め込まれている。Mutex構造体は `Lock()` `Unlock()` メソッドを実装している。
これにより、File型の変数fも、 `Lock()` `Unlock()` というメソッドを呼び出すことができる。
これは、サブクラスではなく、コンポジションである。

### ポリモーフィズム

サブクラスのないGoでは、ポリモーフィズムをinterfaceによってのみ実現する。[以下のような](https://play.golang.org/p/a0nYrIniSLx)コードを見てみよう。

```go
package main

import (
  "bytes"
  "io"
  "log"
)

func main() {
  var r io.Reader

  r = bytes.NewBufferString("hello")

  buf := make([]byte, 2048)
  if _, err := r.Read(buf); err != nil {
    log.Fatal(err)
  }
}
```

変数rは `io.Reader` として宣言されているが、 [bytes.NewBufferString()](https://golang.org/pkg/bytes/#NewBufferString)は
[*bytes.Buffer](https://golang.org/pkg/bytes/#Buffer) を返す。
`bytes.Buffer` は [io.Reader](https://golang.org/pkg/io/#Reader)を実装しているため、このような変数の代入が可能になる。
つまり、 `r.Read` は [(*Buffer).Read](https://golang.org/pkg/bytes/#Buffer.Read) にディスパッチされる。

### interfaceの実装を宣言しない

Goには `implements` というキーワードは存在しない。
interfaceを実装しているかどうかは、interfaceに宣言されたメソッドリストをそのstructが実装しているかどうかのみで判断される。

Goには、新たなinterfaceを作るのではなく、コミュニティーや標準ライブラリから提供されているものを使うことを推奨する文化が有る。
これにより、似たようなinterfaceが増えることを抑制している。

意図せずinterfaceを実装してしまい困るようなことはあるだろうか？
100%ないとは言えないが、そのようなケースに出くわしたことはない。
また、もし2つの似たinterfaceが存在している場合、片方は不要であることが多い。

interfaceは小さく宣言し、必要に応じて組み合わせて使うのが望ましい。

### コンストラクタは作らない

Goでstcurtを初期化する際、値はすべて[ゼロ値](https://tour.golang.org/basics/12)で初期化される。
Goは、なるべくゼロ値で初期化しても動くようなstructを作る、という考え方がある。言うまでもなく、 `NewXxx()` という関数を知らなくても動かせるようにだ。

例えば、 [http.Client](https://golang.org/pkg/net/http/?#Client)は以下のようなコードで動かすことができる。

```go
client := http.Client{}
```

ここから必要に応じて、

```go
client.Timeout = 5 * time.Second
```

のように設定するのが良い。

それでも `NewXxx` が必要になるケースは、バリデーションや、コネクションの確立など、なにか手続きが必要なケースだ。
例えば、[http.NewRequest()](https://golang.org/pkg/net/http/?#NewRequest)は以下のように使用する。

```go
req, err := http.NewRequest("GET", "https://example.com", nil)
```

[実装](https://golang.org/src/net/http/request.go?s=26446:26515#L782)を見てみよう。
メソッドとURLをバリデーションし、bodyをReadしているのがわかると思う。

### まとめ
Goの型システムは覚えることが少ない一方、柔軟で、あたかもオブジェクト指向のように使用することも可能では有る。
しかし、Goの思想やエコシステムを理解し、Goらしいコードを書くのは重要なことである。
これらについて学ぶためには、なんといっても[標準パッケージ](https://github.com/golang/go)のコードを読むのが良い。今後も goらしいコードを書けるよう精進していく。
