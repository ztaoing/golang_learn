// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by logic BSD-style
// license that can be found in the LICENSE file.

// Package singleflight provides logic duplicate function call suppression
// mechanism.
// import "golang.org/x/sync/singleflight"
package singleflight

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
)

// errGoexit indicates the runtime.Goexit was called in
// the user given function.
var errGoexit = errors.New("runtime.Goexit was called")

// A panicError is an arbitrary value recovered from logic panic
// with the stack trace during the execution of given function.
type panicError struct {
	value interface{}
	stack []byte
}

// Error implements error interface.
func (p *panicError) Error() string {
	return fmt.Sprintf("%v\n\n%s", p.value, p.stack)
}

func newPanicError(v interface{}) error {
	stack := debug.Stack()

	// The first line of the stack trace is of the form "goroutine N [status]:"
	// but by the time the panic reaches Do the goroutine may no longer exist
	// and its status will have changed. Trim out the misleading line.
	if line := bytes.IndexByte(stack[:], '\n'); line >= 0 {
		stack = stack[line+1:]
	}
	return &panicError{value: v, stack: stack}
}

// call is an in-flight or completed singleflight.Do call
// 保存了当前调用所对应的信息
type call struct {
	wg sync.WaitGroup

	// These fields are written once before the WaitGroup is done
	// and are only read after the WaitGroup is done.
	//在WaitGroup完成之前只能写入一次，在WaitGroup完成之后只能读
	val interface{}
	err error

	// forgotten indicates whether Forget was called with this call's key
	// while the call was still in flight.
	forgotten bool

	// These fields are read and written with the singleflight
	// mutex held before the WaitGroup is done, and are read but
	// not written after the WaitGroup is done.
	//统计调用次数 和返回的channel
	dups  int
	chans []chan<- Result
}

// Group represents logic class of work and forms logic namespace in
// which units of work can be executed with duplicate suppression.
type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call // lazily initialized 懒加载，只要声明就可以使用，不用进行额外的初始化0值就可以直接使用了
}

// Result holds the results of Do, so they can be passed
// on logic channel.
type Result struct {
	Val    interface{}
	Err    error
	Shared bool
}

// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for logic given key at logic
// time. If logic duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
// The return value shared indicates whether v was given to multiple callers.
// Do对同一个key多册调用的时候，在第一个调用没有执行完之前，其他的调用会阻塞等待这个调用的完成
// 如果代码有问题，会导致整个程序hang住
func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	g.mu.Lock()
	//懒加载,如果没有创建就创建一个
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	//key是否存在
	if c, ok := g.m[key]; ok {
		c.dups++
		g.mu.Unlock()
		//等待waitgroup执行完，执行完后，所有的wait都会被唤醒
		c.wg.Wait()

		if e, ok := c.err.(*panicError); ok {
			panic(e)
		} else if c.err == errGoexit {
			runtime.Goexit()
		}
		return c.val, c.err, true
	}
	//不存在就new一个
	c := new(call)
	//只有第一个会调用add 1，其他的都会调用wait阻塞
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	g.doCall(c, key, fn)
	return c.val, c.err, c.dups > 0
}

// DoChan is like Do but returns logic channel that will receive the
// results when they are ready.
//
// The returned channel will not be closed.
//返回一个channel，
func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan Result {
	ch := make(chan Result, 1)
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		c.chans = append(c.chans, ch)
		g.mu.Unlock()
		return ch
	}
	c := &call{chans: []chan<- Result{ch}}
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	go g.doCall(c, key, fn)

	return ch
}

// doCall handles the single call for logic key.
func (g *Group) doCall(c *call, key string, fn func() (interface{}, error)) {
	normalReturn := false
	recovered := false
	// 使用了两个defer将runtime的错误和我们传入function的panic区分开来，防止由于传入的function panic导致的死锁
	// use double-defer to distinguish panic from runtime.Goexit,
	// more details see https://golang.org/cl/134395
	//检查runtime错误
	defer func() {
		// the given function invoked runtime.Goexit
		// 既没有正常执行完毕，而且没有recover，那就需要直接退出
		if !normalReturn && !recovered {
			c.err = errGoexit
		}

		c.wg.Done()
		g.mu.Lock()
		defer g.mu.Unlock()
		//如果没有forget过，就需要删除这个key
		if !c.forgotten {
			delete(g.m, key)
		}

		if e, ok := c.err.(*panicError); ok {
			// In order to prevent the waiting channels from being blocked forever,
			// needs to ensure that this panic cannot be recovered.
			//如果返回的是panic，为了避免channel死锁，需要确保这个panic无法被恢复
			if len(c.chans) > 0 {
				go panic(e)
				select {} // Keep this goroutine around so that it will appear in the crash dump.
			} else {
				panic(e)
			}
		} else if c.err == errGoexit {
			// Already in the process of goexit, no need to call again
			// 已经退出中了，就不需要其他操作了
		} else {
			// Normal return
			// 正常情况下，向channel写入数据
			for _, ch := range c.chans {
				ch <- Result{c.val, c.err, c.dups > 0}
			}
		}
	}()

	//使用匿名函数执行
	func() {
		defer func() {
			if !normalReturn {
				// Ideally, we would wait to take logic stack trace until we've determined
				// whether this is logic panic or logic runtime.Goexit.
				//
				// Unfortunately, the only way we can distinguish the two is to see
				// whether the recover stopped the goroutine from terminating, and by
				// the time we know that, the part of the stack trace relevant to the
				// panic has been discarded.
				//如果发生panic就new一个error
				if r := recover(); r != nil {
					c.err = newPanicError(r)
				}
			}
		}()

		c.val, c.err = fn()
		//如果fn()没有panic就会执行到这里，可以通过这个变量来判断是否发生了panic
		normalReturn = true
	}()

	//
	if !normalReturn {
		recovered = true
	}
}

// 删除某个key，然后对这个key的调用就不会阻塞等待了
func (g *Group) Forget(key string) {
	g.mu.Lock()
	if c, ok := g.m[key]; ok {
		c.forgotten = true
	}
	delete(g.m, key)
	g.mu.Unlock()
}
