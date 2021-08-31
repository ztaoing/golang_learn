package server_proxy

import (
	"net/rpc"

	"golang_learn/golang_learn/gop/new_helloworld/handler"
)

type HelloServicer interface {
	Hello(request string, reply *string) error
}

// RegisterHelloService 使用接口的方式来解耦
func RegisterHelloService(srv HelloServicer) error {
	return rpc.RegisterName(handler.HelloServiceName, srv)
}
