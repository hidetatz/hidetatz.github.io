Microservicesデザインパターン---2018-11-11 16:06:03

Microservicesにおいて名前のついている設計パターンについて書いていく。
会社では、これらの言葉を使って会話している。

![](/images/microservices-designpattern/1.png)
(画像は[Design patterns for microservices](https://azure.microsoft.com/ja-jp/blog/design-patterns-for-microservices/)より)

### Ambassadorパターン

各サービスにAmbassadorを付加してデプロイし、TLS、モニタリング、ロギング(分散トレーシング)、ルーティングなどに責任を負う。
Istioなどの[Service Mesh](/post/2018/11/11/microservices-service-mesh/)は、このパターンの実装である。

### Anti-corruption layerパターン

モノリスをMicroservicesなどに移行する場合、DBはモノリスで見ていたものをしばらくは使う、というケースが有る。
この時、DBアクセスはこのAnti-corruption layerを必ず経由するようにする。
DDDにおける、Bounded Context間でモデルを変換したりするのに使う。

### Backends for Frontendsパターン

Microservicesは多くの場合APIなので、デスクトップ/モバイルなど異なるタイプのクライアントから呼び出される。
ひとつのサービスは、異なるタイプのクライアントから呼び出される方法を知っている必要はない。
そこを知っているのがBackends for Frontends レイヤーである。

### Bulkheadパターン

一つのノードにふたつのサービスが乗っていて、同じコネクションプールを共有していた場合、
サービスAが不具合でコネクションを使い果たしてしまうと、サービスBも止まってしまう。
これを避けるため、コネクションプールや、CPUリソース等をサービスごとに分離する。

### Gateway Aggregationパターン

複数の別のサービスへのリクエストをひとつにまとめることで、通信を効率化する。

### Gateway Offloadingパターン

API GatewayでSSL認証などをまとめて引き受けることで、各サービスのロードを減らす。

### Gateway Routingパターン

ユーザーが多くのエンドポイントを管理する必要がないように、複数のリクエストをひとつにまとめる。

### Sidecarパターン

アプリケーションのコンポーネントを別のコンテナに載せてデプロイしてカプセル化する。

### Stranglerパターン

サービスのバージョンアップ時などに、徐々にリクエストを新しいサービスにフォワードする。
これにより、新しいバージョンのサービスにバグが有ったときに影響を最小限にできる。

これらのパターンはすべて、サービス規模や特性によって使い分ける必要がある。
