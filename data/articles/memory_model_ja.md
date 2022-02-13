title: メモリモデル
timestamp: 2022-01-20 12:00:00
lang: ja
---

## 概要

マルチスレッドプログラミングにおいては、プログラム1行1行の実行順序がプログラマの直感に反することがある。以下の例を見てみよう。スレッド1と2はそれぞれ別のスレッド (コア) で同時に実行される。またここでは、すべての変数は初期値をゼロとする。
xとyはメモリ上のあるアドレスを指す変数で、r1、r2はスレッドローカルな、レジスタのような領域を示していることにする。

```
# スレッド1
x = 1
y = 1

# スレッド2
r1 = y
r2 = x
```

この時、このプログラムの各行の実行順序は次のように考えられるだろう。

1. スレッド1 ->  スレッド2と実行されたケース

```
x = 1
y = 1
	r1 = y // r1は1
	r2 = x // r2は1
```

2. スレッド2 ->  スレッド1と実行されたケース。

```
	r1 = y // r1は0
	r2 = x // r2は0
x = 1
y = 1
```

3. 各行がインターリーブされたケース。インターリーブのパターンは4通りある。

```
# パターン1
x = 1
	r1 = y // r1は0
y = 1
	r2 = x // r2は1
```

```
# パターン2
x = 1
	r1 = y // r1は0
	r2 = x // r2は1
y = 1
```

```
# パターン3
	r1 = y // r1は0
x = 1
	r2 = x // r2は1
y = 1
```

```
# パターン4
	r1 = y // r1は0
x = 1
y = 1
	r2 = x // r2は1
```

r1とr2の組み合わせは、1のケースで `{1, 1}` 、2のケースで `{0, 0}` 、3のケースは複数のパターンがあるがいずれも `{0, 1}` となることがわかる。

では、r1とr2が `{1, 0}` となることはあり得ないのだろうか？実は場合によっては、 `{1, 0}` になることがある。

上のいずれもケースも、スレッド1・2は、それぞれが自スレッド内での命令実行順を変えない (プログラムで指定された通りの順序で命令が実行される) ことを前提としていた。実はこの前提は必ずしも正しくない。なぜならこの前提は、プロセッサとコンパイラによる最適化を考慮していないためである。プロセッサとコンパイラは、最適化のために命令の順序を変えることがある。しかし、これではプログラマが困ってしまうので、「メモリ上に保存されたデータの可視性と一貫性について定められたルール」が設けられる。このルールのことを「メモリ一貫性モデル」あるいは単に「メモリモデル」などと呼ぶ。

すなわちメモリモデルには、プロセッサを話題とする「ハードウェア・メモリモデル」と、コンパイラを話題とする「ソフトウェア・メモリモデル」の2種類が存在する。筆者は個人的に、メモリモデルという言葉がこれらをあまり区別せず使われていることが、メモリモデルのわかりにくさの理由の一つではないかと考えている。

この記事では、メモリモデルの理解を獲得するために必要な一貫性に関する知識をまず解説し、その後ハードウェア及びソフトウェア、そしてGoのメモリモデルについて書いていく。

## なぜメモリモデルを学ぶのか

JavaやC++、Goといった高級言語でマルチスレッドプログラミングを書く時、メモリモデルの知識は必ずしも必要ない。Goではチャネルやsync/atomic、あるいはsync.Mutexなどの仕組みが既に用意されているため、これらを適切に使えればプログラムはプログラマが普通に考える通り動作する。したがって、プログラマに必要な知識はメモリモデルのような低レイヤな話題ではなく、ライブラリの適切な使用方法や並列処理設計といった高レイヤな部分である。

