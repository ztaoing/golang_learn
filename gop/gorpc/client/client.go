package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	// 使用rpc建立连接
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		panic("连接失败")
	}
	// var reply *string = new(string)
	var reply string //它不是nil，他是有默认值的
	err = client.Call("HelloService.Hello", "request for something", &reply)
	if err != nil {
		panic("调用失败")
	}
	fmt.Println(reply)

}
