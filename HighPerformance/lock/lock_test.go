/**
* @Author:zhoutao
* @Date:2021/1/2 上午10:38
* @Desc:
 */

package lock

import (
	"sync"
	"testing"
	"time"
)

type RW interface {
	Write()
	Read()
}

const cost = time.Microsecond

//const cost = time.Nanosecond * 10
//const cost = time.Microsecond * 10

type Lock struct {
	count int
	mu    sync.Mutex
}

func (l *Lock) Write() {
	l.mu.Lock()

	l.count++
	time.Sleep(cost)

	l.mu.Unlock()
}

func (l *Lock) Read() {
	l.mu.Lock()

	time.Sleep(cost)
	_ = l.count

	l.mu.Unlock()
}

type RWlock struct {
	count int
	mu    sync.RWMutex
}

func (l *RWlock) Write() {
	l.mu.Lock()

	l.count++
	time.Sleep(cost)

	l.mu.Unlock()
}

func (l *RWlock) Read() {
	l.mu.RLock()

	time.Sleep(cost)
	_ = l.count

	l.mu.RUnlock()
}

func benchmark(b *testing.B, rw RW, read, write int) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for k := 0; k < read*100; k++ {
			wg.Add(1)
			go func() {
				rw.Read()
				wg.Done()
			}()
		}

		for k := 0; k < write*100; k++ {
			wg.Add(1)
			go func() {
				rw.Write()
				wg.Done()
			}()
		}
	}
}

//写锁：   读多写少：读9，写1
func BenchmarkReadMore(b *testing.B) { benchmark(b, &Lock{}, 9, 1) }

//读写锁： 读多写少：读9，写1
func BenchmarkReadMoreRW(b *testing.B) { benchmark(b, &RWlock{}, 9, 1) }

//写锁：   写多：读1，写9
func BenchmarkWriteMore(b *testing.B) { benchmark(b, &Lock{}, 1, 9) }

//读写锁:  写多：读1，写9
func BenchmarkWriteMoreRW(b *testing.B) { benchmark(b, &RWlock{}, 1, 9) }

//写锁：   读写相同
func BenchmarkEqual(b *testing.B) { benchmark(b, &Lock{}, 5, 5) }

//读写锁： 读写相同
func BenchmarkEqualRw(b *testing.B) { benchmark(b, &RWlock{}, 5, 5) }

/**

//1秒杀
tao@taodeMacBook-Pro lock % go test lock_test.go -bench='^Benchmark'
goos: darwin
goarch: amd64

//读多写少： 读写锁 > 写锁
BenchmarkReadMore-4      	     421	   2654042 ns/op
BenchmarkReadMoreRW-4    	    1759	   1023469 ns/op

// 和预期不符: 写锁 > 读写锁 ，预期为：读锁和写锁相当
BenchmarkWriteMore-4     	    2469	   3861325 ns/op
BenchmarkWriteMoreRW-4   	     417	   2544585 ns/op

//读写相同：读写锁 > 写锁
BenchmarkEqual-4         	     385	   2907275 ns/op
BenchmarkEqualRw-4       	     870	   1278216 ns/op
PASS
ok  	command-line-arguments	83.635s

//1纳秒
tao@taodeMacBook-Pro lock % go test lock_test.go -bench='^Benchmark'
goos: darwin
goarch: amd64
BenchmarkReadMore-4      	     999	   1397176 ns/op
BenchmarkReadMoreRW-4    	    1622	    733912 ns/op
BenchmarkWriteMore-4     	    1641	   1037246 ns/op
BenchmarkWriteMoreRW-4   	    1281	    857281 ns/op
BenchmarkEqual-4         	    1570	    646676 ns/op
BenchmarkEqualRw-4       	    1880	    877343 ns/op
PASS
ok  	command-line-arguments	14.372s

//10毫秒
tao@taodeMacBook-Pro lock % go test lock_test.go -bench='^Benchmark'
goos: darwin
goarch: amd64
BenchmarkReadMore-4      	     459	   2475461 ns/op
BenchmarkReadMoreRW-4    	    1640	   1030786 ns/op
BenchmarkWriteMore-4     	     810	   3056037 ns/op
BenchmarkWriteMoreRW-4   	    1206	   2392134 ns/op
BenchmarkEqual-4         	     370	   2958725 ns/op
BenchmarkEqualRw-4       	     997	   1638529 ns/op
PASS
ok  	command-line-arguments	73.240s
*/

//互斥锁的两种状态：正常状态和饥饿状态
// 正常状态有很好的性能表现，饥饿状态也非常重要，因为它能阻止尾部延迟的现象
