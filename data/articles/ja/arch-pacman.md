pacmanのオプションについて---2018-08-03 19:09:52

最近はArchLinuxを使っているので、pacmanのオプションについて勉強している。
以下にまとめておく。

* リポジトリの同期(emerge --sync)

```shell
$ sudo pacman -Syy
```

* ローカルのソフトウェアのアップデート(emerge -uDN world)

```shell
$ sudo pacman -Syyu
```

* ソフトウェアの検索(emerge --search)

```shell
$ pacman -Ss
```

* ローカルのソフトウェアの検索

```shell
$ pacman -Qii
```

* ソフトウェアのインストール(emerge)

```shell
$ pacman -S
```

* ローカルのソフトウェアのリスト(equery list "*")

```shell
$ pacman -Ql
```

* 依存関係の表示

```shell
$ pctree
```

* orphanパッケージの表示

```shell
$ pacman -Qdt
```

* orphanパッケージの削除(emerge --depclean)

```shell
$ sudo pacman -Rs $(pacman -Qdtq)
```

* パッケージ削除

```shell
$ sudo pacman -R
```

* 依存関係も含めてパッケージ削除

```shell
$ sudo pacman -Rs
```

* パッケージのローカルキャッシュを削除

```shell
$ sudo pacman -Scc
```

* その他

`/etc/pacman.conf` の `MiscOptions` に `ILoveCandy` を追加することで、 `pacman -S` のプログレスバーにパックマンが現れる。
