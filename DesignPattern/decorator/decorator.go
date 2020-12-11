package main

import "fmt"

//装饰器模式

type Decoer func(i int, s string) bool

func foo(i int, s string) bool {
	fmt.Printf("=== foo ===\n")
	return true
}

func withTx(fn Decoer) Decoer {
	return func(i int, s string) bool {
		fmt.Printf("=== start tx ===\n")
		result := fn(i, s)
		fmt.Printf("=== commit tx ===\n")
		return result
	}
}
func main() {
	foo := withTx(foo)
	foo(1, "hello")
}
