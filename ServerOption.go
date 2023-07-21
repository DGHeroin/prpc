package prpc

import (
	"crypto/tls"
	"google.golang.org/grpc/credentials"
	"net"
)

type serverOption struct {
	Address     string
	Credentials credentials.TransportCredentials
	Listener    net.Listener
}
type ServerOption func(o *serverOption)

func WithListenAddress(addr string) ServerOption {
	return func(o *serverOption) {
		o.Address = addr
	}
}
func WithListener(listener net.Listener) ServerOption {
	return func(o *serverOption) {
		o.Listener = listener
	}
}

func WithServerTLS(c *tls.Config) ServerOption {
	return func(o *serverOption) {
		o.Credentials = credentials.NewTLS(c)
	}
}
