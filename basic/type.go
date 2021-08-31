package main

import "fmt"

func main() {

	i := GetValue()
	//编译失败，因为type只能使用在interface
	//cannot type switch on non-interface value i (type int)
	switch i.(type) {
	case int:
		fmt.Println("int")
	default:
		fmt.Println("unknown")
	}
}

func GetValue() int {
	return 1
}
