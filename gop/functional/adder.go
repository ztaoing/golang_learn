package main

import "fmt"

func add() func(int) int {
	// sum是自由变量,会连接所有的层级的sum，就构成了闭包
	sum := 0
	// return返回的不是function的代码，返回的是这个函数，和对sum的引用
	return func(i int) int {
		//i是局部变量
		sum += i
		return sum
	}
}
func main() {
	//a这个函数体中存储了sum，每次叠加后的结果都会被保存，即闭包
	//
	a := add()

	for i := 0; i < 10; i++ {
		fmt.Println(a(i))
	}
}
