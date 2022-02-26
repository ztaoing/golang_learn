package main

import "fmt"

type Course struct {
	name  string
	price int
	url   string
}

// 接口一般以er结尾
type printer interface {
	printInfo() string
}

func (c Course) printInfo() string {
	return "logic string"
}

// 2、用途：可以用来传递参数
/*func print(x interface{}) {
	// 根据不同的类型做不同的打印

	//判断x是不是int类型
	// 如果是int类型的话，ok就是true否则就是false。如果是bool值，则命名为OK
	// 还有另外的一种返回情况时error，如果是error的话，就定义成error。

	// 此时的ok就变为了局部变量
	if v, ok := x.(int); ok {
		// x是int类型
		fmt.Printf("%d", "(整数)\n", v)
	}

	if s, ok := x.(string); ok {
		fmt.Printf("%s", "(字符串)\n", s)
	}

	// 这里使用switch更合适

}*/

func print(x interface{}) {
	switch v := x.(type) {
	case string:
		fmt.Printf("%s(字符串)\n", v)
	case int:
		fmt.Printf("%d(整数)\n", v)
	}
}
func main() {
	// 空接口
	var i interface{}
	// 空接口类似于Java和python中的object

	// 1、空接口的第一个用途：可以把任何类型都赋值给空接口变量

	// 把struct赋值给空接口
	/*i = Course{}
	print(i)*/
	i = 10
	print(i)
	i = "boddy"
	print(i)

	// 2、参数传递
	// print()

	// 3、空接口可以作为map的值
	var teacherInfo = make(map[string]interface{})
	teacherInfo["name"] = "boodu"
	teacherInfo["age"] = 18
	teacherInfo["weight"] = 76.2
	teacherInfo["courses"] = []string{"django", "sanic"}

	fmt.Printf("%v", teacherInfo)

	// Course实现了printer接口
	// 如果实现的是指针类型，
	// 如果实现的是值类型,
	// var c printer = Course{}

	// 断言
	// i.(int) 意思是看下i是不是int类型

}
