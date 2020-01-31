go1.11のmodulesの使い方について---2018-10-04 19:15:39

### Modulesとは

goのパッケージ管理ツールはglide -> depと進化を遂げたが、go 1.11の世界では新たに[Modules](https://github.com/golang/go/wiki/Modules)という機構が生まれ、`go mod` コマンドでmodulesを管理することができるようになった。(言語レベルでモジューリングがサポートされたことを意味する)

goの思想としてもdepではなくmodulesに置き換えていくことを推進してるっぽいので、早めに直しておくと良いと思う。
今回使ってみて、今後は必ず使うと思う。そのくらい便利だった。

### Requirements

* Go 1.11以上

### go modの使い方

## すでにdepを使っていない場合(新規プロジェクトなど)

ローカルに[mattn/todo](https://github.com/mattn/todo)をgo getしてやってみた。
このリポジトリでは、

* github.com/gonuts/commander
* github.com/gonuts/flag

をimportしている。

```shell
$ export GO111MODULE=on
$ cd $GOPATH/src/github.com/mattn/todo
$ go mod init
$ cat go.mod
module github.com/mattn/todo
```
go.modが作られる。

```shell
$ go build

### buildするとgo.modが編集され、またgo.sumが作られる
$ cat go.mod
require (
        github.com/gonuts/commander v0.1.0
        github.com/gonuts/flag v0.1.0
)

$ cat go.sum
github.com/gonuts/commander v0.1.0 h1:EcDTiVw9oAVORFjQOEOuHQqcl6OXMyTgELocTq6zJ0I=
github.com/gonuts/commander v0.1.0/go.mod h1:qkb5mSlcWodYgo7vs8ulLnXhfinhZsZcm6+H/z1JjgY=
github.com/gonuts/flag v0.1.0 h1:fqMv/MZ+oNGu0i9gp0/IQ/ZaPIDoAZBOBaJoV7viCWM=
github.com/gonuts/flag v0.1.0/go.mod h1:ZTmTGtrSPejTo/SRNhCqwLTmiAgyBdCkLYhHrAoBdz4=
```

* ここでは、vendorディレクトリは作られない
* パッケージが、GOPATH配下にダウンロードされてしまうこともない

*  パッケージを増やすには？go getする

```shell
$ go get github.com/golang/mock
go: finding github.com/golang/mock v1.1.1
go: downloading github.com/golang/mock v1.1.1

$ cat go.sum
github.com/golang/mock v1.1.1 h1:G5FRp8JnTd7RQH5kemVNlMeyXQAztQ3mOWV95KxsXH8=
github.com/golang/mock v1.1.1/go.mod h1:oTYuIxOrZwtPieC+H1uAHpcLFnEyAGVDL/k47Jfbm0A=
github.com/gonuts/commander v0.1.0 h1:EcDTiVw9oAVORFjQOEOuHQqcl6OXMyTgELocTq6zJ0I=
github.com/gonuts/commander v0.1.0/go.mod h1:qkb5mSlcWodYgo7vs8ulLnXhfinhZsZcm6+H/z1JjgY=
github.com/gonuts/flag v0.1.0 h1:fqMv/MZ+oNGu0i9gp0/IQ/ZaPIDoAZBOBaJoV7viCWM=
github.com/gonuts/flag v0.1.0/go.mod h1:ZTmTGtrSPejTo/SRNhCqwLTmiAgyBdCkLYhHrAoBdz4=

$ cat go.mod
module github.com/mattn/todo

require (
        github.com/golang/mock v1.1.1 // indirect
        github.com/gonuts/commander v0.1.0
        github.com/gonuts/flag v0.1.0
)
```

* 増えている
* 当然、GOPATH配下にはダウンロードされないし、vendorが作られるわけでもない

go getせずに、直接ソースコードに書くとどうなるか？以下の行をいい感じに追加

```go
import "github.com/kataras/golog"
golog.Println("This is a sample log message.")
```

```shell
$ go build #これは必要
$ cat go.mod
module github.com/mattn/todo

require (
        github.com/golang/mock v1.1.1 // indirect
        github.com/gonuts/commander v0.1.0
        github.com/gonuts/flag v0.1.0
        github.com/kataras/golog v0.0.0-20180321173939-03be10146386
        github.com/kataras/pio v0.0.0-20180511174041-a9733b5b6b83 // indirect
)

$ cat go.sum
github.com/golang/mock v1.1.1 h1:G5FRp8JnTd7RQH5kemVNlMeyXQAztQ3mOWV95KxsXH8=
github.com/golang/mock v1.1.1/go.mod h1:oTYuIxOrZwtPieC+H1uAHpcLFnEyAGVDL/k47Jfbm0A=
github.com/gonuts/commander v0.1.0 h1:EcDTiVw9oAVORFjQOEOuHQqcl6OXMyTgELocTq6zJ0I=
github.com/gonuts/commander v0.1.0/go.mod h1:qkb5mSlcWodYgo7vs8ulLnXhfinhZsZcm6+H/z1JjgY=
github.com/gonuts/flag v0.1.0 h1:fqMv/MZ+oNGu0i9gp0/IQ/ZaPIDoAZBOBaJoV7viCWM=
github.com/gonuts/flag v0.1.0/go.mod h1:ZTmTGtrSPejTo/SRNhCqwLTmiAgyBdCkLYhHrAoBdz4=
github.com/kataras/golog v0.0.0-20180321173939-03be10146386 h1:VT6AeCHO/mc+VedKBMhoqb5eAK8B1i9F6nZl7EGlHvA=
github.com/kataras/golog v0.0.0-20180321173939-03be10146386/go.mod h1:PcaEvfvhGsqwXZ6S3CgCbmjcp+4UDUh2MIfF2ZEul8M=
github.com/kataras/pio v0.0.0-20180511174041-a9733b5b6b83 h1:NoJ+fI58ptwrPc1blX116i+5xWGAY/2TJww37AN8X54=
github.com/kataras/pio v0.0.0-20180511174041-a9733b5b6b83/go.mod h1:NV88laa9UiiDuX9AhMbDPkGYSPugBOV6yTZB1l2K9Z0=
```

### ポイント

* 最大のポイントは、importするリポジトリのソースコードはローカルのどこにもないということ
  - downloading...みたいなログは出るので、たぶんダウンロードしてバイナリ作って消してるみたいな気がする
* vendorディレクトリはもういらない!!!これはすごいと思った。
* go getするか、ソースコードに依存を書き込めば、あとは普通にgo buildするだけで、goが勝手に依存を組み込んでバイナリを生成してくれる
* じゃあ普通にソースコード取ってくるときのgo getはもうできないの？というと、そうではない
  - そこを切り替えているのが、一番最初に設定した `export GO111MODULE=on`
  - Dockerfileに書くのがおすすめ
