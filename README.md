# prpc
protobuf rpc

```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative rpc.proto
```

``` go
// 连接 grpc
cli := prpc.NewRPCClient(prpc.WithConnectAddress("localhost:50051"))

// ssl
cli := prpc.NewRPCClient(prpc.WithConnectAddress("prpc.test.com:443"),prpc.WithConnectTLSSkipCheck())
```

``` nginx

// nginx conf

server {
    listen                  443 ssl http2;
    listen                  [::]:443 ssl http2;
    server_name             prpc.test.com;

    # SSL
    ssl_certificate         /etc/ssl/wildcard.test.com.crt;
    ssl_certificate_key     /etc/ssl/wildcard.test.com.key;

    underscores_in_headers on;

    # reverse proxy
    location /  {
        grpc_read_timeout 300s;
        grpc_send_timeout 300s;
        grpc_socket_keepalive on;
        grpc_pass grpc://127.0.0.1:50051;
    }
}

```


## 使用方式

### 1. 定义pb文件

```
syntax = "proto3";

option go_package = "pb/";

message PingMessage {
  string name = 1;
}
```

### 2. 编译pb

```
protoc --go_out=./pb --go_opt=paths=source_relative message.proto
```

### 3. server

```
type H struct{}
func (h *H) HandleStreamConnected(context prpc.RPCStreamContext) {}
func (h *H) HandleStreamDisconnected(context prpc.RPCStreamContext) {}
func (h *H) HandleRequest(ctx prpc.RPCRequestContext) {
    ctx.Reply(200, ctx.Message())
}
func main() {
    handler := &H{}
    s := prpc.NewServer(handler, prpc.WithListenAddress(":50051"))
    s.Start()
}
```

### 4. client

```
func main() {
    cli := prpc.NewRPCClient(prpc.WithConnectAddress("127.0.0.1:50051"))
    err := cli.Connect()
    if err != nil {
        fmt.Println(err)
        return
    }
    for {
        call := cli.Request("math.add", &pb.PingMessage{
            Name: "Hello",
        }, time.Hour)
        if call.Error != nil {
            fmt.Println(call.Error)
        } else {
            fmt.Println(call.Status, call.Message)
        }
        time.Sleep(time.Second)
    }
}
```
