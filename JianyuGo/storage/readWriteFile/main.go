/**
* @Author:zhoutao
* @Date:2022/2/16 14:14
* @Desc:
 */

package main

import (
	"fmt"
	"log"
	"os"
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
	// 强转类型
	bytes := Struct2Bytes(unsafe.Pointer(&t), 16)

	fmt.Printf("t -> []byte\t: %v\n", bytes)

	fd, err := os.OpenFile("test_bytes.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("create failed, err:%v\n", err)
	}

	// 结构体写入文件
	_, err = fd.Write(bytes)
	if err != nil {
		log.Fatalf("write failed, err:%v\n", err)
	}

	t1 := Test{}
	// 强转出一个 16 字节的内存 buffer 来
	t1Bytes := Struct2Bytes(unsafe.Pointer(&t1), 16)
	// 从文件中把数据读出来
	_, err = fd.ReadAt(t1Bytes, 0)
	if err != nil {
		log.Fatalf("read failed, err:%v\n", err)
	}

	fmt.Printf("t1 -> []byte\t: %v\n", t1Bytes)
}
