Dockerを使った開発環境---2018-02-02 18:22:42

自分のdockerの使い方について書いていく。

筆者はコマンドラインツールをGoやRubyで書くことが多い。
しかし、Rubyは実行にRubyが必要だ。Goはバイナリひとつで動かせるが、
バイナリのビルドには結局Goが必要になる。
しかし、ローカルにRubyやGoはなるべく入れたくないと思っている。
これは単に潔癖なのではなく、実行環境のポータブル化を意識している。

## 手順
筆者は、何かコードを書き始めるときはまずDockerfileから書く。
例として、nodeが動く環境を作りたいときは、自分はこんな感じで書いている。

### Dockerfileを書く

Dockerfileは、あまり書いたことがない頃は悩みながら書いていたが、慣れてくると、全然難しくないことがわかる。
とにかくnodeさえ動けばいいので、これだけ書く。

```
FROM node:8.9
WORKDIR /app
```

`/app` については後述する。
FROMにはRubyならruby、Nodeならnodeを書けばいい。
汎用的なものを作りたいなら、 `Ubuntu` とかを書けばいい。
自分は、 `node docker` と検索して出て来る[Docker Hub](https://hub.docker.com/_/node/)のページを見て、
バージョンを選んでいる。

### docker-compose.ymlを書く

docker-composeは絶対書かないといけないものではない(基本は複数コンテナのオーケストレーションに使うものなので。)が、
自分は書くことが多い。

```
version: '2.2'
services:
  app:
    build: .
    command: 'tail -f /dev/null'
    volumes:
      - .:/app
```

上記例では、 `app` という識別子で、 `.` (カレントディレクトリ)のDockerfileを基に コンテナを作ります、ということを示している。

`command` は、本来コンテナ起動時に行うべき処理を書く。(例えば、Railsコンテナを立ち上げた際に `rails s` するなど。)
が、ちょっと言語が動く環境を作りたいだけなので、起動時の処理は特にないことが多い。
commandはなんでもいい(書かなくてもだぶん大丈夫)だが、自分はコンテナに入りたいことがある(デバッグなどのために)ので、
そういうときは `tail -f /deb/null` を書いておくと何も起きないので便利。
というのは、ここに `echo 'Hello World'` とかを書くと、echoした後にコンテナが停止してしまうので。

また、 `volumes` で ホスト側の `.` (カレントディレクトリ)とコンテナ側の `/app` をマウントしている。
これにより、ホストでコードを書いて、それを即コンテナで実行、とできるようになる。
前述の `WORKDIR /app` のおかげで、コンテナ側のカレントディレクトリが `/app` になるので、実行ファイルのディレクトリも意識する必要はない。

これらの処理は、docker-composeなしで(docker単体で)もできる。が、実行がちょっとめんどくさい。
dockerのみだと、

```
$ docker build . -t nodeapp
$ docker run -d --name nodeapp \
  -w /app \
  -v `pwd`:/app \
  nodeapp tail -f /dev/null
```

こんな感じで、実行時のオプションをたくさん渡す必要がありめんどくさい。
docker-composeだと、簡単にできる。

```
$ docker-compose build
$ docker-compose up
```

### コードを実行

コードの実行もdocker-composeだと簡単で、

```
docker-compose run app node main.js
```

みたいに実行が可能。
ホスト側のセットアップ済みのvimで開発し、実行はコンテナでできる。

### tips

筆者は以下のように設定して、ちょっと楽にしている。

```
alias dr='docker'
alias drc='docker-compose'
```

### まとめ

dockerがないと開発がつらい。
