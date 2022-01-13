【以下内容来自煎鱼的微信公众号】
* 值为nil能调用函数吗？
    func(p *sometype)Somemethod(a int){} 本质上是func Somemethod(p *sometype,a int){}
    所以参数为nil，不影响方法的调用
  
* go有哪几种无法恢复的致命场景
    

* 动手实现一个localcache：高效的并发访问；减少GC
    1、高效并发访问：【减小锁的粒度】
                本地缓存的本地实现可以使用map[string]interface{}+sync.RWMutex组合
                使用sync.RWMutex对读进行了优化，但是当并发量上来以后，哈市编程了串行读，等待锁的goroutine
                就会被阻塞住，为了解决这个问题我们可以进行分片。
                每一个分片使用一把锁，减少竞争：根据他的key做hash(key),然后进行分片：hash(key)%N；
               
                分片数量的选择，分片并不是越多越好，根据经验，我们的分片数可以选择N的2次幂，
                分片时为了提高效率可以使用位运算代替取余操作。
   2、 减少GC：
                BigCache如何加速并发访问以及避免高额的GC开销： https://pengrl.com/p/35302/
    

[runtime]
* 什么是go runtime.KeepAlive
go 官方文档: https://pkg.go.dev/runtime#KeepAlive
文档: https://medium.com/a-journey-with-go/go-keeping-a-variable-alive-c28e3633673a
  


[类型的比较] golang.org/ref/spec#comparison_operators
* 可比较类型和不可不叫类型。对于不可比较类型，如何比较他们包含的值是否相等呢？使用reflect.DeepEqual


【json.unmarshal】pkg.go.dev/encoding/json#unmarshal
json.unmarshal的类型转换
bool, for JSON booleans
float64, for JSON numbers
string, for JSON strings
[]interface{}, for JSON arrays
map[string]interface{}, for JSON objects
nil for JSON null

* for range :是获取切片的长度，然后执行n次

【编译】总结两个go程序编译的重要知识
*交叉编译，条件编译

[go程序自己监控自己]

[多路复用] channel.multiplex
[atomic.Value为什么不加锁也能保证数据线程安全]

[go中的零值，它有什么作用？] 官方：https://golang.org/ref/spec#the_zero_value
布尔型为false；数字型为0；字符串型为""；指针、函数、接口、切片、通道和映射都为nil

[go是如何实现启动参数的加载的？]
go汇编为了简化汇编代码的编写，引入了PC\FP\SP\SB色哥伪寄存器。
四个伪寄存器加上其他的通用寄存器就是go汇编语言对CPU的重新抽象。该抽象的结构也适用于非x86类型的体系结构

