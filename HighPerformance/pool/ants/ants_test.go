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
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	_ "net/http/pprof"

	"github.com/stretchr/testify/assert"
)

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
	// GiB // 1073741824
	// TiB // 1099511627776             (超过了int32的范围)
	// PiB // 1125899906842624
	// EiB // 1152921504606846976
	// ZiB // 1180591620717411303424    (超过了int64的范围)
	// YiB // 1208925819614629174706176
)

const (
	Param    = 100
	AntsSize = 1000
	TestSize = 10000
	n        = 100000
)

var curMem uint64

// TestAntsPoolWaitToGetWorker 用来测试等待获取worker
func TestAntsPoolWaitToGetWorker(t *testing.T) {
	w := make(chan struct{})
	var wg sync.WaitGroup
	p, _ := NewPool(6)
	defer p.Release()

	for i := 0; i < n; i++ {
		wg.Add(1)
		// 提交n个任务
		_ = p.Submit(func() {
			demoPoolFunc(Param)
			wg.Done()
		})
	}

	wg.Wait()
	t.Logf("pool, running workers number:%d", p.Running())
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)

	_ = http.ListenAndServe(":6060", nil)
	<-w
	/**
	=== RUN   TestAntsPoolWaitToGetWorker
	    ants_test.go:74: pool, running workers number:1000
	    ants_test.go:78: memory usage:2 MB
	--- PASS: TestAntsPoolWaitToGetWorker (10.31s)
	PASS
	*/
}

// 预分配
func TestAntsPoolWaitToGetWorkerPreMalloc(t *testing.T) {
	var wg sync.WaitGroup
	// 开启预分配
	p, _ := NewPool(AntsSize, WithPreAlloc(true))
	defer p.Release()

	for i := 0; i < n; i++ {
		wg.Add(1)
		_ = p.Submit(func() {
			demoPoolFunc(Param)
			wg.Done()
		})
	}
	wg.Wait()
	t.Logf("pool, running workers number:%d", p.Running())
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
	/**
	=== RUN   TestAntsPoolWaitToGetWorkerPreMalloc
	    ants_test.go:103: pool, running workers number:1000
	    ants_test.go:107: memory usage:2 MB
	--- PASS: TestAntsPoolWaitToGetWorkerPreMalloc (10.26s)
	PASS
	*/
}

// TestAntsPoolWithFuncWaitToGetWorker is used to test waiting to get worker.
func TestAntsPoolWithFuncWaitToGetWorker(t *testing.T) {
	var wg sync.WaitGroup
	p, _ := NewPoolWithFunc(AntsSize, func(i interface{}) {
		demoPoolFunc(i)
		wg.Done()
	})
	defer p.Release()

	for i := 0; i < n; i++ {
		wg.Add(1)
		_ = p.Invoke(Param)
	}
	wg.Wait()
	t.Logf("pool with func, running workers number:%d", p.Running())
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
	/**
	=== RUN   TestAntsPoolWithFuncWaitToGetWorker
	    ants_test.go:131: pool with func, running workers number:1000
	    ants_test.go:135: memory usage:1 MB
	--- PASS: TestAntsPoolWithFuncWaitToGetWorker (10.29s)
	PASS
	*/
}

func TestAntsPoolWithFuncWaitToGetWorkerPreMalloc(t *testing.T) {
	var wg sync.WaitGroup
	p, _ := NewPoolWithFunc(AntsSize, func(i interface{}) {
		demoPoolFunc(i)
		wg.Done()
	}, WithPreAlloc(true))
	defer p.Release()

	for i := 0; i < n; i++ {
		wg.Add(1)
		_ = p.Invoke(Param)
	}
	wg.Wait()
	t.Logf("pool with func, running workers number:%d", p.Running())
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
	/**
	=== RUN   TestAntsPoolWithFuncWaitToGetWorkerPreMalloc
	    ants_test.go:158: pool with func, running workers number:1000
	    ants_test.go:162: memory usage:1 MB
	--- PASS: TestAntsPoolWithFuncWaitToGetWorkerPreMalloc (10.35s)
	PASS
	*/
}

