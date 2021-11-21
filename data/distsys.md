分散システムについて学ぶためのリソースを列挙する。全てを読む必要はおそらくなく、一通りさらって見て気になったのを読み込むのが良いと思う。

## 書籍: 

* 大学で使う教科書。かなり基礎的で良い。
  * https://www.amazon.co.jp/dp/4320124499
  * https://www.amazon.co.jp/dp/4627810717
* 所謂タネンバウム本。個人的にはさらっと読むくらいで良いと思う
  * https://www.amazon.co.jp/dp/4894714981
* 個人的に一番おすすめはこれ (ただし英語)
  * https://www.amazon.co.jp/dp/3642112935
* ZooKeeper本なんだけど、ZooKeeperを通して分散システムについて学べるので良いと思う
  * https://www.amazon.co.jp/dp/4873116937
* Designing Data-Intensive Applicationsもかなり良い本
  * https://www.amazon.co.jp/dp/1449373321
* 後は以下のブログを読んで気になった本を読めば良いと思う 
  * http://home.att.ne.jp/sigma/satoh/diary/diary100331.html#20100102

## Webサイト

* system-design-primer、分散システムに特化はしていないが、一通り読んでおくのがおすすめ。
  * https://github.com/donnemartin/system-design-primer/blob/master/README-ja.md
* AmazonがAWS作る時に得た知見をシェアするサイトで、かなりおすすめ。
  * https://aws.amazon.com/jp/builders-library
* Designs, Lessons and Advice from Building Large Distributed Systems by Jeff Dean
  * http://www.cs.cornell.edu/projects/ladis2009/talks/dean-keynote-ladis2009.pdf
* jepsen.ioのconsistency modelsの章。これだけは最低限読んでおいたほうが良いと思う
  * https://jepsen.io/consistency
* MQTT作った人の講演を文字に興したもの (たぶん) 。とっかかりはこれが良いかも。
  * https://postd.cc/learning-about-distributed-systems/
* たぶん学生の人が書いてるんだと思うけど、絵が可愛いので良い。
  * https://medium.com/baseds
* jepsen.ioが提供する分散システムの講義のマテリアルだと思う。基礎的な内容
  * https://github.com/aphyr/distsys-class
* Dropboxのアーキテクチャについての動画。
  * https://www.youtube.com/watch?v=PE4gwstWhmc&ab_channel=Stanford
* PingCapが提供してる分散システムを学ぶためのやつ
  * https://github.com/pingcap/talent-plan

## Paper

* Facebookのアーキテクチャ
  * https://cs.uwaterloo.ca/~brecht/courses/854-Emerging-2014/readings/key-value/fb-memcached-nsdi-2013.pdf
* AWS Auroraのアーキテクチャ、難しい
 * https://spring-mt.hatenablog.com/entry/2021/03/01/123934
* Hekatonの話、かなり難しい
  * https://okachimachiorz.hatenablog.com/entry/20170918/1505729999

## OSS
* comma.ai が本番で使っている、1000行以下で書かれた分散KVS
  * https://github.com/geohot/minikeyvalue
* 極めて基本的な構成で分散KVS
  * https://github.com/chrislusf/vasto
  * 分散システムのことわかってくると分散KVS作りたくなるっぽくて、調べるとたくさん出てくる
* Redisで分散ロックするやつ
  * https://github.com/go-redsync/redsync
  * 分散ロックとはなにかというのはこれを読む
    * https://kumagi.hatenablog.com/entry/distributed_lock

## その他色々
* Tail Latencyの話
  * https://agnozingdays.hatenablog.com/entry/2019/04/25/210721
* Paxos
  * https://qiita.com/kumagi/items/535c9b7a761d2ed52bc0
  * この方の記事は全部読むと良いと思う
* Randezvos Hash
  * https://medium.com/@reduls/about-rendezvous-hashing-76aaf1a3f705
* CRDT
  * https://qiita.com/everpeace/items/bb73ec64d3e682279d26
* replicated log
  * http://mesos.apache.org/documentation/latest/replicated-log-internals/
* Cosmos DBの一貫性モデル
  * https://docs.microsoft.com/en-us/azure/cosmos-db/consistency-levels

