/**
* @Author:zhoutao
* @Date:2021/1/2 下午1:14
* @Desc:
 */

package main

import "fmt"

func foo(a [2]int) {
	a[0] = 200

}

func foo1(a *[2]int) {
	(*a)[0] = 200
}

func foo2(a []int) {
	a[0] = 300
}
func main() {
	// array
	a := [2]int{1, 2}
	// slice.md
	a1 := []int{1, 2}
	// foo
	foo(a)
	fmt.Println(a)
	// output : [1 2]

	//foo1
	foo1(&a)
	fmt.Println(a)
	// output : [200 2]

	//foo2
	// 将切片作为参数时，拷贝了一个新切片，即拷贝了构成切片的三个值：*ptr len cap
	// 对切片中某个元素的修改，实际上是修改了底层数组中的值，因此原切片也发生了改变
	foo2(a1)
	fmt.Println(a1)
	// output : [300 2]

	// deto中使用的a是对切片a1的拷贝
	deto(a1)
	// 打印的原切片
	fmt.Println(a1)
	// output : [300 2]  没有变化
}

func deto(a []int) {
	// 新切片增加了8个元素，原切片对应的底层数字不够放置8个元素，因此申请了新的空间来放置扩充后的底层数组。
	//这个时候新切片的底层数组和原切片的底层数组就不是同一个了
	a = append(a, 1, 2, 3, 4, 5, 6, 7, 8)
	a[0] = 600
}