// TestAntsPoolGetWorkerFromCache 从sync.Pool中获取worker.
func TestAntsPoolGetWorkerFromCache(t *testing.T) {
	p, _ := NewPool(TestSize)
	defer p.Release()

	for i := 0; i < AntsSize; i++ {
		_ = p.Submit(demoFunc)
	}
	time.Sleep(2 * DefaultCleanIntervalTime)
	_ = p.Submit(demoFunc)
	t.Logf("pool, running workers number:%d", p.Running())
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
	/**
	=== RUN   TestAntsPoolGetWorkerFromCache
	    ants_test.go:182: pool, running workers number:1
	    ants_test.go:186: memory usage:0 MB
	--- PASS: TestAntsPoolGetWorkerFromCache (2.00s)
	PASS
	*/
}

// TestAntsPoolWithFuncGetWorkerFromCache is used to test getting worker from sync.Pool.
func TestAntsPoolWithFuncGetWorkerFromCache(t *testing.T) {
	dur := 10
	p, _ := NewPoolWithFunc(TestSize, demoPoolFunc)
	defer p.Release()

	for i := 0; i < AntsSize; i++ {
		_ = p.Invoke(dur)
	}
	time.Sleep(2 * DefaultCleanIntervalTime)
	_ = p.Invoke(dur)
	t.Logf("pool with func, running workers number:%d", p.Running())
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
	/**
	=== RUN   TestAntsPoolWithFuncGetWorkerFromCache
	    ants_test.go:207: pool with func, running workers number:1
	    ants_test.go:211: memory usage:0 MB
	--- PASS: TestAntsPoolWithFuncGetWorkerFromCache (2.00s)
	PASS
	*/
}

func TestAntsPoolWithFuncGetWorkerFromCachePreMalloc(t *testing.T) {
	dur := 10
	p, _ := NewPoolWithFunc(TestSize, demoPoolFunc, WithPreAlloc(true))
	defer p.Release()

	for i := 0; i < AntsSize; i++ {
		_ = p.Invoke(dur)
	}
	time.Sleep(2 * DefaultCleanIntervalTime)
	_ = p.Invoke(dur)
	t.Logf("pool with func, running workers number:%d", p.Running())
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
	/**
	=== RUN   TestAntsPoolWithFuncGetWorkerFromCachePreMalloc
	    ants_test.go:231: pool with func, running workers number:1
	    ants_test.go:235: memory usage:1 MB
	--- PASS: TestAntsPoolWithFuncGetWorkerFromCachePreMalloc (2.00s)
	PASS
	*/
}

//-------------------------------------------------------------------------------------------
// 对比 在不使用pool和使用ants pool的两种情况
//-------------------------------------------------------------------------------------------

func TestNoPool(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			demoFunc()
			wg.Done()
		}()
	}

	wg.Wait()
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
	/*
		=== RUN   TestNoPool
		    ants_test.go:263: memory usage:16 MB
		--- PASS: TestNoPool (0.17s)
		PASS
	*/
}

func TestAntsPool(t *testing.T) {
	defer Release()
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		_ = Submit(func() {
			demoFunc()
			wg.Done()
		})
	}
	wg.Wait()

	t.Logf("pool, capacity:%d", Cap())
	t.Logf("pool, running workers number:%d", Running())
	t.Logf("pool, free workers number:%d", Free())

	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
	/**
	=== RUN   TestAntsPool
	    ants_test.go:284: pool, capacity:2147483647
	    ants_test.go:285: pool, running workers number:13130
	    ants_test.go:286: pool, free workers number:2147470517
	    ants_test.go:291: memory usage:12 MB
	--- PASS: TestAntsPool (0.18s)
	PASS
	*/
}

//-------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------

