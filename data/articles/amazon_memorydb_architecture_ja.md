type: blog
title: Amazon MemoryDB for Redisはどのように耐久性を保証しているか
timestamp: 2021-11-19 19:00:00
lang: ja
---

Note: This article is available in [English](https://hidetatz.io/articles/2021/11/19/amazon_memorydb_architecture/) also.

## はじめに

[Amazon MemoryDB for Redis](https://aws.amazon.com/memorydb/)は2021年にAWSが発表したRedis互換の分散型マネージドデータベースである。
これまでAWSはマネージドRedisとしてAmazon Elasticacheを提供していたが、Amazon MemoryDBはElasticacheにはない分散環境におけるデータの一貫性・耐久性を機能として提供している。
この記事ではRedisの持つ一貫性・耐久性の仕組みを概観しながら、Amazon MemoryDBがどのように分散環境での耐久性保証を実現しているかについて、現状提供されている情報を基にまとめていく (参考とした情報は本記事末尾で列挙する) 。

この記事ではレプリケーションにおける「マスター・レプリカ」モデルを扱うが、一貫してそれぞれを「プライマリ」「レプリカ」と呼称する。

## Redisクラスタにおけるデータの一貫性

まず、Amazon MemoryDBではなく普通のRedisクラスタの話から。

Redisクラスタは[強い一貫性を保証しない](https://redis.io/topics/cluster-tutorial#redis-cluster-consistency-guarantees)。
Redisのレプリケーションは[結果整合性](http://antirez.com/news/36)で行われる。
すなわち、プライマリが書き込みを受け付けると、まずはクライアントにレスポンスが返り、その後非同期でレプリカに変更が反映される。
従って、プライマリへの書き込みとレプリケーションの間にプライマリに障害が発生しフェイルオーバーした場合、プライマリへの書き込みはロストする。
これによって、Redisは素早くレスポンスすることができるが、その分障害が発生した際に強い一貫性は保証できなくなる。

<img width="764" alt="1" src="https://user-images.githubusercontent.com/60682957/142567317-2dc28ff9-371e-4aaa-859a-806ff32e587f.png">

> [Amazon MemoryDB for Redis – Where speed meets consistency](https://www.allthingsdistributed.com/2021/11/amazon-memorydb-for-redis-speed-consistency.html)より引用

### Redisにおけるデータ永続化
Redisにはネイティブに[RDBやAOF](https://redis.io/topics/persistence)といった永続化の仕組みが備わっているが、そのどれもが強い一貫性を保証するものではない。
RDBは単なる定期的なスナップショット取得であるため、最後のスナップショット取得から現在時刻までの書き込みは常に永続化されていない。

AOFとは、「Append Only File」のイニシャリズムで、これはいわゆるトランザクションログである。
すなわち、書き込みは全てオンメモリのデータにコミットされつつ、その操作自体がログファイルに記録される。
重要なのは、データのコミットとログの記録は同期で行われることである。
ノードに障害が発生した時は当然メモリ上のデータはロストするが、ログファイルはファイルシステムに永続化されているので (ディスクがクラッシュしない限りは) 消えることはない。
すなわち、Redisのプロセスを再起動し、AOFの中身を先頭から末尾まで順に再適用していけば、データロストなく障害の直前までのデータが復旧可能である。

<img width="683" alt="2" src="https://user-images.githubusercontent.com/60682957/142567324-65aa25ee-99ef-4e0a-b3ee-5722d1267510.png">

> [Amazon MemoryDB for Redis – Where speed meets consistency](https://www.allthingsdistributed.com/2022/11/amazon-memorydb-for-redis-speed-consistency.html)より引用

このやり方にもデメリットがある。Redisが保持するデータセットが極めて大きい場合、そもそもAOFに記録された操作を全て再適用すること自体に時間がかかる。
また、単純に書き込み時に、AOFに追記する分のパフォーマンス上のオーバーヘッドが (ディスクへのランダムアクセスは発生しないにせよ) 避けられない。

AOFは、クラスタリングとの相性が良くない。Redisクラスタはプライマリ・レプリカモデルであり、フェイルオーバーの際はレプリカがプライマリに昇格する。
しかし前述の通り、フェイルオーバー時はいくつかの書き込みがレプリカに反映されておらずデータがロストするかもしれない。
じゃあAOFからデータを復旧すればいいかというと、AOFはあくまでノードのローカルに持つものなので、レプリカがプライマリのAOFを参照することはできない。
データのロストを避けたければクラスタリングを辞めて障害時にAOFから復旧するようにすればよいが、これではAOFの再適用に時間がかかった場合可用性を阻害する。

Redisでは、 (クラスタリングを行うかどうかに関わらず) 高い可用性とデータの一貫性、耐久性を同時に保証することができない。
クラスタリングは可用性をもたらすが、一部のデータのロストを許容しなければならない。
AOFはデータの耐久性を保証するが、場合によってはRedisの可用性を落とす。
しかし、あくまでキャッシュストアであるというRedisの特性上、こういった欠点は歴史的に許容されてきた。
Redisはあくまでデータのキャッシュを持つだけであり、本来のデータは (データの元となるデータベースなどから) 復旧が可能なためである。

さて、ここまで読むと、AOFをレプリカからアクセス可能にすれば色々解決するじゃん？ということに気づくだろう。
Amazon MemoryDBが実現したのはまさにこれである。

## Amazon MemoryDBのアーキテクチャ

<img width="1097" alt="3" src="https://user-images.githubusercontent.com/60682957/142567327-26d34c4f-2495-4b0e-a422-baac5ba90863.png">

> [Getting Started with Amazon MemoryDB for Redis - AWS Online Tech Talks - YouTube](https://www.youtube.com/watch?v=Jbq_XZMZEKY&ab_channel=AWSOnlineTechTalks)より引用

Amazon MemoryDBは、プライマリが書き込みをオンメモリにコミットするのと同期で、「Multi-AZ Transaction Log」を書き込んでいる。
Multi-AZ Transaction Logとは、詳細は明かされていないが、おそらくプライマリとはネットワーク的に分離したサーバに存在するログであると思われる。
これはマルチAZでかつ同期で書き込みされるため、高い耐久性を実現しながらオンメモリデータとMulti-AZ Transaction Logの間の一貫性は保証される。

レプリカへの書き込みは通常のRedisクラスタと同様に非同期で行われる。実際は、Multi-AZ Transaction Logへの書き込みが伝播するようである。

### フェイルオーバー

<img width="1094" alt="4" src="https://user-images.githubusercontent.com/60682957/142567334-559b4050-7b93-4837-8369-305ae902b149.png">
<img width="1093" alt="5" src="https://user-images.githubusercontent.com/60682957/142567342-cb3bb220-1066-4fbe-a621-c276cdabf2f2.png">

> どちらも[Getting Started with Amazon MemoryDB for Redis - AWS Online Tech Talks - YouTube](https://www.youtube.com/watch?v=Jbq_XZMZEKY&ab_channel=AWSOnlineTechTalks)より引用

フェイルオーバーの際は、レプリカがプライマリに昇格する。しかし、フェイルオーバープロセスの中で、プライマリに昇格するレプリカはMulti-AZ Transaction Logのうち自身に適用されていないオペレーションだけを適用する。
ここの詳細は公開されていないが、例えば以下のようなアルゴリズムが考えられる。

* Amazon MemoryDB的に、全ての書き込みオペレーションはインクリメンタルなIDを持つ (例えば `uint64` )
* Multi-AZ Transaction Logはすべての行にオペレーションIDを保持する
* Multi-AZ Transaction Logからレプリカに非同期で書き込まれる際、オペレーションIDも受け取っており、各レプリカは自分がどのオペレーションIDまで適用しているかを知っている
* フェイルオーバー時は、プライマリに昇格するレプリカはMulti-AZ Transaction Logの最後の行 (最も最近の書き込みオペレーション) のオペレーションIDを確認し、自身のオペレーションIDとの差分を計算する。例えば差分がN行であれば、Multi-AZ Transaction Logの末尾のN行を取得し、N行のみを自身に適用する
* 適用が完了したら、プライマリとしてクライアントからのリクエストの受付を開始する

AOFではファイルの先頭から末尾までを全て再適用しないと完全なデータが再現できなかった。上記のアプローチでは、適用されていないオペレーションのみを適用可能なため、データの復旧にかかる時間が大きく削減できると思われる。これによって、Amazon MemoryDBは高い可用性と一貫性・耐久性を同時に実現できる。

## MySQLクラスタとの類似点

さて、こういった技術は特に新規性のあるものではない。
例えば、MySQLクラスタをフェイルオーバーするとき、MySQLのレプリケーションのBinlogを使って非同期で行われるので、結局はソースとレプリカでの差分が発生する。
RDBは通常、キャッシュストアよりも一貫性にはシビアである。フェイルオーバーの際は結局はレプリカへの差分データのコピーを行う必要があり、これを自動で行うのが[MHA](https://github.com/yoshinorim/mha4mysql-manager/wiki)や[mysqlfailover](https://www.percona.com/blog/2014/07/03/failover-mysql-utilities-part-2-mysqlfailover/)である。 ([orchestrator](https://github.com/openark/orchestrator)もおそらく同じ？)
これはまさしくAmazon MemoryDBが実現していることである。

## 終わりに

AOFはノードローカルに保有するので、ネットワーク通信が不要な分書き込みの失敗の確率がかなり低い。Multi-AZ Transaction LogはAWS内のネットワークとはいえ書き込み失敗の確率が比較的高いはずだが、リトライなどで解決可能だと思われる。
筆者は[Amazon MemoryDB for Redis – Where speed meets consistency](https://www.allthingsdistributed.com/2021/11/amazon-memorydb-for-redis-speed-consistency.html)を読んでいてMemoryDBのアーキテクチャを知った。一通り読んで、「これMySQLクラスタのフェイルオーバー自動化ツールがやってるやつじゃん」と思うと同時に、こういった既に世の中にある仕組みを流用して何かを少し便利にするようなエンジニアリングはクールだな、と感じこの記事を書いた。Amazon MemoryDBはElasticacheよりもコストがかかりそうだが、今後は基本的にAWSでRedis使う場合はMemoryDBを採用するのが良いではないだろうか。

## 参考

* [Amazon MemoryDB for Redis – Where speed meets consistency](https://www.allthingsdistributed.com/2021/11/amazon-memorydb-for-redis-speed-consistency.html)
* [Getting Started with Amazon MemoryDB for Redis - AWS Online Tech Talks - YouTube](https://www.youtube.com/watch?v=Jbq_XZMZEKY&ab_channel=AWSOnlineTechTalks)
* [AWS Announces the General Availability of Amazon MemoryDB for Redis](https://www.infoq.com/news/2021/08/amazon-memorydb-for-redis-ga/)
* [Redis Persistence – Redis](https://redis.io/topics/persistence)
* [Minimizing downtime in ElastiCache for Redis with Multi-AZ - Amazon ElastiCache for Redis](https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/AutoFailover.html)
* [Redis data model and eventual consistency - <antirez>](http://antirez.com/news/36)
* [Yoshinori Matsunobu's blog: Announcing MySQL-MHA: "MySQL Master High Availability manager and tools"](http://yoshinorimatsunobu.blogspot.com/2011/07/announcing-mysql-mha-mysql-master-high.html)
* [openark/orchestrator: MySQL replication topology management and HA](https://github.com/openark/orchestrator)
* [MySQL :: MySQL 8.0 Reference Manual :: 17.1.3.5 Using GTIDs for Failover and Scaleout](https://dev.mysql.com/doc/refman/8.0/en/replication-gtids-failover.html)
* [Failover with the MySQL Utilities: Part 2 - mysqlfailover](https://www.percona.com/blog/2014/07/03/failover-mysql-utilities-part-2-mysqlfailover/)
