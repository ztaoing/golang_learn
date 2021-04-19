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
	"sync"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2/internal"
)

// PoolWithFunc 接收来自client的任务
type PoolWithFunc struct {
	// capacity of the pool.
	capacity int32

	// running is the number of the currently running goroutines.
	running int32

	// workers 存储了可用的worker
	workers []*goWorkerWithFunc

	// state is used to notice the pool to closed itself.
	state int32

	// lock for synchronous operation.
	lock sync.Locker

	// cond for waiting to get a idle worker.
	cond *sync.Cond

	// poolFunc 是用来处理任务的方法
	poolFunc func(interface{})

	// workerCache speeds up the obtainment of the an usable worker in function:retrieveWorker.
	workerCache sync.Pool

	// blockingNum is the number of the goroutines already been blocked on pool.Submit, protected by pool.lock
	blockingNum int

	options *Options
}

// purgePeriodically 定期清除过期的workers，它会单独运行一个goroutine作为清理者
func (p *PoolWithFunc) purgePeriodically() {
	heartbeat := time.NewTicker(p.options.ExpiryDuration)
	defer heartbeat.Stop()

	var expiredWorkers []*goWorkerWithFunc
	for range heartbeat.C {
		if p.IsClosed() {
			break
		}
		currentTime := time.Now()
		p.lock.Lock()
		idleWorkers := p.workers
		n := len(idleWorkers)
		var i int
		for i = 0; i < n && currentTime.Sub(idleWorkers[i].recycleTime) > p.options.ExpiryDuration; i++ {
		}
		expiredWorkers = append(expiredWorkers[:0], idleWorkers[:i]...)
		if i > 0 {
			m := copy(idleWorkers, idleWorkers[i:])
			for i = m; i < n; i++ {
				idleWorkers[i] = nil
			}
			p.workers = idleWorkers[:m]
		}
		p.lock.Unlock()

		// Notify obsolete workers to stop.
		// This notification must be outside the p.lock, since w.task
		// may be blocking and may consume a lot of time if many workers
		// are located on non-local CPUs.
		for i, w := range expiredWorkers {
			w.args <- nil
			expiredWorkers[i] = nil
		}

		// There might be a situation that all workers have been cleaned up(no any worker is running)
		// while some invokers still get stuck in "p.cond.Wait()",
		// then it ought to wakes all those invokers.
		if p.Running() == 0 {
			p.cond.Broadcast()
		}
	}
}

// NewPoolWithFunc generates an instance of ants pool with a specific function.
func NewPoolWithFunc(size int, pf func(interface{}), options ...Option) (*PoolWithFunc, error) {
	if size <= 0 {
		return nil, ErrInvalidPoolSize
	}

	if pf == nil {
		return nil, ErrLackPoolFunc
	}

	opts := loadOptions(options...)

	if expiry := opts.ExpiryDuration; expiry < 0 {
		return nil, ErrInvalidPoolExpiry
	} else if expiry == 0 {
		opts.ExpiryDuration = DefaultCleanIntervalTime
	}

	if opts.Logger == nil {
		opts.Logger = defaultLogger
	}

	p := &PoolWithFunc{
		capacity: int32(size),
		poolFunc: pf,
		lock:     internal.NewSpinLock(),
		options:  opts,
	}
	p.workerCache.New = func() interface{} {
		return &goWorkerWithFunc{
			pool: p,
			args: make(chan interface{}, workerChanCap),
		}
	}
	if p.options.PreAlloc {
		p.workers = make([]*goWorkerWithFunc, 0, size)
	}
	p.cond = sync.NewCond(p.lock)

	// 使用一个goroutine来清理过期的workers
	go p.purgePeriodically()

	return p, nil
}

//---------------------------------------------------------------------------

// Invoke 提交一个任务到pool中
func (p *PoolWithFunc) Invoke(args interface{}) error {
	if p.IsClosed() {
		return ErrPoolClosed
	}
	var w *goWorkerWithFunc
	if w = p.retrieveWorker(); w == nil {
		return ErrPoolOverload
	}
	w.args <- args
	return nil
}

// Running returns the number of the currently running goroutines.
func (p *PoolWithFunc) Running() int {
	return int(atomic.LoadInt32(&p.running))
}

// Free returns a available goroutines to work.
func (p *PoolWithFunc) Free() int {
	return p.Cap() - p.Running()
}