func TestPanicHandler(t *testing.T) {
	var panicCounter int64
	var wg sync.WaitGroup
	p0, err := NewPool(10, WithPanicHandler(func(p interface{}) {
		defer wg.Done()
		atomic.AddInt64(&panicCounter, 1)
		t.Logf("catch panic with PanicHandler: %v", p)
	}))
	assert.NoErrorf(t, err, "create new pool failed: %v", err)
	defer p0.Release()
	wg.Add(1)
	_ = p0.Submit(func() {
		panic("Oops!")
	})
	wg.Wait()
	c := atomic.LoadInt64(&panicCounter)
	assert.EqualValuesf(t, 1, c, "panic handler didn't work, panicCounter: %d", c)
	assert.EqualValues(t, 0, p0.Running(), "pool should be empty after panic")
	p1, err := NewPoolWithFunc(10, func(p interface{}) { panic(p) }, WithPanicHandler(func(p interface{}) {
		defer wg.Done()
		atomic.AddInt64(&panicCounter, 1)
	}))
	assert.NoErrorf(t, err, "create new pool with func failed: %v", err)
	defer p1.Release()
	wg.Add(1)
	_ = p1.Invoke("Oops!")
	wg.Wait()
	c = atomic.LoadInt64(&panicCounter)
	assert.EqualValuesf(t, 2, c, "panic handler didn't work, panicCounter: %d", c)
	assert.EqualValues(t, 0, p1.Running(), "pool should be empty after panic")
}

func TestPanicHandlerPreMalloc(t *testing.T) {
	var panicCounter int64
	var wg sync.WaitGroup
	p0, err := NewPool(10, WithPreAlloc(true), WithPanicHandler(func(p interface{}) {
		defer wg.Done()
		atomic.AddInt64(&panicCounter, 1)
		t.Logf("catch panic with PanicHandler: %v", p)
	}))
	assert.NoErrorf(t, err, "create new pool failed: %v", err)
	defer p0.Release()
	wg.Add(1)
	_ = p0.Submit(func() {
		panic("Oops!")
	})
	wg.Wait()
	c := atomic.LoadInt64(&panicCounter)
	assert.EqualValuesf(t, 1, c, "panic handler didn't work, panicCounter: %d", c)
	assert.EqualValues(t, 0, p0.Running(), "pool should be empty after panic")
	p1, err := NewPoolWithFunc(10, func(p interface{}) { panic(p) }, WithPanicHandler(func(p interface{}) {
		defer wg.Done()
		atomic.AddInt64(&panicCounter, 1)
	}))
	assert.NoErrorf(t, err, "create new pool with func failed: %v", err)
	defer p1.Release()
	wg.Add(1)
	_ = p1.Invoke("Oops!")
	wg.Wait()
	c = atomic.LoadInt64(&panicCounter)
	assert.EqualValuesf(t, 2, c, "panic handler didn't work, panicCounter: %d", c)
	assert.EqualValues(t, 0, p1.Running(), "pool should be empty after panic")
}

func TestPoolPanicWithoutHandler(t *testing.T) {
	p0, err := NewPool(10)
	assert.NoErrorf(t, err, "create new pool failed: %v", err)
	defer p0.Release()
	_ = p0.Submit(func() {
		panic("Oops!")
	})

	p1, err := NewPoolWithFunc(10, func(p interface{}) {
		panic(p)
	})
	assert.NoErrorf(t, err, "create new pool with func failed: %v", err)
	defer p1.Release()
	_ = p1.Invoke("Oops!")
}

func TestPoolPanicWithoutHandlerPreMalloc(t *testing.T) {
	p0, err := NewPool(10, WithPreAlloc(true))
	assert.NoErrorf(t, err, "create new pool failed: %v", err)
	defer p0.Release()
	_ = p0.Submit(func() {
		panic("Oops!")
	})

	p1, err := NewPoolWithFunc(10, func(p interface{}) {
		panic(p)
	})

	assert.NoErrorf(t, err, "create new pool with func failed: %v", err)

	defer p1.Release()
	_ = p1.Invoke("Oops!")
}

func TestPurge(t *testing.T) {
	p, err := NewPool(10)
	assert.NoErrorf(t, err, "create TimingPool failed: %v", err)
	defer p.Release()
	_ = p.Submit(demoFunc)
	time.Sleep(3 * DefaultCleanIntervalTime)
	assert.EqualValues(t, 0, p.Running(), "all p should be purged")
	p1, err := NewPoolWithFunc(10, demoPoolFunc)
	assert.NoErrorf(t, err, "create TimingPoolWithFunc failed: %v", err)
	defer p1.Release()
	_ = p1.Invoke(1)
	time.Sleep(3 * DefaultCleanIntervalTime)
	assert.EqualValues(t, 0, p.Running(), "all p should be purged")
}

