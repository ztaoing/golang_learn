[什么是内联函数]
    在很多讲 Go 语言底层的技术资料和博客里都会提到内联函数这个名词，也有人把内联函数说成代码内联、函数展开、展开函数等等，
    其实想表达的都是 "Go 语言编译器对函数调用的优化"，
    "编译器会把一些函数的调用直接替换成被调函数的函数体内的代码在调用处展开"，以减少函数调用带来的时间消耗。
    它是Go语言编译器对代码进行优化的一个常用手段。

    内联函数并不是 Go 语言编译器独有的，很多语言的编译器在编译代码时都会做内联函数优化，维基百科对内联函数的解释如下 (我把重点需要关注的信息特意进行了加粗)：
    
    在计算机科学中，内联函数（有时称作在线函数或编译时期展开函数）是一种编程语言结构，
    用来建议编译器对一些特殊函数进行内联扩展（有时称作在线扩展）；也就是说建议编译器将指定的函数体插入并取代每一处调用该函数的地方（上下文），
    从而节省了每次调用函数带来的额外时间开支。
    但在选择使用内联函数时，必须在程序占用空间和程序执行效率之间进行权衡，因为过多的比较复杂的函数进行内联扩展将带来很大的存储资源开支。
    另外还需要特别注意的是对递归函数的内联扩展可能引起部分编译器的无穷编译。

    Note：内联优化一般用于能够快速执行的函数，因为在这种情况下函数调用的时间消耗显得更为突出，同时内联体量小的函数也不会明显增加编译后的执行文件占用的空间。
    
    举个例子：
    package main
    import "fmt"
    func main() {
        x := 20
        y := 5
        res := add(x, y)
        fmt.Println(res)
    }
    
    func add(x int, y int) int {
        return x + y
    }
    
    上面的函数非常简单，add 函数对两个参数进行加和，编译器在编译上面的 Go 代码时会做内联优化，把 add 函数的函数体直接在调用处展开，等价于上面的 Go 代码是这么编写的。
    package main
    import "fmt"

    func main() {
        x := 20
        y := 5
        // 内联函数， 或者叫把函数展开在调用处
        res := x + y
        fmt.Println(res)
    }
    
    func add(x int, y int) int {
        return x + y
    }

[告诉编译器不对函数进行内联]
在源码编译的汇编代码里我们也看不到对 add 函数的调用，不过我们可以通过在 add 函数上增加一个特殊的注释来告诉编译器不要对它进行内联优化

    // 注意下面这行注释，"//"后面不要有空格
    //go:noinline
    func add(x int, y int) int {
        return x + y
    }
    怎么验证这个注释真实有效，能让编译器不对add函数做内联优化呢？
    我们可以用 go tool compile -S scratch.go 打印出的 Go 代码被编译成的汇编代码，在汇编代码里我们可以发现对add函数的调用。
    0x0053 00083 (scratch.go:6)  CALL    "".add(SB)
    这也反向证明了，正常情况下 Go 编译器会对 add 函数进行内联优化。
    
[编译器会对代码做哪些优化]
除了分析编译后的汇编源码外，我们还可以通过给 go build 命令传递  -gcflags -m 参数
    
        $ go build -gcflags --help
        [.......]
        // 传递 -m 选项会输出编译器对代码的优化
        -m    print optimization decisions
    
    让编译器告诉我们它在编译 Go 代码对代码都做了哪些优化。
    接下用 -gcflags -m 构建一下我们的小程序
    
    $ go build -gcflags -m scratch.go

    ./scratch_16.go:10:6: can inline add
    ./scratch_16.go:6:12: inlining call to add
    ./scratch_16.go:7:13: inlining call to fmt.Println
    ./scratch_16.go:7:13: res escapes to heap
    ./scratch_16.go:7:13: main []interface {} literal does not escape
    ./scratch_16.go:7:13: io.Writer(os.Stdout) escapes to heap

    通过终端的输出可以了解到，编译器判断 add函数可以进行内联，也对 add 函数进行了内联，
    除此之外还对fmt.Println 进行了内联优化。
    还有一个 io.Writer(os.Stdout) escapes to heap 的输出代表的是 io 对象逃逸到了堆上，
    堆逃逸是另外一种优化，在先前 Go内存管理系列的文章 -- Go内存管理之代码的逃逸分析 有详细说过。
    
[哪些函数不会被内联]
    那么 Go 的编译器是不是会对所有的体量小，执行快的函数都会进行内联优化呢？我查查了资料发现 Go 在决策是否要对函数进行内联时有一个标准：
    函数体内包含：闭包调用，select ，for ，defer，go 关键字的的函数不会进行内联。
    并且除了这些，还有其它的限制。当解析AST时，Go申请了80个节点作为内联的预算。
    每个节点都会消耗一个预算。比如，a = a + 1这行代码包含了5个节点：AS, NAME, ADD, NAME, LITERAL。
    以下是对应的SSA 输出：
    当一个函数的开销超过了这个预算，就无法内联。
    
    以上描述翻译自文章：https://medium.com/a-journey-with-go/go-inlining-strategy-limitation-6b6d7fc3b1be