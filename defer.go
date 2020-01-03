package main

import "fmt"

/**
多个defer语句的执行顺序为：(defer,return,返回值三者的执行逻辑)
defer最先执行一些收尾工作--》然后return执行，return负责将结果写入返回值中--》最后函数携带当前返回值退出
*/

func main() {
	fmt.Println(a())

	fmt.Println(b())
}

//a函数的返回值没有被命名
//特别注意：匿名函数的特点是可以继承变量的值，所以defer2语句继承了defer1语句的值。
//按理说a()的返回值是2才对，但是返回值确实0，这是因为返回值没有被生命，所以函数a()的返回值还是0
func a() int {
	var i int

	defer func() {
		i++
		fmt.Println("a_defer2:", i) //a_defer2:2
	}()

	defer func() {
		i++
		fmt.Println("a_defer1:", i) //a_defer1:1
	}()
	return i //返回：0
}

//b()的返回值被声明,因此defer在return复制返回值i之后，再一次改成了i，最终函数退出后的返回值才是defer修改过的值
func b() (i int) {
	defer func() {
		i++
		fmt.Println("b_defer2:", i) //b_defer2:2
	}()

	defer func() {
		i++
		fmt.Println("b_defer1:", i) //b_defer1:1

	}()
	return i //返回：2
}

/**
defer 原理
*/
