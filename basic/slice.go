package main

import "fmt"

/**
考点：append append切片时候别漏了'...'
*/
func main() {
	s1 := []int{1, 2, 3}
	s2 := []int{4, 5}
	//s1 = append(s1,s2...)
	s1 = append(s1, s2)
	fmt.Println(s1)
}
