UNIXドメインソケットについて---2018-07-30 19:09:52

ISUCONの文脈で、UNIXドメインソケットという言葉がよく出てくる。
詳しく知らなかったので、調べた。

## Abstract
ISUCONの文脈では、通信をTCPからUNIXドメインソケットを利用するように切り替えることで高速化が図れる、ということを言う。
通常、プロセス間通信はポートを通じてやり取りを行う。
しかし、UNIXドメインソケット通信はポートを使わず、ファイルシステムを利用して通信する。

## ファイルシステム
UNIXドメインソケットでは、ファイルをインタフェースに使ってプロセス間通信を行う。
この時、実際にファイルを生成するわけではなく、ソケットファイルという特殊なファイルが使われている。

## 見る

以下のように確認できる。

```
$ netstat -lt --protocol=unix

Active UNIX domain sockets (servers and established)
Proto RefCnt Flags       Type       State         I-Node   PID/Program name     Path
unix  2      [ ACC ]     STREAM     LISTENING     37610    2367/gnome-session-  @/tmp/.ICE-unix/2367
unix  2      [ ACC ]     STREAM     LISTENING     29969    -                    @/tmp/.ICE-unix/1378
unix  2      [ ACC ]     STREAM     LISTENING     14378499 -                    /run/cups/cups.sock
unix  2      [ ]         DGRAM                    37368    2323/systemd         /run/user/1000/systemd/notify
unix  2      [ ]         DGRAM                    29920    -                    /run/user/120/systemd/notify
### 後略
```

## Go

Goでは、unixドメインソケットで使用するソケットファイルを、以下のようなコードで作ることができる。

```
package main

import (
  "fmt"
  "net"
)

func main() {
  listener, err := net.Listen("unix", "socketfile")
  if err != nil {
    panic(err)
  }
  defer listener.Close()
  _, err := listener.Accept()
   if err != nil {
    panic(err)
   }  
}  
```
 `net.Listen` の第一引数を `unix` にする。

上記を `go run` して、別シェルから以下のように実行することで、ソケットファイルが作成されるのがわかる。

```
$ netstat -lt --protocol=unix | grep socketfile
unix  2      [ ACC ]     STREAM     LISTENING     16241224 socketfile
```

逆に、上記で作られた `socketfile` を、使用するクライアントを作るには、 `net.Dial` を使用する。

```
conn, err := net.Dial("unix", "socketfile")
if err != nil {
  panic(err)
}
```

## 使い所
ISUCONでは、同じインスタンス内にアプリケーションサーバとMySQLが同居している、みたいなことがある。
そのような場合は、ひとつのOSの中なのでUNIXドメインソケット通信が使用できる。
nginxの場合は
https://github.com/shirokanezoo/isucon5/commit/0d0f8c6792a0289f97f778d8d1dd7b44d0c3dc48
な感じ、
MySQLの場合は、my.cnfでソケットのパスを見て、アプリ側で指定する。
