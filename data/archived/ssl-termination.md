ELBの動きとSSLターミネーションについて学んだ---2017-11-01 20:40:29

ELBはEC2へリクエストを負荷分散してくれる。
その際の暗号化について学んだ。

本来、HTTPS証明書はサーバに配置する。
ELB配下にEC2が冗長化されているケースでも、EC2に証明書をデプロイすることは可能。
しかし、実際はSSL(TLS)の終端はELBとし、証明書もELBに持たせる。
そして、ELBからEC2への接続は暗号化されない。
これは、いくつかの理由がある。

* 証明書を一元管理するため
* 盗聴の心配がないため
* セッション維持のため
* 盗聴の心配がない

以下のような記載があった。

```
Whenever you are connecting to an RDS instance from an EC2 instance within the same region, that communication is going through AWS private network.

However traffics between RDS DB instances and EC2 instances which are located in different regions will go through the internet.

If you are concerned about security, you may want to use SSL for communication between RDS DB instances and EC2 instances.
```

RDSについての記載だが、要点は同じ。
同一リージョンのインスタンス間のトラフィックは、AWSのプライベートなネットワークを通るため、暗号化は不要とのこと。

### セッション維持

ELBにはスティッキーセッションという、Cookieベースでのセッション維持機能がある。
そして、スティッキーセッションはSSLの終端をEC2にする場合は 使えない 。
これは、ELBからEC2に通信する際のトラフィックが暗号化されると、ELBがCookieの中身を読み取れないため。(読み取れないと当然トラフィックをスティッキーに分散できない)
