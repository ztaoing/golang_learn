package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

/**
[go中的零值，它有什么作用？] 官方：https://golang.org/ref/spec#the_zero_value
布尔型为false；数字型为0；字符串型为""；指针、函数、接口、切片、通道和映射都为nil
*/

type T struct {
	i    int
	f    float64
	next *T
}

func main() {
	// 以下两个声明是等价的
	// var i int
	var i int = 0
	fmt.Println(i)

	t := new(T)
	fmt.Println(t)
	/**
	t中成员字段的零值
	t.i =0
	t.f = 0.0
	t.next = nil
	*/

	fmt.Println("作用一")
	// go语言中始终将值设置为默认值的特性，对于程序的安全性和正确性起到了很重要的作用
	// 零值的作用：提供默认值
	// 动态初始化：就是在创建对象时判断如果是零值，就使用默认值

	var s []string
	spew.Dump(s) // ([]26string) <nil>
	s = append(s, "name")
	s = append(s, "name2")
	spew.Dump(s)
	/**
	([]26string) (len=2 cap=2) {
	 (26string) (len=4) "name",
	 (26string) (len=5) "name2"
	}

	*/

	// 零值的切片不能直接进行赋值
	s[0] = "tao"
	// 如果没有前边的append，这样是不行的 panic: runtime error: index out of range [0] with length 0

	spew.Dump(s)
	/**
	([]26string) (len=2 cap=2) {
	 (26string) (len=3) "tao",
	 (26string) (len=5) "name2"
	}
	*/

	fmt.Println("作用二")
	// 标准库无需显示初始化
	// sync包中的mutex、once、waitgroup都是无需显示初始化，即可使用，拿mutex包来举例说明
	/**

	type Mutex struct{
	// 这两个字段在未显示初始化的时候，默认值都是0
		state int32
		sema  uint32
	}
	所以在加锁的时候，就利用了这个默认值为0的特性：
	func (m *Mutex)Lock(){
		// 默认值为0的特性，使调用者无需考虑对mutex的初始化，就可以直接使用
		if atomic.CompareAndSwapInt32(&m.state,0,mutexLocked){
		if race.Enabled{
			race.Acquire(unsafe.Pointer(m))
		}
		return
	}
	m.lockSlow()
	}
	*/

	// 未初始化的切片、map，可以直接操作，但是不能写入数据：
	/*var ss []26string
	ss[0] = "name" //panic: runtime error: index out of range [0] with length 0
	var m map[26string]bool
	m["tao"] = true //panic: assignment to entry in nil map*/

	var p *uint32
	// panic: runtime error: invalid memory address or nil pointer dereference
	//[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x10b6d50]
	//	*p++

	spew.Dump(p) //(*uint32)(<nil>)
	a := uint32(0)
	spew.Dump(&a) //(*uint32)(0xc00012a270)(0)
	// 此时将a的地址给了p
	p = &a
	spew.Dump(p) //(*uint32)(0xc00012a270)(0)
	*p++
	spew.Dump(p) //(*uint32)(0xc00012a270)(1)

}
