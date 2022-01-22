package main

import "fmt"

func GetValue(m map[int]string, id int) (string, bool) {
	if _, exist := m[id]; exist {
		return "存在数据", true
	}
	//编译出错:cannot use nil as type string in return argument
	//考点：函数返回值类型 nil 可以用作 interface、function、pointer、map、slice.md 和 channel 的“空值”。
	//但是如果不特别指定的话，Go 语言不能识别类型，所以会报错。
	//通常编译的时候不会报错，但是运行是时候会报:cannot use nil as type string in return argument.
	return nil, false
}
func main() {
	intmap := map[int]string{
		1: "a",
		2: "b",
		3: "ccc",
	}
	v, err := GetValue(intmap, 3)
	fmt.Println(v, err)
}
