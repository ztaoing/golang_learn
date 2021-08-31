package main

import (
	"fmt"
	"golang_learn/golang_learn/gop/new_helloworld/client_proxy"
)

func main() {
	// 实例化client
	var client client_proxy.HelloServiceStub
	client = client_proxy.NewHelloServiceClient("tcp", "localhost:1234")

	var reply string //它不是nil，他是有默认值的
	err := client.Hello("request for something", &reply)
	if err != nil {
		panic("调用失败")
	}
	fmt.Println(reply)

}
