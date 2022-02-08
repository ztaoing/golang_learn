[为什么 panic 值得思考？]
初学 Go 的时候，心里常常很多疑问，有时候看似懂了的问题，其实是是而非。

panic 究竟是啥？看似显而易见的问题，但是却回答不出个所以然来。奇伢分两个章节来彻底搞懂 panic 的知识：
* 姿势篇：摸清楚 panic 的诞生，它不是石头里蹦出来的，总结有三种姿势；
* 原理篇：彻底搞明白 panic 的内部原理，理解 panic 的深层原理；

[panic 的三种姿势]
什么时候会产生 panic ？
我们先从“形”来学习。从程序猿的角度来看，可以分为主动和被动方式:
主动方式：程序猿主动调用 panic( ) 函数；
被动的方式：
    * 编译器的隐藏代码触发；
    * 内核发送给进程信号触发 ；

[编译器的隐藏代码]
Go 之所以简单又强大，编译器居功至伟。非常多的事情是编译器帮程序猿做了的，逻辑补充，内存的逃逸分析等等。
包括 panic 的抛出！

举个非常典型的例子：整数算法除零会发生 panic，怎么做到的？

看一段极简代码：
    
    func divzero(a, b int) int {
        c := a/b
        return c
    }
上面函数就会有除零的风险，当 b 等于 0 的时候，程序就会触发  panic，然后退出，如下：
    
    root@ubuntu:~/code/gopher/src/panic# ./test_zero 

    panic: runtime error: integer divide by zero
    
    goroutine 1 [running]:
    main.zero(0x64, 0x0, 0x0)
    /root/code/gopher/src/panic/test_zero.go:6 +0x52

问题来了：程序怎么触发的 panic ？代码面前无秘密。
可代码看不出啥呀，不就是一行 c := a/b 嘛？

奇伢说的是汇编代码。因为这段隐藏起来的逻辑，是编译器帮你加的。

用 dlv 调试断点到 divzero 函数，然后执行 disassemble ，你就能看到秘密了。奇伢截取部分汇编，并备注了下：

(dlv) disassemble
TEXT main.zero(SB) /root/code/gopher/src/panic/test_zero.go

    // 判断 b 是否等于 0 
    test_zero.go:6  0x4aa3c1    4885c9          test rcx, rcx
    // 不等于 0 就跳转到 0x4aa3c8 执行指令，否则就往下执行
    test_zero.go:6  0x4aa3c4    7502            jnz 0x4aa3c8
    // 执行到这里，就说明 b 是 0 值，就跳转到 0x4aa3ed ，也就是 call $runtime.panicdivide
    =>  test_zero.go:6  0x4aa3c6    eb25            jmp 0x4aa3ed
    test_zero.go:6  0x4aa3c8    4883f9ff        cmp rcx, -0x1
    test_zero.go:6  0x4aa3cc    7407            jz 0x4aa3d5
    test_zero.go:6  0x4aa3ce    4899            cqo
    test_zero.go:6  0x4aa3d0    48f7f9          idiv rcx
    // ...
    test_zero.go:7  0x4aa3ec    c3              ret
    // 看到神奇的函数了嘛 ！
    test_zero.go:6  0x4aa3ed    e8ee27f8ff      call $runtime.panicdivide
编译器偷偷加上了一段 if/else 的判断逻辑，并且还给加了 runtime.panicdivide  的代码。

1、如果 b == 0 ，那么跳转执行函数 runtime.panicdivide ；
    再来看一眼 panicdivide 函数，这是一段极简的封装： 
    
    // runtime/panic.go
    func panicdivide() {
        panicCheck2("integer divide by zero")
        panic(divideError) // 这里面调用的就是 panic() 函数。
    }
除零触发的 panic 就是这样来的，它不是石头里蹦出来的，
而是编译器多加的逻辑判断保证了除数为 0 的时候，触发 panic 函数。

划重点：编译器加的隐藏逻辑，调用了抛出 panic 的函数。Go 的编译器才是真大佬！

