// MIT License

// Copyright (c) 2018 Andy Pan

// Permission is hereby granted, free of charge, to any person obtaining a copy
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
	"runtime"
	"time"
)

// goWorker 是实际的执行任务的人，它使用一个goroutine接收任务，然后使用指定的方法处理这个任务
type goWorker struct {
	// 拥有当前worker的指针
	pool *Pool

	// 需要被执行的任务
	task chan func()

	// 回收时间  will be update when putting a worker back into queue.
	recycleTime time.Time
}

// run starts a goroutine to repeat the process
// that performs the function calls.
// run 开启了一个goroutine执行指定的方法来处理任务
func (w *goWorker) run() {
	// 增加运行的goroutine数量
	w.pool.incRunning()
	go func() {
		// 在任务处理完成后，
		defer func() {
			w.pool.decRunning()
			// 将worker归还到workerCache中
			w.pool.workerCache.Put(w)
			//处理异常
			if p := recover(); p != nil {
				// 使用定制的PanicHandler
				if ph := w.pool.options.PanicHandler; ph != nil {
					ph(p)
				} else {
					w.pool.options.Logger.Printf("worker exits from a panic: %v\n", p)
					var buf [4096]byte
					// 获取此时的运行栈
					n := runtime.Stack(buf[:], false)
					w.pool.options.Logger.Printf("worker exits from panic: %s\n", string(buf[:n]))
				}
			}
			// 没有发生panic：
			// 调用 Signal()通知那些等待获取可用goroutine的被阻塞的调用者
			// here in case there are goroutines waiting for available workers.
			w.pool.cond.Signal()
		}()

		for f := range w.task {
			if f == nil {
				return
			}
			// 执行每一个任务
			f()
			// 执行完，将worker归还到pool中
			if ok := w.pool.revertWorker(w); !ok {
				return
			}
		}
	}()
}
