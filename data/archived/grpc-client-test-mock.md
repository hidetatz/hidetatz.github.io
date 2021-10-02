gRPCサーバーをモックしてテストする---2019-02-11 00:36:30

gRPC越しに別のマイクロサービスをコールするクライアントの層をテストしたいことがある。
例えば、次のようなprotoを使ってみよう。

```
syntax = "proto3";

package user.v1;

service Users {
  rpc Get(GetRequest) returns (GetResponse);
}

message GetRequest {
  string id = 1;
}

message GetResponse {
  User user = 1;
}

message User {
  string id    = 1; 
  string name  = 2;
  string email = 3;
}
```

UsersをGetするRPCを作っている。

これを、 `protoc --go_out=plugins=grpc:pb ./users.proto` と実行してみる。
すると、[このような](https://github.com/yagi5/grpc-test/blob/master/pb/users.pb.go)ファイルが自動生成される。

これを素朴に使ってクライアント、サーバを実装するとこんな感じになる。

* Server

```
package main

import (
  "context"
  "fmt"
  "log"
  "net"

  userpb "github.com/yagi5/grpc-test/pb"
  "google.golang.org/grpc"
)

type server struct{}

func (s *server) Get(ctx context.Context, req *userpb.GetRequest) (*userpb.GetResponse, error) {
  if req.Id == "1" {
    return &userpb.GetResponse{User: &userpb.User{
      Id:    "1",
      Name:  "bob",
      Email: "email@example.com",
    }}, nil
  }
  return nil, fmt.Errorf("id invalid")
}

func main() {
  listenPort, err := net.Listen("tcp", ":19003")
  if err != nil {
    log.Fatalln(err)
  }
  grpcServer := grpc.NewServer()
  userpb.RegisterUsersServer(grpcServer, &server{})
  grpcServer.Serve(listenPort)
}
```

Getの実装はサンプルであり、実際はRDBやKVSから取得する。

* Client

```
package main

import (
  "context"
  "fmt"

  userpb "github.com/yagi5/grpc-test/pb"
  "google.golang.org/grpc"
)

type user struct {
  ID    string
  Name  string
  Email string
}

func main() {
  conn, err := grpc.Dial("127.0.0.1:19003", grpc.WithInsecure())
  if err != nil {
    panic(err)
  }
  defer conn.Close()
  client := userpb.NewUsersClient(conn)

  c := usersClient{client: client}
  res, err := c.Get(ctx, &userpb.GetRequest{Id: "1"})
  if err != nil {
    panic(err)
  }

  fmt.Printf("user: %#v\n", u)
}
```

マイクロサービスのコンテキストにおいては、このgRPCサーバーが実態としては別のマイクロサービスであるということがほとんどである。
このとき、ユニットテストでこのサーバー(Userマイクロサービス)を実際に起動するのはなかなか実装がめんどくさい。
telepresenceなどを使えば動作確認は可能だが、それが本番で同じようにつながる保証もないし、ユニットテストではモックするという選択を筆者はよくする。
どうモックするかというと、userpb.UsersClientが[interfaceになっている](https://github.com/yagi5/grpc-test/blob/master/pb/users.pb.go#L192-L194)ことを利用する。
以下のようにリファクタする。

```
package main

import (
  "context"
  "fmt"

  userpb "github.com/yagi5/grpc-test/pb"
  "google.golang.org/grpc"
)

type usersClient struct {
  grpcClient userpb.UsersClient
}

type user struct {
  ID    string
  Name  string
  Email string
}

func main() {
  conn, err := grpc.Dial("127.0.0.1:19003", grpc.WithInsecure())
  if err != nil {
    panic(err)
  }
  defer conn.Close()
  client := userpb.NewUsersClient(conn)

  c := usersClient{grpcClient: client}
  u, err := c.getUser(context.Background(), "1")
  if err != nil {
    panic(err)
  }

  fmt.Printf("user: %#v\n", u)
}

func (c *usersClient) getUser(ctx context.Context, id string) (*user, error) {
  res, err := c.grpcClient.Get(ctx, &userpb.GetRequest{Id: id})
  if err != nil {
    return nil, err
  }
  return &user{ID: res.User.Id, Name: res.User.Name, Email: res.User.Email}, nil
}
```

getUser関数をusersClientのメソッドとする。usersClientは内部でuserpb.UsersClientを持つ。
userpb.UsersClientはinterfaceになっているため、モックで差し替えて、getUser内のc.grpcClient.Get()で任意の結果を返せるようにする。

テストコードは以下のように書ける。

```
package main

import (
  "context"
  "fmt"
  "reflect"
  "testing"

  userpb "github.com/yagi5/grpc-test/pb"
  "google.golang.org/grpc"
)

type mockUsersClient struct {
  MockGetFunc func(context.Context, *userpb.GetRequest, ...grpc.CallOption) (*userpb.GetResponse, error)
}

func (c *mockUsersClient) Get(ctx context.Context, req *userpb.GetRequest, opts ...grpc.CallOption) (*userpb.GetResponse, error) {
  return c.MockGetFunc(ctx, req, opts...)
}

func TestGetUser(t *testing.T) {
  var tests = []struct {
    name    string
    getMock func(context.Context, *userpb.GetRequest, ...grpc.CallOption) (*userpb.GetResponse, error)
    id      string
    wantErr bool
    res     *user
  }{
    {
      name: "returns error",
      getMock: func(context.Context, *userpb.GetRequest, ...grpc.CallOption) (*userpb.GetResponse, error) {
        return nil, fmt.Errorf("")
      },
      id:      "1",
      wantErr: true,
    },
    {
      name: "returns non-error",
      getMock: func(context.Context, *userpb.GetRequest, ...grpc.CallOption) (*userpb.GetResponse, error) {
        return &userpb.GetResponse{User: &userpb.User{Id: "1", Name: "bob", Email: "email@example.com"}}, nil
      },
      id:  "1",
      res: nil,
    },
  }
  for _, tt := range tests {
    tt := tt
    t.Run(tt.name, func(t *testing.T) {
      c := &mockUsersClient{MockGetFunc: tt.getMock}
      cl := usersClient{grpcClient: c}
      u, err := cl.getUser(context.Background(), tt.id)
      if (err != nil) != tt.wantErr {
        t.Fatalf("failed expected: %#v but got %#v", tt.wantErr, err)
      }
      if reflect.DeepEqual(tt.res, u) {
        t.Errorf("failed expected: %#v but got %#v", tt.res, u)
      }
    })
  }
}
```

mockUsersClientを作成し、中にはモックさせる関数を持つ。
mockUsersClientはGet()を実装しているため、userpb.UsersClientを実装していることになる。

```
{
  name: "returns error",
  getMock: func(context.Context, *userpb.GetRequest, ...grpc.CallOption) (*userpb.GetResponse, error) {
    return nil, fmt.Errorf("")
  },
  id:      "1",
  wantErr: true,
},
```

このテストは、errを返す。

```
{
  name: "returns non-error",
  getMock: func(context.Context, *userpb.GetRequest, ...grpc.CallOption) (*userpb.GetResponse, error) {
    return &userpb.GetResponse{User: &userpb.User{Id: "1", Name: "bob", Email: "email@example.com"}}, nil
  },
  id:  "1",
  res: nil,
},
```

このテストは正しい結果を返す。
あとはgetMockを使ってusersClientを生成し、getUserをテストできる。

このように、mockClientを作り、そこにどのような関数を渡すかを各テストで定義するのは、特にinfraレイヤーでかなり使えるテクニックである。
