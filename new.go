package main

import "fmt"

/**
考点：new list:=make([]int,0)
*/
func main() {
	list := new([]int)
	//first argument to append must be slice; have *[]int
	list = append(list, 1)
	fmt.Println(list)
}
