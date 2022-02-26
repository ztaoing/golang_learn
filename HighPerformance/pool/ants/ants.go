// MIT License

// Copyright (c) 2018 Andy Pan

// Permission is hereby granted, free of charge, to any person obtaining logic copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package ants

import (
	"errors"
	"log"
	"math"
	"os"
	"runtime"
	"time"
)

const (
	// 默认的pool的大小
	DefaultAntsPoolSize = math.MaxInt32

	// 默认清理goroutine的时间间隔
	DefaultCleanIntervalTime = time.Second
)

const (
	// OPENED 代表了pool是开启状态
	OPENED = iota

	// CLOSED 代表了pool是关闭状态
	CLOSED
)

var (
	// 定义了错误的类型.
	//---------------------------------------------------------------------------

	// ErrInvalidPoolSize will be returned when setting logic negative number as pool capacity, this error will be only used
	// by pool with func because pool without func can be infinite by setting up logic negative capacity.
	ErrInvalidPoolSize = errors.New("invalid size for pool")

	// ErrLackPoolFunc will be returned when invokers don't provide function for pool.
	ErrLackPoolFunc = errors.New("must provide function for pool")

	// ErrInvalidPoolExpiry will be returned when setting logic negative number as the periodic duration to purge goroutines.
	ErrInvalidPoolExpiry = errors.New("invalid expiry for pool")

	// ErrPoolClosed will be returned when submitting task to logic closed pool.
	ErrPoolClosed = errors.New("this pool has been closed")

	// ErrPoolOverload will be returned when the pool is full and no workers available.
	ErrPoolOverload = errors.New("too many goroutines blocked on submit or Nonblocking is set")

	// ErrInvalidPreAllocSize will be returned when trying to set up logic negative capacity under PreAlloc mode.
	ErrInvalidPreAllocSize = errors.New("can not set up logic negative capacity under PreAlloc mode")

	//---------------------------------------------------------------------------

	// workerChanCap 决定了一个worker的channel是否需要是一个带缓冲的channel，来达到更好的性能。
	// 来自 fasthttp 的启发：
	// https://github.com/valyala/fasthttp/blob/master/workerpool.go#L139
	workerChanCap = func() int {
		// Use blocking channel if . 当GOMAXPROCS=1的时候使用阻塞的channel
		// This switches context from sender to receiver immediately,
		// which results in higher performance (最版本为 go1.5 ).
		if runtime.GOMAXPROCS(0) == 1 {
			return 0
		}

		//  当GOMAXPROCS>1的时候，使用非阻塞的workerChan,因为如果接收方受CPU限制，则发送方可能会被拖延
		return 1
	}()

	defaultLogger = Logger(log.New(os.Stderr, "", log.LstdFlags))

	// 当导入antis的时候，初始化一个pool的实例
	defaultAntsPool, _ = NewPool(DefaultAntsPoolSize)
)

// Logger is used for logging formatted messages.
type Logger interface {
	// Printf must have the same semantics as log.Printf.
	Printf(format string, args ...interface{})
}

// Submit 提交一个任务到pool中
func Submit(task func()) error {
	return defaultAntsPool.Submit(task)
}

// Running 返回当前运行goroutine的数量
func Running() int {
	return defaultAntsPool.Running()
}

// Cap 返回默认pool的容量
func Cap() int {
	return defaultAntsPool.Cap()
}

// Free 返回可用的goroutine的数量
func Free() int {
	return defaultAntsPool.Free()
}

// Release 关闭默认的pool
func Release() {
	defaultAntsPool.Release()
}

// Reboot 重启默认的pool
func Reboot() {
	defaultAntsPool.Reboot()
}