go可以利用os.Args解析程序启动时的命令行参数，他的实现过程是怎样的？
    func main(){
        for i,v:=range os.Args{
            fmt.Printf("arg[%d]:%v\n,i,v)
        }
    }
    输出：
        $ go build main.go
        $ ./main foo bar sss ddd
        arg[0]: ./main
        arg[1]: foo
        arg[2]: bar
        arg[3]: sss
        arg[4]: ddd

[select机制] 
每个case 上的操作例如方法的结果给channel， 每次循环，所有的方法都会执行，但是只会选择其中一个case，其他的case的操作就被丢失了
    
这个问题很常见：最多的就是time.After导致内存泄漏问题，网上有很多的文章解释原因，如何避免，其实最根本原因是select这个机制导致的

以下代码会内存泄漏：
    func main(){
        ch:= make(chan int ,10)
        go func(){
            var i = 1
            for {
                i++
                ch<-i
            }
        }()
        for {
            select{
                case x:=<-ch:
                    println(x)
                case <-time.After(30*time.Second):
                    println(time.now().Unix())
            }
        }

    }
    为什么会内存泄漏？
    答： 每次循环都会执行time.After(30*time.Second)，导致堆内存不断升高，最后泄漏。哪怕没有选择这个case
        time.After(30*time.Second)也会执行。

[如何保存go程序崩溃的现场] 
错误日志能力有限：
    1、日志是开发者在代码中定义的打印信息，我们没法保证日志信息能包含所有的错误情况
    2、在go程序发生panic时，我们也并不能总是能通过recover捕获（没法插入日志代码）
那线上go程序突然崩溃后，当日志记录没有覆盖到错误场景时，还有别的方法排查码？
    
    core dump 又即核心转储，简单来说它就是程序意外终止时产生的内存快照。我们可以通过 core dump 文件来调式程序，找出其崩溃原因。
    在 linux 平台上，可通过ulimit -c命令查看核心转储配置，系统默认为 0，表明未开启 core dump 记录功能。
    
    $ ulimit -c
    输出：0
    
    可以使用ulimit -c [size]命令指定记录 core dump 文件的大小，即是开启 core dump 记录。当然，如果电脑资源足够，避免 core dump 丢失或记录不全，也可执行ulimit -c unlimited而不限制 core dump 文件大小。
    那在 Go 程序中，如何开启 core dump 呢？

    string 是不可以被修改的，当我们将 string 类型通过黑魔法转为 []byte 后，企图修改其值，程序会发生一个不能被 recover 捕获到的错误。
    
    $ go run main.go
    unexpected fault address 0x106a6a4
    fatal error: fault
    [signal SIGBUS: bus error code=0x2 addr=0x106a6a4 pc=0x105b01a]
    
    goroutine 1 [running]:
    runtime.throw({0x106a68b, 0x0})
    /usr/local/go/src/runtime/panic.go:1198 +0x71 fp=0xc000092ee8 sp=0xc000092eb8 pc=0x102bad1
    runtime.sigpanic()
    /usr/local/go/src/runtime/signal_unix.go:732 +0x1d6 fp=0xc000092f38 sp=0xc000092ee8 pc=0x103f2f6
    main.Modify(...)
    /Users/slp/github/PostDemo/coreDemo/main.go:21
    main.main()
    /Users/slp/github/PostDemo/coreDemo/main.go:25 +0x5a fp=0xc000092f80 sp=0xc000092f38 pc=0x105b01a
    runtime.main()
    /usr/local/go/src/runtime/proc.go:255 +0x227 fp=0xc000092fe0 sp=0xc000092f80 pc=0x102e167
    runtime.goexit()
    /usr/local/go/src/runtime/asm_amd64.s:1581 +0x1 fp=0xc000092fe8 sp=0xc000092fe0 pc=0x1052dc1
    exit status 2
    
    这些堆栈信息是由 GOTRACEBACK 变量来控制打印粒度的，它有五种级别。
        1、none，不显示任何 goroutine 堆栈信息
        2、single，默认级别，显示当前 goroutine 堆栈信息
        3、all，显示所有 user （不包括 runtime）创建的 goroutine 堆栈信息
        4、system，显示所有 user + runtime 创建的 goroutine 堆栈信息
        5、crash，和 system 打印一致，但会生成 core dump 文件（Unix 系统上，崩溃会引发 SIGABRT 以触发core dump）
    注：如果我们将 GOTRACEBACK 设置为 system ，我们将看到程序崩溃时所有 goroutine 状态信息
        
    $ GOTRACEBACK=system go run main.go
        unexpected fault address 0x106a6a4
        fatal error: fault
        [signal SIGBUS: bus error code=0x2 addr=0x106a6a4 pc=0x105b01a]
        
        goroutine 1 [running]:
        runtime.throw({0x106a68b, 0x0})
        ...
        
        goroutine 2 [force gc (idle)]:
        runtime.gopark(0x0, 0x0, 0x0, 0x0, 0x0)
        ...
        created by runtime.init.7
        /usr/local/go/src/runtime/proc.go:294 +0x25
        
        goroutine 3 [GC sweep wait]:
        runtime.gopark(0x0, 0x0, 0x0, 0x0, 0x0)
        ...
        created by runtime.gcenable
        /usr/local/go/src/runtime/mgc.go:181 +0x55
        
        goroutine 4 [GC scavenge wait]:
        runtime.gopark(0x0, 0x0, 0x0, 0x0, 0x0)
        ...
        created by runtime.gcenable
        /usr/local/go/src/runtime/mgc.go:182 +0x65
        exit status 2

    如果想获取 core dump 文件，那么就应该把 GOTRACEBACK 的值设置为 crash 。
    当然，我们还可以通过 runtime/debug 包中的 SetTraceback 方法来设置堆栈打印级别。

    delve 调试：
    delve 是 Go 语言编写的 Go 程序调试器，我们可以通过 dlv core 命令来调试 core dump。
    首先，通过以下命令安装 delve：go get -u github.com/go-delve/delve/cmd/dlv
    
    还是以上文中的例子为例，我们通过设置 GOTRACEBACK 为 crash 级别来获取 core dump 文件
    
    $ tree
    .
    └── main.go
    $ ulimit -c unlimited
    $ go build main.go
    $ GOTRACEBACK=crash ./main
    ...
    Aborted (core dumped)
    $ tree
    .
    ├── core  //core dump 文件
    ├── main
    └── main.go
    $ ls -alh core
    -rw------- 1 slp slp 41M Oct 31 22:15 core
    
    此时，在同级目录得到了 core dump 文件 core（文件名、存储路径、是否加上进程号都可以配置修改）。
      1、  通过 dlv 调试器来调试 core 文件，执行命令格式 dlv core 可执行文件名 core文件:
        $ dlv core main core
            Type 'help' for list of commands.
            (dlv)
        
      2、  命令 goroutines 获取所有 goroutine 相关信息:
        (dlv) goroutines
         *Goroutine 1 - User: ./main.go:21 main.main (0x45b81a) (thread 18061)
          Goroutine 2 - User: /usr/local/go/src/runtime/proc.go:367 runtime.gopark (0x42ed96) [force gc (idle)]
          Goroutine 3 - User: /usr/local/go/src/runtime/proc.go:367 runtime.gopark (0x42ed96) [GC sweep wait]
          Goroutine 4 - User: /usr/local/go/src/runtime/proc.go:367 runtime.gopark (0x42ed96) [GC scavenge wait]
          [4 goroutines]
          (dlv)

      3、  Goroutine 1 是出问题的 goroutine （带有 * 代表当前帧），通过命令 goroutine 1 切换到其栈帧:
        (dlv) goroutine 1
        Switched from 1 to 1 (thread 18061)
        (dlv)
    
      4、 执行命令 bt（breakpoints trace） 查看当前的栈帧详细信息:
        (dlv) bt
        0  0x0000000000454bc1 in runtime.raise
        at /usr/local/go/src/runtime/sys_linux_amd64.s:165
        1  0x0000000000452f60 in runtime.systemstack_switch
        at /usr/local/go/src/runtime/asm_amd64.s:350
        2  0x000000000042c530 in runtime.fatalthrow
        at /usr/local/go/src/runtime/panic.go:1250
        3  0x000000000042c2f1 in runtime.throw
        at /usr/local/go/src/runtime/panic.go:1198
        4  0x000000000043fa76 in runtime.sigpanic
        at /usr/local/go/src/runtime/signal_unix.go:742
        5  0x000000000045b81a in main.Modify
        at ./main.go:21
        6  0x000000000045b81a in main.main
        at ./main.go:25
        7  0x000000000042e9c7 in runtime.main
        at /usr/local/go/src/runtime/proc.go:255
        8  0x0000000000453361 in runtime.goexit
        at /usr/local/go/src/runtime/asm_amd64.s:1581
        (dlv)
        
      5、  通过 5 0x000000000045b81a in main.Modify 发现了错误代码所在函数，执行命令 frame 5 进入函数具体代码:
        (dlv) frame 5
            > runtime.raise() /usr/local/go/src/runtime/sys_linux_amd64.s:165 (PC: 0x454bc1)
            Warning: debugging optimized function
            Frame 5: ./main.go:21 (PC: 45b81a)
            16: }
            17:
            18: func Modify() {
            19:  a := "hello"
            20:  b := String2Bytes(a)
            =>  21:  b[0] = 'H'
            22: }
            23:
            24: func main() {
            25:  Modify()
            26: }
            (dlv)
        自此，破案了，问题就出在了擅自修改 string 底层值。
  
    注：有一点需要注意，上文 core dump 生成的例子，我是在 linux 系统下完成的，mac amd64 系统没法弄（很气，害我折腾了两个晚上）。
    这是由于 mac 系统下的 Go 限制了生成 core dump 文件，这个在 Go 源码 src/runtime/signal_unix.go 中有相关说明。
    
    //go:nosplit
    func crash() {
    // OS X core dumps are linear dumps of the mapped memory,
    // from the first virtual byte to the last, with zeros in the gaps.
    // Because of the way we arrange the address space on 64-bit systems,
    // this means the OS X core file will be >128 GB and even on a zippy
    // workstation can take OS X well over an hour to write (uninterruptible).
    // Save users from making that mistake.
    if GOOS == "darwin" && GOARCH == "amd64" {
    return
    }
    
    dieFromSignal(_SIGABRT)
    }
    
    当然，core dump 文件的生成也是有弊端的。core dump 文件较大，如果线上服务本身内存占用就很高，
    那在生成 core dump 文件上的内存与时间开销都会很大。另外，我们往往会布置服务守护进程，
    如果我们的程序频繁崩溃和重启，那会生成大量的 core dump 文件（设定了core+pid 命名规则），产生磁盘打满的风险（如果放开了内核限制 ulimit -c unlimited）。
    最后，如果担心错误日志不能帮助我们定位 Go 代码问题，我们可以为它开启 core dump 功能，在 hotfix 上增加奇兵。
    对于有守护进程的服务，建议设置好 ulimt -c 大小限制。
    

