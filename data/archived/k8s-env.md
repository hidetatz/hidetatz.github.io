Kubernetesを動かすまで---2018-12-16 18:55:58

この記事は[Kubernetes完全ガイド](https://www.amazon.co.jp/dp/B07HFS7TDT)の内容(抜粋)の個人的まとめです。

---

Kubernetesを動かす方法はいくつかある。
筆者が所属しているメルカリではGKE上でKubernetesクラスタを構築しているが、
この記事で、それ以外の方法もまとめる。

Kubernetes環境は大きく分けて以下の3種類が存在する。

* ローカルKubernetes
* Kubernetes構築ツール
* マネージドKubernetesサービス

それぞれ紹介する。
なお、ローカルKubernetesについては、Mac上で動作させることを前提とする。

## ローカルKubernetes

ローカルKubernetesは、ローカルマシン上に構築するものだ。
ネットワークに繋がっていなくてもクラスタが構築できるなど、手元ならではのメリットがいくつかある。
しかし、冗長化はしていないため、プロダクション利用などは適していない。

### Minikube

Minikubeは手元のマシンにKubernetesを簡単に構築するためのツールだ。
[Kubernetesの公式チュートリアル](https://kubernetes.io/docs/tutorials/kubernetes-basics/create-cluster/cluster-intro/)でも使われている。
実行されるKubernetesはシングルノードになる。
また、MacではVirtualBoxなどのハイパーバイザーを必要とする。

### Docker for Mac

MacでDockerを利用する際は、[Docker for Mac](https://docs.docker.com/docker-for-mac/install/)を利用するが、
Docker for MacでもKubernetesが利用できる。
PreferencesからKubernetesをEnableとすることで、Kubernetesがローカルにデプロイされる。

## Kubernetes構築ツール

Kubernetes構築ツールは、kubeadmというKubernetesクラスタ構築用コマンドを利用して構築するものだ。
オンプレミス上や、なんらかの理由でマネージドを使わずにパブリッククラウド上で構築したいケースなどで使用する。

### kubeadm

[kubeadm](https://docs.docker.com/docker-for-mac/install/)は、Kubernetesが公式に提供しているクラスタ構築ツールである。
Masterノード/Slaveノード上でコマンドを実行することでクラスタが構築できる。

### Rancher

[Rancher](https://github.com/rancher/rancher)はRancher社が主に開発している
オープンソースのコンテナプラットフォームである。
Kubernetesクラスタの構築、デプロイだけでなく、モニタリングやWebUIなども提供している。

## マネージドKubernetesサービス

AWS、GCP、Azureなどのパブリッククラウドは、マネージドなKubernetes as a Serviceを提供している。

### GKE(Google Kubernetes Engine)

GCPにおけるKubernetesのマネージドサービスである。
メルカリでも採用している。
歴史的に、他のパブリッククラウドがKubernetesをサポートする前から存在しており、実績が多い。
Stackdriver Loggingとの連携などのめりっとがある。

また、GKEには `NodePool` という、仮想インスタンスのグループ機能がある。
マシン性能が異なるノードをクラスタ内に混在させ、コンテナの特性に合わせたデプロイができるなどの
利点がある。

### AKS (Azure Kubernetes Service)

AKSはMicrosoft AzureでのマネージドKubernetesサービスである。
使ったことがないので省略。

### EKS (Elastic Container Service for Kubernetes)

EKSはAWSにおけるマネージドKubernetesサービスである。
IAMベースでの認証認可や、CloudWatch連携などの利点がある。

## Kubernetesプレイグラウンド

[Play with Kubernetes](https://labs.play-with-k8s.com/)という、WebからKubernetesを試すことのできる環境もある。

## まとめ

Kubernetesをステージング/プロダクションで利用する上では、今だと可能な限りマネージドサービスを利用するのがいいと思う。
が、状況に応じて適切な手段を選択する必要がある。
