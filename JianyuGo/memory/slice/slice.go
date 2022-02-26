package main

import "fmt"

var a []int
var c []int //第三者

func main() {
	//f([]int{0, 1, 2, 3, 4, 5})
	out()
}

func f(b []int) []int {
	a = b[:2]
	//新的切片append导致切片扩容
	c = append(c, b[:2]...)
	fmt.Printf(" logic:%p\n c:%p\n b:%p\n", &a[0], &c[0], &b[0])
	return a
}

/**
输出结果:
logic: 0xc000102060
c: 0xc000124010
b: 0xc000102060

这段程序，新增了一个变量 c，他容量为 0。此时将期望的数据，追加过去。自然而然他就会遇到容量空间不足的情况，也就能实现申请新底层数据。

我们再将原本的切片置为 nil，就能成功实现两者分手的目标了。
*/

func out() {

	sl := make([]int, 0, 10)
	var appenFunc = func(s []int) {
		s = append(s, 10, 20, 30)
		fmt.Println(s)
	}
	fmt.Println(sl)
	appenFunc(sl)
	fmt.Println(sl)
	fmt.Println(sl[:10])

	/**
	你认为程序的输出结果是什么？
	[]
	[10 20 30]
	[]
	[10 20 30 0 0 0 0 0 0 0]

	*/
	/**

	fmt.Println(sl)
	fmt.Println(sl[:10])
	上述代码中，为什么第一个 sl 打印结果是空的，第二个 sl 给索引位置就能打印出来？

		也有小伙伴不断在尝试 sl[:10] 以外的输出，有没有因为一些边界值改变而导致不行。
		fmt.Println(sl[:])你认为这个对应的输出结果是什么？
		[10 20 30 0 0 0 0 0 0 0] 错了
		正确的结果是：[]
		是没有任何元素输出，这下大家更懵了。为什么 sl[:] 的输出结果为空？
	*/
	/**
	再看看变量 sl 的长度和容量：
	fmt.Println(len(sl), cap(sl))
	0 10

	*/

	/**
			[挖掘原因]三个问题
			请思考如下三个问题：

			1、为什么打印 sl[:10] 时，结果包含了 10 个元素，还包含了函数闭包中插入的 10, 20, 30，之间有什么关系？
			2、为什么打印 sl 变量时，结果为空？
			3、为什么打印 sl[:] 时，结果为空。但打印 sl[:10] 就正常输出？

			【了解底层】
			要分析起源，我们就必须要再提到 slice（切片）的底层实现，slice 底层存储的数据结构指向了一个 array（数组）。

			type SliceHeader struct {
			 Data uintptr 指向具体的底层数组
			 Len  int 切片的长度
			 Cap  int 切片的容量
			}
			核心要记住的是：slice 真正存储数据的地方，是一个数组。
			slice 的结构中存储的是指向所引用的数组指针地址。

			【分析原因】
			我们关注到 appenFunc 变量，他其实是一个函数，并且结果中我们所看到的 10, 20, 30，
			也只有这里有插入的动作。因此这是需要分析的。

			func main() {
				 sl := make([]int, 0, 10)
				 var appenFunc = func(s []int) {
				  s = append(s, 10, 20, 30)
				 }

				 appenFunc(sl)
				 fmt.Println(sl)
				 fmt.Println(sl[:10])
			}
			但为什么在 appenFunc 函数中所插入的 10, 20, 30 元素，就跑到外面的切片 sl 中去了呢？
			这其实结合 slice 的底层设计和函数传递就明白了，
			在 Go 语言中，只有值传递：也就是说一个函数总是得到传参的副本，

			实质上在调用 appenFunc(sl) 函数时，实际上修改了底层所指向的数组（虽然传递的是副本，但是副本中存储值的数组地址是不变的，
			所以在操作这个副本的存储数据的地址的时候，就是操作了原有的数组），
			自然也就会发生变化，也就不难理解为什么 10, 20, 30 元素会出现了。

			那为什么 sl 变量的长度是 0，甚至有人猜测是不是扩容了，这其实和上面的问题还是一样，因为是值传递，自然也就不会发生变化。
			（此时的sl成了底层数组的一个切片,原有的底层容量的内容增加了，但是sl这个切片的管理范围并没有变动）
			要记住一个关键点：如果传过去的值是指向内存空间的地址，是可以对这块内存空间做修改的。反之，你也改不了。
			至此，也就解决了我们的第一个大问题。

		【切片小优化】
		还剩下两个大问题，这似乎用上面的结论没法完整解释。虽说程序是诱因，但这块最直接的影响是和切片访问的小优化有关。
		常用的访问切片我们会用：s[low : high]
		注意这里是：low、high。可没有用 len、cap 这种定性的词语，也就代表着这里取的值是可变的。

		当是切片（slice）时，表达式 s[low : high] 中的 high，最大的取值范围对应着切片的容量（cap），不是单纯的长度（len）
		因此调用 fmt.Println(sl[:10]) 时可以输出容量范围内的值，不会出现越界。

		相对的 fmt.Println(sl) 因为该切片 len 值为 0，没有指定最大索引值，high 则取 len 值，导致输出结果为空。
		注：访问元素在 Go 编译期就确定的了，相关逻辑可以在 compile 相关的代码中看到。

	【总结】
	在今天这篇文章中，我们结合了 Go 语言中切片的
	基本底层原理、
	值传递、
	边界值取值
	等进行了多轮探讨。

	所谓的最大取值范围，除非官方给你写定 len 或 cap，否则不要过于主观的认为，因为他会根据访问的数据类型和访问定位等改变。

	*/
}
