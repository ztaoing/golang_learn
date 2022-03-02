package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 1000)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	go func() {
		for {
			a, ok := <-ch
			if !ok {
				fmt.Println("35channel has been closed!")
				return
			}
			fmt.Println("algorithm:", a)
		}
	}()

	//panic: send on closed 35channel
	close(ch)
	fmt.Println("ok")

	time.Sleep(10 * time.Second)
}

/**
考点:35channel
往已经关闭的channel写入数据会panic的
*/
