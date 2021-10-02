https://www.alexedwards.net/blog/configuring-sqldb---2021-10-01 10:00:00


Configuring sql.DB for Better Performance
https://www.alexedwards.net/blog/configuring-sqldb

* sql.DBはデータベースコネクションのプールである
    * 使用中のコネクションと、アイドルのコネクションが含まれている
* sql.DBを使う時、まずアイドルなコネクションがプールに存在するかを調べる
    * あれば、それを使う
    * なければ、コネクションを新たに張ってそれを使う
* SetMaxOpenConnsとは？
    * プールの中のコネクション数にはデフォルトでは制限がない
        * SetMaxOpenConnsはin-useとidleのコネクションの数の合計を制限する
            * なので、例えばSetMaxOpenConns(5)して、5つコネクションが全てin-useな状態で新たなコネクションが必要な場合、どれかがidleになるのを待つ。新たにコネクションを張らない。
            * 上の記事によればこう (コードはこれ https://gist.github.com/alexedwards/5d1db82e6358b5b6efcb038ca888ab07)

```
BenchmarkMaxOpenConns1-8                 500       3129633 ns/op         478 B/op         10 allocs/op
BenchmarkMaxOpenConns2-8                1000       2181641 ns/op         470 B/op         10 allocs/op
BenchmarkMaxOpenConns5-8                2000        859654 ns/op         493 B/op         10 allocs/op
BenchmarkMaxOpenConns10-8               2000        545394 ns/op         510 B/op         10 allocs/op
BenchmarkMaxOpenConnsUnlimited-8        2000        531030 ns/op         479 B/op          9 allocs/op
PASS
```

* SetMaxIdleConnsとは？
    * プール内のアイドルなコネクションの数を制限する。デフォルトでは2
    * 理論上、アイドルコネクションの数が多ければパフォーマンスは向上するはずである
    * 上の記事によればこう

```
BenchmarkMaxIdleConnsNone-8          300       4567245 ns/op       58174 B/op        625 allocs/op
BenchmarkMaxIdleConns1-8            2000        568765 ns/op        2596 B/op         32 allocs/op
BenchmarkMaxIdleConns2-8            2000        529359 ns/op         596 B/op         11 allocs/op
BenchmarkMaxIdleConns5-8            2000        506207 ns/op         451 B/op          9 allocs/op
BenchmarkMaxIdleConns10-8           2000        501639 ns/op         450 B/op          9 allocs/op
PASS
```

* MaxIdleConnを0にするといちいちコネクションを貼り直すので極めて遅い
* 1にするだけでもずいぶん変わった
* じゃあアイドルコネクションをたくさん確保すれば良いのか？
    * 場合による
    * アイドルコネクションを確保するコストはメモリを食ってしまうこと
    * アイドルコネクションは使われていないとはいえ結局データベースとつながりっぱなしではあるので
* コネクションがあまりに長い時間アイドルだったら切る、ということもできる
    * 例えばMySQLのwait_timeoutがそれを示す (デフォルトでは8h)
    * sql.DBの場合、コネクションはgracefullyにハンドルされる。コネクションがサーバーサイドから切られてもGoはそれを知るはずがないので、死んだコネクションを使おうとしたら2回リトライしたあとそれをプールから取り除き、新たにコネクションを張るような動きになる
    * なので、アイドルコネクションを持ちすぎるとより多くのリソースが使われてしまうかも。
    * 本当に使う分のアイドルコネクションだけ持つのが望ましい
* MaxIdleConnsはMaxOpenConnsより小さくなければならない (当たり前のことだが)
* SetConnMaxLifetime(d time.Duration)は、コネクションがどれだけの間使用可能かを定義する
* db.SetConnMaxLifetime(time.Hour)とすると、作られてから1時間でコネクションは 'expire' する
    * コネクションが1時間残り続けることを保証するものではない
    * コネクションは1時間経っても使われ続けていることはあり得る。1時間以上経ってから使用開始されることはないが。
    * アイドルになってから1時間ではなく、作られてから1時間
    * expireしたコネクションをプールから除外するのは1秒に1回実行される
* ConnMaxLifetimeが短いと、すぐにexpireして、より多くのコネクションが新たに張られなければならなくなる
* 上の記事によればこう

```
BenchmarkConnMaxLifetime100-8               2000        637902 ns/op        2770 B/op         34 allocs/op
BenchmarkConnMaxLifetime200-8               2000        576053 ns/op        1612 B/op         21 allocs/op
BenchmarkConnMaxLifetime500-8               2000        558297 ns/op         913 B/op         14 allocs/op
BenchmarkConnMaxLifetime1000-8              2000        543601 ns/op         740 B/op         12 allocs/op
BenchmarkConnMaxLifetimeUnlimited-8         3000        532789 ns/op         412 B/op          9 allocs/op
PASS
```

* MaxConnLifetimeを設定する場合、コネクションがどのくらい頻繁にexpire -> recreateされるのかを意識する
* 100のコネクションがあって1分のlifetimeだと、1秒に1.67のコネクションがexpireする。これがあまりに大きいとパフォーマンスが悪化しうる

http://dsas.blog.klab.org/archives/2018-02/configure-sql-db.html

ポイント

* MySQLではwait_timeoutで接続がサーバから切られた場合、MySQLドライバはそれに気づかず、クエリを終わったコネクションで送ろうとしてしまう
* SetConnMaxLifetimeを短めに設定しておけばそういうことが起こらない
* 副次的な効果は:
    * コネクションがすぐ切れるのでサーバの増減しやすくなる
    * DBのフェイルオーバーがしやすくなる
    * MySQLをオンラインで設定変更しても古い設定のコネクションが残りにくく鳴る
* SetMaxOpenConnsは必ず設定する。負荷が大きい時に新規コネクションを作らない程度の値にする
* SetMaxIdleConnsはおそらくSetMaxOpenConnsと同じにしておけば良いと思う、減らす理由がない
* SetConnMaxLifetimeは、何秒に1回再接続するか？で考える。1秒に1回程度なら殆どの場合負荷は問題にならない
