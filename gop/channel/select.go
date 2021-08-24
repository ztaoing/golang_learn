package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generator() chan int {
	out := make(chan int)
	//负责不停的向channel中发送数据
	go func() {
		i := 0
		for {
			time.Sleep(time.Duration(rand.Intn(15)) * time.Second)
			out <- i
			i++
		}

	}()
	return out
}

func createWorker(id int) chan<- int {
	c := make(chan int)
	go worker(id, c)
	return c
}
func worker(id int, c chan int) {
	for n := range c {
		time.Sleep(5 * time.Second)
		fmt.Printf("Worker %d receied %d\n", id, n)
	}
}
func main() {
	var c1, c2 = generator(), generator()
	w := createWorker(0)

	// nil channel在select中是肯定可以正确运行的，但是肯定不会被select到，永远是阻塞的
	// 非阻塞式的接收
	hasValue := true
	n := 0
	for {
		var activeWorker chan<- int
		if hasValue {
			activeWorker = w
		}

		select {
		case n = <-c1:
			hasValue = true
		case n = <-c2:
			hasValue = true
		case activeWorker <- n:
			hasValue = false

		}
	}

}