func TestPurgePreMalloc(t *testing.T) {
	p, err := NewPool(10, WithPreAlloc(true))
	assert.NoErrorf(t, err, "create TimingPool failed: %v", err)
	defer p.Release()
	_ = p.Submit(demoFunc)
	time.Sleep(3 * DefaultCleanIntervalTime)
	assert.EqualValues(t, 0, p.Running(), "all p should be purged")
	p1, err := NewPoolWithFunc(10, demoPoolFunc)
	assert.NoErrorf(t, err, "create TimingPoolWithFunc failed: %v", err)
	defer p1.Release()
	_ = p1.Invoke(1)
	time.Sleep(3 * DefaultCleanIntervalTime)
	assert.EqualValues(t, 0, p.Running(), "all p should be purged")
}

func TestNonblockingSubmit(t *testing.T) {
	poolSize := 10
	p, err := NewPool(poolSize, WithNonblocking(true))
	assert.NoErrorf(t, err, "create TimingPool failed: %v", err)
	defer p.Release()
	for i := 0; i < poolSize-1; i++ {
		assert.NoError(t, p.Submit(longRunningFunc), "nonblocking submit when pool is not full shouldn't return error")
	}
	ch := make(chan struct{})
	ch1 := make(chan struct{})
	f := func() {
		<-ch
		close(ch1)
	}
	// p is full now.
	assert.NoError(t, p.Submit(f), "nonblocking submit when pool is not full shouldn't return error")
	assert.EqualError(t, p.Submit(demoFunc), ErrPoolOverload.Error(),
		"nonblocking submit when pool is full should get an ErrPoolOverload")
	// interrupt f to get an available worker
	close(ch)
	<-ch1
	assert.NoError(t, p.Submit(demoFunc), "nonblocking submit when pool is not full shouldn't return error")
}

func TestMaxBlockingSubmit(t *testing.T) {
	poolSize := 10
	p, err := NewPool(poolSize, WithMaxBlockingTasks(1))
	assert.NoErrorf(t, err, "create TimingPool failed: %v", err)
	defer p.Release()
	for i := 0; i < poolSize-1; i++ {
		assert.NoError(t, p.Submit(longRunningFunc), "submit when pool is not full shouldn't return error")
	}
	ch := make(chan struct{})
	f := func() {
		<-ch
	}
	// p is full now.
	assert.NoError(t, p.Submit(f), "submit when pool is not full shouldn't return error")
	var wg sync.WaitGroup
	wg.Add(1)
	errCh := make(chan error, 1)
	go func() {
		// should be blocked. blocking num == 1
		if err := p.Submit(demoFunc); err != nil {
			errCh <- err
		}
		wg.Done()
	}()
	time.Sleep(1 * time.Second)
	// already reached max blocking limit
	assert.EqualError(t, p.Submit(demoFunc), ErrPoolOverload.Error(),
		"blocking submit when pool reach max blocking submit should return ErrPoolOverload")
	// interrupt f to make blocking submit successful.
	close(ch)
	wg.Wait()
	select {
	case <-errCh:
		t.Fatalf("blocking submit when pool is full should not return error")
	default:
	}
}

func TestNonblockingSubmitWithFunc(t *testing.T) {
	poolSize := 10
	ch1 := make(chan struct{})
	p, err := NewPoolWithFunc(poolSize, func(i interface{}) {
		longRunningPoolFunc(i)
		close(ch1)
	}, WithNonblocking(true))
	assert.NoError(t, err, "create TimingPool failed: %v", err)
	defer p.Release()
	for i := 0; i < poolSize-1; i++ {
		assert.NoError(t, p.Invoke(nil), "nonblocking submit when pool is not full shouldn't return error")
	}
	ch := make(chan struct{})
	// p is full now.
	assert.NoError(t, p.Invoke(ch), "nonblocking submit when pool is not full shouldn't return error")
	assert.EqualError(t, p.Invoke(nil), ErrPoolOverload.Error(),
		"nonblocking submit when pool is full should get an ErrPoolOverload")
	// interrupt f to get an available worker
	close(ch)
	<-ch1
	assert.NoError(t, p.Invoke(nil), "nonblocking submit when pool is not full shouldn't return error")
}

