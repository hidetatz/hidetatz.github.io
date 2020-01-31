equeryのオプション---2018-02-25 11:12:20

equeryでよく使うコマンドオプションをメモっておく。

### インストール

equery自体は以下のようにインストール。

```
$ emerge --ask app-portage/gentoolkit
```

なくてもportageは使えるが、便利なので入れるのが良いと思う。

### ファイルがどのパッケージに入っていたのか

```
$ equery belongs -e /usr/bin/glxgears
```

`/usr/bin/glxgears` がどのパッケージによってインストールされたのかを検索できる。

### パッケージが壊れていないか検査する

```
$ equery check gentoolkit
```

### 依存されているパッケージをリストアップ

```
$ equery depends pygtk
```

*pygtkを必要としている* パッケージをリストアップする。

### 依存しているパッケージをリストアップ

```
$ equery depgraph mozilla-firefox
```

`depends` の逆。

### ファイルをリストアップ

```
$ equery files --tree gentoolkit
```

### パッケージ自体をリストアップ

```
$ equery list '*'
```

よく使う。
