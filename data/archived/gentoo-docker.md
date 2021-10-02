GentooにDockerを入れた---2017-12-31 08:48:51

Gentooを使うにあたり、Dockerがないとなにもできない。インストールしていく。

### dockerインストール

先にカーネルの設定をしてもいいが、筆者はとりあえず入れた。

```
$ emerge app-emulation/docker
```

筆者のUSEフラグはデフォルトのままにしてあり、なにもいじってない。
注意点としてはファイルシステム周り。overlayを有効にすることだけ気をつければいい。
この状態で試しに、

```
$ rc-service docker start
```

でdockerをスタートする。
カーネルがいい感じに設定されていれば、これでも起動するはず。 `sudo docker ps` などして確認する。

もし起動しなければ、カーネルの再設定が必要。必要に応じて `/var/log/docker.log` を確認しておく。

### カーネル設定

カーネル設定については、[wiki](https://wiki.gentoo.org/wiki/Docker)の通りやる。
筆者の場合、特にモジュールを追加でロードするような設定もしなかった。
注意点としては、このwikiはファイルシステムにoverlayを採用するようなやり方が書いてある。
AUFSを採用する場合は[こちら](https://wiki.gentoo.org/wiki/Aufs) を参考にパッチを当てる必要がある。
また、AUFSを使う場合はdockerインストール時にaufsのUSEフラグを立てておくことも必要。

カーネル設定したら `grub-mkconfig` して再起動する。

### 使用準備

これでdockerがちゃんと起動することを確認する。
また、

```
$ rc-update add docker default
```

でdockerがOS起動時に上がってくるように設定しておく。
最後に、

```
$ sudo gpasswd -a ユーザ名 docker
```

で、ユーザをdockerグループに入れておく。こうしないと、 `docker ps` すら `sudo` が必要になってしまう。
