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
