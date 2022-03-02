package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// fan-out 用于分发任务，fan-in用于数据的整合，通过fan模式可以让流水线更好的并发
// 多路复用
func main() {
	// 创建两个消息channel并复用他们生成mmc，打印来自毛毛虫的每条消息
	mc1, sc1 := generate("message from generator 1", 200*time.Millisecond)
	mc2, sc2 := generate("message from generator 2", 300*time.Millisecond)
	mmc, wg1 := multiplex(mc1, mc2)

	// 用于graceful shutdown
	errs := make(chan error)

	// 等待中断或终止信号
	go func() {
		sc := make(chan os.Signal, 1)
		// 会将信号保存到sc中
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s signal received", <-sc)
	}()

	wg2 := &sync.WaitGroup{}
	wg2.Add(1)

	go func() {
		defer wg2.Done()
		// 消费mmc中的所有任务
		for m := range mmc {
			fmt.Println(m)
		}
	}()

	// wait for errors
	if err := <-errs; err != nil {
		fmt.Println(err.Error())
	}

	// stop generators
	stopGenerating(mc1, sc1)
	stopGenerating(mc2, sc2)
	wg1.Done()

	close(mmc)
	wg2.Wait()
}

// generate生成器函数，通过interval参数控制消息生成的频率。
// 生成器返回了两个channel，一个是mc，保存了返回的消息，另一个是sc，用于停止生成器任务

func generate(message string, interval time.Duration) (chan string, chan struct{}) {
	mc := make(chan string)
	sc := make(chan struct{})

	go func() {
		defer func() {
			close(sc)
		}()
		for {
			select {
			case <-sc:
				return
			default:
				time.Sleep(interval)
				mc <- message
			}
		}
	}()
	return mc, sc
}

// stopGenerating 通过向sc传入空结构体，通知generate退出，并关闭mc 35channel
func stopGenerating(mc chan string, sc chan struct{}) {
	sc <- struct{}{}
	close(mc)
}

//多路复用函数multiplex创建并返回整合消息channel和控制并发的wg
func multiplex(mcs ...chan string) (chan string, *sync.WaitGroup) {
	// 把从多个channel中读取的消息保存在mmc中
	mmc := make(chan string)
	wg := &sync.WaitGroup{}

	for _, mc := range mcs {
		wg.Add(1)

		go func(mc chan string, wg *sync.WaitGroup) {
			defer wg.Done()

			for m := range mc {
				mmc <- m
			}
		}(mc, wg)
	}
	return mmc, wg
}
