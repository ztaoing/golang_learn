/**
* @Author:zhoutao
* @Date:2021/1/1 下午3:02
* @Desc: 内存对齐
 */

package main

import (
	"fmt"
	"unsafe"
)

//可以使用unsafe.Sizeof 计算出一个数据类型实例需要占用的字节数

type Args struct {
	num1 int
	num2 int
}

type Flag struct {
	num1 int16
	num2 int32
}

func main() {
	fmt.Println(unsafe.Sizeof(Args{}))
	fmt.Println(unsafe.Sizeof(Flag{}))

	/**
	16  // Args 由连个int组成，在64位的机器上，一个int占8字节，因此存储一个Args，需要16字节
	8	// Flag int32 占4字节，int16占2字节，8-4-2 = 2，多出来的字节是内存对齐的结果
	因此一个结构体实例所占的空间等于各个字段占据空间之和+ 内存对齐的空间大小
	*/

	fmt.Println(unsafe.Alignof(Args{}))
	fmt.Println(unsafe.Alignof(Flag{}))
	/**
	8  Args的对齐值
	4  Flag的对齐值

	*/

	fmt.Println(unsafe.Sizeof(demo1{}))
	fmt.Println(unsafe.Sizeof(demo2{}))

	/**
	8
	12
	在对内存特别敏感的结构体的设计上，我们可以通过调整资金端的顺序，减少内存的占用
	*/

	fmt.Println(unsafe.Sizeof(demo3{}))
	fmt.Println(unsafe.Sizeof(demo4{}))
	/**
	8 额外补充了4字节的内存空间
	4 与c占用内存空间一致

	*/

}

//cpu 访问内存是以字长为单位访问。比如32的CPU，字长为4字节，那么CPU访问内存也是4字节
//内存对齐是每次的内存访问都是原子的，如果变量的大小不超过字长，那么内存对齐后，对该变量的访问就是原子的，这个特性在并发场景下至关重要

//struct 的对齐

type demo1 struct {
	a int8  // 1字节  从第0个位置开始占据一个字节
	b int16 // 2字节  对齐倍数是2，从第二个位置占据2个字节
	c int32 // 4字节  内存已对齐，从第4个位置开始占据4个字节
}

type demo2 struct {
	a int8  // 1字节 从第0个位置开始占据一个字节
	c int32 // 4字节 对齐倍数是4，必须跳过3个字节，偏移量才是4的倍数，从第4个位置开始占据4个字节
	b int16 // 对齐倍数是2，从第8个位置开始占据2个字节
}

//空struct的对齐
//当struct作为结构体最后一个字段时，需要内存对齐，因为如果有指针指向该字段，返回的地址将在结构体之外
//如果此指针一直存活不释放对应的内存，就会有内存泄漏的问题（该内存不因结构体释放而释放）

type demo3 struct {
	c int32
	a struct{}
}

type demo4 struct {
	a struct{}
	c int32
}
