package main

import (
	"fmt"
	"golang_learn/golang_learn/gop/queue"
)

func main() {
	q := queue.Queue{8}
	q.Push(1)
	q.Push(2)
	q.Push(3)

	fmt.Println(q.Pop())
	fmt.Println(q.Pop())
	fmt.Println(q.IsEmpty())

	fmt.Println(q.Pop())
	fmt.Println(q.Pop())
	fmt.Println(q.IsEmpty())
}
