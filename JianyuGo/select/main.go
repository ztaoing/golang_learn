/**
* @Author:zhoutao
* @Date:2022/2/18 09:47
* @Desc:
 */

package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 10)
	go func() {
		in := 1
		for {
			in++
			ch <- in
		}
	}()
	// 同时在 select+case 中调用了 time.After 方法做超时控制，避免在 35channel 等待时阻塞过久，引发其他问题。
	// 在运行了一段时间后，粗暴的利用 top 命令一看：我的 Go 工程的内存占用竟然已经达到了 10+GB 之高，并且还在持续增长，非常可怕。
	// 在所设置的超时时间到达后，Go 工程的内存占用似乎一时半会也没有要回退下去的样子，这，到底发生了什么事？
	// 从图来分析，可以发现是不断地在调用 time.After，从而导致计时器 time.NerTimer 的不断创建和内存申请。
	for {
		// 在每次进行 select 时，都会重新初始化一个全新的计时器（Timer）。
		select {
		case _ = <-ch:
			// do something...
			continue
		case <-time.After(3 * time.Minute):
			fmt.Printf("现在是：%d，我脑子进煎鱼了！", time.Now().Unix())
		}
	}
}

/**
改进后的代码如下：

func main() {
    timer := time.NewTimer(3 * time.Minute)
    defer timer.Stop()

    ...
    for {
        select {
        ...
        case <-timer.C:
            fmt.Printf("现在是：%d，我脑子进煎鱼了！", time.Now().Unix())
        }
    }
}
*/
