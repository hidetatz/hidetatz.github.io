カーネル構築時のオプションについて---2017-12-29 10:40:43

Gentooを普段使いしているとカーネルをビルドすることが半年に1回くらいある。
その差異、 `make *config` を実行するが、よく調べてみると種類が多くあり、
知らないことも多かったため、まとめた。

以下の内容は `/usr/src/linux/Documentation/admin-guide/README.rst` にあったものを訳している。

```
カーネルの設定
----------------------

   このステップは、たとえマイナーバージョンのアップデートであっても省略することはできません。
   どのリリースでも新たな設定オプションが追加されるし、もし設定ファイルが期待される通りでなければ、
   奇妙な問題が起こるでしょう。もし最低限の手順で既存の設定を新しいバージョンに持っていきたいならば、
   `make oldconfig` を使用します。これは、新たなカーネル設定に対する回答のみをあなたに尋ねます。

 - 設定のコマンド::

     "make config"      プレーンテキストの設定インタフェース。

     "make menuconfig"  テキストベースの色付きのメニュー、ラジオリスト、ダイアログ。

     "make nconfig"     改良されたテキストベースの色付きのメニュー。

     "make xconfig"     Qtベースの設定ツール。

     "make gconfig"     Gtk+ベースの設定ツール。

     "make oldconfig"   全ての質問の答えを既存の `.config` ファイルのコンテンツをベースとし、
                        新たな設定シンボルについては質問する。

     "make silentoldconfig"
                        上と似ているが、すでに回答済みの質問については画面に出さない。

     "make olddefconfig"
                        上と似ているが、新たなシンボルは自動的にそれらのデフォルトにする。

     "make defconfig"   `arch/$ARCH/defconfig` か、 `arch/$ARCH/configs/${PLATFORM}_defconfig` のいずれかから、
                        アーキテクチャに依存して、デフォルトのシンボルを用いて新たな `.config` ファイルを作る。

     "make ${PLATFORM}_defconfig"
                        `arch/$ARCH/configs/${PLATFORM}_defconfig` のデフォルト値を使って、新たな `.config` を作る。
                        あなたのアーキテクチャで使用可能な全てのプラットフォームのリストを得るためには、
                        `make help` を使う。

     "make allyesconfig"
                        可能な限り全てのシンボルを `y(カーネル組み込み)` にした `.config` を作る。

     "make allmodconfig"
                        可能な限り全てのシンボルを `m(カーネルモジュール)` にした `.config` を作る。

     "make allnoconfig" 可能な限り全てのシンボルを `n(ビルドしない)` にした `.config` を作る。

     "make randconfig"  シンボルをランダムに設定した `.config` を作る。

     "make localmodconfig" 今の設定とロード中のモジュール(`lsmod`)をベースに設定を作る。
                           ロード済みのモジュールに不要なオプションは無効にする。

                           別のマシンのための `localmodconfig` を作るには、そのマシンの `lsmod` をファイルにして、
                           `LSMOD` パラメータと一緒に渡す。

                   target$ lsmod > /tmp/mylsmod
                   target$ scp /tmp/mylsmod host:/tmp

                   host$ make LSMOD=/tmp/mylsmod localmodconfig

                           これはクロスコンパイル時にも動作する。

     "make localyesconfig" localmodconfigと似ているが、全てのモジュールオプションをカーネル組み込み(=y)とする。

   Linuxカーネル設定ツールについての情報は、 `Documentation/kbuild/kconfig.txt` で参照できる。
```


普段は `make menuconfig` をよく使う。
`make randconfig` はテスト用なのだろうか。
