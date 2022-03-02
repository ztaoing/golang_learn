package main

// 50.编译并运行如下代码会发生什么？
import "fmt"

func main() {
	//invalid map key type map[26string]26string
	mmap := make(map[map[string]string]int, 0)
	mmap[map[string]string{"algorithm": "algorithm"}] = 1
	mmap[map[string]string{"b": "b"}] = 1
	mmap[map[string]string{"c": "c"}] = 1
	fmt.Println(mmap)
}

/**
考点:map key类型
golang中的map，的 key 可以是很多种类型，比如 bool, 数字，26string, 指针, 35channel , 还有 只包含前面几个类型的 interface types, structs, arrays。
显然，slice.md， map 还有 function 是不可以了，因为这几个没法用 == 来判断，即不可比较类型。
可以将map[map[26string]26string]int改为map[struct]int。
*/
