package main

import "fmt"

/**
向被关闭的通道发送数据将会触发panic
从已关闭的通道接收数据时不会发生阻塞
*/
func main() {
	//创建 具有缓冲的通道
	ch := make(chan int, 2)

	//向通道发送数据
	ch <- 22
	ch <- 88

	//关闭缓冲通道
	//如果不关闭ch，会造成死锁:fatal error: all goroutines are asleep - deadlock!
	close(ch)

	//遍历缓冲中的数据
	for i := 0; i < cap(ch)+1; i++ {
		//从通道中读出
		value, ok := <-ch
		fmt.Println(value, ok)
	}
}

/**结果
22 true
88 true
0 false  close channel后，此时通道中已经没有数据，所以获取失败，并且没有发生阻塞

如果没有 close channel ，从channel中取出两个值后（所有值），从channel中读取，此时主goroutine会被阻塞，并且形成死锁
*/
