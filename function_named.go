package main

import "fmt"

/**
有名函数，可以直接使用函数名调用函数，也可以直接赋值给函数类型变量，后续通过该变量来调用该函数
*/
func sum(a, b int) int {
	fmt.Println(a + b)
	return a + b
}
func main() {
	//直接调用
	sum(3, 4) //7
	//有名函数可以直接赋值给变量
	f := sum
	f(1, 2) //3
}
