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
			time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
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
		time.Sleep(1 * time.Microsecond)
		fmt.Printf("Worker %d receied %d\n", id, n)
	}
}
func main() {
	var c1, c2 = generator(), generator()
	w := createWorker(0)
	tm := time.After(10 * time.Second)
	// 每两秒显示队列的长度
	tick := time.Tick(2 * time.Second)

	// 非阻塞式的接收
	var values []int
	for {
		var activeWorker chan<- int
		var activeValue int
		if len(values) > 0 {
			// 当有值的时候，把值发送给已创建的channel
			// 当没有值的时候，activeWorker这个channel就是nil的
			// nil channel在select中是肯定可以正确运行的，但是肯定不会被select到，永远是阻塞的
			activeWorker = w
			activeValue = values[0]
		}

		// 如果生成数据的速度很快，c1和c2都会将收到的值发送给n，这样前边收到的n，就会被后边的n覆盖掉。
		// 所以需要把收到的n，保存起来，排队。
		select {
		case n := <-c1:
			values = append(values, n)

		case n := <-c2:
			values = append(values, n)
		case activeWorker <- activeValue:
			// 把已经发送的值，从队列中拿掉
			values = values[1:]
		case <-time.After(800 * time.Millisecond):
			//如果两次select时间差，超过800毫秒没有收到数据，就显示超时
			fmt.Println("time out")
		case <-tick:
			// 定时显示队列的长度
			fmt.Println("the length of queue %d", len(values))
		case <-tm:
			// 运行10秒之后，退出
			fmt.Println("bye ")
			return

		}
	}

}
