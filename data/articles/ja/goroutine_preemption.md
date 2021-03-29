Goroutineのプリエンプション---2021-03-28 09:00:00

Goにおけるgoroutineのプリエンプションについて調べていたのでメモ。間違いがあれば指摘いただけると助かります。

Goにおけるプリエンプションは、Go1.14以前とそれ以降で挙動が異なる。Go1.14では、[リリースノート](https://golang.org/doc/go1.14#runtime) にもある通り、goroutineは "asynchronously preemptible" になった。これは何を意味するのか？

まず、簡単な例を見てみよう。
次のようなGoプログラムを考える。

```
package main

import (
	"fmt"
)

func main() {
	go fmt.Println("hi")
	for {
	}
}

```

main関数の中では "hi" と出力するだけのgoroutineをひとつ起動している。また、`for {}` で無限ループしている。

このプログラムを `GOMAXPROCS=1` にして動かすとどうなるだろうか？感覚的には "hi" と出力され、その後何も起きない (無限ループがあるため) というような挙動をしそうだ。
実際、このプログラムをGo1.14以降で (筆者が手元で動かしたときは Go1.16 (on Ubuntu on WSL2)で) 動かすと、その通り動く。
このプログラムをその通り動かさないためには方法はふたつある。ひとつは1.14より前のバージョンのGoで実行すること。もうひとつは `GODEBUG=asyncpreemptoff=1` にして動かすことだ。

筆者の手元で試してみると、以下のように動いた。

```
$ GOMAXPROCS=1 GODEBUG=asyncpreemptoff=1 go run main.go
# ここで止まる
```

"hi"が出ない。なぜこうなるのか？を話す前に、このプログラムを期待通りの挙動にさせる方法もいくつかあるので説明しておく。

ひとつは、ループの中に次のように処理を追加するやり方だ。

```
*** main.go.org	2021-03-26 20:03:16.840000000 +0900
--- main2.go	2021-03-26 20:03:58.970000000 +0900
*************** package main
*** 2,11 ****
--- 2,13 ----
  
  import (
  	"fmt"
+ 	"runtime"
  )
  
  func main() {
  	go fmt.Println("hi")
  	for {
+ 		runtime.Gosched()
  	}
  }

```

`runtime.Gosched()` は、POSIXの [`sched_yield`](https://man7.org/linux/man-pages/man2/sched_yield.2.html) のようなもの (たぶん) だ。 `sched_yield` は、他のスレッドが動けるように当該スレッドにCPUを手放させる。Goの場合はスレッドではなくgoroutineなため、 `Gosched` という名前なのかと思われる (これは予想) 。
すなわち、 `runtime.Gosched()` を明示的にコールすることで強制的にgoroutineの再スケジュールが行われ、別のgoroutineにスイッチされることを期待できる。

また、[GOEXPERIMENT=preemptibleloops](https://github.com/golang/go/blob/87a3ac5f5328ea0a6169cfc44bdb081014fcd3ec/src/cmd/internal/objabi/util.go#L257)を使う方法もある。これは「ループ」の際にプリエンプションさせるためのものだ。これを使えばコードの変更は不要となる。

## GoにおけるCooperative vs. Preemptive スケジューリング

そもそも、マルチタスクのスケジューリングには大きく2つの方式がある。それは「Cooperative(協調的)」と「Preemptive(プリエンプティブ)」だ。協調的マルチタスクは「ノンプリエンプティブ」とも呼ばれる。
協調的マルチタスクは、プログラムのスイッチがどのように行われるかは、そのプログラム自身に依存する。「協調型」と呼ばれるのは、プログラムが相互動作可能に設計されていることを意図した呼び方なのだと思われる。
プリエンプティブ型のマルチタスクでは、プログラムのスイッチはOSに委ねられる。優先度を基にしたものや、FCSV・ラウンドロビンなど、なんらかのアルゴリズムに基づいてOSによってスイッチされるスケジューリング方式を言う。

さて、Goにおけるgoroutineのスケジューリングは協調的か、それともプリエンプティブだろうか？

こうと言い切るのはなかなか難しいが、少なくともGo1.13までは協調的だと言ってよいだろう。

オフィシャルなドキュメントを見つけられなかったが、いろいろ調べたところgoroutineのスイッチは以下のようなときに起こるらしい(網羅的ではない。);

* バッファされていないチャネルへの読み書きによる待ち
* システムコールの呼び出しによる待ち
* time.Sleep()の呼び出し
* mutexの解放待ち

また、Goでは「sysmon」という関数を実行し続けるコンポーネントが動いており、プリエンプション (以外にもネットワーク処理の待ち状態をノンブロッキングにしてあげるとか、いろいろ) をやっている。
sysmonの実体としてはM (Machine) だが、P (Processor) 無しで動く。MとかPとかは様々な解説記事 ([これ](https://developpaper.com/gmp-principle-and-scheduling-analysis-of-golang-scheduler/)とか) を参照することをお薦めする。

sysmonはMが同じG (Goroutine) を10ms以上実行し続けているのを見つけると、そのGの内部パラメータである `preempt` フラグをtrueにする。その後、そのGが関数コールした際のfunction prologueで、Gは自身の `preempt` フラグを確認し、trueだった場合は自身をMから切り離し、グローバルキューというキューにプッシュする。すなわち、無事プリエンプションが実行されたというわけだ。
ちなみに、グローバルキューとはPが持つGのキュー (=ローカルキュー) とは異なるキューである。グローバルキューの目的はいくつかある。

* ローカルキューはキャパシティが256であり、それを超えるGを格納するため
* 様々な要因で待ちになっているGを格納するため
* preemptフラグで切り離されたGを格納するため

ここまでがGo1.13までの実装であるが、ここまでを理解すれば前述の無限ループするコードが何故期待通りに動かなかったのかがわかるだろう。 `for {}` は単なるビジーループなので、先に書いたようなgoroutineのスイッチの契機には特にならない。「10ms以上実行されているからsysmonによってpreemptフラグが立てられるのでは？」と思うが、 **preemptフラグが立っても関数コールがなければそのフラグのチェックが発生しない** のである。先に書いたように、preemptフラグのチェックはfunction prologueで起こるから、何の処理もしないビジーループではプリエンプションの実行まで到達できなかったのである。

そして、Go1.14で導入された "non-cooperative preemption" (asynchronous preemption) によって、この挙動は変わった。

## asynchronously preemptibleとはなにか？

ここまでを整理しよう。Goは10ms以上実行されているgoroutineをsysmonで監視し、適宜強制的にプリエンプションするよう計らう仕組みがそもそも存在した。しかし、その動作の仕組み上、 `for {}` のような場合は実際はプリエンプションが発生しなかった。

Go1.14で導入されたnon-cooperative preemptionによって、Goroutineのスケジューラはプリエンプティブと呼んで差し支えないようになった。それは、シグナルを使ったシンプルながら効果的なアルゴリズムである。

まず、sysmonは今まで通り、10ms以上動き続けているG (goroutine)を検知する。すると、sysmonはそのGを動かしているスレッド (P) にシグナル (SIGURG) を送信する。
Goのsignal handlerはシグナルをハンドリングするためにそのPに対して `gsignal` という別のgoroutineを起動し、それまで実行していたGの代わりにMと対応付け、gsignalにシグナルを確認させる。gsignalはプリエンプションが命じられたことをわかり、それまで実行していたGを停止する。

すなわち、Go1.13までは関数コールがないと仕組み上プリエンプションしなかったが、Go1.14では明示的なシグナルの送信によってプリエンプションが実行されるようになった。言い換えると、プリエンプションをgoroutine自身でなくシグナルを契機とした外的要因で実行できるようになったのだ。

この、シグナルを用いた非同期のプリエンプションの仕組みによって、先述のコードは期待通り動くようになった。それでも、 `GODEBUG=asyncpreemptoff=1` にすることでasynchrnous preemptionはオフにすることが可能だ。

ちなみに、SIGURGを使う理由は、SIGURGが既存のデバッガなどのシグナルの使用を妨げないことや、libcで使われていないことなどから選んだらしい。 ([参考](https://github.com/golang/proposal/blob/master/design/24543-non-cooperative-preemption.md#other-considerations))

## 終わりに

何もしない無限ループが他のgoroutineに処理を渡さないからと言って、Go1.13までの仕組みがダメかというとそうでもないと思われる。 [@davecheney氏](https://github.com/golang/go/issues/11462#issuecomment-116616022)も発言しているように、通常これは特に問題にならないと考えられる。そもそもasynchronous preemptionはこの無限ループの問題を解決するために導入されたのではない。

asynchronous preemptionの導入によってスケジューリングがプリエンプティブになったものの、GCの際の「アンセーフ・ポイント」の取り扱いにさらに注意が必要となった。この辺の実装の考慮が大変面白くて話したかったのだが、力尽きたのでここでこの記事は終わる。気になる読者は自分で[Proposal: Non-cooperative goroutine preemption](https://github.com/golang/proposal/blob/master/design/24543-non-cooperative-preemption.md)を読んでほしい。

## 参考

* [Proposal: Non-cooperative goroutine preemption](https://github.com/golang/proposal/blob/master/design/24543-non-cooperative-preemption.md)
* [runtime: non-cooperative goroutine preemption](https://github.com/golang/go/issues/24543)
* [runtime: tight loops should be preemptible](https://github.com/golang/go/issues/10958)
* [runtime: golang scheduler is not preemptive - it's cooperative?](https://github.com/golang/go/issues/11462)
* [Source file src/runtime/preempt.go](https://golang.org/src/runtime/preempt.go)
* [Goroutine preemptive scheduling with new features of go 1.14](https://developpaper.com/goroutine-preemptive-scheduling-with-new-features-of-go-1-14/)
* [Go: Goroutine and Preemption](https://medium.com/a-journey-with-go/go-goroutine-and-preemption-d6bc2aa2f4b7)
* [At which point a goroutine can yield?](https://stackoverflow.com/questions/64113394/at-which-point-a-goroutine-can-yield)
* [Go: Asynchronous Preemption](https://medium.com/a-journey-with-go/go-asynchronous-preemption-b5194227371c)
* [go routine blocking the others one [duplicate]](https://stackoverflow.com/questions/17953269/go-routine-blocking-the-others-one)
* [Golangのスケジューラあたりの話](https://qiita.com/takc923/items/de68671ea889d8df6904)
* [goroutineがスイッチされるタイミング](https://qiita.com/umisama/items/93333ffe4d9fc7e4ba1f)
