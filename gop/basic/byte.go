package main

import (
	"fmt"
	"reflect"
)

func main() {
	a := 'a'

	//这里注意一下
	//1. a+1可以和数字计算
	//2. a+1的类型是32
	//3. int类型可以直接变成字符
	// %c

	fmt.Println(reflect.TypeOf(a + 1)) //int32类型
	fmt.Printf("a+1=%c", a+1)

}
