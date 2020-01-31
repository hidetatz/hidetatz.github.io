GentooのターミナルにRictyを導入した---2018-02-27 22:05:46

Rictyのいれ方。

### 前提

* WMはGNOME3
* ターミナルはGnome Terminal

### 手順

[こちら](http://www.sa-sa-ki.jp/blog/2016/11/ubuntu-16-04-lts-ricty/)の通り。

```
$ emerge fontforge

### 一時的な作業ディレクトリ作成
$ mkdir ricty
$ cd ricty

### 必要なフォントファイルをDL
$ wget https://github.com/google/fonts/raw/master/ofl/inconsolata/Inconsolata-Bold.ttf
$ wget https://github.com/google/fonts/raw/master/ofl/inconsolata/Inconsolata-Regular.ttf

### Migu DL
$ wget https://osdn.jp/projects/mix-mplus-ipa/downloads/63545/migu-1m-20150712.zip
$ unzip migu-1m-20150712.zip
$ cp migu-1m-20150712/migu-1m-*.ttf ./

$ wget http://www.rs.tus.ac.jp/yyusa/ricty/ricty_generator.sh
$ wget http://www.rs.tus.ac.jp/yyusa/ricty/os2version_reviser.sh
$ chmod 755 ricty_generator.sh os2version_reviser.sh
$ ./ricty_generator.sh auto

$ ./os2version_reviser.sh Ricty*.ttf
```

これでRictyは作れる。
が、生成したRictyにPowerline用のパッチを当てたかったので、以下も行った。

```
$ git clone https://github.com/Lokaltog/powerline-fontpatcher
### Ricty-Regularか、Boldか、Diminishedか、好きなのを選ぶ
$ fontforge -lang=py -script ./powerline-fontpatcher/scripts/powerline-fontpatcher Ricty-Regular.ttf
### ファイル名にスペースが入ったttfができあがるのが気持ち悪いのでリネーム
$ mv Ricty\ Discord\ Regular\ for\ Powerline.ttf Ricty_Discord_Regular_for_Powerline.ttf
```

これでpowerline用のフォントが出来上がる。
あとは、

```
$ cp -f Ricty*.ttf ~/.fonts/
$ fc-cache -vf
```

で、フォントの準備はOK。

あとはgnome-terminal側で、Preferecnes -> Profile -> Custom fontで設定する。
