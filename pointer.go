package main

import "fmt"

func main() {
	//顶一个一个指针，指向nil，即指针的地址为0
	var a *int

	//空指针a的地址
	aP := &a

	fmt.Printf("a-->nil:%x\n", a)
	fmt.Printf("aP-->a:%x\n", aP)
	fmt.Printf("aP-->a-->nil(指针aP指向的指针a的内存地址):%x\n", *aP)
	fmt.Printf("&aP-->aP(表示aP在内存中地址):%x\n", &aP)
}

/**
输出结果
a-->nil(0):0
aP-->a:c00000e028
aP-->a-->nil(指针aP指向的指针a的内存地址):0
&aP-->aP(表示aP在内存中地址):c00000e030

*/
