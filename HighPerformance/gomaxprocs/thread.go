/**
* @Author:zhoutao
* @Date:2021/1/2 下午2:13
* @Desc:go运行程序中的线程数
 */

package main

import (
	"fmt"
	"net"
	"runtime"
	"runtime/pprof"
	"sync"
)

// 实际的thread可能不受GOMAXPROCS限制，go代码进行系统调用的时候被block的，线程数不受这个变量限制
// 也就是，如果并发的blocking的系统调用很多的话，go就会创建大量的线程，但是当系统调用完成后，这些线程因为go运行时的设计
// 却不会被回收：https://github.com/golang/go/issues/14592

// 什么是blocking的系统调用：阻塞的系统调用就是系统调用执行时，在完成之前，调用者必须等待。
// read就是个很好的例子，如果没有数据可读，调用者就一直等到一些数据可读
// 那，go从网络I/O中read数据岂不是每个读取的goroutine都会占用一个系统线程吗？不会的。go使用netpoll而处理网络读写。它使用epoll（Linux）、kqueue(BSD，Darwin)的方式可以轮询network I/O的状态
// 一旦接受了一个链接，连接的文件描述符就被设置为non-blocking,也就是一旦连接中没有数据了，从其中read数据，并不会被阻塞。而是返回一个特定的错误。因此go标准库的网络读写不会产生大量的线程（即为非阻塞）
// 单cgo或者其他一些阻塞的系统调用可能会导致线程大量增加并无法回收

//线程数暴涨:

var threadProfile = pprof.Lookup("threadcreate")

func main() {
	fmt.Printf("thread in starting:%d\n", threadProfile.Count())

	var wg sync.WaitGroup

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < 100; j++ {
				net.LookupHost("www.google.com")
			}
		}()
	}
	wg.Wait()

	fmt.Printf("thread after lookupHost:%d\n", threadProfile.Count())
}

// runtime.LockOSThread(), 会把当前goroutine绑定到当前的系统线程上，这个goroutine总是在这个线程中执行，而且也不会有其他goroutine在这个线程中执行
// 只有这个goroutine调用了相同次数的UnlockOSThread函数之后，才会进行解绑
// 如果goroutine在退出的时候没有unlock这个线程，那么这个线程会被终止。我们可以利用这个特性将绑定的线程杀掉。

func KillOne() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		runtime.LockOSThread()
		return
	}()
	wg.Wait()
}
