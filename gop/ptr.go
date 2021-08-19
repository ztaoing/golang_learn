package main

import "fmt"

func main() {
	var a int = 2
	var pa *int = &a
	*pa = 3
	fmt.Println(a)
}
