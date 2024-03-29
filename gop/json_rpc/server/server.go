package main

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloService struct {
}

func (h *HelloService) Hello(request string, reply *string) error {
	*reply = "get request " + request
	return nil
}
func main() {
	// 1 实例化一个server
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		return
	}

	// 2 注册处理逻辑
	err = rpc.RegisterName("HelloService", &HelloService{})
	if err != nil {
		return
	}

	// 3 启动服务
	for {

		conn, err := listen.Accept() // 当有一个请求到来，就会创建一个套接字
		if err != nil {
			return
		}
		// 将套接字交给rpc处理
		// rpc 解决了call id 和序列化和反序列化的问题
		// 使用ServeCodec来指定编码方式
		// ServeCodec is like ServeConn but uses the specified codec to
		// decode requests and encode responses.

		// 当多个连接同时到达时，采用goroutine的方式异步处理
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}

}
