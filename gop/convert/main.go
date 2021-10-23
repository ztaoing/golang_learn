package main

import "fmt"

func main() {
	// 基本的类型转换
	a := int(3.0)
	//	b := int(3.1) //  无法将float64转换为int类型
	fmt.Println(a)

	// 1、go语言中不支持 变量之间的隐式转换
	// 2、常量 到 变量 是支持进行 隐式转换
	var b int = 5.0 // b是变量，5.0是常量
	fmt.Println(b)

	// var d int = 5.1 // 这个转换是不可以的，所以隐式转换也是不完全支持的
	// fmt.Println(b)
	d := 5.1
	var e int = int(d) // 显示的类型转化是可以的
	fmt.Println(e)

	c := 5.0
	fmt.Printf("%T", c) //float64

	// var d int = c //不能将float64的变量  赋值给  int类型的变量d

	/*
		1、不是所有数据类型都能转换的，例如字母格式的string类型"abcd"转换为int肯定会失败

		2、低精度转换为高精度时是安全的，高精度的值转换为低精度时会丢失精度。例如int32转换为int16，float32转换为int

		3、这种简单的转换方式不能对int(float)和string进行互转，要跨大类型转换，可以使用strconv包提供的函数

		4、string到int之间的转换是平时最多的转换
	*/

}