[在go容器里设置gomaxprocs的正确姿势：]
gomaxprocs是go提供的非常重要的一个环境变量。通过设置gomaxprocs，用户可以调整调度器中processor（即P）
的数量，由于每个系统线程必须要绑定P，P才能把G交给M执行。
    所以，p的数量会很大程度上影响go runtime的并发表现。
gomaxprocs在go1.5之后默认值是机器的CPU核数（runtime.NumCPU）。通过下面的代码可以获取当前机器的核心数和给
gomaxprocs设置值。
    
    func getGOMAXPROCS()int{
       _:= runtime.NumCPU() //获取机器的CPU核心数
        return runtime.GOMAXPROCS(0) //参数为0时用于获取给gomaxprocs设置的值
    }
    func main(){
        fmt.Printf("GOMAXPROCS:%d\n",getGOMAXPROCS())
    }
    
    但是，以docker为代表的容器虚拟化技术，会通过cgroup等技术对CPU资源进行隔离。
    以k8s为代表的基于容器虚拟化实现的资源管理系统，也支持这样的特性，比如在podTemplate的容器定义里：
    limits.cpu=1000m就嗲表给这个容器分配1个核心的使用时间。

    这类隔离技术，导致runtime.NumCPU()无法正确的获取到容器被分配CPU资源数。
    runtime.NumCPU()获取的是物理机的实际核心数。

    设置gomaxprocs高于真正可以使用的核心数会导致go调度器不停的进行os线程切换，
    从而给调度器增加很多不必要的工作。
    
    目前go官方没有好的方式来规避容器里获取到实际使用的核心数，而Uber的uber-go/automaxprocs这个包，
    可以在运行时根据cgroup为容器分配的CPU资源限制来修改稿gomaxprocs

    import _ "go.uber.org/automaxprocs"
    func main(){
        //逻辑
    }

