Kubernetesのオススメ勉強方法---2018-12-23 16:07:32

この記事は[Kubernetes3 Advent Calendar](https://qiita.com/advent-calendar/2018/kubernetes3)23日目の記事です。

筆者は、最近退っ引きならない事情でKubernetesに精通する必要性が生じて勉強中の初心者です。
Kubernetesは世界で見ればそのコミュニティは非常に大きく、様々な勉強方法のtipsがあり、この記事で私のたどった道のりをまとめます。
なお、読む人はDockerなどコンテナの知識を持っていることを想定しています。

---

## 書籍

### Kubernetes完全ガイド

[Kubernetes完全ガイド](https://www.amazon.co.jp/dp/B07HFS7TDT)

体系的な知識の構築に最も役立ったのはKubernetes完全ガイドだった。
この本は、Kubernetesを構成する技術要素(リソースという)をひとつずつ説明している。
また、リソース以外にも周辺のOSSやエコシステムなどにも触れられており、全体像を把握できるようになっている。
実際プロダクションですべてのリソースを使うことはないかもしれないが、全体像を把握する上で読んでおくのがおすすめ。

### 入門Kubernetes

[入門Kubernetes](https://www.amazon.co.jp/dp/4873118409)

Kelsey Hightower氏の本。
Kubernetes完全ガイドとの違いは、実際プロダクションでよく使うリソースに絞って書かれていること。
例えば、CronJobやネットワーキングについては書かれていない。
しかし、ReplicaSetやDeploymentなど、一般的によく使うリソースについては詳しい書かれている。
2冊目として読むのがおすすめ。

## 実践

実際に動かさないと厳しいので、実際に何をしていくかを書いていく。

### Play with Kubernetes

[Play with Kubernetes](https://labs.play-with-k8s.com/)

Dockerが提供する、web上からKubernetesに触れるプレイグラウンド。
初期状態ではクラスタ自体存在しないため、kubeadmでマスターノードを構築するところから
操作できる。
インスタンスの台数も自分で決める。
ローカルにファイルを持っていくことなどはできず、
yamlはkubectl -fでURLを指定する必要がある。

非常に便利だが、4時間でセッションが切れてすべてのリソースが削除されてしまう。
また、単純に動作が重い(たぶん海外のVMにつないでいると思う)。

### Kubernetes The Hard Way

[Kubernetes The Hard Way](https://github.com/kelseyhightower/kubernetes-the-hard-way)

入門Kubernetes作者のKelsey Hightower氏がGitHubで公開しているプログラム。
GCP上で、Kubernetesクラスタをいちから作る手順が書かれている。
kube-controller-managerやkube-proxyなども自分で構築するようになっている。
GKEやEKSを使っていると隠蔽されている部分を自分で構築することで
より理解が深まるため、一度やってみるといいと思う。

## その他読んだほうがいいもの

### what-happens-when-k8s

 [what-happens-when-k8s](https://github.com/jamiehannaford/what-happens-when-k8s)

kubectl runしたときに何が起こっているのかを文章で解説したドキュメント。
The Hard Wayよりも細かい部分を解説してくれている。

### Kubernetes at GitHub

[Kubernetes at GitHub](https://githubengineering.com/kubernetes-at-github/)

GitHubのKubernetesの事例。

## まとめ

あまりいろいろやるのも大変なので、
ここに書いたものを頑張るのがいい気がしている。
