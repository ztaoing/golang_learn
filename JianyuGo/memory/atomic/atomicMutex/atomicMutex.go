package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

/**
可以在本地运行以上这两段代码，可以观察到计数器的结果都最后都是1000000，都是线程安全的。

需要注意的是，所有原子操作方法的"被操作数形参必须是指针类型"，通过指针变量可以获取被操作数在内存中的地址，
从而施加特殊的CPU指令，确保同一时间只有一个goroutine能够进行操作。
*/
func mutexAdd() {
	var a int32 = 0
	var wg sync.WaitGroup
	var mu sync.Mutex
	start := time.Now()
	for i := 0; i < 100000000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			a += 1
			mu.Unlock()
		}()
	}
	wg.Wait()
	timeSpends := time.Now().Sub(start).Nanoseconds()
	fmt.Printf("use mutex algorithm is %d, spend time: %v\n", a, timeSpends)
}

func AtomicAdd() {
	var a int32 = 0
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//被操作数必须是指针类型
			atomic.AddInt32(&a, 1)
		}()
	}
	wg.Wait()
	timeSpends := time.Now().Sub(start).Nanoseconds()
	fmt.Printf("use atomic algorithm is %d, spend time: %v\n", atomic.LoadInt32(&a), timeSpends)
}
