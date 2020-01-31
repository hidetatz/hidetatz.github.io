Kubernetesについて学ぶ 2---2018-01-14 11:01:57

前回は[こちら](https://yagi5.com/2018/01/13/kubernetes-learn-1/)。

今日からModule2の、 `Deploy an App` について学んでいく。

### Kubernetes Deployments

```
Kubernetesクラスタが稼動したら、コンテナ化されたアプリケーションをその上に配置できます。
これを行うには、KubernetesのDeployment Configurationを作成します。
デプロイメントは、Kubernetesにアプリケーションのインスタンスの作成と更新の方法を指示します。
デプロイメントを作成したら、Kubernetesのマスタースケジュールはアプリケーションインスタンスを
クラスター内の個々のノードにメンションしました。

アプリケーションインスタンスが作成されると、
Kubernetes Deployment Controllerはそれらのインスタンスを継続的に監視します。
インスタンスをホストしているノードが停止したり削除されたりすると、
Deployment Controllerがそれを置き換えます。
これにより、マシンの障害またはメンテナンスに対処するための自己修復メカニズムが提供されます。

プリオーケストレーションの世界では、
アプリケーションを起動するためにインストールスクリプトが使用されることがよくありましたが、
マシンの障害からの回復は許可されませんでした。
Kubernetes Deploymentsは、アプリケーションインスタンスを作成し、
それらのノード間でアプリケーションインスタンスを実行することによって、
アプリケーション管理に対する根本的に異なるアプローチを提供します。
```

* Deployment Configurationを作成する
* masterがnodeにアプリケーションを配備する
* Deployment Controllerがインスタンスを監視する

Kubernetesがアプリケーションインスタンス(実体はたぶんコンテナ)を、ノードに配備する。そして管理もやってくれる。
という話らしい。

### Deploying your first app on Kubernetes

```
デプロイメントを作成および管理するには、KubernetesコマンドラインインターフェイスKubectlを使用します。
KubectlはKubernetes APIを使用してクラスタと対話します。
このモジュールでは、Kubernetesクラスタでアプリケーションを実行するデプロイメントを
作成するために必要な最も一般的なKubectlコマンドについて学習します。

デプロイメントを作成するときは、
アプリケーションのコンテナイメージと実行するレプリカの数を指定する必要があります。
後でその情報を変更することができます。
ブートキャンプのモジュール5と6では、配備の拡張と更新の方法について説明しています。

最初のデプロイでは、Dockerコンテナにパッケージ化されたNode.jsアプリケーションを使用します。
ソースコードとDockerfileは、Kubernetes BootcampのGitHubリポジトリにあります。

デプロイメントが何であるかを知ったので、オンラインチュートリアルに行き、最初のアプリをデプロイしましょう！
```

### [Interactive Tutorial - Deploying an App](https://kubernetes.io/docs/tutorials/kubernetes-basics/deploy-interactive/)

### kubectl basics

```
minikubeのように、kubectlはオンライン端末にインストールされています。
ターミナルにkubectlと入力して、その使用方法を確認します。
kubectlコマンドの一般的な形式は次のとおりです。
```

```
$ kubectl action resource
```

```
これは、指定されたリソース（node、containerなど）に対して指定されたアクション（create、describeなど）を実行します。
コマンドの後に--helpを使用すると、可能なパラメータに関する追加情報を取得できます（kubectl get nodes --help）。

kubectl versionコマンドを実行して、kubectlがクラスタと通信するように設定されていることを確認します。
```

```
$ kubectl version
Client Version: version.Info{Major:"1", Minor:"8", GitVersion:"v1.8.0", GitCommit:"6e937839ac04a38cac63e6a7a306c5d035fe7b0a", GitTreeState:"clean", BuildDate:"2017-09-28T22:57:57Z", GoVersion:"go1.8.3", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"5", GitVersion:"v1.5.2", GitCommit:"08e099554f3c31f6e6f07b448ab3ed78d0520507", GitTreeState:"clean", BuildDate:"1970-01-01T00:00:00Z", GoVersion:"go1.7.1", Compiler:"gc", Platform:"linux/amd64"}
```

```
OK、kubectlがインストールされており、クライアントとサーバの両方のバージョンが表示されます。

クラスタ内のノードを表示するには、kubectl get nodesコマンドを実行します。
```

```
$ kubectl get nodes
NAME      STATUS    ROLES     AGE       VERSION
host01    Ready     <none>    6m        v1.5.
```

```
ここでは、利用可能なノード（ここでは01）が表示されます。
Kubernetesは、ノードの利用可能なリソースに基づいてアプリケーションをどこに展開するかを選択します。
```

### Deploy our app

```

kubectlの実行コマンドを使ってKubernetesで最初のアプリケーションを実行しましょう。
runコマンドは、新しいデプロイメントを作成します。
デプロイメント名とアプリケーションイメージの場所
（Dockerハブの外部にホストされているイメージの完全なリポジトリURLを含む）
を提供する必要があります。
特定のポートでアプリケーションを実行したいので--portパラメータを追加します：
```

```
$ kubectl run kubernetes-bootcamp --image=docker.io/jocatalin/kubetes-bootcamp:v1 --port=8080
deployment "kubernetes-bootcamp" created
```

```
Great！デプロイメントを作成するだけで、最初のアプリケーションをデプロイしました。
これはあなたのためにいくつかのことを行いました：

* アプリケーションのインスタンスを実行できる適切なノードを検索しました（利用可能なノードは1つしかありません）
* そのノードでアプリケーションを実行するようにスケジュールしました
* 必要に応じて新しいノードでインスタンスを再スケジュールするようにクラスタを構成しました

デプロイメントをリストするには、get deploymentsコマンドを使用します。
```

```
$ kubectl get deployments
NAME                  DESIRED   CURRENT   UP-TO-DATE   AVAILABLE AGE
kubernetes-bootcamp   1         1         1            1 1m
```

```
アプリの1つのインスタンスを実行する1つのデプロイメントがあることがわかります。
インスタンスは、ノード上のDockerコンテナ内で実行されています。
```

ためしに、 `kubectl run --help` してみた。

```
$ kubectl run --help
Create and run a particular image, possibly replicated.

Creates a deployment or job to manage the created container(s).

Examples:
  # Start a single instance of nginx.
  kubectl run nginx --image=nginx

  # Start a single instance of hazelcast and let the container expose port 5701 .
  kubectl run hazelcast --image=hazelcast --port=5701

  # Start a single instance of hazelcast and set environment variables "DNS_DOMAIN=cluster" and "POD_NAMESPACE=default" in the container.
  kubectl run hazelcast --image=hazelcast --env="DNS_DOMAIN=cluster" --env="POD_NAMESPACE=default"

  # Start a single instance of hazelcast and set labels "app=hazelcast" and "env=prod" in the container.
  kubectl run hazelcast --image=nginx --labels="app=hazelcast,env=prod"

  # Start a replicated instance of nginx.
  kubectl run nginx --image=nginx --replicas=5

  # Dry run. Print the corresponding API objects without creatingthem.
  kubectl run nginx --image=nginx --dry-run

  # Start a single instance of nginx, but overload the spec of the deployment with a partial set of values parsed from JSON.
  kubectl run nginx --image=nginx --overrides='{ "apiVersion": "v1", "spec": { ... } }'

  # Start a pod of busybox and keep it in the foreground, don't restart it if it exits.
  kubectl run -i -t busybox --image=busybox --restart=Never

  # Start the nginx container using the default command, but use custom arguments (arg1 .. argN) for that command.
  kubectl run nginx --image=nginx -- <arg1> <arg2> ... <argN>

  # Start the nginx container using a different command and custom arguments.
  kubectl run nginx --image=nginx --command -- <cmd> <arg1> ... <argN>

  # Start the perl container to compute π to 2000 places and print it out.
  kubectl run pi --image=perl --restart=OnFailure -- perl -Mbignum=bpi -wle 'print bpi(2000)'

  # Start the cron job to compute π to 2000 places and print it out every 5 minutes.
  kubectl run pi --schedule="0/5 * * * ?" --image=perl --restart=OnFailure -- perl -Mbignum=bpi -wle 'print bpi(2000)'

(Optionsは省略)

Usage:
  kubectl run NAME --image=image [--env="key=value"] [--port=port] [--replicas=replicas] [--dry-run=bool] [--overrides=inline-json][--command] -- [COMMAND] [args...] [options]

Use "kubectl options" for a list of global command-line options (applies to all commands).
```

最低限必要なのは、 `NAME` と コンテナイメージらしい。ローカルでビルドしたイメージを使うにはどうするのか？が気になる。

### View our app

```
Kubernetes内で動作しているポッドは、プライベートな独立したネットワーク上で実行されています。
デフォルトでは、同じkubernetesクラスタ内の他のポッドやサービスからは見えますが、
そのネットワーク外では表示されません。
kubectlを使用するときは、アプリケーションと通信するためにAPIエンドポイントを介して対話しています。
```

Podとは、 `group of one or more containers (such as Docker containers), with shared storage/network, and a specification for how to run the containers` とのこと。
作成したデプロイメントによってあがるコンテナ郡のことをPodという単位で呼称しているのだと思う(たぶん)。
重要な概念っぽい。

```
モジュール4のkubernetesクラスターの外部にアプリケーションを公開する方法に関するその他のオプションについても説明します。

kubectlコマンドは、通信をクラスタ全体のプライベートネットワークに転送するプロキシを作成できます。
プロキシはcontrol-Cを押して終了することができ、実行中は出力を表示しません。

プロキシを実行するための2番目の端末ウィンドウを開きます。
```

```
$ kubectl proxy
Starting to serve on 127.0.0.1:8001
```

```
私たちは現在、ホスト（オンライン端末）とKubernetesクラスタとの間の接続を持っています。
プロキシは、これらの端末からAPIへの直接アクセスを可能にします。

プロキシエンドポイント経由でホストされているすべてのAPIを見ることができます。
これはhttp://localhost:8001から入手できます。
たとえば、curlコマンドを使用してAPIを介して直接バージョンを問い合わせることができます。
```

```
$ curl http://localhost:8001/version
{
  "major": "1",
  "minor": "5",
  "gitVersion": "v1.5.2",
  "gitCommit": "08e099554f3c31f6e6f07b448ab3ed78d0520507",
  "gitTreeState": "clean",
  "buildDate": "1970-01-01T00:00:00Z",
  "goVersion": "go1.7.1",
  "compiler": "gc",
  "platform": "linux/amd64"
}
```

つまり、ポッドはプライベートなネットワーク上で稼働するので、そこに外部からアクセスするためのプロキシを立てられるということっぽい。


```
APIサーバーは、ポッド名に基づいて各ポッドのエンドポイントを自動的に作成します。ポッド名は、プロキシを介してアクセスすることもできます。

まず、Pod名を取得する必要があります。それを、環境変数POD_NAMEに保存します。
```

```
$ export POD_NAME=$(kubectl get pods -o go-template --template '{ge .items}}{{.metadata.name}}{{"\n"}}{{end}}')
$ echo Name of the Pod: $POD_NAME
Name of the Pod: kubernetes-bootcamp-390780338-7bv66
```

```
これで、そのポッドで実行されているアプリケーションに対してHTTPリクエストを行うことができます。
```

```
$ curl http://localhost:8001/api/v1/proxy/namespaces/default/pods/$POD_NAME/
Hello Kubernetes bootcamp! | Running on: kubernetes-bootcamp-390780338-7bv66 | v=1
```

```shell
urlは、PodのAPIへのルートです。

注：端末の上部を確認してください。プロキシは新しいタブ（ターミナル2）で実行され、最近のコマンドは元のタブ（ターミナル1）で実行されました。プロキシはまだ2番目のタブで実行されており、curlコマンドはlocalhost：8001を使用して動作することができました。
```
