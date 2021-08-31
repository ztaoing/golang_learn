package client_proxy

import (
	"net/rpc"

	"golang_learn/golang_learn/gop/new_helloworld/handler"
)

type HelloServiceStub struct {
	*rpc.Client
}

func NewHelloServiceClient(protocol, address string) HelloServiceStub {
	client, err := rpc.Dial(protocol, address)
	if err != nil {
		panic("connect error")
	}

	return HelloServiceStub{
		client,
	}
}
func (h HelloServiceStub) Hello(request string, reply *string) error {
	err := h.Call(handler.HelloServiceName+".Hello", "request for something", &reply)
	return err
}
