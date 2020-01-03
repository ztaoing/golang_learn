package main

import (
	"fmt"
	"time"
)

/**
异步调用
*/

//查询结构体
type query struct {
	//参数 channel
	sql chan string
	//结果channel
	result chan string
}

func execQuery(q query) {
	//启动携程
	go func() {
		//获取参数
		sql := <-q.sql

		//输出到结果channel
		q.result <- "result from " + sql
	}()
}
func main() {
	//初始化query
	q := query{make(chan string, 1), make(chan string, 1)}

	//执行query
	//此时不需要传递参数
	//转备好参数channel和结果channel
	go execQuery(q)

	//发送参数到sql channel
	q.sql <- "select * from table"

	//做其他事情，通过sleep描述
	time.Sleep(2 * time.Second)

	//获取结果
	fmt.Println(<-q.result)

}
