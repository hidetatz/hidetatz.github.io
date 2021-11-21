分散システムについて学ぶためのリソースを単純に列挙する。全てを読む必要はなくて (無理だと思う) 一通りさらって見て気になったのを読み込むのが良いと思う。
随時更新。

## 書籍: 

* [分散システム (第2版) (未来へつなぐデジタルシリーズ) | 忠則, 水野 |本 | 通販 | Amazon](https://www.amazon.co.jp/dp/4320124499)
* [分散処理システム (情報工学レクチャーシリーズ) | 真鍋 義文 |本 | 通販 | Amazon](https://www.amazon.co.jp/dp/4627810717)
  * 大学で使う教科書。かなり基礎的で良い。
* [分散システム　第二版 | アンドリュー・S・タネンバウム, マールティン・ファン・スティーン, Andrew S. Tanenbaum, Maarten van Steen, 水野 忠則, 佐藤 文明, 鈴木 健二, 竹中 友哉, 西山 智, 峰野 博史, 宮西 洋太郎 |本 | 通販 | Amazon](https://www.amazon.co.jp/dp/4894714981)
  * 所謂タネンバウム本。個人的にはさらっと読むくらいで良いと思う
* [Amazon | Replication: Theory and Practice (Lecture Notes in Computer Science / Theoretical Computer Science and General Issues) (Lecture Notes in Computer Science, 5959) | Charron-Bost, Bernadette | Cryptography](https://www.amazon.co.jp/dp/3642112935)
  * 個人的に一番おすすめはこれ (ただし英語)
* [ZooKeeperによる分散システム管理 | Flavio Junqueira, Benjamin Reed, 中田 秀基 |本 | 通販 | Amazon](https://www.amazon.co.jp/dp/4873116937)
  * ZooKeeper本なんだけど、ZooKeeperを通して分散システムについて学べるので良いと思う
* [Amazon | Designing Data-Intensive Applications: The Big Ideas Behind Reliable, Scalable, and Maintainable Systems | Kleppmann, Martin | Accounting](https://www.amazon.co.jp/dp/1449373321)
  * Designing Data-Intensive Applicationsもかなり良い本
* [佐藤一郎: Web日記 (2010年)](http://home.att.ne.jp/sigma/satoh/diary/diary100331.html#20100102)
  * 後はこちらのブログを読んで気になった本を読むとか

## Webサイト

* [Readings in distributed systems](http://christophermeiklejohn.com/distributed/systems/2013/07/12/readings-in-distributed-systems.html)
* [papers-we-love/distributed_systems at master · papers-we-love/papers-we-love](https://github.com/papers-we-love/papers-we-love/tree/master/distributed_systems)
* [system-design-primer/README-ja.md at master · donnemartin/system-design-primer](https://github.com/donnemartin/system-design-primer/blob/master/README-ja.md)
  * system-design-primer、分散システムに特化はしていないが、一通り読んでおくのがおすすめ。
* [Patterns of Distributed Systems](https://martinfowler.com/articles/patterns-of-distributed-systems/)
  * リファレンス的に使う。
* [The Amazon Builders' Library](https://aws.amazon.com/jp/builders-library)
  * AmazonがAWS作る時に得た知見をシェアするサイトで、かなりおすすめ。
* [Designs, Lessons and Advice from Building Large Distributed Systems by Jeff Dean](http://www.cs.cornell.edu/projects/ladis2009/talks/dean-keynote-ladis2009.pdf)
  * GoogleのJeff Deanさんが書いたスライド
* [Consistency Models](https://jepsen.io/consistency)
  * jepsen.ioのconsistency modelsの章。これだけは最低限読んでおいたほうが良いと思う
* [分散システムについて語るときに我々の語ること ― 分散システムにまつわる重要な概念について | POSTD](https://postd.cc/learning-about-distributed-systems/)
  * MQTT作った人の講演を文字に興したもの (たぶん) 。とっかかりはこれが良いかも。
* [baseds – Medium](https://medium.com/baseds)
  * たぶん学生の人が書いてるんだと思うけど、絵が可愛いので良い。
* [aphyr/distsys-class: Class materials for a distributed systems lecture series](https://github.com/aphyr/distsys-class)
  * jepsen.ioが提供する分散システムの講義のマテリアルだと思う。基礎的な内容
* [How We've Scaled Dropbox - YouTube](https://www.youtube.com/watch?v=PE4gwstWhmc&ab_channel=Stanford)
  * Dropboxのアーキテクチャについての動画。かなりおすすめ
* [pingcap/talent-plan: open source training courses about distributed database and distributed systemes](https://github.com/pingcap/talent-plan)
  * PingCapが提供してる分散システムを学ぶためのやつ
* [分散システムについて語らせてくれ](https://www.slideshare.net/kumagi/ss-78765920)
  * 故障モデルのところわかりやすい

## Paper

* [Scaling Memcache at Facebook](https://cs.uwaterloo.ca/~brecht/courses/854-Emerging-2014/readings/key-value/fb-memcached-nsdi-2013.pdf)
  * FacebookのMemcachedアーキテクチャ
* [Amazon Aurora: Design Considerations for High Throughput Cloud-Native Relational Databasesを読む(その1 Introduction) - CubicLouve](https://spring-mt.hatenablog.com/entry/2021/03/01/123934)
  * AWS Auroraのアーキテクチャ、難しい
* [SQLServer 2014 “Hekaton”再考 - 急がば回れ、選ぶなら近道](https://okachimachiorz.hatenablog.com/entry/20170918/1505729999)
  * Hekatonの話、かなり難しい

## OSS

* [geohot/minikeyvalue: A distributed key value store in under 1000 lines. Used in production at comma.ai](https://github.com/geohot/minikeyvalue)
  * comma.ai が本番で使っている、1000行以下で書かれた分散KVS
* [chrislusf/vasto: A distributed key-value store. On Disk. Able to grow or shrink without service interruption.](https://github.com/chrislusf/vasto)
  * 極めて基本的な (たぶん) 構成で分散KVS
  * 分散システムのことわかってくると分散KVS作りたくなるっぽくて、調べるとたくさん出てくる
* [go-redsync/redsync: Distributed mutual exclusion lock using Redis for Go](https://github.com/go-redsync/redsync)
  * Redisで分散ロックするやつ
  * 分散ロックとはなにかというのはこれを読む
    * https://kumagi.hatenablog.com/entry/distributed_lock

## その他色々

* [Tail Latencyに関する論文読み - 勘と経験と読経](https://agnozingdays.hatenablog.com/entry/2019/04/25/210721)
* [今度こそ絶対あなたに理解させるPaxos - Qiita](https://qiita.com/kumagi/items/535c9b7a761d2ed52bc0)
  * この方の記事は全部読むと良いと思う
* [Rendezvous Hashingについての雑感. 先日、 Rendezvous… | by sile | Medium](https://medium.com/@reduls/about-rendezvous-hashing-76aaf1a3f705)
* [CRDT (Conflict-free Replicated Data Type)を15分で説明してみる - Qiita](https://qiita.com/everpeace/items/bb73ec64d3e682279d26)
* [Apache Mesos - The Mesos Replicated Log](http://mesos.apache.org/documentation/latest/replicated-log-internals/)
* [Consistency levels in Azure Cosmos DB | Microsoft Docs](https://docs.microsoft.com/en-us/azure/cosmos-db/consistency-levels)
* [実践TLA+（Hillel Wayne 株式会社クイープ 株式会社クイープ）｜翔泳社の本](https://www.shoeisha.co.jp/book/detail/9784798169163)
  * TLA+の本。最初に読む本ではない

## 並行並列プログラミング

たぶんやってると並行並列プログラミングについてだんだん学びたくなってくる (と思う) ので、そのへんも

* [O'Reilly Japan - 並行プログラミング入門](https://www.oreilly.co.jp/books/9784873119595/)
  * 素晴らしい本
* [一人トランザクション技術 Advent Calendar 2016 - Qiita](https://qiita.com/advent-calendar/2016/transaction)
  * トランザクション
* [research!rsc: Hardware Memory Models (Memory Models, Part 1)](https://research.swtch.com/hwmm)
  * rscさんが書いたメモリモデルのハードウェアの方
* [「強いメモリモデル」と「弱いメモリモデル」 - yamasaのネタ帳](https://yamasa.hatenablog.jp/entry/2020/11/29/171322)
  * これもハードウェアメモリモデルの話
* [メモリモデル？なにそれ？おいしいの？ - yohhoyの日記（別館）](https://yohhoy.hatenablog.jp/entry/2014/12/21/171035)
  * これはソフトウェアメモリモデル。C++をベースにしている
* [Lock-free入門](https://kumagi.com/lockfree.pdf)
  * Lockfree
* [Ruby に Software Transactional Memory (STM) を入れようと思った話 - クックパッド開発者ブログ](https://techlife.cookpad.com/entry/2020/11/20/110047)
  * STM
* [Survey of Transactional Memory - Speaker Deck](https://speakerdeck.com/ytakano/survey-of-transactional-memory)
  * これもSTM
