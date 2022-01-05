package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// 在go语言，甚至大部分语言中，一条普通的复制语句其实不是一个原子操作
// 例如：在32位机器上写int64类型的变量就会有中间状态，因为他会被拆分成两次写操作：写低32位和高32位

// 原子操作由底层硬件支持，对于一个变量更新的保护，原子操作通常会更有效率，并且更能利用计算机多核的优势，
// 如果要更新的是一个复合对象，则应当使用atomic.Value封装好的实现。
func main() {
	/**
	atomic.Value类型提供了两个读写方法：
	v.Sotre(c) :写，将原始的变量c存放到一个atomic.Value类型的v里
	c:=v.Load() :读，从v中，线程安全的读取上一步存放的内容

	*/

	wg := sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			update(i, i+5)
		}()
	}
	wg.Wait()
	// 从v中线程安全的读取
	r := rect.Load().(*Rectangle)
	fmt.Printf("rect.width=%d\nrect.length=%d\n", r.width, r.length)
	/**
	rect.width=10
	rect.length=15
	atomic.Value在不加锁的情况下提供了读写操作的线程安全保证，内部是如何实现的呢？

	1、atomic.Value 是interface{}类型，所以可以用来存储任意类型的数据，
	atomic包顶一个ifaceWords类型，其实是interface{}的内部表示（runtime.eface），他的作用是将interface{}类型分解
	得到其原始类型(typ)和真正的值(data)
	type ifaceWords struct{
	//ifaceWords is interface{} internal representation
		typ  unsafe.Pointer
		data unsafe.Pointer
	}

	2、unsafe.Pointer
	出于安全考虑，go语言不支持直接操作内存，但他的标准库中又提供了一种不安全（不保证向后兼容）的指针类型unsafe.Pointer，让程序可以灵活的操作内存
	unsafe.Pointer可以与任意的指针类型互相转换。
	也就说，如果两种类型具有相同的内存结构(layout)，我们可以将unsafe.Pointer作为桥梁，然后让这两种类型的指针相互转换，从而是实现
	同一个内存，拥有两种不同的解读方式。

	例如：[]byte和string的内部存储结构其实是一样的，他们在运行时类型分别为reflect.SliceHeader 和reflect.StringHeader
	type SliceHeader {
		Data uintptr
		Len int
		Cap int
	}
	type StringHeader {
		Data uintptr
		Len int
	}
	但是go语言的类型系统禁止他俩相互转换。但是可以借助unsafe.Pointer ,就可以在实现零拷贝的情况下，将[]byte数组直接转换成string类型

	bytes:=[]byte{104,101,110,111}
	p:=unsafe.Pointer(&bytes) // 将*[]byte 指针强制转换为unsafe.Pointer
	// 注意这个写法的含义
	str:=*(*string)(p)        // 将unsafe.Pointer转换为string类型的指针，再将这个指针的值当做string类型取出来
	fmt.Println(str)          // 输出"hello"


	func (v *Value) Store(x interface{}) {
	  if x == nil {
	    panic("sync/atomic: store of nil value into Value")
	  }
	// 通过unsafe.Pointer将现有的和要写入的值分别转成ifaceWords类型，
	// 这样我们下一步就可以得到这两个interface{}的原始类型（typ）和真正的值（data）。
	  vp := (*ifaceWords)(unsafe.Pointer(v))  // Old value
	  xp := (*ifaceWords)(unsafe.Pointer(&x)) // New value

	// 1、开始就是一个无限 for 循环。配合CompareAndSwap使用，可以达到乐观锁的效果。
	  for {

		//通过LoadPointer这个原子操作拿到当前Value中存储的类型。下面根据这个类型的不同，分3种情况处理。
	    typ := LoadPointer(&vp.typ)

		//  typ == nil: 第一次写入 - 一个atomic.Value实例被初始化后，它的typ字段会被设置为指针的零值 nil，
	    // 所以先判断如果typ是nil 那就证明这个Value实例还未被写入过数据。那之后就是一段初始写入的操作：
	    if typ == nil {
			// 那之后就是一段初始写入的操作：
	      // Attempt to start first store.
	      // Disable preemption so that other goroutines can use
	      // active spin wait to wait for completion; and so that
	      // GC does not see the fake type accidentally.

	      // runtime_procPin()这是runtime中的一段函数，一方面它禁止了调度器对当前 goroutine 的抢占（preemption），
		  // 使得它在执行当前逻辑的时候不被打断，以便可以尽快地完成工作，因为别人一直在等待它。
	      // 另一方面，在禁止抢占期间，GC 线程也无法被启用，这样可以防止 GC 线程看到一个莫名其妙的指向^uintptr(0)的类型（这是赋值过程中的中间状态）。
	      runtime_procPin()

		  // 使用CAS操作，先尝试将typ设置为^uintptr(0)这个 中间状态 。如果失败，则证明已经有别的线程抢先完成了赋值操作，那它就解除抢占锁，然后重新回到 for 循环第一步。
	      if !CompareAndSwapPointer(&vp.typ, nil, unsafe.Pointer(^uintptr(0))) {
	        runtime_procUnpin()
	        continue
	      }

	      // Complete first store. 使用乐观锁的思想
		  // 如果设置成功，那证明当前线程抢到了这个"乐观锁”，它可以安全的把v设为传入的新值了。
		  // 注意，这里是先写data字段，然后再写typ字段。因为我们是以typ字段的值作为写入完成与否的判断依据的。
	      StorePointer(&vp.data, xp.data)
	      StorePointer(&vp.typ, xp.typ) // 我们是以typ字段的值作为写入完成与否的判断依据的。
	      runtime_procUnpin()
	      return
	    }

		// 第一次写入还未完成- 如果看到typ字段还是^uintptr(0)这个中间类型，
		// 证明刚刚的第一次写入还没有完成，所以它会继续循环，一直等到第一次写入完成。
	    if uintptr(typ) == ^uintptr(0) {
	      // First store in progress. Wait.
	      // Since we disable preemption around the first store,
	      // we can wait with active spinning.
	      continue
	    }

		// 第一次写入已完成 - 首先检查上一次写入的类型与这一次要写入的类型是否一致，如果不一致则抛出异常。
		// 反之，则直接把这一次要写入的值写入到data字段。
		// 这个逻辑的主要思想就是，为了完成多个字段的原子性写入，我们可以抓住其中的一个字段，以它的状态来标志整个原子写入的状态。
	    // First store completed. Check type and overwrite data.
	    if typ != xp.typ {
	      panic("sync/atomic: store of inconsistently typed value into Value")
	    }

		// 直接把这一次要写入的值写入到data字段。
	    StorePointer(&vp.data, xp.data)
	    return
	  }
	}




	func (v *Value) Load() (x interface{}) {
	  // 原值
	  vp := (*ifaceWords)(unsafe.Pointer(v))
	  typ := LoadPointer(&vp.typ)

	  // 如果当前的typ是 nil 或者^uintptr(0)，那就证明第一次写入还没有开始，或者还没完成，那就直接返回 nil （不对外暴露中间状态）。
	  if typ == nil || uintptr(typ) == ^uintptr(0) {
	    // First store not yet completed.
	    return nil
	  }
	  // 否则，根据当前看到的typ和data构造出一个新的interface{}返回出去。

	  // 通过LoadPointer这个原子操作拿到当前Value中存储的类型
	  data := LoadPointer(&vp.data)
	  // 注意这个写法
	  // 新值
	  xp := (*ifaceWords)(unsafe.Pointer(&x))
	  xp.typ = typ
	  xp.data = data
	  return
	}










	*/
}

// 长方形
type Rectangle struct {
	length int
	width  int
}

var rect atomic.Value

func update(width, length int) {
	rectLocal := new(Rectangle)
	rectLocal.length = length
	rectLocal.width = width
	//更新,原子更新
	rect.Store(rectLocal)
}