// 设置了最大的阻塞值
func TestMaxBlockingSubmitWithFunc(t *testing.T) {
	poolSize := 10
	p, err := NewPoolWithFunc(poolSize, longRunningPoolFunc, WithMaxBlockingTasks(1))
	assert.NoError(t, err, "create TimingPool failed: %v", err)
	defer p.Release()
	for i := 0; i < poolSize-1; i++ {
		assert.NoError(t, p.Invoke(Param), "submit when pool is not full shouldn't return error")
	}
	ch := make(chan struct{})
	// p is full now.
	assert.NoError(t, p.Invoke(ch), "submit when pool is not full shouldn't return error")
	var wg sync.WaitGroup
	wg.Add(1)
	errCh := make(chan error, 1)
	go func() {
		// should be blocked. blocking num == 1
		if err := p.Invoke(Param); err != nil {
			errCh <- err
		}
		wg.Done()
	}()
	time.Sleep(1 * time.Second)
	// already reached max blocking limit
	assert.EqualErrorf(t, p.Invoke(Param), ErrPoolOverload.Error(),
		"blocking submit when pool reach max blocking submit should return ErrPoolOverload: %v", err)
	// interrupt one func to make blocking submit successful.
	close(ch)
	wg.Wait()
	select {
	case <-errCh:
		t.Fatalf("blocking submit when pool is full should not return error")
	default:
	}
}

func TestRebootDefaultPool(t *testing.T) {
	defer Release()
	Reboot()
	var wg sync.WaitGroup
	wg.Add(1)
	_ = Submit(func() {
		demoFunc()
		wg.Done()
	})
	wg.Wait()
	Release()
	assert.EqualError(t, Submit(nil), ErrPoolClosed.Error(), "pool should be closed")
	Reboot()
	wg.Add(1)
	assert.NoError(t, Submit(func() { wg.Done() }), "pool should be rebooted")
	wg.Wait()
}

func TestRebootNewPool(t *testing.T) {
	var wg sync.WaitGroup
	p, err := NewPool(10)
	assert.NoErrorf(t, err, "create Pool failed: %v", err)
	defer p.Release()
	wg.Add(1)
	_ = p.Submit(func() {
		demoFunc()
		wg.Done()
	})
	wg.Wait()
	p.Release()
	assert.EqualError(t, p.Submit(nil), ErrPoolClosed.Error(), "pool should be closed")
	p.Reboot()
	wg.Add(1)
	assert.NoError(t, p.Submit(func() { wg.Done() }), "pool should be rebooted")
	wg.Wait()

	p1, err := NewPoolWithFunc(10, func(i interface{}) {
		demoPoolFunc(i)
		wg.Done()
	})
	assert.NoErrorf(t, err, "create TimingPoolWithFunc failed: %v", err)
	defer p1.Release()
	wg.Add(1)
	_ = p1.Invoke(1)
	wg.Wait()
	p1.Release()
	assert.EqualError(t, p1.Invoke(nil), ErrPoolClosed.Error(), "pool should be closed")
	p1.Reboot()
	wg.Add(1)
	assert.NoError(t, p1.Invoke(1), "pool should be rebooted")
	wg.Wait()
}

func TestInfinitePool(t *testing.T) {
	c := make(chan struct{})
	p, _ := NewPool(-1)
	_ = p.Submit(func() {
		_ = p.Submit(func() {
			<-c
		})
	})
	c <- struct{}{}
	if n := p.Running(); n != 2 {
		t.Errorf("expect 2 workers running, but got %d", n)
	}
	p.Tune(10)
	if capacity := p.Cap(); capacity != -1 {
		t.Fatalf("expect capacity: -1 but got %d", capacity)
	}
	var err error
	p, err = NewPool(-1, WithPreAlloc(true))
	if err != ErrInvalidPreAllocSize {
		t.Errorf("expect ErrInvalidPreAllocSize but got %v", err)
	}
}

