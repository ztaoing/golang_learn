/**
* @Author:zhoutao
* @Date:2022/2/16 13:59
* @Desc:
 */

package main

import (
	"fmt"
	"unsafe"
)

// 按照对齐原则，Test 变量占用 16 字节的内存
type Test struct {
	F1 uint64
	F2 uint32
	F3 byte
}

func Struct2Bytes(p unsafe.Pointer, n int) []byte {
	return ((*[4096]byte)(p))[:n]
}

func main() {
	t := Test{F1: 0x1234, F2: 0x4567, F3: 12}
	bytes := Struct2Bytes(unsafe.Pointer(&t), 16)
	//  %v:相应值的默认格式。
	fmt.Printf("t -> []byte\t: %v\n", bytes)
}
