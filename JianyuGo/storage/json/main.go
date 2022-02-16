/**
* @Author:zhoutao
* @Date:2022/2/16 10:50
* @Desc:
 */

package main

import (
	"encoding/json"
	"fmt"
)

// 按照对齐原则，Test 变量占用 16 字节的内存
type Test struct {
	F1 uint64
	F2 uint32
	F3 byte
}

func main() {
	t := Test{F1: 0x1234, F2: 0x4567, F3: 12}

	//t := Test{F1: 4660, F2: 17767, F3: 12} 同上

	// 测试序列化
	bs, err := json.Marshal(&t)
	if err != nil {
		panic("")
	}
	fmt.Printf("t -> []byte\t: %v\n", bs)

	// 测试反序列化
	t1 := Test{}
	err = json.Unmarshal(bs, &t1)
	if err != nil {
		panic("")
	}
	fmt.Printf("[]byte -> t1\t: %v\n", t1)
}
