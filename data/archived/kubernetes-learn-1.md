Kubernetesについて学ぶ 1---2018-01-12 18:40:08

Kubernetesについて、ほぼ何も知らないので学んでいく。

### 学び方

[チュートリアル](https://kubernetes.io/docs/tutorials/kubernetes-basics/)を読んで、書いてある通りに動かし、このブログに書いていく。

### Overview

[こちら](https://kubernetes.io/docs/tutorials/kubernetes-basics/)をまずは読んでいく。

```
このチュートリアルでは、Kubernetesクラスタオーケストレーションシステムの基礎を説明します。
各モジュールには、主要なKubernetesの機能と概念に関する背景情報が含まれており、
インタラクティブなオンラインチュートリアルも含まれています。
これらのインタラクティブなチュートリアルでは、簡単なクラスタとそのコンテナ化されたアプリケーションを自分で管理できます。

インタラクティブなチュートリアルを使用すると、次のことを学ぶことができます。

* containerizedされたアプリケーションをクラスタにデプロイする
* デプロイメントのスケーリング
* コンテナ化されたアプリケーションを新しいソフトウェアバージョンで更新する
* コンテナ化されたアプリケーションをデバッグする

このチュートリアルではKatacodaを使用して、Webブラウザに仮想端末を実行し、
どこでも実行できる小規模なKubernetesのローカル展開であるMinikubeを実行します。
ソフトウェアをインストールしたり、何かを設定する必要はありません。各インタラクティブチュートリアルは、Webブラウザ自体から直接実行されます。
```

ここまでで。とりあえず、

* クラスタってものがあり
* コンテナ化されたアプリケーションを、クラスタに対してデプロイする

ことを理解した。クラスタとはなにか？については、おいおい出てくることを期待する。

```
現代のWebサービスでは、ユーザーはアプリケーションを24/7で利用できると期待しており、
開発者はこれらのアプリケーションの新しいバージョンを1日に数回デプロイすることを期待しています。
コンテナ化は、パッケージソフトウェアがこれらの目標を達成するのを助け、
アプリケーションをダウンタイムなしで簡単かつ迅速にリリースおよびアップデートできるようにします。
Kubernetesは、コンテナ化されたアプリケーションをいつでもどこでも実行できるようにし、
作業に必要なリソースとツールを見つけるのに役立ちます。
Kubernetesは、コンテナオーケストレーションでのGoogleの蓄積された経験と、
コミュニティから得た最高のアイデアを組み合わせた、プロダクション対応のオープンソースプラットフォームです。
```

つまり、Kubernetesを我々が使うことで得られるメリットは。

* アプリケーションをダウンタイム無しで1日に複数回デプロイすることができる

ことだと言えるだろう。

### [Create a Cluster](https://kubernetes.io/docs/tutorials/kubernetes-basics/cluster-intro/)

## Using Minikube to Create a Cluster

```
Kubernetesは、単一ユニットとして動作するように接続されている可用性の高いコンピュータクラスタを調整します。
Kubernetesの抽象化によって、コンテナ化されたアプリケーションを個々のマシンに具体的に結びつけることなく
クラスタに展開することができます。この新しいデプロイメントモデルを利用するには、
アプリケーションを個々のホストから切り離す方法でパッケージ化する必要があります。つまり:
それらをコンテナ化する必要があります。

コンテナ化されたアプリケーションは、アプリケーションがホストに深く統合されたパッケージとして
特定のマシンに直接インストールされていた過去のデプロイメントモデルよりも柔軟性があり、利用可能です。
Kubernetesは、クラスタ全体のアプリケーションコンテナの配布とスケジューリングをより効率的に自動化します。
Kubernetesはオープンソースのプラットフォームで、プロダクションの準備が整いました。

Kubernetesクラスタは、2種類のリソースで構成されています。

* The Master: クラスタを調整する
* Nodes: アプリケーションを実行するワーカー
```

従来は、アプリケーションとホストが深く結合されており、柔軟性を欠いていた。(ホストに障害が起きたりしたらアプリも死ぬ。逆に、スケールも難しい。)
しかし、アプリケーションをコンテナ化することで、個々のマシンと疎結合な状態を生み出すことができる。
`クラスタ` というものがそのabstractionsを実現してくれているらしい。
しかし、コンテナという物自体はLXCやdockerが提供してくれていたが、
kubernetesが **クラスタ全体のアプリケーションコンテナの配布とスケジューリングをより効率的に自動化** してくれるらしい。
つまり、ホストとアプリを疎結合にすることはkubernetesの仕事ではなく、(それはコンテナで実現される)
クラスタにアプリを配備し、スケジューリングすることが仕事だと言える(たぶん)。

## Cluster Diagram

```
マスターはクラスタの管理を担当します。
マスターは、アプリケーションのスケジューリング、アプリケーションの望ましい状態の維持、
アプリケーションのスケーリング、新しい更新の展開など、クラスター内のすべてのアクティビティーを調整します。

ノードは、Kubernetesクラスタ内のワーカーマシンとして機能するVMまたは物理コンピュータです。
各ノードには、ノードを管理し、Kubernetesマスターと通信するためのエージェントであるKubeletがあります。
ノードには、Dockerやrktなどのコンテナ操作を処理するツールも必要です。
運用トラフィックを処理するKubernetesクラスタには、最低3つのノードが必要です。

Kubernetesにアプリケーションをデプロイするときは、
マスターにアプリケーションコンテナを開始するよう指示します。
マスターは、クラスタのノード上で実行するコンテナのスケジュールを設定します。
ノードは、マスターが公開するKubernetes APIを使用してマスターと通信します。
エンドユーザーは、Kubernetes APIを直接使用してクラスタとやりとりすることもできます。

Kubernetesクラスタは、物理マシンまたは仮想マシンのいずれにも展開できます。
Kubernetes開発を開始するには、Minikubeを使用できます。
Minikubeは、ローカルマシンにVMを作成し、1つのノードだけを含む単純なクラスタを展開する軽量Kubernetesの実装です。
Minikubeは、Linux、macOS、およびWindowsシステムで使用できます。
Minikube CLIは、クラスタを操作するための基本的なブートストラップ操作（開始、停止、ステータス、および削除を含む）を提供します。
ただし、このチュートリアルでは、Minikubeがプリインストールされているオンライン端末を使用します。

Kubernetesが何であるかを知ったので、オンラインチュートリアルに進み、最初のクラスターを開始しましょう！
```

* マスターはクラスタ全体の管理をしている(実体はホスト)
* ノードはサーバーのこと。
* 各ノードには、Kubeletというエージェントがいる
* クラスタには最低3ノード必要

つまり。
masterにはクラスタ全体の情報を管理させ、ノードがマスターに対して、API経由で情報をもらい、
自ノードの上でコンテナを作る。
APIはデベロッパがたたくこともできるらしい。

ここで `Minikube` というものが出てくる。
Minikubeを使って、ローカルにVMを起動し、1つのノードを持つクラスタを動かせるらしい。

### [Interactive Tutorial](https://kubernetes.io/docs/tutorials/kubernetes-basics/cluster-interactive/)

### Interactive Tutorial - Creating a Cluster

チュートリアルはブラウザ上にターミナルがあって、どっかの仮想サーバをいじれるようになっている。

```
$ minikube version
minikube version: v0.17.1-katacoda
```

```
$ minikube start
Starting local Kubernetes cluster...
$
```

`start` でスタートできるっぽい。

```
このブートキャンプ中にKubernetesと対話するために、
コマンドラインインターフェイスkubectlを使用します。
次のモジュールでkubectlについて詳しく説明しますが、今はクラスタ情報をいくつか見ていきます。
kubectlがインストールされているかどうかを確認するには、kubectl versionコマンドを実行します：
```

```
$ kubectl version
Client Version: version.Info{Major:"1", Minor:"8", GitVersion:"v1.8.0", GitCommit:"6e937839ac04a38cac63e6a7a306c5d035fe7b0a", GitTreeState:"clean", BuildDate:"2017-09-28T22:57:57Z", GoVersion:"go1.8.3", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"5", GitVersion:"v1.5.2", GitCommit:"08e099554f3c31f6e6f07b448ab3ed78d0520507", GitTreeState:"clean", BuildDate:"1970-01-01T00:00:00Z", GoVersion:"go1.7.1", Compiler:"gc", Platform:"linux/amd64"}
$
```

Go製らしい。kubectlというインタフェースでminikubeを操作する。

```

OK、kubectlが設定されており、クライアントのバージョンとサーバーのバージョンの両方を見ることができます。
クライアントのバージョンはkubectlのバージョンです。サーバーのバージョンは、
マスターにインストールされたKubernetesバージョンです。ビルドの詳細を表示することもできます。
```

とのこと。


```
$ kubectl cluster-info
Kubernetes master is running at http://host01:8080
heapster is running at http://host01:8080/api/v1/namespaces/kube-system/services/heapster/proxy
kubernetes-dashboard is running at http://host01:8080/api/v1/namespaces/kube-system/services/kubernetes-dashboard/proxy
monitoring-grafana is running at http://host01:8080/api/v1/namespaces/kube-system/services/monitoring-grafana/proxy
monitoring-influxdb is running at http://host01:8080/api/v1/namespaces/kube-system/services/monitoring-influxdb/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

`cluster-info` をたたくことで、クラスタの情報をとれるっぽい。
上記を見るに、master、及びdashboardなどがクラスタ内で動いているらしい。


ちなみに、 `kubectl cluster-info dump` を叩いてみたら、めちゃめちゃ長いjsonが吐かれた。

```
$kubectl get nodes
NAME      STATUS    ROLES     AGE       VERSION
host01    Ready     <none>    8m        v1.5.2
```

これは `docker ps` した時のやつですね。


要するに、 `host01` っていうホストで構成されるクラスタがあり、その上でmasterやdashboardが動いている。その情報にアクセスするためのAPIが
`kubectl` らしい。

ここまでで `Module1(Create a Cluster)` が終了した。

この後はModule2の `Deploy an App` をやっていく。