// 测试覆盖率
func TestRestCodeCoverage(t *testing.T) {
	_, err := NewPool(-1, WithExpiryDuration(-1))
	t.Log(err)
	_, err = NewPool(1, WithExpiryDuration(-1))
	t.Log(err)
	_, err = NewPoolWithFunc(-1, demoPoolFunc, WithExpiryDuration(-1))
	t.Log(err)
	_, err = NewPoolWithFunc(1, demoPoolFunc, WithExpiryDuration(-1))
	t.Log(err)

	options := Options{}
	options.ExpiryDuration = time.Duration(10) * time.Second
	options.Nonblocking = true
	options.PreAlloc = true
	poolOpts, _ := NewPool(1, WithOptions(options))
	t.Logf("Pool with options, capacity: %d", poolOpts.Cap())

	p0, _ := NewPool(TestSize, WithLogger(log.New(os.Stderr, "", log.LstdFlags)))
	defer func() {
		_ = p0.Submit(demoFunc)
	}()
	defer p0.Release()
	for i := 0; i < n; i++ {
		_ = p0.Submit(demoFunc)
	}
	t.Logf("pool, capacity:%d", p0.Cap())
	t.Logf("pool, running workers number:%d", p0.Running())
	t.Logf("pool, free workers number:%d", p0.Free())
	p0.Tune(TestSize)
	p0.Tune(TestSize / 10)
	t.Logf("pool, after tuning capacity, capacity:%d, running:%d", p0.Cap(), p0.Running())

	pprem, _ := NewPool(TestSize, WithPreAlloc(true))
	defer func() {
		_ = pprem.Submit(demoFunc)
	}()
	defer pprem.Release()
	for i := 0; i < n; i++ {
		_ = pprem.Submit(demoFunc)
	}
	t.Logf("pre-malloc pool, capacity:%d", pprem.Cap())
	t.Logf("pre-malloc pool, running workers number:%d", pprem.Running())
	t.Logf("pre-malloc pool, free workers number:%d", pprem.Free())
	pprem.Tune(TestSize)
	pprem.Tune(TestSize / 10)
	t.Logf("pre-malloc pool, after tuning capacity, capacity:%d, running:%d", pprem.Cap(), pprem.Running())

	p, _ := NewPoolWithFunc(TestSize, demoPoolFunc)
	defer func() {
		_ = p.Invoke(Param)
	}()
	defer p.Release()
	for i := 0; i < n; i++ {
		_ = p.Invoke(Param)
	}
	time.Sleep(DefaultCleanIntervalTime)
	t.Logf("pool with func, capacity:%d", p.Cap())
	t.Logf("pool with func, running workers number:%d", p.Running())
	t.Logf("pool with func, free workers number:%d", p.Free())
	p.Tune(TestSize)
	p.Tune(TestSize / 10)
	t.Logf("pool with func, after tuning capacity, capacity:%d, running:%d", p.Cap(), p.Running())

	ppremWithFunc, _ := NewPoolWithFunc(TestSize, demoPoolFunc, WithPreAlloc(true))
	defer func() {
		_ = ppremWithFunc.Invoke(Param)
	}()
	defer ppremWithFunc.Release()
	for i := 0; i < n; i++ {
		_ = ppremWithFunc.Invoke(Param)
	}
	time.Sleep(DefaultCleanIntervalTime)
	t.Logf("pre-malloc pool with func, capacity:%d", ppremWithFunc.Cap())
	t.Logf("pre-malloc pool with func, running workers number:%d", ppremWithFunc.Running())
	t.Logf("pre-malloc pool with func, free workers number:%d", ppremWithFunc.Free())
	ppremWithFunc.Tune(TestSize)
	ppremWithFunc.Tune(TestSize / 10)
	t.Logf("pre-malloc pool with func, after tuning capacity, capacity:%d, running:%d", ppremWithFunc.Cap(),
		ppremWithFunc.Running())
}
