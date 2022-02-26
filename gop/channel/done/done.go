package main

import (
	"fmt"
	"sync"
)

func doWork(id int, w worker) {
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
		for n := range w.in {
			fmt.Printf("Worker %d received :%c\n", id, n)
			// 通过channel来通信
			go func() {
				w.done()
			}()
		}
	}()
}

type worker struct {
	in chan int
	//wg *sync.WaitGroup //他是一个引用，我们必须使用外边的传进来的wg
	done func() // 不直接使用wg，而使用函数式编程的方式, 把wg抽象掉了，done究竟做了什么事情呢？它有CreateworkerC()的时候来配置.这样抽象程度更高!
}

// 告诉使用者这是一个发数据（<-）用的channel
func CreateworkerC(id int, wg *sync.WaitGroup) worker {

	w := worker{
		in: make(chan int),
		done: func() {
			wg.Done()
		},
	}
	// 真正做事情的是在一个goroutine中，所以CreateworkerC()的操作非常快
	go doWork(id, w)
	return w
}

func chanDemo() {
	var workers [10]worker
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		// 把创建的channel保存起来
		workers[i] = CreateworkerC(i, &wg)

	}
	// 有20个任务需要等待
	wg.Add(20)
	//发送
	for i := 0; i < 10; i++ {
		//发
		workers[i].in <- 'logic' + i
		//收
		//<-workers[i].done
	}
	for i := 0; i < 10; i++ {
		workers[i].in <- 'A' + i
		//<-workers[i].done
	}
	wg.Wait()

	// 等待所有

}

func main() {
	chanDemo()

	// 通过channel的方式通知结束，所以不需要使用sleep了。
	//time.Sleep(10 * time.Second)
}
