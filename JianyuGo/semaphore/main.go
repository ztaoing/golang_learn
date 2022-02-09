package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/semaphore"
)

/**
这三个goroutine占用的资源都等于最大资源数，也就是说只能有一个goruotine运行成功，
另外两个goroutine都会被阻塞，因为goroutine是抢占式调度，所以我们不能确定哪个gouroutine会第一个被执行，


这里我们假设第一个获取到信号量的是gouroutine:2，阻塞等待的调用者列表顺序是：goroutine:1 -> goroutine:0，
因为在goroutine:2中有一个3s的延时，所以会触发ctx的超时，ctx会下发Done信号，

因为goroutine:2和goroutine:1都是被ctx控制的，所以就会把goroutine:1从等待者队列中取消，
但是因为goroutine:1属于队列的第一个队员，并且因为goroutine:2已经释放资源，
那么就会唤醒goroutine:0继续执行，画个图表示一下：

使用这种方式可以避免goroutine永久失眠。
*/
func main() {
	s := semaphore.NewWeighted(3)
	// goroutine:0 使用ct对象来做控制，超时时间为2s
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// goroutine:1和goroutine:2对象使用ctx对象来做控制，超时时间为3s

	for i := 0; i < 3; i++ {
		if i != 0 {
			// goroutine:2对象使用ctx对象来做控制，超时时间为3s
			go func(num int) {
				if err := s.Acquire(ctx, 3); err != nil {
					fmt.Printf("goroutine： %d, err is %s\n", num, err.Error())
					return
				}
				time.Sleep(2 * time.Second)
				fmt.Printf("goroutine： %d run over\n", num)
				s.Release(3)

			}(i)
		} else {
			// goroutine:1对象使用ctx对象来做控制，超时时间为3s
			go func(num int) {
				ct, cancel := context.WithTimeout(context.Background(), time.Second*2)
				defer cancel()
				if err := s.Acquire(ct, 3); err != nil {
					fmt.Printf("goroutine： %d, err is %s\n", num, err.Error())
					return
				}
				time.Sleep(3 * time.Second)
				fmt.Printf("goroutine： %d run over\n", num)
				s.Release(3)
			}(i)
		}

	}
	time.Sleep(10 * time.Second)
}
