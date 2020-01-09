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

输出结果：
[5 7 9] [5 7 9 13] [5 7 9 13]

*/
