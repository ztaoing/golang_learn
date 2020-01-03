package main

import "fmt"

/**
使用type定义函数的类型
可以使用%T格式化参数打印函数的类型
*/
func add(a, b int) int {
	return a + b
}

func sub(a, b int) int {
	return a - b
}

//定义一个函数类型，输入的是两个int类型，返回值是一个int类型
type Op func(int, int) int

//定义一个函数类型，第一个参数是函数类型Op，函数类型变量可以直接用来进行函数调用
func do(f Op, a, b int) int {
	return f(a, b)
}

func main() {
	//函数名add可以当做相同函数类型形参，不需要强制类型转换
	a := do(add, 1, 2)
	fmt.Println(a)              //结果为：3
	fmt.Printf("函数类型为：%T\n", a) //函数类型为：int

	s := do(sub, 2, 1)
	fmt.Println(s)              //结果为：1
	fmt.Printf("函数类型为：%T\n", s) //函数类型为：int

}
