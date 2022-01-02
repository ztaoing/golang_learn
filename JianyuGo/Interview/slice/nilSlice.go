package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	var s1 []int
	s2 := make([]int, 0)
	s3 := make([]int, 0)

	// s1 s2 s3的地址
	fmt.Printf("s1 pointer:%v,"+
		" s2 pointer:%v,"+
		"s3 pointer:%v",
		*(*reflect.SliceHeader)(unsafe.Pointer(&s1)),
		*(*reflect.SliceHeader)(unsafe.Pointer(&s2)),
		*(*reflect.SliceHeader)(unsafe.Pointer(&s3)))
	//s1 pointer:{0 0 0}的data=0,
	//s2 pointer:{824634355408 0 0}的data=824634355408,
	//s3 pointer:{824634355408 0 0}的data=824634355408
	// s2和s3的地址相同，说明都指向了相同的空地址
	fmt.Printf("%v\n", (*(*reflect.SliceHeader)(unsafe.Pointer(&s1))).Data == (*(*reflect.SliceHeader)(unsafe.Pointer(&s2))).Data) //false
	fmt.Printf("%v\n", (*(*reflect.SliceHeader)(unsafe.Pointer(&s2))).Data == (*(*reflect.SliceHeader)(unsafe.Pointer(&s3))).Data) //true
	// nil切片和空切片的地址不一样，nil切片引用数组指针地址为（没有任何实际地址）
	// 空切片的引用数组指针地址是一个固定值
}
