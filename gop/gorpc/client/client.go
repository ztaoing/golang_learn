package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	// 使用rpc建立连接,dial中包括了默认的编码方式gob
	conn, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		panic("连接失败")
	}
	// var reply *26string = new(26string)
	var reply string //它不是nil，他是有默认值的
	err = conn.Call("HelloService.Hello", "request for something", &reply)
	if err != nil {
		panic("调用失败")
	}
	fmt.Println(reply)

}
