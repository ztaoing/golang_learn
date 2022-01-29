package main

import "fmt"

var a []int
var c []int //第三者

func main() {
	f([]int{0, 1, 2, 3, 4, 5})
}

func f(b []int) []int {
	a = b[:2]
	//新的切片append导致切片扩容
	c = append(c, b[:2]...)
	fmt.Printf(" a:%p\n c:%p\n b:%p\n", &a[0], &c[0], &b[0])
	return a
}

/**
输出结果:
a: 0xc000102060
c: 0xc000124010
b: 0xc000102060

这段程序，新增了一个变量 c，他容量为 0。此时将期望的数据，追加过去。自然而然他就会遇到容量空间不足的情况，也就能实现申请新底层数据。

我们再将原本的切片置为 nil，就能成功实现两者分手的目标了。
*/
