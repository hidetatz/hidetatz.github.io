Microservicesにおけるサービス間通信とService Mesh---2018-11-11 13:00:44

Microservicesにおいて、サービス間の通信は高速でかつ堅牢である必要がある。
MercariのようにgRPCを採用する場合もあれば、HTTPを使う場合もある。
また、多くの場合、[ActiveMQ](http://activemq.apache.org/)や[GCP Cloud Pub/Sub](https://cloud.google.com/pubsub/)、[AWS SQS](https://aws.amazon.com/jp/sqs/)などの非同期メッセージングプロトコルも併せて採用するだろう。

しかし、Microservicesにおけるサービス間通信には下記のような課題がある。

* 回復性
  - 大規模なMicroservicesでは、数百のインスタンスが存在する可能性がある。ノードレベルでの障害や、過負荷によるタイムアウトなど、様々な理由で通信が失敗する可能性がある。
  - 回復性を高めるためには リトライ、[Circuit Breaker](/post/2018/11/10/about-circuit-breaker/) などのデザインパターンがある。
* Load Balancing
  - Kubernetesでは各サービスはPodのグループという形で抽象化される。この時、待機時間や負荷によって、負荷分散のアルゴリズムを賢くしたいケースが有る。
* 分散トレーシング
  - 1トランザクションで複数のサービスを(依存サービスなどの形で)利用するケースが有る。その場合、各サービスがロギングを行っても、それらを1トランザクションとして結び付けないと、人間に読むことは難しい。
* サービスのバージョニング
  - サービスをデプロイするときに、すでに流れている依存サービスへのトラフィックなどは、正しいバージョンにルーティングされなければならない。リクエストを特定のバージョンにルーティングしなければならないケースがある。
* TLS暗号化と認証
  - セキュア化のために通信内容をTLSで暗号化したいケースが有る。

これらの課題に解決策を提供するのが、 `Service Mesh` である。
KubernetesではServish Meshとして、 [Linkerd](https://linkerd.io/)、[Istio](https://istio.io/)を選択できる。
これらは、各サービスの[Side Car](https://docs.microsoft.com/ja-jp/azure/architecture/patterns/sidecar)として稼働する。

Service Meshは以下のような機能を提供する。

* セッションレベルでのLoad Balancing。
* URL/Host Header/API VersionなどのルールベースでのL7 Load Balancing。
* リトライ
* Circuit Breaker
* メトリクスのキャプチャと、リクエストの各ホップを使った分散トレーシング
* 相互TLS認証

Service Meshを使用しない場合、冒頭の課題を開発者が考慮する必要がある。
まだ枯れていない技術ではあるが、大規模な分散システム(Microservices)を設計する際には考慮したいと思う。