// Cap returns the capacity of this pool.
func (p *PoolWithFunc) Cap() int {
	return int(atomic.LoadInt32(&p.capacity))
}

// Tune changes the capacity of this pool.
func (p *PoolWithFunc) Tune(size int) {
	if size <= 0 || size == p.Cap() || p.options.PreAlloc {
		return
	}
	atomic.StoreInt32(&p.capacity, int32(size))
}

// IsClosed indicates whether the pool is closed.
func (p *PoolWithFunc) IsClosed() bool {
	return atomic.LoadInt32(&p.state) == CLOSED
}

// Release 关闭pool
func (p *PoolWithFunc) Release() {
	atomic.StoreInt32(&p.state, CLOSED)
	p.lock.Lock()
	idleWorkers := p.workers
	for _, w := range idleWorkers {
		w.args <- nil
	}
	p.workers = nil
	p.lock.Unlock()
	//  这里可能有一些调用者等待在retrieveWorker()，所以我们需要唤醒他，以防这些调用者永久的阻塞
	p.cond.Broadcast()
}

// Reboot 重启一个已经释放的pool
func (p *PoolWithFunc) Reboot() {
	if atomic.CompareAndSwapInt32(&p.state, CLOSED, OPENED) {
		go p.purgePeriodically()
	}
}

//---------------------------------------------------------------------------

// incRunning 递增当前运行的goroutine的数量
func (p *PoolWithFunc) incRunning() {
	atomic.AddInt32(&p.running, 1)
}

// decRunning 递减当前运行的goroutine的数量
func (p *PoolWithFunc) decRunning() {
	atomic.AddInt32(&p.running, -1)
}

// retrieveWorker 返回一个可用的worker来运行任务
func (p *PoolWithFunc) retrieveWorker() (w *goWorkerWithFunc) {
	spawnWorker := func() {
		w = p.workerCache.Get().(*goWorkerWithFunc)
		w.run()
	}

	p.lock.Lock()
	//可使用的worker
	idleWorkers := p.workers
	n := len(idleWorkers) - 1
	if n >= 0 {
		w = idleWorkers[n]
		idleWorkers[n] = nil
		p.workers = idleWorkers[:n]
		p.lock.Unlock()
	} else if p.Running() < p.Cap() {
		//当前运行的goroutine的数量少于容量
		p.lock.Unlock()
		spawnWorker()
	} else {
		if p.options.Nonblocking {
			p.lock.Unlock()
			return
		}
	Reentry:
		if p.options.MaxBlockingTasks != 0 && p.blockingNum >= p.options.MaxBlockingTasks {
			// MaxBlockingTasks已经设置并且不等于0 && 阻塞的个数 大于等于 允许的最大的阻塞数，就直接返回
			p.lock.Unlock()
			return
		}
		// 加入等待队列
		p.blockingNum++
		p.cond.Wait()
		p.blockingNum--
		// 当前没有goroutine在运行
		if p.Running() == 0 {
			p.lock.Unlock()
			// pool没有关闭
			if !p.IsClosed() {
				spawnWorker()
			}
			return
		}
		// 可使用的worker的数量
		l := len(p.workers) - 1
		// 没有可用的worker时
		if l < 0 {
			// 可以直接取用goroutine，因为容量没有被用完
			if p.Running() < p.Cap() {
				p.lock.Unlock()
				spawnWorker()
				return
			}
			// 再去尝试
			goto Reentry
		}
		// 有可用的worker时：
		w = p.workers[l]
		p.workers[l] = nil
		p.workers = p.workers[:l]
		p.lock.Unlock()
	}
	return
}

// revertWorker puts a worker back into free pool, recycling the goroutines.
func (p *PoolWithFunc) revertWorker(worker *goWorkerWithFunc) bool {
	if capacity := p.Cap(); (capacity > 0 && p.Running() > capacity) || p.IsClosed() {
		// 运行的goroutine数超过了容量 或者 pool已经关闭了
		return false
	}
	worker.recycleTime = time.Now()
	p.lock.Lock()

	// 避免内存的泄漏，在锁范围内增加了一个双检测
	// Issue: https://github.com/panjf2000/ants/issues/113
	// 如果已经关闭
	if p.IsClosed() {
		p.lock.Unlock()
		return false
	}
	// 将worker归还
	p.workers = append(p.workers, worker)

	// 提醒： 卡在了'retrieveWorker()' 的调用者，现在有一个可用的worker了
	p.cond.Signal()

	p.lock.Unlock()
	return true
}
