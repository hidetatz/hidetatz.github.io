Goにおけるミューテックスのスタベーション---2021-04-06 20:00:00

## はじめに

ミューテックスはマルチスレッドプログラミングにおける同期プリミティブのひとつだ。ミューテックスオブジェクトは通常 `Lock()` `Unlock()` のインタフェースを持ち、複数のスレッドから同時にロックされることがないことを保証する。利用する側のプログラムは、クリティカルセクションに入る前には必ずロックを獲得し、クリティカルセクションが終わり次第速やかにロックを手放すようにすれば、マルチスレッド環境でも排他ができるというわけだ。

ミューテックスは極めてシンプルなロック機構であるがゆえに、そのアルゴリズムには様々なバリエーションが存在する。しかし、そのアルゴリズムは慎重に選択されないと時にスタベーションを引き起こし、スループットの低下やシステムの停止状態などを引き起こす。
この記事では、Go1.8までに存在したある特定の状況でスタベーションを引き起こしてしまうミューテックスのアルゴリズムについて、実際のプログラムを実例に取りながら、ミューテックスの動作やスタベーション、ロックのアルゴリズムについての理解を深める。そして、Goが1.9で施した変更を見ていく。

## Go1.8に存在したミューテックスの不公平性

Go1.9の[リリースノート](https://golang.org/doc/go1.9#minor_library_changes) の `Minor changes to the library` 内 `sync` の項には、 `Mutex is now more fair.` とある。

まず、Go1.8までのミューテックスに存在したある問題について見ていこう。これは、「2つのgoroutineがロックを獲得しようと争っていて、片方はロックを長時間保持して短時間解放、もう片方は短時間保持して長時間解放する場合」に発生する。それは、ちょうどこんなコードだ。

```go
package main

import (
	"sync"
	"time"
)

func main() {
	done := make(chan bool, 1)
	var mu sync.Mutex

	// goroutine 1
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				mu.Lock()
				time.Sleep(100 * time.Microsecond)
				mu.Unlock()
			}
		}
	}()

	// goroutine 2
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Microsecond)
		mu.Lock()
		mu.Unlock()
	}
	done <- true
}
```

このプログラムは、1つgoroutineを起動し、その中ではfor-loopを起動している。ループの中ではミューテックスをロックし、ロックした後100µsスリープし、ロックを手放している。
また、別のgoroutineでは、10回のみ実行するfor-loopを起動し、そちらでは100μsスリープした後、ロック獲得、および速やかなロック解除をしている。

このプログラムがどう動くかを予想してみると、goroutine2は最終的には10回ロックを獲得するが、その間にgoroutine1にもロックを何度か取られるだろう、そういった取りつ取られつを繰り返しつつもすぐに、おそらくは1秒か2秒でプログラムは正常終了する、といった感じに筆者には見える。
このプログラムを、手元のGo1.8で動かしてみると実際はどうなったかというと、いつまで経ってもプログラムが終わらない。時間がかかりすぎたので途中でプログラムを終了してしまったが、少なくとも3分ほどはプログラムが終了しなかった。 (興味があれば試してみてほしい。)

このプログラムは、Go1.9以降であれば期待した通りすぐに終了する。なぜこうなるのかを説明する前に、いくつか説明すべきことがあるので先に説明する。

## スタベーション (Starvation)

スタベーションとは「飢餓」を意味する。マルチタスクの文脈でスタベーションとは、あるプロセスが必要なリソースを長期間獲得できないことを言う。ミューテックスの話題のときはより具体的に、READ/WRITEロックをなかなか獲得できない状況のことをスタベーションという。
上のプログラムではスタベーションが起こっている。goroutine2はロックを獲得したいにもかかわらず、そのロックは後述するミューテックスのアルゴリズムの関係でロックを獲得できていない。一般的にスレッド間でのロックの獲得は、スループットを犠牲にしない範囲でなるべく公平であることが期待される。

## ミューテックスの動作を決めるアルゴリズム

前述のとおり、ミューテックスはシンプルな同期プリミティブである。その動作をもう一度確認しよう。

* ロックされていないミューテックスは `Lock()` でロックすることができる
* ロックされているミューテックスは `Unlock()` でロックを解除することができる
* 同時にロックできるのは1スレッドだけである
* ロックされているミューテックスを `Lock()` しようとすると、ブロックする

このようなインタフェースを守りながら動くロックのアルゴリズムは複数存在する。

それらロックのアルゴリズムが規定するのは結局、ロックの獲得がかち合った時に、次のロックを獲得できるのはどのスレッドか？という点だ。
以下でいくつか説明するが、これは網羅的ではない。

## ロックのアルゴリズム #1. Thunder lock

Thunder Lockではロックを獲得しようとしたが既にロックされていた場合、そのスレッドはロック解放待ちのプールのようなものに入りスリープする。これをよく「パーキング (Parking、駐車)」という。
ロック解放時にはパーキング中の全てのスレッドを一斉に起こし、それらのスレッドはまたロックを獲得しようとする。当然ながら、その後ロックをひとつのスレッドだけが獲得し、他のスレッドは再度スリープする。
Thunder Lockの性質上、ミューテックスを解放する時にはThundering herdが発生する。
Thunder lockは実装がシンプルというメリットがあるが、筆者としては特にこれを採用する理由はない印象ではある。

## ロックのアルゴリズム #2. Barging

Bargingはいろいろ調べたところおそらくJavaで初めて生まれたタームである。Bargeは日本語で「はしけ」というものらしいが、これは要するに「ある点とある点を行き来して何かを運ぶ橋渡しのようなもの」と理解すれば良いはず (たぶん) 。

Bargingでは多くの場合ミューテックスを待つスレッドのキューがある。スレッドはtry Lock時にそのロックが既にLockedであることを確認すると、そのキューに入りスリープする。ここまではThunder lockと似ている。

ミューテックスがリリースされる時には、そのキューからポップしたスレッドをスリープから起こす。ロックはその起こされたスレッドか、新たにロックを取ろうとしたスレッドの**どちらか**に、ロックを渡す。どちらに渡すかはその時次第である。
ロック獲得を待つスレッドが既にいたにもかかわらず新たにロックを取りに来たスレッドがロックを取ってしまう事を `Lock stealing` 、ロックの盗難という。

BargingはLock stealingがあるがゆえに比較的高いスループットを実現可能である。そもそも、スリープしているスレッドを起こすことは通常ハイコストな操作で非常に時間がかかる。
Bargingの持つデメリットは、ロック獲得の公平性を犠牲にすること、そして、スタベーションが起こりうることだ。

## ロックのアルゴリズム #3. Handoff

Handoffは文字通り「手渡し」を意味する。Handoffでは、ミューテックスがリリースされる時にキューからスレッドをポップし、そのスレッドを起こす。Bargingのように新たに来たロック獲得リクエストにロックを渡すことはなく、必ずロックの待ち行列の先頭にいるスレッドが次にロックを獲得できるようになる。すなわち、Handoffロックは厳密なFIFOロックと言える。

Handoffでは一般的に、スループットは犠牲になる。何故なら、キューの先頭にいるスレッドの他に、ロックを獲得しようとしたスレッドが仮にいても、そのスレッドにはロックを決して渡さないからだ。ロックを獲得しようとしたスレッドは必ずキューに入り、順番を待つ必要がある。
Linuxカーネルで、Handoffなミューテックスを実装したパッチは[こちら](https://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git/commit/?id=9d659ae14b545c4296e812c70493bfdc999b5c1c)である。

## Adaptive lockとスピン

いくつかのロックアルゴリズムを簡単に見てきたが、これらは別のなにかと組み合わされて使用することがあり、これをAdaptive lockということが多い。
Adaptive lockという言葉にははっきりした定義がなく、様々な「適応型」ロックアルゴリズムの総称としてそう言われることが多い。例えば、[Adaptive Locks: Combining Transactions and Locks for Efficient Concurrency](https://yanniss.github.io/al-transact09.pdf) ではいわゆるミューテックスとSTMを組み合わせるアプローチを扱っている。

Adaptive lockのパターンのひとつに、スピンとパーキングを組み合わせるものがある。とりあえず何度かスピンしてみて、それでもロックが取れなければスレッドをパーキングするというものだ。
とりあえずスピンしてみることは、「今はロックが取れないけど、実際はすぐに - 具体的には、スレッドをパーキングしてから起こすのにかかる時間よりは短い間に - ロックは取れるようになる」といった場面で、スピンする間にロックが取れることを意図している。

例えば、[WebKit](https://github.com/WebKit/WebKit/blob/88278b55563e5ccdc0b3419c6c391c3becc19e40/Source/WTF/wtf/LockAlgorithmInlines.h#L57-L62)のソースコードを見てみると、40回スピンして、それでもだめならパーキングしているのがわかる。

## Go1.8のsync.Mutexのアルゴリズム

話をGoに戻そう。Go1.8までのミューテックスは、上記のBargingとAdaptive lockを組み合わせていた。ロック済みのミューテックスに対してロックを獲得しようとした時、まず[mutex.Lock()](https://github.com/golang/go/blob/go1.8/src/sync/mutex.go#L61-L72)の中で呼び出す[sync_runtime_canSpin](https://github.com/golang/go/blob/go1.8/src/runtime/proc.go#L4477-L4490)でスピンするかを判定している。それによれば、Goのミューテックスはまず[4回](https://github.com/golang/go/blob/go1.8/src/runtime/lock_futex.go#L30)だけスピンする。

このアルゴリズムは特定のgoroutineがスタベーションに陥っていたとしてもそれを考慮せずに動く。このため、先のプログラムは期待通り動かなかった。これが1.9でどう変わったのかを説明する。

## Go1.9のsync.Mutexのアルゴリズム

修正のコミットは[これ](https://github.com/golang/go/commit/0556e26273f704db73df9e7c4c3d2e8434dec7be#diff-f6dc3e83d9b4548fbba149aca4d4307b8d4551951978fd9c1b98dff9c1ada149)である。

1.9のsync/mutex.goでは、新たにスタベーションモード (starvation mode、飢餓モード) という概念が導入された。

1.9では、[1e6ns (= 1ms)](https://github.com/golang/go/commit/0556e26273f704db73df9e7c4c3d2e8434dec7be#diff-f6dc3e83d9b4548fbba149aca4d4307b8d4551951978fd9c1b98dff9c1ada149R43-R67) の間ロックを獲得できなかったgoroutineは[スタベーションモード](https://github.com/golang/go/commit/0556e26273f704db73df9e7c4c3d2e8434dec7be#diff-f6dc3e83d9b4548fbba149aca4d4307b8d4551951978fd9c1b98dff9c1ada149R136)に入っていると判定される。以下のようなコードだ。

```go
	starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
```

スタベーションモードに入ったgoroutineは、自身が1ms以上待ち続けていることを検知すると、ミューテックスが内部で使用しているセマフォの獲得に [`lifo=true` を併せて渡すようになる](https://github.com/golang/go/commit/0556e26273f704db73df9e7c4c3d2e8434dec7be#diff-f6dc3e83d9b4548fbba149aca4d4307b8d4551951978fd9c1b98dff9c1ada149R130-R135)。

```go
	// If we were already waiting before, queue at the front of the queue.
	queueLifo := waitStartTime != 0
	if waitStartTime == 0 {
		waitStartTime = runtime_nanotime()
	}
	runtime_SemacquireMutex(&m.sema, queueLifo)
```

セマフォ (runtimeパッケージ内) 側では `lifo=true` の場合、待ち状態のgoroutineを格納しているTreap (キューの実体となるデータ構造) で、そのgoroutineをキューの先頭に配置する。これによってほかに待っているgoroutineをごぼう抜きすることができる。

また、ミューテックスの解放時にはセマフォのリリースに[handoff=true](https://github.com/golang/go/commit/0556e26273f704db73df9e7c4c3d2e8434dec7be#diff-f6dc3e83d9b4548fbba149aca4d4307b8d4551951978fd9c1b98dff9c1ada149R208-R212)を併せて渡すようになった。
セマフォ側ではリリース時に待ち状態のgoroutineのキューからあるgoroutineをdequeueするが、 `handoff=true` の場合はデキューしたgoroutineの `ticket` を直接 `1` に、すなわち[キューから出る権利を与えられる](https://github.com/golang/go/commit/0556e26273f704db73df9e7c4c3d2e8434dec7be#diff-bb07e8d113e0257192c87f8b6153be1bcb547aa7826db102178ce2e6b7fd98d8R178-R195)。これによって、ミューテックスのアンロックの際に新たにやってきたgoroutineにロックをstealされることなく、直接そのオーナーシップをキューの中のgoroutineに渡すようになったことを意味する。

```go
	s, t0 := root.dequeue(addr)
	if s != nil {
		atomic.Xadd(&root.nwait, -1)
	}
	unlock(&root.lock)
	if s != nil { // May be slow, so unlock first
		// ...
		if handoff && cansemacquire(addr) { // here
			s.ticket = 1
		}
		// ...
	}
```

## 終わりに

Go1.8までに潜在的に存在したミューテックスのスタベーションの問題について見てきた。コードのコメントなどにもある通り、HandoffはBargingと比較するとパフォーマンスで劣るがより高い公平性をもたらす。
[大本となったissue](https://github.com/golang/go/issues/13086) では `sync.FairMutex` みたいなのがあったほうがいいのではみたいな話もあって面白いので、興味のある読者は読んでみてほしい。

## 参考

* [runtime: fall back to fair locks after repeated sleep-acquire failures](https://github.com/golang/go/issues/13086)
* [Go 1.9 Release Notes - Minor changes to the library](https://golang.org/doc/go1.9#minor_library_changes)
* [Locking in WebKit](https://webkit.org/blog/6161/locking-in-webkit/)
* [Fine-grained Adaptive Biased Locking](http://www.filpizlo.com/papers/pizlo-pppj2011-fable.pdf)
* [GNU/Linux でのスレッドプログラミング NPTL (Native POSIX Thread Library) Programming.](http://www.tsoftware.jp/nptl/)

## 参考にしたソースコード

* https://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git/commit/?id=9d659ae14b545c4296e812c70493bfdc999b5c1c
* https://github.com/WebKit/WebKit/blob/88278b55563e5ccdc0b3419c6c391c3becc19e40/Source/WTF/wtf/LockAlgorithmInlines.h#L57-L62
* https://trac.webkit.org/browser/trunk/Source/WTF/benchmarks/LockSpeedTest.cpp?rev=200444
* https://trac.webkit.org/browser/webkit/trunk/Source/WTF/wtf/ParkingLot.cpp?rev=200444
* https://trac.webkit.org/browser/trunk/Source/WTF/wtf/Lock.cpp?rev=200444
