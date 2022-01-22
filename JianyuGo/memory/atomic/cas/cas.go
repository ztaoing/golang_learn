package main

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

/**
该操作简称CAS (Compare And Swap)。这类操作的前缀为 CompareAndSwap :
func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)

func CompareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool)

该操作在进行交换前首先确保被操作数的值未被更改，即仍然保存着参数 old 所记录的值，满足此前提条件下才进行交换操作。
CAS的做法类似操作数据库时常见的乐观锁机制。

需要注意的是，当有大量的goroutine 对变量进行读写操作时，可能导致CAS操作无法成功，这时可以利用for循环多次尝试。

上面我只列出了比较典型的int32和unsafe.Pointer类型的CAS方法，主要是想说除了读数值类型进行比较交换，还支持对指针进行比较交换。

unsafe.Pointer提供了绕过Go语言指针类型限制的方法，unsafe指的并不是说不安全，而是说官方并不保证向后兼容。
*/

// 定义一个struct类型P
type P struct{ x, y, z int }

// 执行类型P的指针
var pP *P

//示例并不是在并发环境下进行的CAS，只是为了演示效果，先把被操作数设置成了Old Pointer。
func main() {

	// 定义一个执行unsafe.Pointer"值的指针"变量
	var unsafe1 = (*unsafe.Pointer)(unsafe.Pointer(&pP))

	// Old pointer
	var sy P

	// 为了演示效果先将unsafe1设置成Old Pointer
	px := atomic.SwapPointer(
		unsafe1, unsafe.Pointer(&sy))

	// 执行CAS操作，交换成功，结果返回true
	y := atomic.CompareAndSwapPointer(
		unsafe1, unsafe.Pointer(&sy), px)

	fmt.Println(y)
}

/**
其实Mutex的底层实现也是依赖原子操作中的CAS实现的，原子操作的atomic包相当于是sync包里的那些同步原语的实现依赖。

比如互斥锁Mutex的结构里有一个state字段，其是表示锁状态的状态位。
type Mutex struct {
	 state int32
	 sema  uint32
}

为了方便理解，我们在这里将它的状态定义为0和1，0代表目前该锁空闲，1代表已被加锁，以下是sync.Mutex中Lock方法的部分实现代码。
func (m *Mutex) Lock() {
   // Fast path: grab unlocked mutex.
// 在atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked)中，m.state代表锁的状态，通过CAS方法，
// 判断锁此时的状态是否空闲（m.state==0），是，则对其加锁（mutexLocked常量的值为1）。
   if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
       if race.Enabled {
           race.Acquire(unsafe.Pointer(m))
       }
       return
   }
   // Slow path (outlined so that the fast path can be inlined)
    m.lockSlow()
}
*/
