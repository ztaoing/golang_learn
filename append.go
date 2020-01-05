package main

import "fmt"

func main() {
	x := make([]int, 4, 4)
	x1 := append(x, 1)

	fmt.Printf("x1长度为：%d\n", len(x1))
	for i, v := range x1 {
		fmt.Printf("x1[%d]=[%d]\n", i, v)
	}
	/**
	x1长度为：5
	输出结果：
	x1[0]=[0]
	x1[1]=[0]
	x1[2]=[0]
	x1[3]=[0]
	x1[4]=[1]
	*/

	fmt.Printf("x长度为：%d\n", len(x))
	for i, v := range x {
		fmt.Printf("x[%d]=[%d]\n", i, v)
	}
	/**
	x长度为：4
	输出结果：
	  x[0]=[0]
	  x[1]=[0]
	  x[2]=[0]
	  x[3]=[0]

	*/

	//从结果来看，向x中append后，并没有使原来的x容量增加，而是生成了一个新的内存空间保存增容后的数据即x1
}
