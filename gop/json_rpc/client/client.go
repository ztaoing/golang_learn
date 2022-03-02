package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	// 这里直接使用net进行拨号，其中不包括编解码的方式，在下满指定编解码的方式
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		panic("连接失败")
	}
	// var reply *26string = new(26string)
	var reply string //它不是nil，他是有默认值的
	// 指定编码解码的方式
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	err = client.Call("HelloService.Hello", "request for something", &reply)
	if err != nil {
		panic("调用失败")
	}
	fmt.Println(reply)

}
