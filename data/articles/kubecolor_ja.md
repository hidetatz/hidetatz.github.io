title: kubecolorでkubectlの出力に色を付ける
timestamp: 2020-10-14 00:00:00
lang: ja
---

![](https://storage.googleapis.com/zenn-user-upload/tugxpw8uu3bg8c10jounker54xls)

# 1. はじめに

kubectlはKubernetesクラスタを操作するためのCLIクライアントです。
筆者は仕事でほぼ毎日kubectlを使用しますが、その使い勝手に対して一点改善したい点がありました。それは、 _出力が一色のみで見辛い_　という点です。
この記事では、筆者が開発したkubectlの出力に色を付けるためのOSS "[kubecolor](https://github.com/hidetatz/kubecolor)" について紹介します。筆者はkubecolorは非常に便利で、もっと多くの人に使ってもらいたいと考えているため、この記事を書くことにしました。
次章からは、既存のkubectlに対して筆者が改善したいと感じていた点についてもう具体例を上げながら説明し、その後kubecolorではそれをどう改善できているのかを解説します。そして、kubecolorの使い方について説明します。

また本記事では、「見やすい」「見辛い」「わかりやすい」などの表現を使用しますが、筆者は色彩学やデザインの専門家では一切なく、単なる主観であることをご了承ください。

# 2. kubectlを改善したかった点

初めに、素のkubectlの出力を見てみます。（画像内で使用しているクラスタはテスト用のものです）

![](https://storage.googleapis.com/zenn-user-upload/krqr1ud7rwhr8bcturrrqwbbjiz2)

kubectlは出力に色を一切つけません。上の画像はまだ出力結果が十分に短いため視認性に大きな問題はありませんが、　"kubectl describe" や "kubectl get -o json" など、大きめの出力を行うコマンドは、次のように見えます。

![](https://storage.googleapis.com/zenn-user-upload/7sb702m4htks09bq1mzgn00qg3cm)

![](https://storage.googleapis.com/zenn-user-upload/lfm3wjluvleivppudndy5d5am8pv)

出力が大きくなると、先程よりも「見辛く」感じないでしょうか。もちろん個人差があるとは思うのですが、筆者は欲しいデータがどこにあるかを瞬時に認識することが難しいなと以前から感じていました。これが、筆者がkubecolorを開発した理由です。
次章では、kubecolorがどのようなアプローチで「見やすさ」に貢献できるかを解説します。

# 3. kubecolorを使うと何が変わるか

kubecolorは筆者が開発したOSSで、以下のGitHubリポジトリで開発しています。

https://github.com/hidetatz/kubecolor

ライセンスはMITです。

さて、冒頭の画像で既にお見せしてしまっていましたが改めて、kubecolorを使うとkubectlの出力は以下のように見えます。

* kubectl get pods

![](https://storage.googleapis.com/zenn-user-upload/0axnlztv2wusigzjpcg4craxxhmk)

* kubectl describe pod

![](https://storage.googleapis.com/zenn-user-upload/3uucjklexol3n7s5vlvbjchwb0it)

* kubectl get pods -o json

![](https://storage.googleapis.com/zenn-user-upload/8yi2f4hmiu1nqhkoc7mstelqy0od)

* kubectl get pods -o yaml

![](https://storage.googleapis.com/zenn-user-upload/zr0i8bvgzrgo5bbu7uym6lkjgh2s)

また、コマンドがエラーを出力した場合、kubecolorはそれを赤で出力することでエラーが起きたことをユーザーにわかりやすく伝えることもできます。

![](https://storage.googleapis.com/zenn-user-upload/6gjh1kums4rwp7cu9h5togy5uedz)

色がつくことで、かなり見やすくなったのではないでしょうか。筆者は元々の一色のみの出力よりも、目的のものが見つけやすくなったなと感じます。
次章では、kubecolorのインストール方法及び使い方などについて解説します。

# 4. kubecolorの使い方

本章の内容は、kubecolorの開発が進むにつれ最新の情報でなくなる可能性があります。

## kubecolorは何をするのか

kubecolorはkubectlコマンドに渡すべきコマンドラインオプションを受け取り、内部でkubectlコマンドを実行し、帰ってきた出力に色付けしてその後出力します。色付け以外のことは何もしません。
kubecolorコマンドはまだ開発中であり、kubectlの全てのサブコマンドの出力に対して色付けをサポートしているわけではありません。また、出力に用いる色が今後変更される可能性もあります。サポートされている機能や今後実装されるものについて、[README](https://github.com/hidetatz/kubecolor)をご覧ください。

## インストール

kubecolorはGoで実装されており、インストールにはgoコマンドを使用します。goコマンドがインストールされていない場合はまず[公式ドキュメント](https://golang.org/doc/install)などを参考にインストールする必要があります。

インストールは次のコマンドで行うことができます。

```shell
go get github.com/hidetatz/kubecolor/cmd/kubecolor
```

アップデートは次のコマンドで行うことができます。

```shell
go get -u  github.com/hidetatz/kubecolor/cmd/kubecolor
```

(上記コマンドを実行してもkubecolorコマンドが実行できない場合、 `$GOPATH/bin` が `$PATH` に入っていない、`Go modules`が有効になっていて `$GOPATH/bin` にバイナリができていないなどの原因が考えられます。)

## 使い方

kubectlコマンドに渡したいコマンドをそのままkubecolorに渡すことで、出力が色付けされます。例えば、次のようなコマンドが考えられます。

```shell
kubecolor --context my-context get pods
kubecolor edit deployment
kubecolor exec -it pod-a bash
```

kubecolorはkubectlコマンドを完全に代替するものとして設計してあります。
これはすなわち、以下のようなaliasをシェルに設定しても機能するようになっていることを意味します。

```shell
alias kubectl=kubecolor
```

しかし、単にこれを設定するとkubectlのネイティブの補完機能が動作しなくなります。対処法は後述します。

## 有効なコマンドラインオプション

### --plainで色をつけなくする

`--plain` フラグをkubecolorに渡すと、kubecolorは色をつけずに出力します。結果をファイルに書き込みたい時などに便利です。

### --light-backgroundで使用する色を設定する

kubecolorが色付けに使用する色のプリセットは、デフォルトでは「背景が暗い色に設定してあるターミナルエミュレータ」上で見やすくなるように設定されています。これは逆に言えば、ターミナルの背景を白などの明るい色に設定している環境では、kubecolorの使う色が見えにくいと感じる可能性があります。 `--light-background` フラグをkubecolorに渡すことで、明るい背景に最適な色のプリセットを使って出力に色付けを行うようになります。
実際に背景色を白に設定したターミナルで、kubecolor describe podコマンドを実行します。 `--light-background` なしでは次のように見えます。

![](https://storage.googleapis.com/zenn-user-upload/4n9qrv259imkmxe9cior4hz5tacx)

`--light-background` をつけると、次のように見えます。

![](https://storage.googleapis.com/zenn-user-upload/h0y2cffxga86lr9pse466ee8php9)

白 -> 黒に、シアン -> 濃い青にそれぞれ変わったことで、読みやすくなります。

## kubectlの自動補完機能と組み合わせる

kubectlにはネイティブに自動補完機能があります。これは、kubectlコマンドに対して補完を効かせるようになっているため、kubecolorコマンドでも同様に補完を効かせるためには次の設定が必要です。

```shell
complete -o default -F __start_kubectl kubecolor
```

また、もし `k=kubecolor` などのaliasを貼る場合は、次のようにする必要があります。

```shell
complete -o default -F __start_kubectl k
```

なお、この設定はbash上で期待通りに動作することを確認済みですが、fishなどでは事情が違います。kubectlの補完機能はfishをオフィシャルにサポートしていないためです。fish上での自動補完動作は調査中のステータスです。（fishに詳しい方がいれば、どうすれば補完可能か情報提供いただけると助かります！）

## kubectlコマンドを指定する

kubecolorは内部でkubectlコマンドを実行しますが、そこで利用するコマンドを変更したいことがたまにあります。例えば、gcloudを使えばkubectlをバージョン違いでインストールできますが、(kubectl.1.18など) これをkubecolorに使わせたいときは、次のようにします。

```shell
KUBECTL_COMMAND="kubectl.1.18" kubecolor get pods
```

`KUBECTL_COMMAND` に何もセットしていない場合、kubecolorはデフォルトで `kubectl` を使用します。

# おわりに

筆者が開発したOSS kubecolorについて紹介しました。色がつくとkubectlライフがより便利で、そして楽しくなりますね！ぜひ使ってみてください。
良さそうだなと思ってもらえたら、GitHub上でスターを付けてもらえると嬉しいです！また、issueやプルリクエストも歓迎です。
最後までお読みいただきありがとうございました。
