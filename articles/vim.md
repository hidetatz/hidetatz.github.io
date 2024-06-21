title: なるべくプラグインを使わずにどうvimでプログラミングしているか
timestamp: 2023-01-02 22:58:12
lang: ja
---

筆者はvimを使い始めて10年弱くらい経ちますが、未だにvimへの習熟度は、感覚的には50%くらいではないかなと思っています。もしかしたらもっと低いのかもしれない。使い始めた頃は、インターネット上のvimrcをコピペし、プラグインマネージャーの使い方を頑張って読んだりしていましたが、最近はvim本来のパワフルな機能をもっと活用することに興味があります。プラグインが良いものであることは知っていますが、ラップトップのvim周辺をシンプルに保つことは私の嗜好に合っていると感じています。
筆者のvimrcは[こんな感じ](https://gist.github.com/hidetatz/32fba337b62e81953056da6a028bd919#file-vimrc)です。

```vimscript
filetype plugin indent on
set autowrite
set belloff=all
set smartindent
set path+=**
set wildmenu
set ignorecase
set smartcase
set hlsearch
set incsearch
```

普段どんなふうにvimを使ってプログラミングしているか、ちょっと書いてみようと思います。今から書くことは、上記のvimrcがホームディレクトリにあることを前提としてます。また、筆者はUbuntu20.04をOSとして使っていて、 `/etc/vim/vimrc` は特に変更してないです。vim自体は8.1のHugeバージョン ( `apt-get install vim` でインストールされるもの ) を使ってます。

諸々説明のために、下記のようなディレクトリ構成で考えます。

```
hidetatz@fox:~/tmp$ tree
.
├── dir
│   ├── a.py
│   ├── b.py
│   ├── c.py
│   └── dir2
│       ├── d.py
│       ├── e.py
│       └── f.py
├── main.py
└── test
    └── tests.py

3 directories, 8 files

```

## ファイル (バッファ) を開く / 閉じる

vimにおけるバッファとは、オープンしてメモリに載っているファイルのことです。例えば、 `vim main.py` などとコマンドを実行すると、main.pyはメモリに展開され、 `main.py` のバッファが開かれます。ここで `:ls` すると、下の方にバッファリストが表示され、バッファ `main.py` が、おそらくバッファ番号1として表示されると思います。

```
:ls
  1 %a   "main.py"              line 1
```

ちなみに、 `%a` はアクティブバッファであることを示します。他にもいろいろ種類があり、[ドキュメント](https://vim-jp.org/vimdoc-ja/windows.html#:ls)を見ると意味がわかります。

vimを開いた状態で別のファイルをバッファとしてオープンするには、 `:e main.py` とか、 `:e test/tests.py` のように実行します。ただし、ファイル名やディレクトリ名はたいてい正確に記憶していないので、実際はvimに補完してもらいます。どうやって補完するかというと、 `:e ` まで打って、ここでTabキーを打ちます。vimrcに `:set wildmenu` があると、保管候補の一覧が表示されるようになり見やすいのでオススメです。Tabキーを押すとこんな感じになります。

```
a.py  b.py  c.py  dir2/
:e dir/a.py
```

ここでTab / Shift-Tabや、← / →で候補を移動し、Enterで選択してバッファをオープンできます。候補に出てくるディレクトリ内のファイル名を補完したいときは、ディレクトリを選択して矢印下を押すことでディレクトリが掘れます。

ここで、 `:bd 1` とかすると、バッファ番号1のバッファ ( `main.py` ) がメモリからパージされ、バッファリストから削除されます。 `:bd` は `:bdelete` でも大丈夫。バッファリストが0件になると、ウインドウは空となり無が表示されます。

## バッファ間を移動する

複数のバッファを開いているときにバッファ間を移動するのは、 `:b main.py` でできます。ただ、これも `:b` してTabを押すと補完が効きます。ちなみに実は、 `:b m` だけ打ってTabを押しても `main.py` のバッファが開きます。ただし、mで始まるファイルが他にもバッファに載っていると失敗します。ファイルがひとつに特定できるまでは入力が必要です。まあ、Tabを使うのが大抵の場面で無難かなという印象です。

## 簡易Fuzzy-Find (あいまい検索) する

ディレクトリをいくつも下に降りていったところにあるファイルを開くときはどうするか？上記で言う `f.py` のような。Tab補完で頑張るのもいいですが、一つの方法としては `path` をいじるというのがあります。vimrcに以下のように書きます。

```
set path+=**
```

で、 `:e` ではなく `:fin` を使います。 `:fin f` などと打ってTabを打つと、 `f.py` が表示され、そこでEnterを打つと `dir/dir2/f.py` が開かれてバッファに載ります。

なぜこうなるかというと、まず `fin` というのはfindの意味です。 `e` は引数として渡されたファイルを開くときに使うのですが、 `fin` はpathオプションに指定されているディレクトリからファイルを探してきてそれを開きます。pathオプションというのは実際のところただの文字列のリストです。これは `find` 系のコマンド ( `:sfind` や `:tabfind` など) でファイルの検索パスとして機能します。デフォルトでは何が入っているか見てみましょう。 `:set path?` すると、次のように表示されます。

```
path=.,/usr/include,,
```

`/usr/include` が入っているようですね。これはCプログラマのために設定されています。Cでは `/usr/include/` ディレクトリの配下にあるファイルはプログラムからインクルードして標準ライブラリとして使えます。なので、例えばCのコードを読んでいるときに `#include <stdio.h>` などとインクルードが書いてあって、このヘッダファイルの中身が読みたい!というときは、 `:fin stdio.h` って実行すれば、 `/usr/include/stdio.h` をvimが開いてくれます。これが便利だという話です。なので、デフォルトのオプションはCプログラマ以外にはあまり役に立ちません。また、最近は [macOSだとそもそも `/usr/include` がなかったりする](https://qiita.com/yoya/items/c0b26cba3c040c581643) らしいです。
pathはただのオプションなので、ユーザーが適当に設定することができます。例えばGoプログラマはローカルマシンにGoが入っていると思います。[Goのオフィシャルなインストール方法](https://go.dev/doc/install) でインストールしている場合、Goのソースコードは `/usr/local/go/src` に置かれているはずです。なので、vimrcに `path+=/usr/local/go/src/*` などと書けば、Goのソースコードも `:fin` できるわけです。これがpathオプションの使い方のひとつです。

話を戻すと、我々はvimrcに `set path+=**` と書きました。 `+=` はpathにアペンドする演算子です。 `**` は、ディレクトリの区切りを無視してカレントディレクトリ配下全てにマッチするような動きになります。そのため、dirやdir2を入力しなくても `f.py` の補完ができたわけです。

## ファイルエクスプローラを使う

netrwとはvimにビルトインのファイルエクスプローラです。というか、名前の通り本来はリモートのファイルをネットワーク越しに編集するための機能だと思います (筆者はこの用途でnetrwを使ったことがない) 。ファイルエクスプローラとしてはnerdtreeなどのプラグインが有名かなと思うのですが、個人的にはnetrwでも十分です。netrwは、ファイル名でなくディレクトリ名 (あるいは `.` など) をオープンしようとすると開きます。例えば `vim .` としたり、 `:e dir` など。

```
" ============================================================================
" Netrw Directory Listing                                        (netrw v165)
"   /home/hidetatz/tmp
"   Sorted by      name
"   Sort sequence: [\/]$,\<core\%(\.\d\+\)\=\>,\.h$,\.c$,\.cpp$,\~\=\*$,*,\.o$,\.obj$,\.info$,\.swp$,\.bak$,\~$
"   Quick Help: <F1>:help  -:go up dir  D:delete  R:rename  s:sort-by  x:special
" ==============================================================================
../
./
dir/
test/
main.py
```

netrwもvimのウインドウの中なので、j/kで移動できます。 `D` で削除、 `R` でリネームなども可能です。 `-` を押すとひとつ上のディレクトリに飛べます。 `/` での検索も可能です。Enterでファイルを開けますが、 `v` や `o` で開けばウインドウを分割して開くことができます。

## 行内検索

ノーマルモードで `f` すると、カーソルのある行内で飛びたい場所にジャンプできます。具体的に言うと、

```
Netrw Directory Listing
```

今カーソルがこの行の一番左にあるとします。ここで、Listingの `L` にジャンプしたいとき、`fL` と押すとカーソルが `L` の上にジャンプします。仮にそれよりも後ろにLがあって、飛びたいのはそっちだった!というときは、 `;` を押すと次の `L` に移動します。行き過ぎたら `,` で戻れます。
`f` のポイントは、なるべく一発で行ける文字を選ぶことです。例えば上記だと、 `i` は3回登場しているので、一番右のiに飛びたい場合は、 `fi;;` と押す必要があります。それよりは、 `fn` でnに飛び、hでiに移動することでキーを1回分節約できるわけです。
確か英語では `e` `t` `a` `n` あたりが登場しやすかったはずなので、これらはfの入力として意識的に避けることが多いです。 `z` や `x` がジャンプしたい文字の近くにあると嬉しくて笑顔になります。
注意点としては、fはカーソルのある行の中でしか移動できません。

## コードフォーマット

外部のプログラムを使ったコードフォーマットをvimからかけたい、みたいなのはどうするか。一つのやり方としては、 `formatprg` や `equalprg` が使えます。例えばC++を書いていて、コードフォーマッターにclang-formatを使いたいとしましょう。このとき、vimrc (またはvimのコマンドラインで) `set formatprg=clang-format` とします。で、 `gg` で一番上に移動、 `gq` を押し、 `G` で一番下に移動します。これで、clang-formatがかかります。 `formatprg` ではなく `equalprg` を使えば `gq` ではなく `=` でおなじことができます。

これは何が起きているのか。まず、gqコマンドと=コマンドですが、これらはvimで編集している文章の見栄えを調整するものです。微妙に用途が異なっていて、gqは `textwidth` というオプションに設定された各行の最大文字数に合わせて調整を行うコマンドです。長すぎる行は分割し、短すぎる行はマージします。=は、これはコードにフィルタをかけるコマンドです。ただし、デフォルトではC言語用のフォーマッタとして機能します。おそらく新しいプログラミング言語とかではちゃんと動かないのだろうと思います。

で、gqと=は、それぞれ `formatprg` と `equalprg` を設定することで、外部のプログラムを起動することができます。上記でclang-formatを設定したような感じです。
実際のところどう動くかですが、まず前提としてgqも=も[オペレーターコマンド](https://vim-jp.org/vimdoc-ja/motion.html#operator)です。オペレーターコマンドについて説明していると記事がもう1本書けてしまうので省略しますが、これは移動コマンドを引数のような形で受け取り、移動コマンドで示した範囲のテキストになんらかの処理を施すものです。例えば、ノーマルモードで `diw` と打つとカーソルが載っている単語を削除しますが、これは `d` が「削除」のオペレーターコマンドで、 `iw` が「単語 (inner word)」を意味する移動コマンドだからです。移動コマンドのことはモーションとも言います。ちなみにオペレーターを2回続ける (モーションを受け取らずに) と、「カーソルのある1行」がオペレーターの対象となります。 `dd` で1行削除できるのはvimの基本ですが、裏側の理屈はこの削除オペレータが効いているわけです。

で、話を戻します。formatprgとequalprgはどちらもオペレーターコマンドなので、モーションを受け取ります。ggで先頭に移動後 gq/= してGで一番下に移動することで、formatprg/equalprgの処理対象はバッファ全体になります。ここからの動きですが、次のようになります。

* formatprg/equalprgに移動コマンドの示すテキストが標準入力として渡る
  * つまり、 `clang-format` コマンドにカレントバッファのテキストが標準入力として渡る
* clang-formatは標準入力を受け取ると、それをフォーマットして標準出力に出力する (これはvim関係なく、clang-formatコマンドの動作)
* vimはformatprg/equalprgの標準出力に出力されたテキストをカレントバッファに設定 (上書き) する

このような感じで、clang-formatをvimから実行できます。

ただ、筆者はformatprg/equalprgはコードフォーマッタとして使うにはやや使い勝手が悪いと感じています。
例えばGoプログラマはフォーマッタにgoimportsを使うと思いますが、goimportsはシンタックスに違反したプログラムを受け取ると、フォーマットに失敗してエラーメッセージを出力します。vimは出力されたものをそのままカレントバッファに上書きしてしまうので、書いていたプログラムは消えてgoimportsのエラーメッセージだけが残ってしまいます。もちろん `u` で戻せますが、これはまああまり使い勝手として良くはないですよね。なので筆者は、フォーマットはvimからは実施しておらず、gitにコミットする前などにターミナルから実行しています。
また、エラーメッセージなどが出ずフォーマットに成功したときも、 `G` を押しているのでカーソルは一番下に移動します。 Ctrl-oで元の場所に戻れますが、これもやや面倒です。

ちなみに、formatprgとequalprgはどちらを使うのか。これは筆者の結論としてはどっちでもいいです。formatprgという名前なのだからフォーマット用であるわけなので、formatprgを使っておくと間違いないのかなとは思います。equalprgは本来フィルタ用途なのになぜこの話題に出てくるかというと、=コマンドがなぜかデフォルトではインデントを整えてくれるので、なんとなくフォーマッタに使うのが良さそうに思えてしまうからです。

さらに余談ですが、 gqではなくgwというコマンドもあります。これはgqと何が違うかというと、gwはカーソルの位置がもとに戻ります。gqではGで一番下に移動してしまったのですが、gwはggした一番上に戻ってくれます。しかし、gwはformatprgを見てくれないので、フォーマット用途では使えません。なぜ、、

## makeとエラー修正

コードを書いているとき、少し書いては `make` を実行しビルドが通ることを確かめ、通らなければ修正する、通ればさらにコードを書く、というようなやり方を取る人も多いと思います。vimには組み込みの `:make` コマンドがあり、これはそのまま、makeコマンドを実行するのと同じことをしてくれます。引数も、 `:make build` のように書けばそのとおり実行できます。
vimのmakeコマンドの利点は、その出力をvimが解釈し、エラーがある箇所に簡単にジャンプができることです。
例えば、次のようなmain.cを書いてみます。

```
#include <stdio.h>
int main() {
   printf("Hello, World!");
   a
   return 0;
}
```

printfの下に `a` とだけ書かれた行を置きました。これは当然コンパイルエラーになります。次のようなMakefileを書いてみます。

```
build:
	gcc main.c
```

makeを実行すると、こんな感じでエラーが出ます。

```
$ make
gcc main.c
main.c: In function ‘main’:
main.c:4:4: error: ‘a’ undeclared (first use in this function)
    4 |    a
      |    ^
main.c:4:4: note: each undeclared identifier is reported only once for each function it appears in
main.c:4:5: error: expected ‘;’ before ‘return’
    4 |    a
      |     ^
      |     ;
    5 |    return 0;
      |    ~~~~~~
make: *** [Makefile:2: build] Error 1
```

次に、これをvimから実行してみます。 `vim main.c` で開き、 `:make` を実行します。コンソールにエラーが表示されたらEnterを押し、vimのウインドウに戻ったら `cw` でQuickfixリストを開きます。

```
main.c
|| gcc main.c
|| main.c: In function ‘main’:
main.c|4 col 4| error: ‘a’ undeclared (first use in this function)
||     4 |    a
||       |    ^
main.c|4 col 4| note: each undeclared identifier is reported only once for each function it appears in
main.c|4 col 5| error: expected ‘;’ before ‘return’
||     4 |    a
||       |     ^
||       |     ;
[Quickfix List] :make
```

Quickfixリストには、gccが出力したエラーが表示されています。ここで重要なのは、vimはこの出力を解釈できることです。解釈できるとなにが嬉しいかというと、エラーを選択してエラーがあるファイルのエラーがある行を直接vimから開くことができます。Quickfixリストの中をj/kで移動してEnterを押せば、そのエラーの箇所が開かれます。
なぜvimはgccのエラーメッセージを解釈できるのでしょうか？そもそもvimがエラーメッセージをどう解釈しているかというと、 `errorformat` というオプションに従っています。errorformatはエラーメッセージのフォーマットを指定するリストです。フォーマットなのでフォーマット指定子 ( `%f` みたいなやつ) で組み立てます。筆者の環境では、デフォルトで以下のようになっていました。

```plaintext
errorformat=%*[^"]"%f"%*\D%l: %m,"%f"%*\D%l: %m,%-G%f:%l: (Each undeclared identifier is reported only once,%-G%f:%l: for each function it appears in.),%-GIn file included from %f:%l:%c:,%-GIn file included from %f:%l:%c\,,%-GIn file included from %f:%l:%c,%-GIn file included from %f:%l,%-G%*[ ]from %f:%l:%c,%-G%*[ ]from %f:%l:,%-G%*[ ]from %f:%l\,,%-G%*[ ]from %f:%l,%f:%l:%c:%m,%f(%l):%m,%f:%l:%m,"%f"\, line %l%*\D%c%*[^ ] %m,%D%*\a[%*\d]: Entering directory %*[`']%f',%X%*\a[%*\d]: Leaving directory %*[`']%f',%D%*\a: Entering directory %*[`']%f',%X%*\a: Leaving directory %*[`']%f',%DMaking %*\a in %f,%f|%l| %m
```

errorformatはリストで、カンマで区切られます。ちょっと見づらいので、カンマで改行してみます。

```plaintext
%*[^"]"%f"%*\D%l: %m,
"%f"%*\D%l: %m,
%-G%f:%l: (Each undeclared identifier is reported only once,
%-G%f:%l: for each function it appears in.),
%-GIn file included from %f:%l:%c:,
%-GIn file included from %f:%l:%c\,
,
%-GIn file included from %f:%l:%c,
%-GIn file included from %f:%l,
%-G%*[ ]from %f:%l:%c,
%-G%*[ ]from %f:%l:,
%-G%*[ ]from %f:%l\,
,
%-G%*[ ]from %f:%l,
%f:%l:%c:%m,
%f(%l):%m,
%f:%l:%m,
"%f"\,
 line %l%*\D%c%*[^ ] %m,
%D%*\a[%*\d]: Entering directory %*[`']%f',
%X%*\a[%*\d]: Leaving directory %*[`']%f',
%D%*\a: Entering directory %*[`']%f',
%X%*\a: Leaving directory %*[`']%f',
%DMaking %*\a in %f,
%f|%l| %m
```

こうしてみると、上の方はCやC++コンパイラのエラーメッセージのフォーマットぽいな、下の方はMakeのエラーメッセージっぽいなということが解ると思います。errorformatは上記のようにリストで、先頭からエラーメッセージがフォーマットと一致するかをチェックし、一致したらそのとおり解釈して、Quickfixリストに表示する、ということをやっています。

Goプログラマの人がMakefileに `go build` を書いたらどうなるでしょうか？Goコンパイラのエラーメッセージは次のような感じです。

```
./main.go:4:2: undefined: a
```

これは、errorformatの中にある `%f:%l:%c:%m` とフォーマットが一致します。つまり、Goのエラーメッセージをvimはそのまま読めるわけですね。
他のプログラミング言語やツールなどで、エラーメッセージがerrorformatのフォーマットと一致しない場合も、自分でerrorformatにフォーマットを追加すればvimから読ませることができます。

しかし、Makefileがないプロジェクトや、Makefile以外のタスクランナーを使っているプロジェクトではどうすればよいでしょうか？こういうときは、 `makeprg` というオプションを設定すれば良いです。例えば、 `set makeprg=scons` とすれば、 `:make` したときにvimが `scons` を実行してくれます。
つまり、 `makeprg` と `errorformat` を自分のプロジェクトに合わせて設定すれば、コーディングのフィードバックループをvimから高速に回せるわけです。でも、 `makeprg` はまだしも、 `errorformat` の設定はけっこう面倒ですよね。こういうときのために、vimには `compiler` というオプションがあります。`compiler`は、要するに `makeprg` と `errorformat` 及び関連するオプションを、予めvim側で用意されたものを簡単に使うための設定です (たぶん) 。vim上で、 `:compier` してみると次のように表示されます。

```
/usr/share/vim/vim81/compiler/ant.vim
/usr/share/vim/vim81/compiler/bcc.vim
/usr/share/vim/vim81/compiler/bdf.vim
/usr/share/vim/vim81/compiler/cargo.vim
/usr/share/vim/vim81/compiler/checkstyle.vim
/usr/share/vim/vim81/compiler/context.vim
/usr/share/vim/vim81/compiler/cs.vim
/usr/share/vim/vim81/compiler/csslint.vim
/usr/share/vim/vim81/compiler/cucumber.vim
/usr/share/vim/vim81/compiler/decada.vim
/usr/share/vim/vim81/compiler/dot.vim
/usr/share/vim/vim81/compiler/erlang.vim
/usr/share/vim/vim81/compiler/eruby.vim
/usr/share/vim/vim81/compiler/fortran_F.vim
/usr/share/vim/vim81/compiler/fortran_cv.vim
/usr/share/vim/vim81/compiler/fortran_elf90.vim
/usr/share/vim/vim81/compiler/fortran_g77.vim
/usr/share/vim/vim81/compiler/fortran_lf95.vim
/usr/share/vim/vim81/compiler/fpc.vim
/usr/share/vim/vim81/compiler/g95.vim
/usr/share/vim/vim81/compiler/gcc.vim
/usr/share/vim/vim81/compiler/gfortran.vim
/usr/share/vim/vim81/compiler/ghc.vim
/usr/share/vim/vim81/compiler/gnat.vim
/usr/share/vim/vim81/compiler/go.vim
/usr/share/vim/vim81/compiler/haml.vim
/usr/share/vim/vim81/compiler/hp_acc.vim
/usr/share/vim/vim81/compiler/icc.vim
/usr/share/vim/vim81/compiler/ifort.vim
/usr/share/vim/vim81/compiler/intel.vim
/usr/share/vim/vim81/compiler/irix5_c.vim
/usr/share/vim/vim81/compiler/irix5_cpp.vim
/usr/share/vim/vim81/compiler/javac.vim
/usr/share/vim/vim81/compiler/jikes.vim
/usr/share/vim/vim81/compiler/mcs.vim
/usr/share/vim/vim81/compiler/mips_c.vim
/usr/share/vim/vim81/compiler/mipspro_c89.vim
/usr/share/vim/vim81/compiler/mipspro_cpp.vim
/usr/share/vim/vim81/compiler/modelsim_vcom.vim
/usr/share/vim/vim81/compiler/msbuild.vim
/usr/share/vim/vim81/compiler/msvc.vim
/usr/share/vim/vim81/compiler/neato.vim
/usr/share/vim/vim81/compiler/ocaml.vim
/usr/share/vim/vim81/compiler/onsgmls.vim
/usr/share/vim/vim81/compiler/pbx.vim
/usr/share/vim/vim81/compiler/perl.vim
/usr/share/vim/vim81/compiler/php.vim
/usr/share/vim/vim81/compiler/pylint.vim
/usr/share/vim/vim81/compiler/pyunit.vim
/usr/share/vim/vim81/compiler/rake.vim
/usr/share/vim/vim81/compiler/rspec.vim
/usr/share/vim/vim81/compiler/rst.vim
/usr/share/vim/vim81/compiler/ruby.vim
/usr/share/vim/vim81/compiler/rubyunit.vim
/usr/share/vim/vim81/compiler/rustc.vim
/usr/share/vim/vim81/compiler/sass.vim
/usr/share/vim/vim81/compiler/se.vim
/usr/share/vim/vim81/compiler/splint.vim
/usr/share/vim/vim81/compiler/stack.vim
/usr/share/vim/vim81/compiler/tcl.vim
/usr/share/vim/vim81/compiler/tex.vim
/usr/share/vim/vim81/compiler/tidy.vim
/usr/share/vim/vim81/compiler/xbuild.vim
/usr/share/vim/vim81/compiler/xmllint.vim
/usr/share/vim/vim81/compiler/xmlwf.vim
```

例えば、 `/usr/share/vim/vim81/compiler/perl.vim` を見てみるとperl用のmakeprgとerrorformatの設定をしています。

筆者はこの `:make` はかなり多用しています。プログラミング言語によっては設定しないと動かないですが、C++やGoを書くときにはよく使います。Makefileの中でdockerを起動したりといった複雑なことをやり始めるとこのへんのインテグレーションがやりにくくなるので注意です。

## grep

コードを書いているとgrepしたくなることは多いですよね。vimにおけるgrepは内部grepと外部grepがあります。内部grepはvimに組み込まれたgrepの処理を実施し、外部grepは外部grep (GNU grepコマンド) を実行します。後者は外部grepの結果をvimにロードするのですが、前者はファイルをバッファにロードするので、外部grepよりも遅いです。じゃあ内部grepの利点はなんなのかという話なのですが、プラットフォーム依存がなく、vimさえあれば動くことが利点です。ただ、筆者は基本常に外部grepを使っています。

内部grepは、細かいオプションを除けば次のような感じで動きます。

```
:vimgrep keyword **
```

`:vimgrep` は `:vim` でも動きます。

外部grepは、デフォルトではGNU grepコマンドを使うので、普段使っているgrepコマンドと同じように使えば良いです。筆者は以下のように普段使います。

```
:grep -ir keyword *
```

例えば拡張子を絞りたいとかも、 `man grep` を見て必要なオプションを渡せばOKです。

外部grepも内部grepも、実行したあとはその結果を見たいわけですが、これには先程も登場したQuickfixを使います。Quickfixは本来エラーメッセージを表示する機能なのだと思うのですが、grepの結果を表示してそこにジャンプするのにも使えます。 `:grep` や `:vim` した後に `:cw` するとQuickfixリストを開けます。

また、grepコマンドには、ackやag、あるいはgit grepなどいろいろAlternativeがありますが、それらは `grepprg` を設定することで使えるようです (筆者はGNU grepしか使っていないです) 。

## タグジャンプ

関数を呼び出している箇所にカーソルを合わせ、その関数の定義元に飛ぶのも、実はプラグイン無しでできるっちゃできます。これには `ctags` という外部コマンドが必要です。ctagsはUbuntuには確かプリインストールされておらず、 `apt-get install universal-ctags` でインストールします。ctagsはもともとExuberant Ctagsというプロジェクトだったのですが、こちらは開発が停止しているので、 [Universal Ctags](https://github.com/universal-ctags/ctags)を使うほうがよいです (自分のctagsがどちらなのかは `ctags --version` でわかります) 。

vimで適当なファイルを開き、 `:!ctags -R` と実行します。 `:!` は外部コマンドの実行です。これをすると `tags` ファイルができます。git管理しているプロジェクトであればこれは管理対象にしなくて良いので、 `echo tags >> .git/info/exclude` などしてignoreしましょう。
vimは標準でこのtagsファイルを読んでジャンプができます。ジャンプしたいシンボルにカーソルを合わせ `Ctrl-]` するとジャンプできます。

注意点としては、ctagsを使ったジャンプはIDEやLSPのようにコンテキストを正しく理解した賢いものではありません。例えば `get()` という関数呼び出しにカーソルを合わせ `Ctrl-]` した時に、importなどで呼び出されるはずのない `get()` 関数もジャンプの対象になります。また、 `get` という名前の変数があったとしてもジャンプの対象になるようです。 `Ctrl-]` の代わりに `gCtrl-]` とすると候補をリストアップしてくれます。
IDEやLSPのように賢くないとは書きましたが、実際それが問題になることは筆者としてはほぼない印象です。ただ、クリーンアーキテクチャ的なことをしていると、同じような名前のシンボルをいろんなところで作ることがあると思うのでそれはちょっと相性が悪いかもしれないです。

## 補完

最後は補完です。前提として、これもctagsのジャンプと同様、importしているライブラリの関数名だけを補完するような「賢い」補完はvimだけではできません。そういうのが必要ならIDEやLSPを使いましょう。

vimの補完は Ctrl-xでトリガーできます。Ctrl-xしたあと、どのタイプの補完をしたいのかをさらに入力します。例えば、次のようなコードを考えます。

```
bool long_func_name() {
	return true;
}


int main() {
	// ここでlong_func_nameを呼び出したい
}
```

この時、 `long_func_name` を呼び出したい場所でインサートモードに入り、最初の一文字である `l` を入力して、そのまま `Ctrl-x Ctrl-n` と押します。ちょっとややこしいのですが、「Ctrlキーを押す」「Ctrlキーを押したまま、xを押す」「Ctrlキーは離さず、押したまま今度はnを押す」という動きです。これで `long_func_name` が勝手に補完されるはずです。
この動きは、Ctrl-nがカレントファイルから単語を補完するという補完タイプなことから、この動きが実現されています。
補完タイプには他にもいくつかあります。よく使うものは:

* Ctrl-f: ファイル名補完
* Ctrl-l: 行補完
* Ctrl-]: タグ補完

などです。これらはすべて、 Ctrl-xに続けて入力する必要があります。また、Ctrl-nは実際はカレントファイルだけではなく別のバッファ、タグファイル、インクルードファイルなどからも補完対象を検索してくれます。

Ctrl-lの行補完はけっこう便利です。Goプログラマの人は `if err != nil` を何度も書くことになると思うのですが、これは行補完すれば `i` と入力し、 `Ctrl-x Ctrl-l` で補完する、という感じで瞬時に入力できます。

## まとめ

長々と書いてきましたが、細かいことはまだ書けてないことも多いと思います。ここで紹介したのは、筆者が実際にプログラミングのときに使っているテクニックです。

筆者がvimを使うときは、大体シェルから `vi main.py` などとファイルを開いて、 `:e` か `:fin` でバッファを開き、開きすぎたなと思ったら `:bd` で消します。バッファ間の移動は `:b` にTabキーでなんとかしてます。検索したいときは `:grep -ir xxx *` と打ちます。補完は Ctrl-x Ctrl-nでなんとかなります。 `:!ctags -R` さえ実行しておけば、ジャンプもできます。コードが書けたら `:make` します。 `:make test` もします。

筆者は以前はLSPを使っていたのですが、補完とジャンプくらいしか使っていないことに気付き、vim本来の機能でそれらを達成できないかいろいろ試した結果、今のスタイルに落ち着きました。IDEやプラグインでもっと便利になることはもちろんあると思うのですが、シンプルさやポータビリティが筆者にとってはとても大きなメリットなので、今の形をとりあえず今後も続けていく予定です。
またこの記事は、IDEやプラグインを使っている人に、IDEやプラグインは使わないほうがいいよ、と主張するための記事ではありません。こういうのは大前提として、好みとか慣れの問題だと思っています。他人がどんな風にコードを書いているのかはけっこう面白い話だと思っていて、それを自分も書いてみたくなったので書きました。
最後まで読んでいただきありがとうございました!

## 参考

* [How to Do 90% of What Plugins Do (With Just Vim)](https://thoughtbot.com/blog/how-to-do-90-of-what-plugins-do-with-just-vim)
* [実践Vim](https://tatsu-zine.com/books/practical-vim)
