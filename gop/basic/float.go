package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println(math.MaxFloat32)
	fmt.Println(math.MaxFloat64)
	fmt.Println(math.MaxInt32)
	/*
		3.4028234663852886e+38
		1.7976931348623157e+308
	*/
	// 为什么64位的float最大值远大于int64,主要是因为float底层和int的存储是不一样的
	// %T 类型
	weight := 72.1
	fmt.Printf("%T\n", weight) //默认选择的是float64
	age := 18
	fmt.Printf("%T\n", age) //默认使用的是int

}
