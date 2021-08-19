package main

import (
	"fmt"
	"golang_learn/golang_learn/gop/infra"
)

func getRetriever() retriever {
	return infra.Retrieve{}
}

type retriever interface {
	Get(string) string
}

func main() {
	url := "https://www.imooc.com"
	// retriever 是一个动态的retriever（可能用到的是testing.Retrieve或者infra.Retrieve），不是固定的
	// 所以我们需要一个接口,来解耦具体的某个retriever
	var r retriever = getRetriever()
	// retriever := getRetrieve() //这样看不出retrieve的类型
	all := r.Get(url)
	fmt.Println(all)
}
