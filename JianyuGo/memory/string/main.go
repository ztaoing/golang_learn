/**
* @Author:zhoutao
* @Date:2022/2/18 08:56
* @Desc:
 */

package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

// 由小结构 向大的结构转换，导致内存占用变大，变大后的结构占用了后边结构的内存，导致后边结构的前边的内存的内容被覆盖了
func StringToByte(key *string) []byte {
	// 把这个地址通过 unsafe 库强转类型，当作 slice 管理结构来用，往后多踩了 8 字节(string的struct是16，slice的struct是24)，
	// 也就是说把变量 key 的头 8 字节踩掉了。
	// key 的头 8 个字节原本是个地址，指向堆上内存字符串所在位置。现在却被无情的踩成了一个整数 16 。
	strPtr := (*reflect.SliceHeader)(unsafe.Pointer(key))
	strPtr.Cap = strPtr.Len

	// 重点来了，且覆盖写 strPtr.Cap 这个字段。这就踩内存了，往后多踩了 8 字节的内存。
	b := *(*[]byte)(unsafe.Pointer(strPtr))

	return b
}

func main() {
	decryptContent := "/AvYEjm4g6xJ3LVrk2/Adk"
	iv := decryptContent[0:16]
	key := decryptContent[2:18]
	fmt.Println(&iv)
	fmt.Println(&key)

	// 直接取了 key，iv 这两个变量本身的地址，作为参数传进了 StringToByte

	// StringToByte 函数把这 24 个字节复制给一个栈上的局部变量 ivBytes
	ivBytes := StringToByte(&iv)

	// 又执行了一次 StringToByte 的函数，这次传入的地址是 key 变量的地址，又是往后踩了 8 字节。
	// 本应该是指针的字段，却活生生被踩成了 16 ，然后把这个值 16 当作指针传递到 slicebytetostring 函数里去转类型，如果这都不出非法地址的 panic ，那才真的是神奇了。
	keyBytes := StringToByte(&key)
	// 这行代码为什么不会报错，因为 ivBytes 变量没事。ivBytes 能够转成 string
	fmt.Println(string(ivBytes))
	// 则会出 panic（还记得吗？我们文章最开始的截图，panic 的位置就是 25 行），但是变量 key 被踩了呀，导致 keyBytes 这个变量也是错的。
	// 明明参数是指针，但是却传了一个 16 进去，这个就是为什么出 panic 的原因了
	fmt.Println(string(keyBytes))
}

/**

0xc000010240
/AvYEjm4g6xJ3LVr
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x10 pc=0x105c030]

*/
/**
	 keyBytes := StringToByte(&key)
	 ivBytes := StringToByte(&iv)

当然了，最后还是顺手把变量 key 的头部踩成 16 了，不过此时已经没有影响，因为下面用的是 keyBytes, ivBytes ，所以程序自然可以正常运行。
*/

/**
怎么才能把程序改正确？

func StringToByte(key *string) []byte {

    strPtr := (*reflect.SliceHeader)(unsafe.Pointer(key))
    strPtr.Cap = strPtr.Len

    b := *(*[]byte)(unsafe.Pointer(strPtr))
    return b
}

func StringToByte(key *string) []byte {
	加了这行代码之后，就不会踩到外面的内存了，因为这样先在栈上分配出 24 字节的局部变量，然后是在这个局部变量上赋值的，
    slic := reflect.SliceHeader{}

    slic = *(*reflect.SliceHeader)(unsafe.Pointer(key))
    slic.Cap = slic.Len

    b := *(*[]byte)(unsafe.Pointer(&slic))
    return b
}
*/
