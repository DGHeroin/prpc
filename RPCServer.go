package prpc

import (
    "context"
    "fmt"
    "github.com/DGHeroin/prpc/pb"
    "google.golang.org/grpc"
    "google.golang.org/protobuf/proto"
    "google.golang.org/protobuf/types/known/anypb"
    "net"
)

var (
    ErrStreamMessageInvalid = fmt.Errorf("stream message invalid")
)

type MessagePayload proto.Message

type RPCServer interface {
    Start() error
}
type RPCRequestContext interface {
    Context() context.Context
    Name() string
    Message() MessagePayload
    Reply(int32, MessagePayload)
}
type RPCStreamContext interface {
    Read() <-chan StreamMessage
    Write(StreamMessage) error
}
type StreamMessage struct {
    Name    string
    Status  int32
    Message interface{}
}
type RPCServerHandler interface {
    HandleRequest(ctx RPCRequestContext)
    HandleStreamConnected(RPCStreamContext)
    HandleStreamDisconnected(RPCStreamContext)
}
type (
    server struct {
        pb.UnimplementedRPCServiceServer
        handler RPCServerHandler
        option  *serverOption
    }
    streamContext struct {
        sendFn func(message StreamMessage) error
        ch     chan StreamMessage
    }
    requestContext struct {
        ctx   context.Context
        data  *pb.RPCRequest
        reply *pb.RPCResponse
    }
)

func (s *streamContext) Read() <-chan StreamMessage {
    return s.ch
}

func (s *streamContext) Write(message StreamMessage) error {
    return s.sendFn(message)
}

func (r *requestContext) Context() context.Context {
    return r.ctx
}

func (r *requestContext) Name() string {
    return r.data.Name
}

func (r *requestContext) Message() MessagePayload {
    return r.data.Data
}

func (r *requestContext) Reply(statusCode int32, message MessagePayload) {
    data, _ := anypb.New(message)
    r.reply = &pb.RPCResponse{
        Status: statusCode,
        Data:   data,
    }
}
func (message *StreamMessage) toRPCStreamMessage() *pb.RPCStreamMessage {
    var (
        data *anypb.Any
    )
    if v2data, ok := message.Message.(proto.Message); ok {
        if v3data, err := anypb.New(v2data); err == nil {
            data = v3data
        }
    }

    return &pb.RPCStreamMessage{
        Data:   data,
        Name:   message.Name,
        Status: message.Status,
    }
}
func NewServer(handler RPCServerHandler, opts ...ServerOption) RPCServer {
    srv := &server{
        handler: handler,
        option: &serverOption{
            Address: "127.0.0.1:5555",
        },
    }
    for _, fn := range opts {
        fn(srv.option)
    }
    return srv
}
func (s *server) Start() error {
    var (
        err error
    )
    if s.option.Listener == nil {
        s.option.Listener, err = net.Listen("tcp", s.option.Address)
        if err != nil {
            return err
        }
    }
    var opts []grpc.ServerOption
    if s.option.Credentials != nil {
        opts = append(opts, grpc.Creds(s.option.Credentials))
    }

    grpcServer := grpc.NewServer(opts...)
    pb.RegisterRPCServiceServer(grpcServer, s)
    return grpcServer.Serve(s.option.Listener)
}

func (s *server) RequestCall(ctx context.Context, in *pb.RPCRequest) (*pb.RPCResponse, error) {
    wrapCtx := &requestContext{
        ctx:  ctx,
        data: in,
    }
    s.handler.HandleRequest(wrapCtx)
    if wrapCtx.reply == nil {
        wrapCtx.reply = &pb.RPCResponse{}
    }
    return wrapCtx.reply, nil
}

func (s *server) Streaming(stream pb.RPCService_StreamingServer) error {
    ctx := &streamContext{
        ch: make(chan StreamMessage, 10),
        sendFn: func(message StreamMessage) error {
            msg := message.toRPCStreamMessage()
            if msg == nil {
                return ErrStreamMessageInvalid
            }
            return stream.Send(msg)
        },
    }
    s.handler.HandleStreamConnected(ctx)
    defer func() {
        close(ctx.ch)
        s.handler.HandleStreamDisconnected(ctx)
    }()
    // 读取消息
    for {
        msg, err := stream.Recv()
        if err != nil {
            return err
        }
        ctx.ch <- StreamMessage{
            Name:    msg.GetName(),
            Status:  msg.GetStatus(),
            Message: msg.GetData(),
        }
    }
}
