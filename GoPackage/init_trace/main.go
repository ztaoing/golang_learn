/**
* @Author:zhoutao
* @Date:2023/1/3 14:42
* @Desc:
 */

package main

import (
	"fmt"
	"golang_learn/golang_learn/GoPackage/init_trace/pkg1"
	"golang_learn/golang_learn/GoPackage/init_trace/pkg2"
	"golang_learn/golang_learn/GoPackage/init_trace/tracer"
)

var M_v1 = tracer.Trace("init main v1", pkg1.P1_v1+10)

// 引入pkg2中的变量
var M_v2 = tracer.Trace("init main v2", pkg2.P2_v1+10)

func init() {
	fmt.Println("main init 1")
}

func init() {
	fmt.Println("main init 2")
}
func main() {
	fmt.Println("main print")
	fmt.Println(M_v1, M_v2)
	/**
	先初始化p2：
	init p2_v1 : 21
	init_p2_v2 : 22
	 init func in pkg2

	再初始化p1：
	p1_v1 : 31
	p1_v2 : 32
	init func in pkg1

	多次引用 只会初始化一次：
	init main v1 : 41
	init main v2 : 31

	main包中的初始化顺序：
	main init 1
	main init 2

	main print

	*/
}
