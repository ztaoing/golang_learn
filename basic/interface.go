package main

import "fmt"

/**
考点：interface内部结构
*/
func main() {
	var x *int = nil
	Foo(x)
}
func Foo(x interface{}) {
	if x == nil {
		fmt.Println("empty interface")
		return
	}
	fmt.Println("non-empty interface")
}

/**
输出结果：
non-empty interface
*/
