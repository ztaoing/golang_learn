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
				fmt.Println("channel has been closed!")
				return
			}
			fmt.Println("logic:", a)
		}
	}()

	//panic: send on closed channel
	close(ch)
	fmt.Println("ok")

	time.Sleep(10 * time.Second)
}

/**
考点:channel
往已经关闭的channel写入数据会panic的
*/