[The Go Memory Model](https://go.dev/ref/mem)の冒頭には、次のようにある。

> If you must read the rest of this document to understand the behavior of your program, you are being too clever.

> Don't be clever.

メモリモデルを理解している必要があるのは、OSカーネルや並行・並列処理ライブラリ、コンパイラなどの開発者である。マルチスレッドプログラミングを書く開発者は必ずしもメモリモデルを学ぶ必要はない。筆者がメモリモデルについて勉強してこのブログまで書いている理由は、単に興味があったからでしかないので、そういう前提でこれより下は読んでいただきたい。

## 逐次一貫性

まずはハードウェア・ソフトウェア関係なく、メモリモデルを理解するために必須な知識である「逐次一貫性」について。
逐次一貫性 ([Sequential Consistency](https://jepsen.io/consistency/models/sequential)) とは並行システムにおける一貫性モデルのひとつである。これ自体は単なる一貫性モデルのひとつなので、メモリモデルやマルチスレッドプログラミングとは独立して理解可能である。これから「プロセッサ」や「スレッド」という言葉を使って逐次一貫性を説明するが、これはわかりやすさのためであって、一貫性モデル自体はプロセッサやOSのスレッドとは独立した概念であることに注意して欲しい。

### 逐次一貫性の定義

逐次一貫性は、1979年のLeslie Lamportの論文「[How to Make a Multiprocessor Computer That Correctly Executes Multiprocess Programs](https://www.microsoft.com/en-us/research/publication/make-multiprocessor-computer-correctly-executes-multiprocess-programs/)」でその定義が与えられている。

> … the result of any execution is the same as if the operations of all the processors were executed in some sequential order, and the operations of each individual processor appear in this sequence in the order specified by its program.

[Wikipedia](https://ja.wikipedia.org/wiki/%E9%80%90%E6%AC%A1%E4%B8%80%E8%B2%AB%E6%80%A7)では以下のように訳されている。

> 「どのような実行結果も、すべてのプロセッサがある順序で逐次的に実行した結果と等しく、かつ、個々のプロセッサの処理順序がプログラムで指定された通りであること」

### 逐次一貫性は何を許し、何を許さないのか

逐次一貫性は、各プロセッサ自身が自分で実行する命令の順序がプログラム上の順序と一致することを保証する。その上で、複数のプロセッサによる全体の実行順序は任意となる。すなわち、各プロセッサが実行する命令はプロセッサごとに任意の順序にインターリーブされうる。以下の例を考えると、

```
# スレッド1
op_1()
op_2()

# スレッド2
op_3()
op_4()
```

__op_1はop_2よりも前に発生すること__ と __op_3はop_4の前に発生すること__ を逐次一貫性は保証する。しかし、op_3がop_1より前か、op_1とop_2の間にインターリーブされるか、あるいはop_2の後かは保証されない。複数のプロセッサによる命令がどのようにインターリーブされても、それは「the result of any execution is the same as if the operations of all the processors were executed in some sequential order」のルールに違反しない。

これら逐次一貫性に関する挙動は特に、プログラマの直感に違反するものではないだろう。書いた通りの順序で実行されるが、マルチスレッドとなるとスレッド間では命令の順序は (当然のごとく) ひとつに定まらないよ、というだけのことである。これらはプログラマにとって自然であり、かつ理想的なモデルであると考えられている (本当にこれが理想的か?というのはまた別の話である) 。

さて、メモリモデルについてのポストにも関わらず逐次一貫性について説明しているのには理由がある。逐次一貫性は、前述したコンパイラとプロセッサによる最適化を目的とした命令の順序の入れ替えに大きく関連している。というのは、逐次一貫性を諦めることで、プロセッサ・コンパイラはプログラムの実行を高速化できるのである。
次章からは、プロセッサ及びコンパイラがそれぞれどのように命令の順序を変更するのか、それに対してどうメモリモデルが関係するのかを見ていく。

### 寄り道: 他の一貫性との比較

こういった一貫性モデルには[多くの種類がある](https://en.wikipedia.org/wiki/Consistency_model)が、逐次一貫性は比較的強いモデルである (Strong consistencyと呼ばれる一貫性モデルのひとつである) 。例えば、weak consistencyに分類される、最近はマイクロサービスなどの文脈でよく言及される[結果整合性](https://en.wikipedia.org/wiki/Eventual_consistency)という一貫性モデルがある。結果整合性が保証するのは「ある操作はいずれ見えるようになる (いつかは不明)」ということだけなので、一度見えた値が巻き戻ったり、あるプロセッサが施した操作がそれとは異なる順番で見えるようになったりすることがあり得る。言い方を変えれば、あり得ないことを特に保証していない一貫性モデルである。

## ハードウェア・メモリモデル

まずはハードウェア、すなわちプロセッサにおけるメモリモデルから。
マルチプロセッサのシステムでは、あるコアによるメモリへのロードやストアが、ほかのコアから可視になる順序が、プログラムの順序と異なることが発生する。なお、この章はハードウェアについて話しているので、ここでいうプログラムとは高級言語ではなくアセンブリ (または機械語) を指している。

[Latency Numbers Every Programmer Should Know](https://gist.github.com/jboner/2841832)によれば、プロセッサのメインメモリ参照のレイテンシは、L1キャッシュ参照の200倍の時間がかかるらしい。プロセッサから見るとメモリアクセスは極めて時間がかかるので、メモリへのアクセスをなるべく不要にすることはプロセッサの設計において重要な観点である。各CPUベンダはこれを実現するために、キャッシュがどのように振る舞うかや命令の順序の入れ替えルールなどを設計した。しかし、メモリへのアクセスを減らそうとするとどうしても、あるプロセッサが他のプロセッサによる書き込みを確実に観測することが保証できなくなる。こういった各プロセッサに設計の違いが結果的に各プロセッサ固有のメモリモデルとなった。

メモリモデルの種類は、[Wikipedia](https://en.wikipedia.org/wiki/Memory_ordering)を見た感じは以下のようになっているがこれの網羅性は不明。

<img width="1458" alt="1" src="https://user-images.githubusercontent.com/60682957/153757264-b32c0900-9fd6-409b-a870-a2cfba4a67fd.png">

次に、実際のハードウェアメモリモデルについて、TSOとRMO/WMOを挙げて見ていく。ただし、これら以外にもメモリモデルの種類はある。

### ハードウェアメモリモデル1. TSO

まずは、有名な (?) メモリモデルであるTSOから。TSOは「Total Store Order」のアクロニムである。TSOは順序変更が発生しにくい、いわゆる「強いメモリモデル」に分類される。
TSOを実装しているプロセッサの典型はx86であるが、[SPARC](https://cr.yp.to/2005-590/sparcv9.pdf)や[RISC-V](https://riscv.org/wp-content/uploads/2018/05/14.25-15.00-RISCVMemoryModelTutorial.pdf)もTSOをサポートしている。

TSOでは、メモリへのストアがロードの後に順序変更されることがあり得る。それ以外の、「ロードがストアの後に順序変更」「ロードがロードの後に順序変更」「ストアがストアの後に順序変更」などは発生しない。

なぜストアがロードの後に順序変更されるのか?というと、例えば以下はx86、SPARC TSOのアーキテクチャ図である ([A Tutorial Introduction to the ARM and POWER Relaxed Memory Models](https://www.cl.cam.ac.uk/~pes20/ppc-supplemental/test7.pdf)より図を引用) 。

<img width="650" alt="2" src="https://user-images.githubusercontent.com/60682957/153757275-1d896011-2dff-4a12-adae-a4ced3ca11ef.png">

このアーキテクチャのポイントは以下の通りである。

* あるスレッド (プロセッサ) のストアは、「Write Buffer (Store Bufferとも言う)」というFIFOキューにまずプッシュされる
* プッシュが完了したらプロセッサは次の命令の実行を始める
* Write Bufferはプロセッサごとに存在し、他のプロセッサと共有されない
* ロードは、まずローカルのWrite Bufferを参照し、そこになければメインメモリを参照する
* すなわち、あるプロセッサによるストアがメモリに到達していれば、ほかのプロセッサは確実にその値を参照する (Write Bufferに書き込みがなければ)

これでなぜストアがロードの後に順序変更されるのか?というと、次のように発生する。

* プロセッサAがストアをWrite Bufferに書き込み (この時点で、プロセッサの利用者から見るとストアは完了している)
* プロセッサBが、Aが書き込んでいるメモリアドレスを参照。この時点でAの書き込みはまだWrite Bufferにあるため、BはAがストアするよりも前の値を読み取る
* Aによる書き込みが完了する

上記のようなシーケンスでは、プロセッサの利用者にとってプロセッサBは古い値を読み込んでいるように見える。これが「ストアがロードの後に順序変更」の内部的な仕組みである。
これ以外の順序変更は発生しないことは詳しくは説明しないが、アーキテクチャ図を見ながら考えれば理解できると思う。

TSOは「ストアがロードの後に順序変更」以外発生しないという点で強いメモリモデルではあるが、逐次一貫性が保証するものをTSOが全て保証するわけではない。これは後ほど説明する。

### ハードウェアメモリモデル2. WMO/RMO

WMOは「Weak Memory Order」、RMOは「Relaxed Memory Order」を意味するがいずれも「弱いメモリモデル」である。
これらはARMv7やIBM POWER、またSPARC RMOや[RISC-V WMO](https://riscv.org/wp-content/uploads/2019/06/16.15-Stefanos-Kaxiras.pdf)でサポートされている。
下はARM、POWERのアーキテクチャ図である  ([A Tutorial Introduction to the ARM and POWER Relaxed Memory Models](https://www.cl.cam.ac.uk/~pes20/ppc-supplemental/test7.pdf)より図を引用) 。

<img width="655" alt="3" src="https://user-images.githubusercontent.com/60682957/153757278-09f43ffc-f23a-4121-a5e4-c8076b6109f5.png">

以下のようなポイントがある。

* 各プロセッサはメモリへのストア・ロードをメモリの「コピー」に対して実行する
* コピーへのストアは非同期でほかのプロセッサのコピーに伝播する
* 伝播の際、順序変更があり得る

TSOと比較するとかなり緩い、「これは何を保証しているんだ?」と思ってしまうようなメモリモデルに見えるが、これは実際のところ、ロード・ストアの順序は全く保証されない。

次に、これらのメモリモデルの違いが現実のプログラムにどう影響するのかを見ていく。

### ハードウェア・メモリモデル in Action

まず、一番上で見た以下のプログラムを再度考える。繰り返すが、今私達はハードウェア・メモリモデルの話をしているので、これらのプログラムは実際にはアセンブリまたは機械語で書かれていると考えて欲しい。

```
# スレッド1
x = 1
y = 1

# スレッド2
r1 = y
r2 = x
```

このプログラムの実行結果を、「Sequential Consistencyを保証するハードウェア」「TSOなハードウェア」「RMOなハードウェア」でそれぞれ比較すると、次のようになる。

|ハードウェア|{r1, r2}が{1, 0}になることはあり得るか？|
|:---|:---|
|Sequential Consistency|No|
|TSO|No|
|RMO|Yes|

Sequential Consistencyなハードウェアにおいては、一番上で説明したように、いかなるインターリーブが発生しても{1, 0}にはならない。

TSOではどうか?TSOでも{1, 0}にはならない。TSOではストアがロードの後に順序変更されることはありえるが、このプログラムではスレッド1はストア -> ストアで、スレッド2はロード -> ロードなので、順序変更の影響を受けないためである。

RMOでは、{1, 0}が発生してしまう。RMOではロード・ストアは一切順序が保証されないので、例えば次のように実行されると結果は{1, 0}になる。

```
y = 1
	r1 = y // r1は1
	r2 = x // r2は0
x = 1
```

上記では、スレッド1のストアの順序が入れ替わっているので、結果が{1, 0}になることがありえてしまう。

TSOとSequential Consistencyなハードウェアの違いを考えるためには、次のようなプログラムが必要だ。

```
# スレッド1
x = 1
r1 = y

# スレッド2
y = 1
r2 = x
```

上記のプログラムを実行後、{r1, r2}が{0, 0}になることはあり得るかどうかを考える。

|ハードウェア|{r1, r2}が{0, 0}になることはあり得るか？|
|:---|:---|
|Sequential Consistency|No|
|TSO|Yes|
|RMO|Yes|

Sequential Consistentなハードウェアでは{0, 0}は発生しない。インターリーブの結果がどうなるかを考えてみれば理解できると思う。
TSOで{0, 0}になるのは、次のようなケースである。

```
r1 = y
	r2 = x
x = 1
	y = 1
```

両方のスレッドにおいて、ストアがロードの後に並び替えられている。TSOでは「ストアがロードの後に順序変更」は発生しうるので、これはTSO的に問題のない順序変更であるが、この場合結果は{0, 0}になってしまう。このプログラムは、TSOは逐次一貫性を満たさないことを示している。
RMOはTSOよりも弱いので、TSOで発生する順序変更はRMOでも起こりうる。すなわち、RMOも同様に{0, 0}になることがあり得る。

### メモリバリア

どういった順序変更が発生しうるのかはプロセッサによって、あるいはプロセッサがサポートするメモリモデルによって異なる。
こういった順序変更は、__メモリバリア__ (あるいは__メモリフェンス__) と呼ばれる命令を使うで明示的に禁止できる。例えば、上記で考えた次のプログラムについて、

```
# スレッド1
x = 1
r1 = y

# スレッド2
y = 1
r2 = x
```

以下のようにメモリバリアを行うことで、順序変更は発生しなくなる。

```
# スレッド1
x = 1
Memory_barrier
r1 = y

# スレッド2
y = 1
Memory_barrier
r2 = x
```

メモリバリア命令は文字通りバリアとして、バリア前後の命令の順序変更が行われないことを保証する。

### 寄り道: Dependent loads reordering

最後に、RMOやWMOでも発生しない問題であるDependent loads reordering (依存関係のあるロードの順序変更) について触れて、ハードウェア・メモリモデルについては終わりにする。

Wikipediaから引用した表によれば、AlphaプロセッサではDependent loads (依存関係のあるロード) の順序変更が許されている。これは、WMOやRMOなプロセッサでも発生しないため、Alphaプロセッサのメモリモデルはこの点でこれらよりも弱いと言える。
Alphaプロセッサ (Alpha21264ベースのプロセッサ) における依存関係のあるロードの順序変更とは、つまり以下のようなことである。

```
# 初期状態: p = &x, x = 1, y = 0

# スレッド1
y = 1
Memory_barrier
p = &y

# スレッド2
i = &p
```

Alpha21264ベースのプロセッサを備えたコンピュータ[^1]では、このプログラムを実行した結果iが0になり得る。
上記のプログラムでは、 `y=1` と `p=&y` はメモリバリアによって順序付けされている (=これらの命令には依存関係がある) 。スレッド2が `i` として `0` を読み取るためには、yが `0` である必要があるが、この依存関係によってそれはありえないように思われる。
これがどのように発生するかと言うと、次の通りである。

* 実行前: p = &x、x = 1、y = 0 かつ、スレッド2は `y=0` をキャッシュしている
* `y=1` が実行される
* スレッド2に対して、yのキャッシュのインバリデーションが送られる
* スレッド2へのインバリデーションは、スレッド2の「Probe queue」にキューイングされる。なお、この時点でスレッド1にはAckが返る
* スレッド1はAckを受け取ったのでメモリバリアを「通過」できる。 `p=&y` に向かう
* スレッド2が `i=&p` を実行。pをデリファレンスするとyが出て来るが、スレッド2はまだy=0だと思っている (インバリデーションはまだキューの中にあるから) ので、 `i=0` になる

これは、メモリバリアがあるにも関わらずスレッド2がスレッド1のyへの書き込みを読み取れないという点で興味深い事象である。この問題を解決するには、スレッド2のロードの前にProbe queueの中身をフラッシュする必要がある。[Reordering on an Alpha processor](http://www.cs.umd.edu/~pugh/java/memoryModel/AlphaReordering.html)によれば、Alpha21264ではメモリバリアのタイミングでキューのフラッシュが実行されるため、スレッド2にメモリバリアを追加するとよい、としている。

## Goにおけるメモリモデル

前章では、プロセッサに関するメモリモデル「ハードウェア・メモリモデル」について書いた。
筆者はGoプログラマなので、Goのメモリモデルはどのようなものかについて説明する。

### ハードウェア・メモリモデル vs ソフトウェア・メモリモデル

ハードウェアメモリモデルが示すのは、あるプロセッサによるメモリへの書き込みや読み取りに関する命令順序に関するルールであった。
Goやその他の高水準言語が提供するメモリモデルは、ハードウェア・メモリモデルと対比して「ソフトウェア・メモリモデル」と呼ばれる。
ソフトウェア・メモリモデルが規定するのは、マルチスレッドプログラムが共有メモリに対してアクセスする際のルールである。マルチスレッドプログラムといっても、Goではgoroutineを話題にしているし、JavaScriptはシングルスレッドなプログラミング言語なため[SharedArrayBuffer](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/SharedArrayBuffer)に対するルールとなっているなど、言語ごとの差異は存在する。

### ハードウェア・メモリモデルとソフトウェア・メモリモデルの接点

昨今のプロセッサは逐次一貫性を保証していないが、代わりにDRF-SCと呼ばれる「同期モデル」をサポートしている。DRF-SCはソフトウェアと密接に関係しており、プログラマが高水準言語をどのように書けばハードウェアにおいても逐次一貫した振る舞いになるのかを示している。
DRF-SCとは「Data-race-free Sequential Consistency」の略で、「ソフトウェアがデータ競合 (Data Race) を回避するならば、ハードウェアは逐次一貫しているように振る舞う」という考え方である。内容の詳細はオリジナルの論文である[Weak Ordering - A New Definition](http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.42.5567)を参照してもらいたい。
DRF-SCにおいて、ソフトウェアはatomicなどの同期プリミティブを用いて複数のスレッドを同期する。同期は「Happens-before (先行発生)」と呼ばれる関係を規定する。例えば、次のようなプログラムを考える。

<img width="276" alt="4" src="https://user-images.githubusercontent.com/60682957/153757270-030d7bb8-90a9-4f16-a3c8-483cfb9816e6.png">

`S(a)` は、 `a` という (アトミックな) 変数への同期命令を指す。	こういった同期がある時、同期の前の命令は同期の後の命令よりも先行して発生することが保証される。
「先行して発生している」とはそのまま、「順序変更が発生しない」ことを意味する。
つまり、プログラマは適切に同期を行えば、コンパイラが適切にメモリバリア命令を発行して順序変更がないことを保証するようになる。このルールがDRF-SCである。

### Goのメモリモデル

[Goのメモリモデル](https://go.dev/ref/mem)においては、次のことが書かれている。

* 「同期」を行うことで、操作と操作の間に「先行発生 (Happens-before)」を関係付けられること
* 同期を行う方法

まず、先行発生は次のように定義されている。

> If event e1 happens before event e2, then we say that e2 happens after e1. Also, if e1 does not happen before e2 and does not happen after e2, then we say that e1 and e2 happen concurrently.

筆者による日本語訳は以下の通り。

> イベントe1がイベントe2よりも先行発生する場合、e2はe1よりも後に発生している。

> また、もしe1がe2よりも先行発生しておらず、e2がe1よりも先行発生していなければ、e1とe2は並行で発生している。

例えば、ある変数 _v_ への書き込み _w_ を、読み取り _r_ が観測するためには、 _w_ が _r_ よりも先行発生してかつ、_w_ と _r_ の間に別の書き込み _w'_ がない ( _w_ よりも後に発生し、 _r_ よりも先行発生する書き込み _w'_ がない) ことが必要である。

さらに、次のようにある。

> When multiple goroutines access a shared variable v, they must use synchronization events to establish happens-before conditions that ensure reads observe the desired writes.

同期イベントを使うことで、先行発生を関係付けることができる。

先行発生を定義したい理由は、先行発生を示すことが機械語プログラムにおける順序変更の回避 (DRF-SCによる) にそのままつながるからである。
つまり私達Goプログラマは、同期イベントを使うことでプログラムの先行発生関係をコンパイラに伝えることができ、その結果コンパイラは適切にメモリバリア命令が発行できるので、意図しない順序変更を回避できる、というわけである。

ここまでくれば、後は同期イベントにどのようなものがあるかを見ていくだけである。同期というとatomicやミューテックスをイメージしてしまうが、実際にはほかにも色々書かれており、いくつか例を挙げる。

* 初期化
  - パッケージpがパッケージqをインポートしている場合、qのinit関数の完了はqのある関数の開始よりも先行発生する
  - main関数の開始は、全てのinit関数の完了よりも後に発生する
* goroutineの作成
  - 新しいgoroutineを開始する `go` 文は、goroutineの実行よりも先行発生する
  - 例えば次のプログラムで、hello関数を実行すると「hello, world」が出力される

```go
var a string

func f() {
	print(a)
}

func hello() {
	a = "hello, world"
	go f()
}
```

これら以外には、goroutineの破棄、チャネル、sync.Mutexとsync.RWMutex、sync.Onceについて書かれている。

これら以外にも、Goのメモリモデルは[アップデートが予定されている](https://research.swtch.com/gomm)。sync.MapやPoolがメモリモデルで言及されていないなどの問題があるが、[proposal: Go Memory Model clarifications · Issue #50590 · golang/go](https://github.com/golang/go/issues/50590)がクローズされているなど、今後どうなるのかは不明である。

## まとめ

ハードウェアメモリモデル及びDRF-SC、そしてGoのメモリモデルについて見てきた。メモリモデルというものは元々ハードウェアがアセンブリプログラマに何を保証するのかを示すものでしかなかった。
しかし、[フリーランチは終わり](http://www.gotw.ca/publications/concurrency-ddj.htm)、人々はマルチスレッドなプログラムを書かなければならなくなった。ここで、コンパイラを対象としたソフトウェア・メモリモデルが生まれた。

現代では、DRF-SCのおかげで私達のような普通のプログラマは適切な同期を行いさえすれば不可解なハードウェアの挙動に悩まされることはない。これは過去の研究者やハードウェアエンジニア、プログラミング言語の開発者などによる高度な抽象化によるものである。
こういった普段は意識することのない領域も調べてみると面白いので、この記事を書いた。最後の方は力尽きて駆け足になってしまったが、まだ理解できていない部分も多いので、引き続き興味を持って調べていきたい。

## 参考文献

* [research!rsc: Hardware Memory Models (Memory Models, Part 1)](https://research.swtch.com/hwmm)
* [research!rsc: Programming Language Memory Models (Memory Models, Part 2)](https://research.swtch.com/plmm)
* [research!rsc: Updating the Go Memory Model (Memory Models, Part 3)](https://research.swtch.com/gomm)
* [A Tutorial Introduction to the ARM and POWER Relaxed Memory Models](https://www.cl.cam.ac.uk/~pes20/ppc-supplemental/test7.pdf)
* [Memory model - cppreference.com](https://en.cppreference.com/w/cpp/language/memory_model)
* [The Free Lunch Is Over: A Fundamental Turn Toward Concurrency in Software](http://www.gotw.ca/publications/concurrency-ddj.htm)
* [concurrency - Why is a program with only atomics in SC-DRF but not in HRF-direct? - Computer Science Stack Exchange](https://cs.stackexchange.com/questions/29043/why-is-a-program-with-only-atomics-in-sc-drf-but-not-in-hrf-direct)
* [What's the relationship between CPU Out-of-order execution and memory order? - Stack Overflow](https://stackoverflow.com/questions/70749012/whats-the-relationship-between-cpu-out-of-order-execution-and-memory-order)
* [x86 - Are memory barriers needed because of cpu out of order execution or because of cache consistency problem? - Stack Overflow](https://stackoverflow.com/questions/63970362/are-memory-barriers-needed-because-of-cpu-out-of-order-execution-or-because-of-c)
* [Memory Barriers Are Like Source Control Operations](https://preshing.com/20120710/memory-barriers-are-like-source-control-operations/)
* [Memory barrier - Wikipedia](https://en.wikipedia.org/wiki/Memory_barrier)
* [Memory ordering - Wikipedia](https://en.wikipedia.org/wiki/Memory_ordering)
* [The Go Memory Model - The Go Programming Language](https://go.dev/ref/mem)
* [Bridging the gap in the RISC-V memory models](https://riscv.org/wp-content/uploads/2019/06/16.15-Stefanos-Kaxiras.pdf)
* [RISC-V Memory Consistency Model Tutorial](https://riscv.org/wp-content/uploads/2018/05/14.25-15.00-RISCVMemoryModelTutorial.pdf)
* [How to Make a Multiprocessor Computer That Correctly Executes Multiprocess Programs - Microsoft Research](https://www.microsoft.com/en-us/research/publication/make-multiprocessor-computer-correctly-executes-multiprocess-programs/)
* [Intel® 64 Architecture Memory Ordering White Paper](http://www.cs.cmu.edu/~410-f10/doc/Intel_Reordering_318147.pdf)
* [Peter Sewell: bibliography by topic](https://www.cl.cam.ac.uk/~pes20/papers/topics.html)
* [Weak Ordering - A New Definition](http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.42.5567)
* [Updating the Go memory model · Discussion #47141 · golang/go](https://github.com/golang/go/discussions/47141)
* [doc: define how sync interacts with memory model · Issue #7948 · golang/go](https://github.com/golang/go/issues/7948)

[^1]: Linus Torvaldsが、Alphaの中でもごく一部のハードウェアでしか起きないよ、とどこかで言っていたらしい https://stackoverflow.com/questions/35115634/dependent-loads-reordering-in-cpu#comment57952162_35115634