[unsafe包]
go标准库中大量使用了unsafe.pointer,想要更好的理解源码实现，就要知道unsafe.pointer到底是什么？

1、什么是unsafe
众所周知，Go语言被设计成一门强类型的静态语言，那么他的类型就不能改变了，静态也是意味着类型检查在运行前就做了。
所以在Go语言中是不允许两个指针类型进行转换的，使用过C语言的朋友应该知道这在C语言中是可以实现的，
Go中不允许这么使用是处于安全考虑，毕竟强制转型会引起各种各样的麻烦，有时这些麻烦很容易被察觉，
有时他们却又隐藏极深，难以察觉。大多数读者可能不明白为什么类型转换是不安全的，这里用C语言举一个简单的例子：

        int main(){
            double pi = 3.1415926;
            double *pv = &pi;
            void *temp = pd;
            int *p = temp;
        }
    在标准C语言中，任何非void类型的指针都可以和void类型的指针相互指派，也可以通过void类型指针作为中介，
    实现不同类型的指针间接相互转换。上面示例中，指针pv指向的空间本是一个双精度数据，占8个字节，
    但是经过转换后，p指向的是一个4字节的int类型。这种发生内存截断的设计缺陷会在转换后进行内存访问是存在安全隐患。
    我想这就是Go语言被设计成强类型语言的原因之一吧。
    
    虽然类型转换是不安全的，但是在一些特殊场景下，使用了它，可以打破Go的类型和内存安全机制，
    可以绕过低效的类型系统，提高运行效率。所以Go标准库中提供了一个unsafe包，之所以叫这个名字，
    就是不推荐大家使用，但是不是不能用，如果你掌握的特别好，还是可以实践的。

    unsafe实现原理：unsafe包中只提供了3种方法：
    func Sizeof(x ArbitraryType) uintptr
    func Offsetof(x ArbitraryType) uintptr
    func Alignof(x ArbitraryType) uintptr

    Sizeof(x ArbitrayType)方法主要作用是用返回类型x所占据的字节数，但并不包含x所指向的内容的大小，
    与C语言标准库中的Sizeof()方法功能一样，比如在32位机器上，一个指针返回大小就是4字节。

    Offsetof(x ArbitraryType)方法主要作用是返回结构体成员在内存中的位置离结构体起始处(结构体的第一个字段的偏移量都是0)的字节数，
    即偏移量，我们在注释中看一看到其入参必须是一个结构体，其返回值是一个常量。
    
    Alignof(x ArbitratyType)的主要作用是返回一个类型的对齐值，也可以叫做对齐系数或者对齐倍数。
    对齐值是一个和内存对齐有关的值，合理的内存对齐可以提高内存读写的性能。一般对齐值是2^n，最大不会超过8(受内存对齐影响).
    获取对齐值还可以使用反射包的函数，也就是说：unsafe.Alignof(x)等价于reflect.TypeOf(x).Align()。
    对于任意类型的变量x，unsafe.Alignof(x)至少为1。对于struct结构体类型的变量x，计算x每一个字段f的unsafe.Alignof(x，f)，
    unsafe.Alignof(x)等于其中的最大值。对于array数组类型的变量x，unsafe.Alignof(x)等于构成数组的元素类型的对齐倍数。
    没有任何字段的空struct{}和没有任何元素的array占据的内存空间大小为0，不同大小为0的变量可能指向同一块地址。

    细心的朋友会发发现这三个方法返回的都是uintptr类型，这个目的就是可以和unsafe.poniter类型相互转换，
    因为*T是不能计算偏移量的，也不能进行计算，但是uintptr是可以的，所以可以使用uintptr类型进行计算，
    这样就可以可以访问特定的内存了，达到对不同的内存读写的目的。三个方法的入参都是ArbitraryType类型，
    代表着任意类型的意思，同时还提供了一个Pointer指针类型，即像void *一样的通用型指针。

    type ArbitraryType int
    type Pointer *ArbitraryType
    // uintptr 是一个整数类型，它足够大，可以存储
    type uintptr uintptr

    *T：普通类型指针类型，用于传递对象地址，不能进行指针运算。
    unsafe.poniter：通用指针类型，用于转换不同类型的指针，不能进行指针运算，不能读取内存存储的值(需转换到某一类型的普通指针)
    uintptr：用于指针运算，GC不把uintptr当指针，uintptr无法持有对象。uintptr类型的目标会被回收。

    三者关系就是：unsafe.Pointer是桥梁，可以让任意类型的指针实现相互转换，
    也可以将任意类型的指针转换为uintptr进行指针运算，也就说uintptr是用来与unsafe.Pointer打配合，用于指针运算。画个图表示一下：
    
    unsafe.Pointer基本使用:
    我们在上一篇分析atomic.Value源码时，看到atomic/value.go中定义了一个ifaceWords结构，
    其中typ和data字段类型就是unsafe.Poniter，这里使用unsafe.Poniter类型的原因是传入的值就是interface{}类型，
    使用unsafe.Pointer强转成ifaceWords类型，这样可以把类型和值都保存了下来，方便后面的写入类型检查。
    截取部分代码如下：

    // ifaceWords is interface{} internal representation.
        type ifaceWords struct {
            typ  unsafe.Pointer
            data unsafe.Pointer
        }
        // Load returns the value set by the most recent Store.
        // It returns nil if there has been no call to Store for this Value.
        func (v *Value) Load() (x interface{}) {
        vp := (*ifaceWords)(unsafe.Pointer(v))
        for {
            typ := LoadPointer(&vp.typ) // 读取已经存在值的类型
            /**
            ..... 中间省略
            **/
            // First store completed. Check type and overwrite data.
            if typ != xp.typ { //当前类型与要存入的类型做对比
                panic("sync/atomic: store of inconsistently typed value into Value")
            }
        }
    
        上面就是源码中使用unsafe.Pointer的一个例子，有一天当你准备读源码时，unsafe.pointer的使用到处可见。
        好啦，接下来我们写一个简单的例子，看看unsafe.Pointer是如何使用的。
        
        func main()  {
            number := 5
            pointer := &number
            fmt.Printf("number:addr:%p, value:%d\n",pointer,*pointer)
            
            float32Number := (*float32)(unsafe.Pointer(pointer))
            *float32Number = *float32Number + 3
            
            fmt.Printf("float64:addr:%p, value:%f\n",float32Number,*float32Number)
        }
        
        运行结果：
        number:addr:0xc000018090, value:5
        float64:addr:0xc000018090, value:3.000000
        由运行可知使用unsafe.Pointer强制类型转换后指针指向的地址是没有改变，只是类型发生了改变。
        这个例子本身没什么意义，正常项目中也不会这样使用。

        注：总结一下基本使用：先把*T类型转换成unsafe.Pointer类型，然后在进行强制转换转成你需要的指针类型即可。
        
        Sizeof、Alignof、Offsetof三个函数的基本使用：
        type User struct {
            Name string
            Age uint32
            Gender bool // 男:true 女：false 就是举个例子别吐槽我这么用。。。。
        }
        
        func func_example()  {
            // sizeof
            fmt.Println(unsafe.Sizeof(true))
            fmt.Println(unsafe.Sizeof(int8(0)))
            fmt.Println(unsafe.Sizeof(int16(10)))
            fmt.Println(unsafe.Sizeof(int(10)))
            fmt.Println(unsafe.Sizeof(int32(190)))
            fmt.Println(unsafe.Sizeof("asong"))
            fmt.Println(unsafe.Sizeof([]int{1,3,4}))

            // Offsetof
            user := User{Name: "Asong", Age: 23,Gender: true}
            userNamePointer := unsafe.Pointer(&user)
            
            nNamePointer := (*string)(unsafe.Pointer(userNamePointer))
            *nNamePointer = "Golang梦工厂"
            
            nAgePointer := (*uint32)(unsafe.Pointer(uintptr(userNamePointer) + unsafe.Offsetof(user.Age)))
            *nAgePointer = 25
            
            nGender := (*bool)(unsafe.Pointer(uintptr(userNamePointer)+unsafe.Offsetof(user.Gender)))
            *nGender = false
            
            fmt.Printf("u.Name: %s, u.Age: %d,  u.Gender: %v\n", user.Name, user.Age,user.Gender)
            // Alignof
            var b bool
            var i8 int8
            var i16 int16
            var i64 int64
            var f32 float32
            var s string
            var m map[string]string
            var p *int32
            
            fmt.Println(unsafe.Alignof(b))
            fmt.Println(unsafe.Alignof(i8))
            fmt.Println(unsafe.Alignof(i16))
            fmt.Println(unsafe.Alignof(i64))
            fmt.Println(unsafe.Alignof(f32))
            fmt.Println(unsafe.Alignof(s))
            fmt.Println(unsafe.Alignof(m))
            fmt.Println(unsafe.Alignof(p))
        }