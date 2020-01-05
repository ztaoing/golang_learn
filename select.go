package main

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(1)
	int_chan := make(chan int, 1)
	string_chan := make(chan string, 1)
	int_chan <- 1
	string_chan <- "hello"

	select {
	case value := <-int_chan:
		fmt.Printf("int_chan:[%d]", value)
	case value := <-string_chan:
		panic(value)
	}
}

/**
输出结果：
int_chan:[1]
第一个case已经满足，就不会再执行第二个
*/
