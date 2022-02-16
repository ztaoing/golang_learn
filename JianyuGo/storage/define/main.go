/**
* @Author:zhoutao
* @Date:2022/2/16 10:55
* @Desc:
 */

package main

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// 按照对齐原则，Test 变量占用 16 字节的内存
type Test struct {
	F1 uint64
	F2 uint32
	F3 byte
}

func (t *Test) Marshal() ([]byte, error) {
	// 创建一个 16 字节的 buffer
	buf := make([]byte, 16)
	// 序列化
	binary.BigEndian.PutUint64(buf[0:8], t.F1)  //0-7
	binary.BigEndian.PutUint32(buf[8:12], t.F2) //8-11
	buf[12] = t.F3                              //12

	return buf, nil
}

func (t *Test) Unmarshal(buf []byte) error {
	if len(buf) != 16 {
		return errors.New("length not match")
	}
	// 反序列化
	t.F1 = binary.BigEndian.Uint64(buf[0:8])
	t.F2 = binary.BigEndian.Uint32(buf[8:12])
	t.F3 = buf[12]
	return nil
}

func main() {
	t := Test{F1: 0x1234, F2: 0x4567, F3: 12}

	// 测试序列化
	bs, err := t.Marshal()
	if err != nil {
		panic("")
	}
	fmt.Printf("t -> []byte\t: %v\n", bs)

	// 测试反序列化
	t1 := Test{}
	err = t1.Unmarshal(bs)
	if err != nil {
		panic("")
	}
	fmt.Printf("[]byte -> t1\t: %v\n", t1)
}
