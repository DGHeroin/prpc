package prpc

import (
    "context"
    "crypto/tls"
    "github.com/DGHeroin/prpc/pb"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/connectivity"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/status"
    "google.golang.org/protobuf/proto"
    "google.golang.org/protobuf/types/known/anypb"
    "sync"
    "time"
)

type (
    RPCClient interface {
        Connect() error
        Request(name string, d proto.Message, timeout time.Duration) *RequestCall
        BuildStreaming() (*StreamCall, error)
    }
    client struct {
        mu     sync.RWMutex
        conn   *grpc.ClientConn
        client pb.RPCServiceClient
        option *clientOption
    }
    RequestCall struct {
        Status  int32
        Message interface{}
        Error   error
    }
    StreamCall struct {
        rpcClient *client
        sendFn    func(message StreamMessage) error
        closeFn   func() error
        closed    bool
    }
)

func (c *client) Connect() error {
    return c.reconnect()
}

func NewRPCClient(opts ...ClientOption) RPCClient {
    c := &client{
        option: &clientOption{},
    }
    for _, fn := range opts {
        fn(c.option)
    }
    return c
}

func (c *client) Request(name string, message proto.Message, timeout time.Duration) *RequestCall {
    var call = &RequestCall{}
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    var data *anypb.Any
    if v3data, err := anypb.New(message); err == nil {
        data = v3data
    }
    r, err := c.client.RequestCall(ctx, &pb.RPCRequest{Name: name, Data: data})
    if err != nil {
        call.Error = err
        return call
    }
    call.Error = nil
    call.Status = r.GetStatus()
    call.Message = r.GetData()
    return call
}

func (c *client) BuildStreaming() (*StreamCall, error) {
    sc := &StreamCall{
        rpcClient: c,
    }
    return sc, nil
}

func (c *client) reconnect() error {
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.conn != nil {
        if c.conn.GetState() == connectivity.Ready {
            return nil
        }
    }
    var options = []grpc.DialOption{grpc.WithBlock()}
    if c.option.Credentials == nil {
        options = append(options, grpc.WithInsecure())
    } else {
        options = append(options, grpc.WithTransportCredentials(c.option.Credentials))
    }
    if c.option.TLSEnable {
        options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
            InsecureSkipVerify: c.option.TLSSkipCheck, // Skip TLS verification for Cloudflare proxy
        })))
    }
    conn, err := grpc.Dial(c.option.Address, options...)
    if err != nil {
        return err
    }
    c.conn = conn
    c.client = pb.NewRPCServiceClient(conn)
    return nil
}
func (c *client) Close() error {
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.conn != nil {
        err := c.conn.Close()
        return err
    }
    return nil
}
func (sc *StreamCall) Send(message StreamMessage) error {
    return sc.sendFn(message)
}
func (sc *StreamCall) Start(cb func(<-chan StreamMessage)) error {
    ch := make(chan StreamMessage)
    ctx := context.Background()
    stream, err := sc.rpcClient.client.Streaming(ctx)
    if err != nil {
        return err
    }
    sc.sendFn = func(message StreamMessage) error {
        msg := message.toRPCStreamMessage()
        if msg == nil {
            return ErrStreamMessageInvalid
        }
        return stream.Send(msg)
    }
    sc.closeFn = func() error {
        return stream.CloseSend()
    }
    go cb(ch)
    go func() {
        defer func() {
            close(ch)
        }()
        for !sc.closed {
            msg, err := stream.Recv()
            if err == nil {
                ch <- StreamMessage{
                    Name:    msg.GetName(),
                    Status:  msg.GetStatus(),
                    Message: msg.GetData(),
                }
                continue
            }
            if err != nil {
                if status.Code(err) == codes.Unavailable {

                }
                sc.reconnectAndStart(cb, time.Second)
                break
            }
        }
    }()

    return nil
}
func (sc *StreamCall) reconnectAndStart(cb func(<-chan StreamMessage), delay time.Duration) {
    if delay != 0 {
        time.Sleep(delay)
    }
    if sc.closed {
        return
    }
    if err := sc.rpcClient.reconnect(); err != nil {
        sc.reconnectAndStart(cb, delay)
        return
    }
    if err := sc.Start(cb); err != nil {
        sc.reconnectAndStart(cb, delay)
        return
    }
}
func (sc *StreamCall) Close() error {
    sc.closed = true
    return sc.closeFn()
}
