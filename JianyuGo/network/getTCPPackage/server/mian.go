package main

import (
	"log"
	"net"
)

// 通过go语言是实现tcp的套接字编程，并结合tcpdump工具，展示他的三次握手、数据传输以及四次回收的过程，帮助读者更好的理解tcp协议与go网络编程
/**
网络编程属于io的范畴，其发展可以简单概括为：多进程-》多线程-》non-block + I/O多路复用

go的网络编程模型是 同步网络编程。它基于协程 + I/O多路复用（Linux下epoll、Darwin下kqueue、windows下iocp，通过网络轮询器netpoller进行封装）
结合网络轮询器与调度器实现。

用户层goroutine中的block socket，实际上是通过netpoller模拟出来的。
注：runtime拦截了底层socket系统调用的错误码，并通过netpoller和goroutine调度让goroutine阻塞在用户层得到的socket fd上。

go将网络编程的复杂性隐藏于runtime中：开发者不用关注socket是否是non-block的，也不用处理回调，只需要在每个连接对应的goroutine中以block I/O的方式对待socket即可
例如： 当用户层针对某个socket fd发起read操作，如果该socket fd中尚无数据，那么runtime会将该socket fd加入到netpoller中监听，同时对应的goroutine将被挂起，直到runtime收到
	socket fd数据ready的通知，runtime才会重新唤醒等待在该socket fd上准备read的那个goroutine。而这个过程从goroutine的视角来看，就像read操作一直block在哪个socket fd上似的。
	（看起来是goroutine阻塞在了read上，实际是因为没有数据而被runtime挂起，待收到socket fd数据ready的通知后，runtime才会重新唤醒等待在该socket fd上准备read的那个goroutine）

总结：go将复杂的网络模型进行封装，放在用户面前的知识阻塞式的goroutine，这让我们可以非常轻松地实现高性能网络编程。


*/
// 在go中网络编程非常容易，我们通过go的net包，就可以轻松实现一个tcp服务器。
func main() {
	// 1、创建listener,监听端口8000
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("error listener returned:%s", err)
	}
	defer l.Close()

	for {
		// 2、等待接收新的连接，accept方法将以阻塞式地等待新的连接到达，并将该连接作为net.Conn接口类型返回
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("error to accpet new connection:%s", err) // os.Exit(1)
		}
		// 3、创建一个goroutine来读取和写入
		go func() {
			log.Printf("TCP session openClose")
			defer conn.Close()

			for {
				d := make([]byte, 100)
				// 从tcp buffer中读取
				_, err := conn.Read(d)
				if err != nil {
					log.Printf("Error reading TCP session:%s", err)
					break // break the loop
				}

				// 将数据写到tcp client中
				_, err = conn.Write(d)
				if err != nil {
					log.Printf("error writing TCP session:%s", err)
					break
				}
			}
		}()
	}
}
