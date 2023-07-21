package prpc

import (
    "crypto/tls"
    "google.golang.org/grpc/credentials"
)

type clientOption struct {
    Credentials  credentials.TransportCredentials
    Address      string
    TLSSkipCheck bool
}
type ClientOption func(option *clientOption)

func WithClientTLS(c *tls.Config) ClientOption {
    return func(o *clientOption) {
        o.Credentials = credentials.NewTLS(c)
    }
}
func WithConnectAddress(addr string) ClientOption {
    return func(o *clientOption) {
        o.Address = addr
    }
}
func WithConnectTLSSkipCheck() ClientOption {
    return func(o *clientOption) {
        o.TLSSkipCheck = true
    }
}
