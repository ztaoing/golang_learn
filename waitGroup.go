package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	//每次只能有一个goroutine运行
	runtime.GOMAXPROCS(1)

	wg := sync.WaitGroup{}
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println("i:", i)
			wg.Done()
		}()
	}

	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("i1:", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

/**
输入结果:
i1: 9
i: 10
i: 10
i: 10
i: 10
i: 10
i: 10
i: 10
i: 10
i: 10
i: 10
i1: 0
i1: 1
i1: 2
i1: 3
i1: 4
i1: 5
i1: 6
i1: 7
i1: 8
*/
