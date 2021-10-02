nginx学び直し---2018-09-23 19:09:52

## user

```nginx.conf
user nobody;
```

ワーカプロセスの実行ユーザ。
デフォルトは `nobody` 。

## worker_processes

```nginx.conf
worker_processes 1;
```

ワーカプロセス数。
デフォルトは1。
nginxのワーカはクライアントからの接続要求から始まる一連の処理をシングルスレッドで、イベント駆動で実行する。
基本的にはCPUのコア数と同じにしておく。
autoに設定する (`worker_processes auto`) ことで自動でコア数にしてくれる。
たまに適当に `1024` とかにしているのを見ることがあるがやめておいたほうが良い。

## worker_rlimit_nofile

```nginx.conf
worker_rlimit_nofile 1024;
```

ワーカプロセスがオープン可能なファイルディスクリプタの数。
Linuxでは、プロセスがオープンできるファイルディスクリプタの上限はデフォルトで1024。(カーネルパラメータの設定で変更可能)
大量の静的ファイルをひã¨つのワーカが返したりするとファイルディスクリプタが枯渇し、 `Too many open files` が発生する。
静的ファイルの配信が大量に発生し、かつ数千のコネクションを同時に処理するケースでは、設定しておくべき。

## events

```nginx.conf
events {
}
```

### worker_connections

```nginx.conf
events {
  worker_connections 1024;
}
```

ワーカが処理するコネクション数。
デフォルトは512。
`コネクション` はクライアントとの通信だけではない。(アップストリームとの通信なども含む。)
静的ファイルの配信であれば、数倍に増やしても問題ない。

### use

```nginx.conf
events {
  use epoll;
}
```

コネクションの処理方式を指定する。
nginxはデフォルトで、システムに最適なメソッドを選択してくれるので、基本的には設定不要。
Linux2.6以上ではepollが最適。

## keepalive_timeout

```nginx.conf
keepalive_timeout 60s;
```

keepaliveのタイムアウト時間。
デフォルトは75s。
0にすれば常時接続を無効化できる。

## sendfile

```nginx.conf
sendfile on;
```

デフォルトはoff。
onにすることで、ファイルの読み込みとレスポンス送信にsendfile()システムコールを使用する。
それにより、ファイルをオープンしているファイルディスクリプタから直接クライアントにファイルを送信してくれる。
基本onにしておく。

## tcp_nopush

```nginx.conf
tcp_nopush on;
```

デフォルトはoff。
onにすることで、TCP_CORKオプションが使用され、最も大きなパケットサイズでレスポンスヘッダとボディを送信でき、送信するパケットを最小化できる。
基本onにしておく。

## open_file_cache

```nginx.conf
open_file_cache max=最大エントリ数 [inactive=有効期間];
```

デフォルトはoff。
onにすることで、一度オープンしたファイルの情報を一定期間キャッシュする。
- ファイルのディスクリプタ、サイズ、更新日時
- ディレクトリが存在するか
- `file not found` `permission denied` などのエラー情報

キャッシュアルゴリズムはLRU。

## pcre_jit

```nginx.conf
pcre_jit on;
```

デフォルトはoff。
onにすることで、正規表現のJITコンパイルを有効にでき、正規表現の処理速度向上が見込める。
ただし、nginxのコンパイル時にPCREライブラリを動的リンクする必要がある(`--enable-jit`)。
