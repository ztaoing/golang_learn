package main

import "fmt"

type Gopher interface {
	WriteGoCode()
}

type person struct {
	name string
}

func (p person) WriteGoCode() {
	fmt.Printf("i am %s,i am writing go code!\n", p.name)
}

func Coding(g Gopher) {
	g.WriteGoCode()
}

func Coding1(gs []Gopher) {
	for _, g := range gs {
		g.WriteGoCode()
	}
}

func main() {
	p := person{
		name: "小菜刀",
	}
	Coding(p)

	p1 := []person{
		{name: "小菜刀1号"},
		{name: "小菜刀2号"},
	}
	//编译不通过:cannot use p1 (type []person) as type []Gopher in argument to Coding1
	// 明明person类型实现了Gopher接口，且当函数入参为Gopher类型时，能够顺利被执行，但是参数变为[]Gopher就编译不了，这是为什么？
	Coding1(p1)
	/**
	这个问题在Stack Overflow上被热议：
	在go中有一个通用的规则，即语法不应该隐藏负载/昂贵的操作。
	转换一个string到interface{}的时间复杂度是O（1），
	转换[]string到interface{}同样也是O(1),因为他还是一个单一值的转换。
	如果要将[]string转换为[]interface{},他是O(N)，因为切片的每一个元素都必须转换为interface{},这违背了go的语法原则

	这个回答，你们同意吗？
	当然，此规则存在一个例外：转换字符串。
	在将身体乳能够转换为[]byte或[]rune时，即使需要O(n)，但是go会运行他的执行,不会报错！

	interfaceSlice问题：
	Ian Lance Taylor（Go 核心开发者）在go官方仓库中也回答了这个问题：
	https://github.com/golang/go/wiki/InterfaceSlice
	他给出了这样做的两个原因：
	原因一：类型为 []interface{} 的变量不是 interface！它仅仅是一个元素类型恰好为 interface{} 的切片。
	原因二：[]interface{} 变量有特定大小的内存布局，在编译期可知。这与 []MyType 是不同的。

	每个 interface{} （运行时通过 runtime.eface 表示）占两个字长（一个字代表所包含内容的类型 _type，另外一个字表示所包含的数据 data 或者指向它的指针 ）
	因此，类型为 []interface{} 的长度为 N 的变量，它是由 N*2 个字长（一个字代表所包含内容的类型 _type，另外一个字表示所包含的数据 data 或者指向它的指针）的数据块支持。
	而这与类型为 []MyType 的长度为 N 的变量的数据块大小是不同的，因为后者的数据块是 N*sizeof(MyType) 字长(只需要考虑MyType所占的字节大小)。

	重要：数据块的不同，造成的结果是编译器无法快速地将 []MyType 类型的内容分配给 []interface{} 类型的内容。
	同理，[]Gopher 变量也是特定大小的内存布局（运行时通过 runtime.iface 表示）。这同样不能快速地将 []MyType 类型的内容分配给 []Gopher 类型。

	因此，Ian Lance Taylor 回答闭环了 Go 的语法通用规则：Go 语法不应隐藏复杂/昂贵的操作，编译器会拒绝它们
	*/

	/**
	代码解决方案：
	再次将文章开头的例子附上，如果我们需要 [] person 类型的 p 能够成功入参 Coding() 函数，应该如何做呢。
	代码方案如下，核心是需要一个 []Gopher 类型的转换变量。
	*/
	var interfaceSlice []Gopher = make([]Gopher, len(p1))
	// 将p1放到interfaceSlice中,即将p1中的多条内容放到类型为Gopher的slice中
	for i, g := range p1 {
		interfaceSlice[i] = g
	}
	Coding1(interfaceSlice)
}
