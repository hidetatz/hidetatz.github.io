static siteをFirebaseでホストする---2018-10-06 21:09:22

このブログをfirebaseでホストすることにした。
手順をdumpしておく。

### プロジェクトの作成
コンソールからプロジェクトを作成する。

### firebase cliの設定

firebaseはコマンドラインでデプロイできるため、ツールの設定を行う。

```shell
### ツールインストール
$ sudo npm install -g firebase-tools

### ログイン。URLが表示されたらコピーしてブラウザからアクセス
$ firebase login --no-localhost

### init いろいろ聞かれるので、各自の設定に合わせる
$ firebase init

### デプロイ。initで指定したディレクトリがデプロイされる
$ firebase deploy
```

siteを更新した際は `firebase deploy` を実行すれば良い。

### ドメインの設定
ProjectのコンソールからConnect Domainを選択し、ドメインを入力する。
Aレコードに指定すべきIPアドレスが表示されるので、DNS側で設定する。

あとはTLS証明書のプロビジョニングをしばらく待てば完了する。
