package main

import (
	"net"
	"net/rpc"

	"golang_learn/golang_learn/gop/new_helloworld/handler"
	"golang_learn/golang_learn/gop/new_helloworld/server_proxy"
)

func main() {
	// 1 实例化一个server
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		return
	}

	// 2 注册处理逻辑
	err = server_proxy.RegisterHelloService(&handler.NewHelloService{})
	if err != nil {
		return
	}
	for {
		// 3 启动服务
		conn, _ := listen.Accept() // 当有一个请求到来，就会创建一个套接字
		go rpc.ServeConn(conn)
	}

}
