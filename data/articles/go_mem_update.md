title: Goのメモリモデルのアップデート
timestamp: 2022-06-12 11:00:00
lang: ja
---

Goのメモリモデルのページは長らく[May 31, 2014](https://web.archive.org/web/20211227220453/https://tip.golang.org/ref/mem)バージョンだったんだけど、つい最近[June 6, 2022](https://tip.golang.org/ref/mem)バージョンに更新されていた。
これまでのざっくりしたものとは違い、かなり厳密にリライトされているようなので、少し中身を詳しく読んでみようと思う。

メモリモデルに馴染みのない人は、筆者が以前書いた記事[メモリモデルとはなにか](/articles/2022/02/14/memory_model_ja/)も併せて読んでもらえると良いと思う。

## これまでのメモリモデル

これまでのGoのメモリモデルは、内容としては上記のMay 31, 2014のものを読むのが一番良いのだけど、かなりざっくりしたものであった。
Happens-Before関係を (改めて) 説明したあとに、Goroutineやチャネル、init関数といったものがランタイムでどのように順序保証されるかと、間違った同期のやり方はどんなものかという例を挙げるに留まっている。
すなわち、データレースに関する詳しい説明や、最も重要 (と筆者は思うのだけど) なデータレース時の挙動についてはあまり書かれていなかった。

Goでデータレース発生時にどうなるかは、[Benign Data Race and Undefined Behaviour](https://groups.google.com/g/golang-nuts/c/EHHMCdcenc8)や[Does Go have `undefined behaviour` ?](https://groups.google.com/g/golang-nuts/c/MB1QmhDd_Rk)、あるいは[Benign data races: what could possibly go wrong?](https://web.archive.org/web/20150604005924/http://software.intel.com/en-us/blogs/2013/01/06/benign-data-races-what-could-possibly-go-wrong)を読むと、「おそらく (C/C++のような) Undefined Behaviorになるのかな?」といった感じに見える。しかし、Goのドキュメントに直接書かれていないので、正直なところ筆者はよくわかっていなかった。ただ、別にデータ競合のあるプログラムを書きたいわけではないので、ロックやatomicはもちろん使っていたし、ThreadSanitizerも積極的に利用する、といった感じであった。

これは、別にGoのメモリモデルが本質的に曖昧だったわけでは別になく、単に古くて十分に書かれていなかったということだと認識している。[Updating the Go Memory Model](https://research.swtch.com/gomm)にもそんなことが書かれているし、[GitHub Discussion](https://github.com/golang/go/discussions/47141)での議論も進んでいたようだ。[doc: define how sync/atomic interacts with memory model #5045](https://github.com/golang/go/issues/5045)では、2013年からこの辺について会話しているようである。

## アップデートされたメモリモデル

アップデートされたメモリモデルでは、いきなりとても重要なことが書かれているのでそのまま引用する (「Informal Overview」より) 。

> While programmers should write Go programs without data races, there are limitations to what a Go implementation can do in response to a data race. An implementation may always react to a data race by reporting the race and terminating the program. Otherwise, each read of a single-word-sized or sub-word-sized memory location must observe a value actually written to that location (perhaps by a concurrent executing goroutine) and not yet overwritten. These implementation constraints make Go more like Java or JavaScript, in that most races have a limited number of outcomes, and less like C and C++, where the meaning of any program with a race is entirely undefined, and the compiler may do anything at all. Go's approach aims to make errant programs more reliable and easier to debug, while still insisting that races are errors and that tools can diagnose and report them.

ここに書かれているのは、「Goではデータレースはエラーであって、Goはそれを報告しプログラムを終了させることができるよ」ということである。
すなわちGoでは、いわゆるC++の「DRF-SC or Catch Fire」のように、プログラマが正しく同期しないと未定義動作が引き起こされなにが起こるか全く不明、ということではないという表明だと思われる。

これについて詳しいことが、「Implementation Restrictions for Programs Containing Data Races」に書かれているのでこちらも引用する。

> First, any implementation can, upon detecting a data race, report the race and halt execution of the program. Implementations using ThreadSanitizer (accessed with “go build -race”) do exactly this.

> Otherwise, a read r of a memory location x that is not larger than a machine word must observe some write w such that r does not happen before w and there is no write w' such that w happens before w' and w' happens before r. That is, each read must observe a value written by a preceding or concurrent write.

> Reads of memory locations larger than a single machine word are encouraged but not required to meet the same semantics as word-sized memory locations, observing a single allowed write w. For performance reasons, implementations may instead treat larger operations as a set of individual machine-word-sized operations in an unspecified order. This means that races on multiword data structures can lead to inconsistent values not corresponding to a single write. When the values depend on the consistency of internal (pointer, length) or (pointer, type) pairs, as can be the case for interface values, maps, slices, and strings in most Go implementations, such races can in turn lead to arbitrary memory corruption.

* データレースが発生すると、Goはそれを報告してプログラムの実行を停止することができる
* データレースがないのであれば、マシンワード以下のサイズのメモリへの読み書きはプログラマの期待通りになる
* マシンワードより大きなメモリ読み取りは、アトミックにはならない可能性がある

1つ目に関して。[Benign data races: what could possibly go wrong?](https://web.archive.org/web/20150604005924/http://software.intel.com/en-us/blogs/2013/01/06/benign-data-races-what-could-possibly-go-wrong)では、「プログラマがデータレースを起こすと未定義動作が引き起こされaccidental nuclear missile launchが発生するかもしれないよ」と冗談ぽく書かれていた。Goではプログラムを停止させることがOKなので、そういったリスクは幾分低減されているようである。ただし、これはあらゆるGoの実装がデータレースに際して必ずプログラムを停止させなければいけないという意味ではないと思われるので、プログラマは依然としてデータレースのないプログラムを書かなければいけないし、 `-race` なども使用すべきである。

2つ目に関して。データレースのないプログラムでは、マシンワードサイズ以下のメモリへの読み書きは、rやwといった言葉で説明されているが、これはすなわち逐次一貫した、プログラマの期待通りの動きとなると思われる。

3つ目に関して。当然のことだが、メモリ読み取りのサイズがシングルワードを超えてしまうと、それらは順序不定な複数回のマシンワードサイズの読み取りになる。例えば64ビットマシンでメモリから128ビット読みたいとき、それは64ビットの読み取り2回で実現される。これまでこのことは簡単にしか明記されていなかったが、このドキュメントから詳しく書かれている。
マルチワード変数へのアクセスは自動ではアトミックにならないことはとても重要で、Goではインタフェースやマップ、スライス、あるいは単なる文字列であっても内部構造はマルチワードになる。このことは[Ice cream makers and data races](https://dave.cheney.net/2014/06/27/ice-cream-makers-and-data-races)に詳しく書かれていて、手元でも簡単に試すことができるのでやってみると良いのではないかなと思う。

## 終わりに

アップデートされたメモリモデルでは、「データレース時の実装の制約によってGoのメモリモデルはJavaやJavascriptに近づいた」とあるが、どちらかというとCやC++っぽさが減ったのがポイントかなと思っている。ここで紹介した内容以外にも、データレースのないプログラムのセマンティクスについてかなり厳密に書かれていたり、コンパイラ開発者向けの最適化に関するガイドも書かれていたりするので、興味があれば最初から最後まで読んでみると良いのではないかなと思う。
