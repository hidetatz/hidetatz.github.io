Fluxboxのキーボード周りの設定について---2017-12-30 17:35:09

いろいろ設定した。

### Capslock -> Ctrl

Capslockは不要なので、Ctrlにする。
`xmodmap` を使う。 `xev` も後で使うので入れておく。

```
$ emerge xmodmap xev
```

まず、 `Capslock` のキーコードを取得する。

```
$ xev | grep keycode

### Capslockを押す
```

筆者の環境では `66` だった。
あとは設定ファイルで66のキーを `Ctrl` に変えていく。

```
$ xmodmap
xmodmap:  up to 4 keys per modifier, (keycodes in parentheses):

shift       Shift_L (0x32),  Shift_R (0x3e)
lock        Caps_Lock (0x42)
control     Control_L (0x25),  Control_R (0x69)
mod1        Alt_L (0x40),  Alt_R (0x6c),  Alt_L (0xcc),  Meta_L (0xcd)
mod2        Num_Lock (0x4d)
mod3  
mod4        Super_L (0x85),  Super_R (0x86),  Super_L (0xce),  Hyper_L (0xcf)
mod5        ISO_Level3_Shift (0x5c),  Mode_switch (0xcb)
```

設定ファイルは `~/.Xmodmap` に書く。

```
vi .Xmodmap

### 以下記載

remove lock = Caps_Lock

keycode 66 = Control_L

add control = Control_L
```

後は `xmodmap ~/.Xmodmap` することで設定をロードする。
`.xinitrc` に書くと良い気がするが、fluxboxの場合は `~/.fluxbox/startup` にすでにロード処理が書いてあったので、
特に何もしていない。

### コロンとセミコロンを入れ替える

USキーボードは、 `L` キーの右のキーにコロンとセミコロンが割り振られており、空打ちすると(デフォルトは)セミコロンが入力され、 Shiftキーと同時に押すとコロンになる。
あんまりやってる人を見たことないが、筆者はこれを逆にして、空打ちでコロンになるようにしている。(主にvim及びSlackでの絵文字を打つ時用。)

これは簡単で、さっき作った `~/.Xmodmap` に以下のように書く。

```shell
### 47かどうかは各自確かめる
keycode 47 = colon semicolon Cyrillic_ZHE
```

### 日本語入力

日本語を入力できるようにしていく。
`fcitx` が良いらしいので、使うことにした。インプットメソッドは `mozc` とする。
まずはツールをインストール。

```shell
### gentooには fcitx-mozcがないので、分けてインストール
$ emerge fcitx
$ emerge mozc
$ emerge fcitx-configtool
```

なおこの際、[wiki](https://wiki.gentoo.org/wiki/Fcitx)によると、gtk2、gtk3、qt4のUSEを立てておくほうがいいらしいのでそのように設定する。

次に日本語フォントをインストール。

```shell
$ emerge media-fonts/ja-ipafonts
```

次に `~/.xinitrc` を以下のように編集。

```shell
### 必要な環境変数を定義
export GTK_IM_MODULE=fcitx
export QT_IM_MODULE=fcitx
export XMODIFIERS="@im=fcitx"
### fcitxを起動しておく
fcitx
```

これで、X起動時にfcitxが上がる。

最後に、fcitxの設定を行う。

```shell
$ fcitx-configtool
```

を叩くと設定のウインドウが上がってくる。
`入力メソッド` に `mozc` を追加する。

画面から設定してもOKだが、筆者は以下の設定をするために、直接設定ファイルをいじった。

* スペースキーの左の `Alt` キーを空打ちしたら日本語入力モードを解除する
* スペースキーの右の `Alt` キーを空打ちしたら日本語入力モードを実行する

Macと同じ。

```shell
$ vi ~/.config/fcitx/config

### 以下をアンコメントしてR_Altに設定
ActivateKey=R_Alt
### 以下をアンコメントしてL_Altに設定
InactivateKey=L_Alt
```

これでいい感じに動いた。

### SandS

最後にSandS。
SandSとは、スペースキーをShiftキーの様に使用し、大文字や記号入力する、ただしスペース空打ち時は通常のスペースが入力される、というもの。
個人的には必須設定。

[こちら](https://qiita.com/ychubachi@github/items/95830219f1bdf912280b)を参考にさせていただいた。

```shell
$ emerge xcape
$ xcape -d

### Spaceキーを押してキーコード確認

$ vi ~/.fluxbox/startup

### 以下を追記
xcape -e '#65=space'
```

```shell
$ vi ~/.Xmodmap

### 以下を追記
keycode 255 = space
keycode 65 = Shift_L
```

これで設定をリロードするかXを上げ直せばイケるはず。

### まとめ
色々やりすぎて若干辛い。キーボード設定はやりすぎないのが良いと思う。
