GentooにSkypeをインストールした---2018-01-02 19:27:43

Skypeをインストールした。

仕事で使うのでSkypeを導入した。[wiki](https://wiki.gentoo.org/wiki/Skype)にある、portageのSkypeは動かなかったので、オーバーレイを使用した。

### layman導入

laymanを使う。

```
$ emerge --ask layman
```

しかし、下記のようなエラーが出る。

```
ERROR: dev-python/cryptography-2.0.2-r1::gentoo failed (compile phase)
```

色々調べたところ、opensslのバージョンが1.1.0以上なら、ダウングレードしろという情報があった。
筆者の環境は、

```
$ openssl version
OpenSSL 1.1.0g  2 Nov 2017
```

だったので、ダウングレードしてみる。  
`/etc/portage/package.mask/openssl` を作成し、以下のように書く。

```
>dev-libs/openssl-1.1.0
```

これで1.1.0は強制的にマスクできるので、この状態で `emerge -uDN world` することで、1.1.0以下の最新を取ってくるようになる。

その後、再び `emerge layman` したら通った。
あとは[公式](https://wiki.gentoo.org/wiki/Layman)の通り、

```
$ layman-updater -R
```

しておいた。

### skypeforlinuxをインストール

いい感じに動いているという情報のあった、 [skypeforlinux-8.13.76.6](http://gpo.zugaina.org/net-im/skypeforlinux) をインストールする。

```
$ emerge -a skypeforlinux
```

あとは

```
$ skypeforlinux
```

で起動する。