[进程信号触发]
最典型的是非法地址访问，比如， nil 指针 访问会触发 panic，怎么做到的？
看一个极简的例子：
    
    func nilptr(b *int) int {
        c := *b
        return c
    }
    当调用 nilptr( nil ) 的时候，将会导致进程异常退出：

    root@ubuntu:~/code/gopher/src/panic# ./test_nil 

    panic: runtime error: invalid memory address or nil pointer dereference
    [signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x4aa3bc]
    
    goroutine 1 [running]:
    main.nilptr(0x0, 0x0)
    /root/code/gopher/src/panic/test_nil.go:6 +0x1c

问题来了：这里的 panic 又是怎么形成的呢？

    在 Go 进程启动的时候会注册默认的信号处理程序（ sigtramp ）
在 cpu 访问到 0 地址会触发 page fault 异常，这是一个非法地址，
内核会发送 SIGSEGV 信号给进程，所以当收到 SIGSEGV 信号的时候，
就会让 sigtramp 函数来处理，最终调用到 panic 函数 ：

        
        // 信号处理函数回调
    sigtramp （纯汇编代码）
    -> sigtrampgo （ signal_unix.go ）
      -> sighandler  （ signal_sighandler.go ）
        -> preparePanic （ signal_amd64x.go ）
          -> sigpanic （ signal_unix.go ）
            -> panicmem 
              -> panic (内存段错误)
    在 sigpanic 函数中会调用到 panicmem  ，在这个里面就会调用 panic 函数，
    从而走上了 Go 自己的 panic 之路。

panicmem 和 panicdivide 类似，都是对 panic( ) 的极简封装：

    func panicmem() {
        panicCheck2("invalid memory address or nil pointer dereference")
        panic(memoryError)
    }

划重点：这种方式是通过信号软中断的方式来走到 Go 注册的信号处理逻辑，从而调用到 panic( )  的函数。
童鞋可能会好奇，信号处理的逻辑什么时候注册进去的？

在进程初始化的时候，创建 M0（线程）的时候用系统调用 sigaction 给信号注册处理函数为 sigtramp ，
调用栈如下：

    mstartm0 （proc.go）
    -> initsig (signal_unix.go:113)
    -> setsig （os_linux.go）
    
    这样的话，以后触发了信号软中断，就能调用到 Go 的信号处理函数，从而进行语言层面的 panic 处理 。
总的来说，这个是从系统层面到特定语言层面的处理转变。

[程序猿主动]
第三种方式，就是程序猿自己主动调用 panic 抛出来的。
    
    func main() {
        panic("panic test")
    }
    
[聊聊 panic 到底是什么？]
现在我们摸透了 panic 产生的姿势，以上三种方式，无论哪一种都归一到 panic( ) 这个函数调用。
所以有一点很明确：panic 这个东西是语言层面的处理逻辑。

panic 发生之后，如果 Go 不做任何特殊处理，默认行为是打印堆栈，退出程序。

现在回到最本源的问题：panic 到底是什么？

这里不纠结概念，只描述几个简单的事实：
* panic( ) 函数内部会产生一个关键的数据结构体 _panic ，并且挂接到 goroutine 之上(panic是单个goroutine上的，不是全局的)；
* panic( ) 函数内部会执行 _defer 函数链条（link），并针对 _panic 的状态进行对应的处理；

什么叫做 panic( ) 的对应的处理？

循环执行 goroutine 上面的 _defer 函数链，

没有recover：如果执行完了都还没有恢复 _panic 的状态（没有被recover重新设置），
那就没得办法了，退出进程，打印堆栈。

有recover：如果在 goroutine 的 _defer 链上，有个朋友 recover 了一下，把这个 _panic 标记成恢复
，那事情就到此为止，就从这个 _defer 函数执行后续正常代码即可，走 deferreturn 的逻辑。

【所以，panic 是什么 ？】
小奇伢认为，它就是个特殊函数调用，仅此而已。

* panic 究竟是啥？是一个结构体？还是一个函数？
* 为什么 panic 会让 Go 进程退出的 ？
* 为什么 recover 一定要放在 defer 里面才生效？
* 为什么 recover 已经放在 defer 里面，但是进程还是没有恢复？
* 为什么 panic 之后，还能再 panic ？有啥影响？
