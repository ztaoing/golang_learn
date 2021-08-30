package main

import (
	"fmt"
	"sync"
	"time"
)

type atomicInt struct {
	value int
	lock  sync.Mutex
}

func (a *atomicInt) increment() {

	// 如果想利用defer的特性，在代码块区内对元素进行保护，在代码区结束时释放锁，可以使用匿名函数
	func() {
		a.lock.Lock()
		defer a.lock.Unlock()
		a.value++
	}()

}
func (a *atomicInt) get() int {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.value
}
func main() {
	var a atomicInt
	a.increment()
	go func() {
		a.increment()
	}()
	time.Sleep(1 * time.Second)
	fmt.Println(a.get())
}
