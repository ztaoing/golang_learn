package main

import "fmt"

/**
匿名函数
*/
//匿名函数被直接赋值给函数变量
var sum = func(a, b int) int {
	return a + b
}

//将匿名函数作为参数
func doInput(f func(int, int) int, a, b int) int {
	return f(a, b)
}

//匿名函数作为返回值
func wrap(op string) func(int, int) int {
	switch op {
	case "add":
		return func(a int, b int) int {
			return a + b
		}
	case "sub":
		return func(a, b int) int {
			return a + b
		}
	default:
		return nil

	}
}
func main() {
	//直接调用匿名函数
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	//匿名函数作为实参
	doInput(func(x int, y int) int {
		return x + y
	}, 1, 2)

	opFunc := wrap("add")
	re := opFunc(2, 3)

	fmt.Printf("%d\n", re)

}
