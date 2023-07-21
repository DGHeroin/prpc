package main

import (
    "github.com/DGHeroin/prpc"
    "log"
)

type H struct{}

func (h *H) HandleRequest(ctx prpc.RPCRequestContext) {
    log.Printf("%p Received: %v", ctx, ctx.Name())
    ctx.Reply(200, ctx.Message())
}

func (h *H) HandleStreamConnected(ctx prpc.RPCStreamContext) {
    log.Println("streaming open", ctx)

    go func() {
        for msg := range ctx.Read() {
            ctx.Write(prpc.StreamMessage{
                Name:    "[server] a stream connected",
                Status:  200,
                Message: msg.Message,
            })
        }
    }()
    ctx.Write(prpc.StreamMessage{
        Name:   "hello",
        Status: 200,
    })
}

func (h *H) HandleStreamDisconnected(ctx prpc.RPCStreamContext) {
    log.Println("streaming closed", ctx)
}

func main() {
    handler := &H{}
    s := prpc.NewServer(handler, prpc.WithListenAddress(":50051"))
    s.Start()
}
