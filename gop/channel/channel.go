package main

import (
	"fmt"
	"time"
)

func workerC(id int, c chan int) {
	go func() {
		/*for {
			// 如果channel已经close了，从channel中读取到的值就是零值，这里chan int，所以这里收到的是0
			n, ok := <-c
			//如果已经没有值了,就结束
			if !ok {
				break
			}
			fmt.Printf("Worker %d received :%c\n", id, n)
		}*/
		// 和ok的方法是一样的
		// 如果没有close channel，range就会一直等待，直到main goroutine结束
		// range会自动检测close
		for n := range c {
			fmt.Printf("Worker %d received :%c\n", id, n)
		}
	}()
}

// 告诉使用者这是一个发数据（<-）用的channel
func CreateworkerC(id int) chan<- int {
	// 先创建一个channel，然后创建一个goroutine（对这个channel要做的事情），然后返回这个channel
	c := make(chan int)
	// 真正做事情的是在一个goroutine中，所以CreateworkerC()的操作非常快
	go workerC(id, c)
	return c
}
func bufferC() {
	// 非缓冲的channel，在发送完数据的时候必须有人接收，这样就需要切换goroutine，切换就会耗费资源
	// 可以使用缓冲型的的channel
	c := make(chan int, 3)
	go workerC(0, c)
	// 向channel中发送3个数,这时候channel没有满，即没有阻塞，所以不会发生goroutine的切换
	// 此时没有人从channel接收，也不会发生deadlock
	c <- 'a'
	c <- 'b'
	c <- 'c'

	// 多发一个，此时goroutine会发生阻塞,触发deadlock
	// fatal error: all goroutines are asleep - deadlock!
	c <- 'd'

}
func chanDemo() {
	var channels [10]chan<- int
	for i := 0; i < 10; i++ {
		// 把创建的channel保存起来
		channels[i] = CreateworkerC(i)
	}
	//发送
	for i := 0; i < 10; i++ {
		channels[i] <- 'a' + i
	}
	for i := 0; i < 10; i++ {
		channels[i] <- 'A' + i
	}
}
func closeC() {
	c := make(chan int)
	go workerC(0, c)
	c <- 'a'
	c <- 'b'
	c <- 'c'
	// 告诉接收方，发送完了
	close(c)
}
func main() {
	chanDemo()
	// bufferC()
	//closeC()

	time.Sleep(10 * time.Second)
}
