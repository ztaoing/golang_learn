package main

import "fmt"

func main() {
	c := 200

	c <<= 2
	fmt.Printf("第 6行  - <<= 运算符实例，c 值为 = %d\n", c)

	c >>= 2
	fmt.Printf("第 7 行 - >>= 运算符实例，c 值为 = %d\n", c)

	c &= 2
	fmt.Printf("第 8 行 - &= 运算符实例，c 值为 = %d\n", c)

	c ^= 2
	fmt.Printf("第 9 行 - ^= 运算符实例，c 值为 = %d\n", c)

	c |= 2
	fmt.Printf("第 10 行 - |= 运算符实例，c 值为 = %d\n", c)

	/*
		第 6行  - <<= 运算符实例，c 值为 = 800
		第 7 行 - >>= 运算符实例，c 值为 = 200
		第 8 行 - &= 运算符实例，c 值为 = 0
		第 9 行 - ^= 运算符实例，c 值为 = 2
		第 10 行 - |= 运算符实例，c 值为 = 2
	*/

}
