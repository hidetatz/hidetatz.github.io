秘匿情報をメモリ上でどう扱うか---2021-04-23 19:28:13

**この記事はメモ書きに近いです。この記事に書かれたような内容を実装する際はセキュリティエンジニアと相談の上行うことをおすすめいたします**

---

パスワード・暗号化に使用する鍵などの秘匿情報は、ディスクに書き残すべきでないだけでなく、メモリにもなるべく残さないことが望ましい。理由は:

* プロセスがクラッシュした際にコアダンプされ結果的にディスクに書かれる可能性がある
* ページアウトやスワップ時にディスクに書かれる可能性がある

などが挙げられる。
ディスクに書き残されてしまうことで、攻撃者に読みだされてしまうおそれがある。

実装したいプログラムにもよるが、これらを完全にメモリに載せないまま目的を達成する (認証を行う、暗号化を行うなど) ことが不可能であることも多いだろう。このような場合にアプリケーション開発者が取れる手段には次のようなものがある。

* コアダンプを無効にする (コアダンプ自体の無効もしくはコアダンプから秘匿情報を除外)
* 機密情報が入っている部分のメモリをスワップさせないようにする
* 機密情報は必要な時以外暗号化しておき、使い終わったらメモリをクリーンアップ (e.g. 0で埋める) する。または即GCさせる

以上が基本的なテクニックだが、これをどう実現するか。

### コアダンプの無効

コアダンプの無効はLinuxやBSDではsetrlimitを使ってRLIMIT_COREに対しsoft limit、hard limit共に0を指定する。これは正確に言えばコアダンプの最大サイズをセットする機能だが、これにゼロを設定することで事実上コアダンプ自体の無効化ができる。

Goだと、次のようなコードになる。 ( `x/sys/unix` を使っているが、もちろんcgoを使っても良いだろう)

```go
import "golang.org/x/sys/unix"

err := unix.Setrlimit(unix.RLIMIT_CORE, &unix.Rlimit{Cur: 0, Max: 0})
```

コアダンプを無効にすることでクラッシュ時にメモリの状況を知ることは当然できなくなるので、デバッグ時に困る可能性がある。そういった場合は、コアダンプを無効にするのではなく、特定のメモリだけスワップアウトの対象から外す (厳密には外すようカーネルに "advise" する) やり方もある。これは、madviseを使って `advise` に `MADV_DONTDUMP` を渡す。Goでは以下のように書ける。

```go
unix.Madvise(some_secret_var, unix.MADV_DONTDUMP)
```

注意すべきは、[madvise(2) - Linux man page](https://linux.die.net/man/2/madvise) にあるように、 `The kernel is free to ignore the advice.` すなわち、カーネルはこのアドバイスを無視することもあるらしい。

また、そもそもコアダンプはアプリケーションからでなくOS事態に設定も可能である。

### スワップの抑止

`mlock` を使えばメモリをロックできるので、メモリがスワップアウト (ページアウト) されるのを防ぐことができる。Goでは次のようなコードだ。

```go
unix.Mlock(some_secret_var)
```

また、[hashicorp/vault](https://github.com/hashicorp/vault/blob/c44f1c9817955d4c7cd5822a19fb492e1c2d0c54/helper/mlock/mlock_unix.go#L15-L17)のように `mlockall` を使ってプロセス上のすべてのメモリをロックする方法もある。この辺はアプリケーションの性質にもよるだろう。

`Mlock` があるということは `Munlock` もあるのか？というと、ある。mlockには `RLIMIT_MEMLOCK` というロック可能なメモリの制限があり、それを超えると追加でのmlockができなくなってしまう。したがって、ロックの必要がなくなったメモリは速やかにアンロックすることが望ましい。

コアダンプ同様、スワップも `swapoff` コマンドでカーネルの設定から無効にすることができる。

### メモリ上での秘匿情報の扱い

HeartbleedやCold Boot Attackなどによる不正なメモリ読み取りから秘匿情報を守る方法のひとつに、秘匿情報をメモリ上で扱う際暗号化したままにしておき、必要な時だけ複合化することが挙げられる。しかしこれは、使用するキーの扱いなどで複雑性を生みやすく、個人的にはあまり採用されないように思われる。

これに対し、使い終わった秘匿情報をGCを待たずに、明示的にゼロで埋めるなどしてパージするのはまだ見る実装だろう。例えばJavaの[Secure Coding Guidelines for Java SE
](https://www.oracle.com/java/technologies/javase/seccodeguide.html)の `Guideline 2-3 / CONFIDENTIAL-3: Consider purging highly sensitive from memory after use` にも同様の手法が紹介されている。

Cではメモリの扱いが良くも悪くも直感的なので、単純に0を書き込むなり、 `free` してしまえばよい。しかし、Goにおいてはこれはそう単純な話ではないようである。例えば、Goの言語仕様によれば、Goのガベージコレクタはメモリのフラグメンテーションを回避するためにメモリを動かしたりコピーしたりすることを禁止していない。Goのランタイムが管理するメモリは、Goのランタイムが自由に操作することができるので、その場合単純な `mlock` は不十分である。例えば変数がランタイムによってコピーされたとき、コピー後のメモリは `mlock` のロック対象にならない。

この場合どうすればよいか。Goでこれをやるには、 `mmap` や `munmap` で自前でメモリを管理する必要がある。これで確保したメモリはガベージコレクタの観測範囲外になるので、これに対して `mlock` を行ったり、使い終わった後のゼロ埋めを行えばよい。また言うまでもないが、使った後は自分でメモリを解放する必要がある。

余談だが、自前でメモリを管理するにあたっては解放を忘れないことだけでなく、 `mprotect` を使ったガードページの実装や、カナリア領域の確保など、バッファオーバーフロー・アンダーランを防ぐ機構もおそらく必要になるだろう。

## 終わりに

メモリ上の秘匿情報の保護は煩雑で、コードを複雑化する。パブリッククラウドでアプリケーションを構築することの多い現代では、きちんとTLS化やストレージ暗号化、ネットワークの設定などでなるべく安全性を保つことがアプリをシンプルにするコツだと思う。それでも、セキュリティ要件によってはこういったテクニックが必要になることもあるので、覚えておくと役に立つことがあるかもしれない。

## 参考

* [memory security in go](https://spacetime.dev/memory-security-go)
* [Secure Coding Guidelines for Java SE](https://www.oracle.com/java/technologies/javase/seccodeguide.html)
* [Java Cryptography Architecture (JCA) Reference Guide](https://docs.oracle.com/javase/8/docs/technotes/guides/security/crypto/CryptoSpec.html)
* [github.com/awnumar/memguard](https://github.com/awnumar/memguard)
* [Memguard must manage memory allocation to work #3](https://github.com/awnumar/memguard/issues/3)
* [Hacker News - Securely Handle Encryption Keys in Go](https://news.ycombinator.com/item?id=14174500)
* [Clearing secrets from memory](https://www.sjoerdlangkemper.nl/2016/05/22/should-passwords-be-cleared-from-memory/)
* [github.com/hashicorp/vault](https://github.com/hashicorp/vault)
