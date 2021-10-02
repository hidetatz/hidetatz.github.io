type: blog
title: sync.Condとはなにか
timestamp: 2021-04-12 11:00:00
lang: ja
---

## はじめに

Goの `sync.Cond` は、マルチスレッド (正確にはgoroutine) プログラミングにおける同期プリミティブだ。
sync.Condはミューテックスと比べユースケースが限定的で、かつ使い方もやや複雑であると感じられる。そのためか、あまり現実世界での使用例を見ないように筆者には思われる。
本記事ではsync.Condについて、単に使い方のリファレンスとしてではなく、筆者の可能な限り体系的に説明する。併せて、なるべく現実世界にありそうで、かつsync.Condを使うとうまくプログラミングできるような題材をベースにさらに説明し、理解を深めていく。

## sync.Condとはなにか？

sync.Condを一言でいうと、「**条件変数 (condition variable)**」をGoプログラムで利用するための仕組みだ。
まず、[sync.Condのgodoc](https://golang.org/pkg/sync/#Cond)には次のように書かれている。

```
Cond implements a condition variable, a rendezvous point for 
goroutines waiting for or announcing the occurrence of an event.
```

ここでも `condition variable` という言葉が出てきた。condition variableとは一体何だろうか？

条件変数とはGo特有の言葉では決してなく、 (筆者の理解では) POSIXの言葉である。実際、[pthread](https://linux.die.net/man/3/pthread_cond_init) のための条件変数を操作するインタフェースも提供されている。
[c++](https://en.cppreference.com/w/cpp/thread/condition_variable) や、[Rust](https://doc.rust-lang.org/std/sync/struct.Condvar.html)、[Ruby](https://docs.ruby-lang.org/ja/2.0.0/class/ConditionVariable.html)など様々な言語で同様のメカニズムが存在する。
また、各言語のリファレンスを見てもらえばわかるように、どの言語でもほぼ同じようなインタフェースになっている。 ( `wait` ,  `signal (またはnotify_one)` ,  `broadcast (またはnotify_all)` )

条件変数はセマフォやミューテックスと同様、マルチスレッドプログラミングにおける同期プリミティブだ。すなわち、共有資源への同時アクセスを避けるための仕組みである。
ただし、セマフォやミューテックスがそれ単体で同期プリミティブとして動作するのに対し、条件変数はミューテックスとセットで使用する。
ミューテックスは同期プリミティブとして非常にシンプルで理解しやすいのに対し、条件変数はやや複雑である。そこで、ミューテックスと条件変数を対比することで条件変数への理解を深めることを目指す。

### 条件変数とミューテックスの違い

ミューテックスは、「一つのスレッドだけがクリティカルセクションに入れること」を保証するための極めてシンプルな仕組みだ。
ミューテックスを使いたい場面とは、排他制御を行いたい時のはずだ。すなわち、「あるクリティカルセクションにアクセスしたくて、自スレッドだけがそこにアクセスしていることを保証したい時」と言い換えることもできるだろう。

条件変数は、上記の「あるクリティカルセクションにアクセスしたい」状況にあるポイントを加えたようなものだ。そのポイントとは、「あるクリティカルセクションにアクセスしたいが、そのクリティカルセクションがxxxという条件になるまではアクセスしたくない」という状況である。やや抽象的な書き方なのでわかりにくいかもしれないが、後で具体的な例を挙げて説明する。

ミューテックスで上記のようなことを達成するためには、一般的にはビジーウェイトとして実装する必要がある。以下の疑似コードを見てみよう。

```
mutex_lock(); // クリティカルセクションをロック
for (!condition) { // conditionがtrueになるまでスピンする (ビジーウェイト)
    mutex_unlock();
    sched_yield();
    mutex_lock();
}
do_something();
mutex_unlock();
```

まずクリティカルセクションをロックで守った後、「condition (条件)」を確認し、条件が満たされるまで無限ループしながら確認し続ける。ループの中では他のスレッドのために一度ロックを解放し、再度ロックが取れたらループの頭に戻っている。

ビジーウェイトはCPUリソースを浪費する点で、通常のアプリケーションプログラミングでは基本的には避けたい実装である。
条件変数を使えば、これをスピンなしで書くことができる。

```
mutex_lock();
for (!condition) {
  cond_wait(); // 内部ではmutexのUnlockとLockが呼ばれる
}
do_something();
cond_notify_all();
mutex_unlock();
```

疑似コードなのでちょっとわかりにくいのだが、 (Goによる動くプログラムは後述) `cond_wait` は内部で以下のように動作する:

* ミューテックスのアンロック
* 自スレッドをサスペンドし、条件変数の通知を待機
* 通知が来たらミューテックスのロックの獲得

(Goの [`sync.Cond.Wait()` の実装](https://golang.org/src/sync/cond.go?s=1353:1374#L42)を見るとよりわかりやすい。)

ロックが取れるとconditionのチェックを行い、真ならfor-loopを抜けクリティカルセクション内で何らかの処理をする。その後、notify_allをコールする。これによって、他のwaitしているすべてのスレッドが起こされる。

`notify_all` の代わりに `notify_one` というインタフェースもあり、こちらでは「すべてのスレッド」ではなく「あるスレッドひとつだけ」が起こされる。

ミューテックスと条件変数の違いは、「条件の確認」をするにあたってミューテックスはスピンする必要があるのに対し、条件変数ではイベントのような仕組みで明示的な通知の送信を待つことができる。

これは例えるなら、クライアント・サーバシステムを作っていて、「サーバ側で行われる何らかの処理の成功を待つ」時に、 **処理のステータスを定期的にポーリングするよりも、成功したらイベントを送ってもらうほうが効率がいい** よね、といった話に似ていると思う。後述するソースコードを見てもらえれば、条件変数ではこのようなコードを**「ある条件の変更があったらイベントをパブリッシュし、その条件の変更を利用したい人たち (スレッドたち) はそれをサブスクライブする」** といった、メッセージングモデルにおけるパブリッシュ・サブスクライブモデルのように実装できることを見て取れると思う。

また、ミューテックスを使った実装ではsched_yieldをコールして他スレッドが走れるように計らっているが、条件変数ではこれは条件変数自身が勝手に行ってくれることもメリットと言える。

次に、sync.Condを使ったプログラム例を見ながらさらに理解を深める。

## sync.Cond in Real World

ここでは例として、マルチスレッドで動くキューをsync.Condを使ってうまく実装してみようと思う。単なるキューだと条件変数など不要なので、特別に以下のような仕様に則るものとする。

* 格納する値はint型のみ
* 最大長が存在し、最大長を超えた数の要素が格納されてはならない
* 最大長を超えて要素をプッシュしようとすると **ブロックする**
* 要素が存在しない状態でポップすると **ブロックする**

なお、この例は[A simple condition variable primitive](http://joeduffyblog.com/2009/07/13/a-simple-condition-variable-primitive/)および[条件変数 Step-by-Step入門](https://yohhoy.hatenablog.jp/entry/2014/09/23/193617)から拝借した。良い記事をありがとうございます。

### ミューテックスを使った実装

まずはミューテックスを使った実装を見てみよう。このソースコードはGitHubでも参照できる。[dty1er/size-limited-queue/mutex_slqueue.go](https://github.com/dty1er/size-limited-queue/blob/381725020d4de089741743523ac2b032ee946767/mutex_slqueue.go)

```go
type MutexQueue struct {
	mu       sync.Mutex
	capacity int
	queue    []int
}

func NewMutexQueue(capacity int) *MutexQueue {
	return &MutexQueue{
		capacity: capacity,
		queue:    []int{},
	}
}

func (s *MutexQueue) Push(i int) {
	s.mu.Lock()
	for len(s.queue) == s.capacity {
		s.mu.Unlock()
		runtime.Gosched()
		s.mu.Lock()
	}

	s.queue = append(s.queue, i)
	s.mu.Unlock()
}

func (s *MutexQueue) Pop() int {
	s.mu.Lock()
	for len(s.queue) == 0 {
		s.mu.Unlock()
		runtime.Gosched()
		s.mu.Lock()
	}

	ret := s.queue[0]
	s.queue = s.queue[1:]
	s.mu.Unlock()

	return ret
}
```

`Push()` から見ていく。このプログラムはマルチスレッド (goroutine) で正しく動かねばならない。また、前述のとおりキューには最大長がある。すなわち、キュー長の確認と要素の追加はアトミックでなければならない。そのアトミック性の担保にミューテックスを使用している。

最大長を超えるプッシュ操作をしようとすると空きができるまでブロックする仕様なため、まずは実体となるスライスがcapacityを超えていないことを確認し、capacityを下回るまでロックを解除しながらスピンする。

`Pop()` もほとんど同じで、アトミックにキュー長の確認と要素のポップを行う。キューが空であれば要素が追加されるまで (ポップ可能になるまで) スピンしながらブロックする。

これはまさしく、上述の「特定の条件になるまでスピンする」、非効率的な実装になっている。これを条件変数を使って書き換えてみよう。

### 条件変数を使った実装

条件変数を使った実装は以下のようなものだ。これもGitHubで参照できる。 [dty1er/size-limited-queue/slqueue.go](https://github.com/dty1er/size-limited-queue/blob/7482018a4aae723aebe80f0ff11f6b4f4fc265bc/slqueue.go)

```go
type SizeLimitedQueue struct {
	cond     *sync.Cond
	capacity int
	queue    []int
}

func New(capacity int) *SizeLimitedQueue {
	return &SizeLimitedQueue{
		cond:     sync.NewCond(&sync.Mutex{}),
		capacity: capacity,
		queue:    []int{},
	}
}

func (s *SizeLimitedQueue) Push(i int) {
	s.cond.L.Lock()
	for len(s.queue) == s.capacity {
		s.cond.Wait()
	}

	s.queue = append(s.queue, i)
	s.cond.Broadcast()
	s.cond.L.Unlock()
}

func (s *SizeLimitedQueue) Pop() int {
	s.cond.L.Lock()
	for len(s.queue) == 0 {
		s.cond.Wait()
	}

	ret := s.queue[0]
	s.queue = s.queue[1:]
	s.cond.Broadcast()
	s.cond.L.Unlock()

	return ret
}
```

条件変数には[syncパッケージにあるCond](https://golang.org/pkg/sync/#Cond)を使う。
前述の通り条件変数はミューテックスとセットで用いられるが、sync.Condは内部でミューテックス (正確には `sync.Locker` ) を持てるので、これをミューテックスとして使う。

`Push()` を見てみよう。まず明示的にミューテックスをロックした後、条件をチェックし `Wait()` している。 前述の通り、[Wait](https://golang.org/src/sync/cond.go?s=1353:1374#L42) の中ではミューテックスのUnlockとLockを勝手に行ってくれる。条件の変更を待つのにスピンは必要なく、ランタイムが勝手にgoroutineのサスペンドおよびウェイクアップをやってくれる。

キューに要素を追加した後は、 `Broadcast()` をコールする。これによって、 `Wait()` しているすべてのgoroutineが起こされ、ロックの獲得およびループ条件の再確認を行う。
 
 `Pop()` でもやってることはほとんど同じである。ロックを取ってから `Wait()` 、その後要素のポップと `Broadcast` のコールを行う。

 条件変数を使うことで無駄なスピンがなくなり効率的な実装となった。しかし、このプログラムはまだ非効率的な部分が残っており改善の余地がある。それらについて検討しながら、さらに条件変数についての理解を深めていく。

### 条件変数を使った実装の改良

ここまでを整理しよう。条件変数でできることは、「あるスレッドが、何らかの条件が満たされることを効率よく待つ」ことと言える。例として挙げた最大長付きキューにおいての「何らかの条件」とは今のところ、「キューの要素数の変更」である。上のプログラムは `Push()` `Pop()` どちらにおいても、「キューの要素数が変わる -> キューの要素数変更イベントを待っているスレッドを起こす -> キューの要素数が自スレッドの希望する条件と一致しているかを確認」というような動きをする。

#### 1. 条件変数への通知タイミングの検討

さて、1点目の改良ポイントは、「条件変数への通知タイミングを絞ることができる」という点である。

現状、 `Push()` では「キュー長が変更された」時に `Wait()` から復帰するが、よく考えるとそもそも必要とするイベントは「キュー長がcapacityと同じだった (満杯だった) のが、それよりも少なくなった」というものだけである。例えばキューのキャパシティが10のとき、キュー長が2とか3とかしかないのであればプッシュがブロックされることはないのだから通知を飛ばす必要はない。

同様に `Pop()` でも、「キュー長が0だったのが1になった」というイベントだけ送ってもらえば良いことになる。

今はこれらを考慮せずにブロードキャストしているため、ここを必要な時のみ行うよう変えることでより効率的にすることができる。

#### 2. 条件変数に意味を持たせる

現状、条件変数 `s.cond` はキューに対して一つで、 `Push()` でも `Pop()` でも同じものを共同利用している。しかし、これらの意味について考えてみると、これはそれぞれで分離することができることがわかる。

意味的に、「 `Push()` の中でコールされる `Broadcast()` 」によって起こされるgoroutineは `Pop()` しようと待っているものであって、 `Push()` しようと待っているいるものではない。 `Push()` しようと待っているしているgoroutineは何を待っているのか？と考えてみると、彼らはキューのキャパシティが空くこと、言い換えれば「キャパシティいっぱいのキューから要素が取り出され (= Popされ) 空きが出る」のを待っているのである。逆も同じであり、「 `Pop()` の中でコールされる `Broadcast()` 」は `Push()` の機会を待つgoroutineのためのものである。

言い換えると、 `Push()` と `Pop()` が通知したいものは意味的に異なっているので、これらがひとつの条件変数 `s.cond` を共有している意味はない。

`1. 条件変数への通知タイミングの検討` と `2. 条件変数に意味を持たせる` を組み合わせると、以下のような整理をすることができる。

* 条件変数 `s.cond` : `s.nonEmptyCond` に変更し、「空だったキューが空でなくなったときに通知する」ように修正。これは `Pop()` の中の `Wait()` の際に使用する。
* 新条件変数 `s.nonFullCond` : 新たに定義。「満杯だったキューが満杯でなくなったときに通知する」ようにする。これは `Push()` の中の `Wait()` の際に使用する。

この修正により、「単にキュー長が変わっただけなのにすべてのgoroutineが起こされる」だった実装が「各goroutineが自分が関心のあるキュー長の変更の時にのみ起こされる」に変わった。

#### 3. `Broadcast()` ではなく `Signal()` を使う

1と2の変更で起こされるスレッドを最小限にすることができたが、まだ改良点はある。例えば `Push()` の際に、空だったキューが要素数1に変更されると、 `Pop()` を待っていたgoroutineは全て起こされることになる。しかし考えてみると、 `Pop()` を待っていたgoroutineをすべて起こしても次にクリティカルセクションに入れるgoroutineはひとつだけなので、すべて起こすのではなく一つのスレッドだけ起こすのが効率的だ。こういう時に使えるのが、 [`sync.Cond.Signal()`](https://golang.org/pkg/sync/#Cond.Signal) である。 `Signal()` は `Wait()` している全てのgoroutineではなく、あるひとつのgoroutineのみを起こす。これによって、待機しているgoroutineが多数あるときのThundering Herd Problemが緩和されるなどより効率的な実装となる。

上記の1, 2, 3を適用した実装を以下に記す。GitHubでも参照可能。 [dty1er/sie-limited-queue/slqueue.go](https://github.com/dty1er/size-limited-queue/blob/e9dd8dc2c2937d6cffe7e784bc5c2d436632b758/slqueue.go) 
diffは[こちら](https://github.com/dty1er/size-limited-queue/commit/e9dd8dc2c2937d6cffe7e784bc5c2d436632b758)

```go
type SizeLimitedQueue struct {
	nonFullCond *sync.Cond

	nonEmptyCond *sync.Cond

	capacity     int
	queue        []int
	mu           *sync.Mutex
}

func New(capacity int) *SizeLimitedQueue {
	mu := &sync.Mutex{}
	return &SizeLimitedQueue{
		nonFullCond:  sync.NewCond(mu),
		nonEmptyCond: sync.NewCond(mu),
		capacity:     capacity,
		queue:        []int{},
		mu:           mu,
	}
}

func (s *SizeLimitedQueue) Push(i int) {
	s.nonFullCond.L.Lock()
	for len(s.queue) == s.capacity {
		s.nonFullCond.Wait()
	}

	wasEmpty := len(s.queue) == 0
	s.queue = append(s.queue, i)

	if wasEmpty {
		s.nonEmptyCond.Signal()
	}
	s.nonFullCond.L.Unlock()
}

func (s *SizeLimitedQueue) Pop() int {
	s.nonEmptyCond.L.Lock()
	for len(s.queue) == 0 {
		s.nonEmptyCond.Wait()
	}

	wasFull := len(s.queue) == s.capacity
	ret := s.queue[0]
	s.queue = s.queue[1:]

	if wasFull {
		s.nonFullCond.Signal()
	}
	s.nonEmptyCond.L.Unlock()

	return ret
}
```

条件変数を使った実装の効率化には、「必要な時だけ通知を送る」「必要なスレッドだけ起こす」を徹底することが重要である。

次章からは条件変数にまつわるいくつかの注意点等を補足的に記していく。

## 条件変数とスプリアス・ウェイクアップ

まずはスプリアス・ウェイクアップ (Spurious wakeup、直訳では「偽りの目覚め」) について。Javaでマルチスレッドプログラミングをやっていた人なら聞いたことのある言葉かもしれない。

スプリアス・ウェイクアップとは、「条件変数利用時の待機処理 (Wait) において、通知が来ていないのにサスペンドされていたスレッドが起動すること」を言う。これは通常OSやハードウェアなどに起因していて、条件変数の利用者ではどうしようもない。

[Javaのドキュメント](https://docs.oracle.com/en/java/javase/11/docs/api/java.base/java/lang/Object.html#wait(long,int) や [C++のドキュメント](https://en.cppreference.com/w/cpp/thread/condition_variable) にはスプリアス・ウェイクアップが発生しうることが明示的に書かれている。

スプリアス・ウェイクアップに対応するのは簡単で、本記事で扱った通り、Wait処理は必ずループで包む、というのを徹底すればよい。条件変数を使いたい状況によってはループ処理はなくてもよいことがあるが、そういう場合でも必ずループの中で待つようにすれば良い。例えば、[chromeの条件変数に関する開発者向けのドキュメント](http://www.chromium.org/developers/lock-and-condition-variable)にも同様の記載がある。

さて、Goのsync.Condではスプリアス・ウェイクアップは起きるのか？というと、これは筆者が参照した[Go1.16.3](https://github.com/golang/go/blob/go1.16.3/src/sync/cond.go) の時点でのGoDocには明記されておらず、正直よくわからなかった。下にいくつか参照した記事やディスカッションを残しておく。ただ、少なくとも我々同期プリミティブの利用者は、きちんとWaitをループ内に置くことを徹底すべきと筆者は考える。

* [GopherCon and dotGo 2019 liveblogs - GopherCon 2018 - Rethinking Classical Concurrency Patterns](https://about.sourcegraph.com/go/gophercon-2018-rethinking-classical-concurrency-patterns/)
* [proposal: sync: mechanism to select on condition variables](https://github.com/golang/go/issues/16620)
* [sync.Cond and spurious wekeups](https://groups.google.com/g/golang-dev/c/Kc1nOjju3zk/discussion)
* [proposal: Go 2: sync: remove the Cond type](https://github.com/golang/go/issues/21165)

## 条件変数と Two step dance

最後にTwo step danceについて。条件変数を利用すると、不要に思えるミューテックスのロック/アンロックがある状況で1回発生する。
上記のソースコードの `Push()` `Pop()` を以下に再度示す。

```go
func (s *SizeLimitedQueue) Push(i int) {
	s.nonFullCond.L.Lock()
	for len(s.queue) == s.capacity {
		s.nonFullCond.Wait()
	}

	wasEmpty := len(s.queue) == 0
	s.queue = append(s.queue, i)

	if wasEmpty {
		s.nonEmptyCond.Signal() // 1
	}
	s.nonFullCond.L.Unlock() // 3
}

func (s *SizeLimitedQueue) Pop() int {
	s.nonEmptyCond.L.Lock()
	for len(s.queue) == 0 {
		s.nonEmptyCond.Wait() // 2
	}

	wasFull := len(s.queue) == s.capacity
	ret := s.queue[0]
	s.queue = s.queue[1:]

	if wasFull {
		s.nonFullCond.Signal()
	}
	s.nonEmptyCond.L.Unlock()

	return ret
}
```

まず順番として、ソースコード内にコメントで示した `2` のところで `nonEmptyCond` に通知が来るのを待っているgoroutineがいる。そして今、あるgoroutineが `Push()` を行っている (goroutineA)としよう。goroutineAはキューの操作を行って `1` の箇所で通知を送る。

ここで、通知が送られると、 `2` で `Wait()` していたgoroutineのうちひとつ (goroutineB) が`Wait()` の待機状態からウェイクアップされる。前述の通り、 (または `sync.Cond.Wait` のソースコードを参照すればわかりやすいが) `Wait()` の中では待機から起こされた後ミューテックスの再ロックを取ろうとする (クリティカルセクションに入るため) 。しかし、 **この時点でそのミューテックスはgoroutineAによってロックが取られている** 。 これがアンロックされるのは `3` の箇所だ。すなわち、これらの過程の中で以下のような不要なブロックが発生することになる:

* goroutineBは条件変数 `nonEmptyCond` でまずブロックされている
* goroutineBは条件変数 `nonEmptyCond` に通知が来たのでブロック解除
* goroutineBが `2` でミューテックスのロックを取ろうとするが **ブロック**
* goroutineAが `3` でミューテックスのロックを解除
* goroutineBが `2` でミューテックスをロック
* ロックが取れたのでここでようやく、goroutineBがクリティカルセクションに入れる

goroutineBは、ブロック解除 -> ブロック -> ブロック解除という状態を辿ることになり、これは [Two-step dance](https://docs.microsoft.com/en-us/archive/msdn-magazine/2008/october/concurrency-hazards-solving-problems-in-your-multithreaded-code#two-step-dance) と呼ばれている。
不要なコンテキストスイッチが発生してしまうことがTwo-step danceのデメリットだ。

コンテキストスイッチを嫌ってTwo-step danceを避ける方法もあり、これは `1` の `Signal()` を `3` の後に配置する、つまり条件変数への通知をクリティカルセクションの後にすればよい。しかし、これは時に間違った実装になってしまうこともある。
[chromiumの条件変数に関する開発者向けのドキュメント](http://www.chromium.org/developers/lock-and-condition-variable)には以下のようにある。

```
In rare cases, it is incorrect to call Signal() after the critical section, 
so we recommend always using it inside the critical section. 
The following code can attempt to access the condition variable after it has been deleted, 
but could be safe if Signal() were called inside the critical section.

  struct Foo { Lock mu; ConditionVariable cv; bool del; ... };
  ...
  void Thread1(Foo *foo) {
    foo->mu.Acquire();
    while (!foo->del) {
      foo->cv.Wait();
    }
    foo->mu.Release();
    delete foo;
  }
  ...
  void Thread2(Foo *foo) {
    foo->mu.Acquire();
    foo->del = true;
                          // Signal() should be called here
    foo->mu.Release();
    foo->cv.Signal();     // BUG: foo may have been deleted
  }
```

同ドキュメントには、Chromeの条件変数の実装はTwo-step danceを検知し、待機スレッドのウェイクアップを遅延させることでコンテキストスイッチを起こらなくしている (だから、条件変数への通知は必ずクリティカルセクション内で行うべき) とも書かれている。

## 終わりに

Goにおいて、条件変数は場合によってはチャネルで代用することも可能なため、それを使う場面は限られている。筆者は「sync.Condがよくわからない」という声をたまに聞くが、そういったことにも起因しているかもしれない。

よくある「sync.Condは複数のgoroutineを一斉に起動させたい時に使う」といった動作から理解するやり方ではsync.Condをマスターするのは難しいかもしれない。なぜなら、「複数のgoroutineを一斉に起動する場面ってどんな時だろうか？」という疑問を抱いてしまうからである。
筆者はそれよりも、「条件変数とは、マルチスレッド環境においてある条件が満たされるのを効率的に待つ方法である」ということをまず理解し、そこからGoではそれをどう使うのか、といった順番のほうがわかりやすいのではないかな、と思う。

スプリアス・ウェイクアップについては正直理解しきれていないので、詳しい人が解説記事などを書いてくれたらいいな～と思っている。ちなみに、この記事を書くにあたり色々調べていたら、以前のGoでは [sync.WaitGroupでスプリアス・ウェイクアップが起きる問題](https://github.com/golang/go/issues/7734)があったらしく興味深かった。

## 参考

* [A simple condition variable primitive](http://joeduffyblog.com/2009/07/13/a-simple-condition-variable-primitive/)
* [Chrome C++ Lock and ConditionVariable](http://www.chromium.org/developers/lock-and-condition-variable)
* [GopherCon 2018 Bryan C .Mills Rethinking Classical Concurrency Patterns](https://youtu.be/5zXAHh5tJqQ?t=801)
* [GopherCon and dotGo 2019 liveblogs - GopherCon 2018 - Rethinking Classical Concurrency Patterns](https://about.sourcegraph.com/go/gophercon-2018-rethinking-classical-concurrency-patterns/)
* [proposal: sync: mechanism to select on condition variables](https://github.com/golang/go/issues/16620)
* [sync.Cond and spurious wekeups](https://groups.google.com/g/golang-dev/c/Kc1nOjju3zk/discussion)
* [proposal: Go 2: sync: remove the Cond type](https://github.com/golang/go/issues/21165)
* [sync: spurious wakeup from WaitGroup.Wait](https://github.com/golang/go/issues/7734)
* [Ruby-Doc.org class ConditionVariable](https://ruby-doc.org/core-2.5.0/ConditionVariable.html)
* [Solving 11 Likely Problems In Your Multithreaded Code](https://docs.microsoft.com/en-us/archive/msdn-magazine/2008/october/concurrency-hazards-solving-problems-in-your-multithreaded-code)
* [sync.Cond／コンディション変数についての解説](https://lestrrat.medium.com/sync-cond-%E3%82%B3%E3%83%B3%E3%83%87%E3%82%A3%E3%82%B7%E3%83%A7%E3%83%B3%E5%A4%89%E6%95%B0%E3%81%AB%E3%81%A4%E3%81%84%E3%81%A6%E3%81%AE%E8%A7%A3%E8%AA%AC-dd2050cdfab7)
* [条件変数 Step-by-Step入門](https://yohhoy.hatenablog.jp/entry/2014/09/23/193617)
* [条件変数とダンス(Two-Step Dance)を](https://yohhoy.hatenadiary.jp/entry/20120504/p1)
* [条件変数とデッドロック・パズル（出題編）](https://yohhoy.hatenadiary.jp/entry/20140926/p1)
* [条件変数とspurious wakeup](https://yohhoy.hatenadiary.jp/entry/20120326/p1)
* [マルチスレッド・プログラミングの道具箱](https://zenn.dev/yohhoy/articles/multithreading-toolbox)

