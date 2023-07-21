package main

import (
    "fmt"
    "github.com/DGHeroin/prpc"
    "log"
    "time"
)

func main() {
    cli := prpc.NewRPCClient(prpc.WithConnectAddress("localhost:50051"))
    err := cli.Connect()
    if err != nil {
        fmt.Println(err)
        return
    }

    s, err := cli.BuildStreaming()
    err = s.Start(func(ch <-chan prpc.StreamMessage) {
        log.Println("[client] stream open")
        s.Send(prpc.StreamMessage{
            Name:   "[client] stream ready",
            Status: 122,
        })
        for msg := range ch {
            log.Println("[client] recv message", msg.Name)
        }
        log.Println("[client] stream closed")
    })
    if err != nil {
        fmt.Println("err:", err)
        return
    }

    for {
        call := cli.Request("hello world", nil, time.Second)
        if call.Error == nil {
            fmt.Println("[client] server response:", call.Status, call.Message)
        } else {
            fmt.Println("[client] request error:", err)
        }
        time.Sleep(time.Second * 5)
    }
}
