Portageのオプション---2018-02-25 05:01:33

portage操作によく使うオプションたちを覚えるためにまとめる。

### 同期

Portageツリーの同期。ebuildファイル(パッケージの設計書)を更新する。

```
### emerge --sync
```

### 検索

パッケージを検索する。

```
### emerge --search <pattern>
```

以下でパッケージの詳細を表示する。(--pretendによりインストールは行われない)

```
### emerge --verbose --pretend(-p) <package>
```

### インストール

パッケージをインストールする。

```
### emerge --ask(-a) --verbose(-v) <package>
```

### 更新

インストールした全てのパッケージを更新する。
あんま使わない。

```
### emerge --update(-u) --deep(-D) world
```

### 更新(--newuse)

USEフラグに変更があった全てのパッケージを更新する。
よく使う。

```
### emerge --update --deep --newuse(-N) world
```

### アンインストール

指定したパッケージを削除する(依存関係は削除されない)

```
### emerge --unmerge(-C) <package>
```

### 不要パッケージ削除

依存関係を持たないパッケージを削除する。unmergeした後はたたくのが推奨。

```
### emerge --depclean
```

### その他

Portageによって使用されている変数を確認する。

```
### emerge --info
```

前回のセッションを再開する。

```
### emerge --resume
```

最初のパッケージを飛ばして、前回のセッションを再開する。

```
### emerge --resume --skipfirst
```

/var/lib/portage/worldにパッケージ名の書き込みを行わない。
ライブラリなど、依存関係でbuildされるパッケージを単独でrebuildする場合などに使う。

```
### emerge --oneshot(-1)
```

* システムの再構築

```
### emerge -e system
### emerge -e world
```
