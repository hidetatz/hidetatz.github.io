GentooにFluxboxを入れた---2017-12-30 11:19:33

LinuxはWindowsやmacOSと異なり、GUI環境はあくまでカーネル上で動くプロセス(アプリケーション)に過ぎない。
GUI環境を構築したい場合、最低限 `X` の導入が必要。また、多くの場合は別途ウインドウマネージャを導入することになる。
Xとは `X Window System` のことで、サーバの役割を果たす。
ウインドウマネージャ上でターミナルやブラウザが動作しており、その際ユーザのキーボード及びマウスからの入力を利用する。
そこをXが中継していて、XはディスプレイやIOをウインドウマネージャに提供していることになる。

### 導入
まずはXを入れる。

```
$ emerge x11-base/xorg-server
```

起動には `startx` する。
筆者環境ではいくつかコマンドが足りなかったので追加した。

```
$ emerge xclock
$ emerge xterm
```

### ウインドウマネージャ導入

```
$ emerge fluxbox
```

### 設定

Xを上げた時に、ウインドウマネージャとしてfluxboxを使用する設定。

```
$ vi ~/.xinitrc
```

```shell
exec startfluxbox
```

これで `startx` することで、fluxboxが起動する。

次に、fluxboxのメニューを作る。
メニューは画面上で右クリックした際に表示され、ここからアプリを立ち上げたり、Xを落としたりなどできる。

```
$ fluxbox-generate_menu
```

後は、 `~/.fluxbox/menu` をいじることでメニューを変えられる。

次に、google-chromeを入れた。

```
$ emerge google-chrome
```

そして、fluxboxのメニューに登録する。

```
$ vi .fluxbox/menu
```

```
### 好きなところに追記
      [exec] (Google Chrome) {google-chrome-stable}
```

`()` がメニュー表示文言で、 `{}` がコマンドを示している。

chromeが起動できることを確認する。

### 終了
終了時は何もないところで右クリックして、 `Fluxbox menu` -> `Exit` すればOK。
