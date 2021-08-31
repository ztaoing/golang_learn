package main

import "fmt"

func main() {
	s := []int{5}
	s = append(s, 7)
	s = append(s, 9)

	x := append(s, 11)
	y := append(s, 13)

	fmt.Println(s, x, y)

}

/**
s为底层数组的切片，x = 5，7，9，11 ； y = 13  append to s（5，7，9），将底层数组的11，替换为了13 。此时s仍然是x = 5，7，9，11
输出结果：
[5 7 9] [5 7 9 13] [5 7 9 13]

*/
