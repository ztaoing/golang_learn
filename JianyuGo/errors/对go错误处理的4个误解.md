[Go 语言中错误处理的机制一直是各大 Gopher 热议的问题。]
甚至一直有人寄望 Go 支持 throw 和 catch 关键字，实现与其他语言类似的特性。社区里的讨论也从未停过。

今天煎鱼带大家了解几个 Go 语言的错误处理中，大家最关心，也是最容易被误解、被嫌弃的问题：
1、为什么不支持 try-catch？
2、为什么不支持全局捕获的机制？
3、为什么要这么设计错误处理？
4、未来的错误处理机制会怎么样？

【落寞的 try-catch】

在 Go1 时，大家知道基本不可能支持。于是打起了 Go2 的主意。为什么 Go 就不能支持 try-catch 组合拳？
上一年宣发了 Go2 的构想，所以在 2020 年就有小伙伴乘机提出过类似 
《proposal: Go 2: use keywords throw, catch, and guard to handle errors[1]》的提案，
这可是其他语言都支持的，Go 语言怎么了？

下面来自该提案的演示，Go1 的错误处理：


    type data struct {}
    
    func (d data) bar() (string, error) {
      return "", errors.New("err")
    }
    
    func foo() (data, error) {
      return data{}, errors.New("err")
    }
    
    func do () (string, error) {
      d, err := foo()
      if err != nil {
        return "", err
      }
    
        s, err := d.bar()
        if err != nil {
            return "", err
        }
    
        return s, nil
    }
新提案所改造的方式：

    type data struct {}
    
    func (d data) bar() string {
      throw "", errors.New("err")
    }
    
    func foo() (d data) {
      throw errors.New("err")
      return
    }
    
    func do () (string, error) {
      catch err {
        return "", err
      }
    
        s := foo().bar()
        return s, nil
    }

不过答复非常明确，@davecheney 在底下回复“以最强烈的措辞，不（In the strongest possible terms, no）”。
这可让人懵了圈，为什么这么硬呢？
其实 Go 官方早在《Error Handling — Problem Overview[2]》提案早已明确提过，
Go 官方在设计上会有意识地选择使用显式错误结果和显式错误检查。

结合《language: Go 2: error handling meta issue[3]》可得知，要拒绝 try-catch 关键字的主要原因是：
* 会涉及到额外的流程控制，因为使用 try 的复杂表达式，会导致函数意外返回。
* 在表达式层面上没有流程控制结构，只有 panic 关键字，它不只是从一个函数返回。

说白了，就是设计理念不合，加之实现上也不大合理。在以往的多轮讨论中早已被 Go 团队拒绝了。
反之 Go 团队倒是一遍遍在回答这个问题，已经不大耐烦了，直接都整理了 issues 版的 FAQ 了。

【想捕获所有 panic】

在 Go 语言中，有一个点，很多新同学会不一样碰到的。那就是在 goroutine 中如果 panic 了，
没有加 recover 关键字（有时候也会忘记），就会导致程序崩溃。

又或是"以为"加了 recover 就能保障一个 goroutine 下所派生出来的 goroutine 所产生的 panic，一劳永逸。

但现实总是会让人迷惑，我经常会看到有同学提出类似的疑惑：
    
    为什么recover要设计成，不能处理更上面的panic

这时候，有其他语言经验的同学中，又有想到了一个利器。能不能设置一个全局的错误处理 handler。
像是 PHP 语言也可以有类似的方法：
    
    set_error_handler();
    set_exception_handler();
    register_shutdown_function();

显然，Go 语言中并没有类似的东西。归类一下，我们聚焦以下两个问题：

1、为什么 recover 不能捕获更上层的 panic？

2、为什么 Go 没有全局的错误处理方法？

【源码层面】
如果是讲设计的话，其实只是通过 Go 的 GMP 模型和 defer+panic+recver 的源码剖析就能知道了。

本质上 defer+panic 都是挂载在 G 上的，可查看我以前写的《深入理解 Go panic and recover[4]》，你会有更多深入的理解。

    深入理解 Go panic and recover: https://eddycjy.com/posts/go/panic/2019-05-21-panic-and-recover/

    【深入理解 Go panic and recover】
    作为一个 gophper，我相信你对于 panic 和 recover 肯定不陌生，但是你有没有想过。
    当我们执行了这两条语句之后。底层到底发生了什么事呢？前几天和同事刚好聊到相关的话题，
    发现其实大家对这块理解还是比较模糊的。希望这篇文章能够从更深入的角度告诉你为什么，它到底做了什么事？

    [思考]
    一、为什么会中止运行
    func main() {
	    panic("EDDYCJY.")
    }

    输出结果：
    $ go run main.go
    panic: EDDYCJY.
    
    goroutine 1 [running]:
    main.main()
    /Users/eddycjy/go/src/github.com/EDDYCJY/awesomeProject/main.go:4 +0x39
    exit status 2
    请思考一下，为什么执行 panic 后会导致应用程序运行中止？（而不是单单说执行了 panic 所以就结束了这么含糊）
    
    二、为什么不会中止运行
    func main() {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("recover: %v", err)
            }
        }()

	    panic("EDDYCJY.")
    }
    输出结果：
    $ go run main.go 
    2019/05/11 23:39:47 recover: EDDYCJY.

    请思考一下，为什么加上 defer + recover 组合就可以保护应用程序？
    
    三、不设置 defer 行不
    上面问题二是 defer + recover 组合，那我去掉 defer 是不是也可以呢？如下：\
    func main() {
        if err := recover(); err != nil {
            log.Printf("recover: %v", err)
        }
    
        panic("EDDYCJY.")
    }
    输出结果：
    $ go run main.go
    panic: EDDYCJY.
    
    goroutine 1 [running]:
    main.main()
    /Users/eddycjy/go/src/github.com/EDDYCJY/awesomeProject/main.go:10 +0xa1
    exit status 2

    竟然不行，啊呀毕竟入门教程都写的 defer + recover 组合 “万能” 捕获。
    但是为什么呢。去掉 defer 后为什么就无法捕获了？

    请思考一下，为什么需要设置 defer 后 recover 才能起作用？
    同时你还需要仔细想想，我们设置 defer + recover 组合后就能无忧无虑了吗，各种 “乱” 写了吗？
    
    四、为什么起个 goroutine 就不行
    func main() {
        // recover在goroutine中
        go func() {
            defer func() {
                if err := recover(); err != nil {
                    log.Printf("recover: %v", err)
                }
            }()
        }()
    
        panic("EDDYCJY.")
    }
    输出结果：
    $ go run main.go 
    panic: EDDYCJY.
    
    goroutine 1 [running]:
    main.main()
    /Users/eddycjy/go/src/github.com/EDDYCJY/awesomeProject/main.go:14 +0x51
    exit status 2

    请思考一下，为什么新起了一个 Goroutine 就无法捕获到异常了？到底发生了什么事…
    [源码]
    接下来我们将带着上述 4+1 个小思考题，开始对源码的剖析和分析，
    尝试从阅读源码中找到思考题的答案和更多为什么
    
    [数据结构]
    type _panic struct {
        argp      unsafe.Pointer //指向 defer 延迟调用的参数的指针
        arg       interface{}    //panic 的原因，也就是调用 panic 时传入的参数
        link      *_panic        //指向上一个调用的 _panic,link 字段，可得知其是一个链表的数据结构
        recovered bool           //panic 是否已经被处理，也就是 是否被 recover
        aborted   bool           //panic 是否被中止
    }
    在 panic 中是使用 _panic 作为其基础单元的，每执行一次 panic 语句，都会创建一个 _panic。
    它包含了一些基础的字段用于存储当前的 panic 调用情况.

    [恐慌 panic]
    func main() {
	    panic("EDDYCJY.")
    }

    输出结果：

    $ go run main.go
    panic: EDDYCJY.
    
    goroutine 1 [running]:
    main.main()
    /Users/eddycjy/go/src/github.com/EDDYCJY/awesomeProject/main.go:4 +0x39
    exit status 2

    我们去反查一下 panic 处理具体逻辑的地方在哪，如下：
    $ go tool compile -S main.go
    "".main STEXT size=66 args=0x0 locals=0x18
    0x0000 00000 (main.go:23)	TEXT	"".main(SB), ABIInternal, $24-0
    0x0000 00000 (main.go:23)	MOVQ	(TLS), CX
    0x0009 00009 (main.go:23)	CMPQ	SP, 16(CX)
    ...
    0x002f 00047 (main.go:24)	PCDATA	$2, $0
    0x002f 00047 (main.go:24)	MOVQ	AX, 8(SP)
    0x0034 00052 (main.go:24)	CALL	runtime.gopanic(SB) (这里)
    
    显然汇编代码直指内部实现是 runtime.gopanic，我们一起来看看这个方法做了什么事，如下（省略了部分）
    
    func gopanic(e interface{}) {
        gp := getg()
        ...
        var p _panic
        p.arg = e
        p.link = gp._panic

        gp._panic = (*_panic)(noescape(unsafe.Pointer(&p)))
        
        for {
            d := gp._defer
            if d == nil {
                break
            }
    
            // defer...
            ...
            d._panic = (*_panic)(noescape(unsafe.Pointer(&p)))
    
            p.argp = unsafe.Pointer(getargp(0))

            // 若当前存在 defer 调用，则调用 reflectcall 方法去执行先前 defer 中延迟执行的代码，
            // 若在执行过程中需要运行 recover 将会调用 gorecover 方法
            reflectcall(nil, unsafe.Pointer(d.fn), deferArgs(d), uint32(d.siz), uint32(d.siz))
            p.argp = nil
    
            // recover...
            if p.recovered {
                ...
                mcall(recovery)
                throw("recovery failed") // mcall should not return
            }
        }
        // 结束前，使用 preprintpanics 方法打印出所涉及的 panic 消息
        preprintpanics(gp._panic)
        // 最后调用 fatalpanic 中止应用程序，实际是执行 exit(2) 进行最终退出行为的
        fatalpanic(gp._panic) // should not return
        *(*int)(nil) = 0      // not reached
    }
    1、获取指向当前 Goroutine 的指针
    2、初始化一个 panic 的基本单位 _panic 用作后续的操作
    3、获取当前 Goroutine 上挂载的 _defer（数据结构也是链表）
    4、若当前存在 defer 调用，则调用 reflectcall 方法去执行先前 defer 中延迟执行的代码，若在执行过程中需要运行 recover 将会调用 gorecover 方法
    5、结束前，使用 preprintpanics 方法打印出所涉及的 panic 消息
    6、最后调用 fatalpanic 中止应用程序，实际是执行 exit(2) 进行最终退出行为的

    通过对上述代码的执行分析，可得知 panic 方法实际上就是处理当前 Goroutine(g) 上
    所挂载的 ._panic 链表（所以无法对其他 Goroutine 的异常事件响应），
    然后对其所属的 defer 链表和 recover 进行检测并处理，最后调用退出命令中止应用程序

    [无法恢复的恐慌 fatalpanic]
        func fatalpanic(msgs *_panic) {
            pc := getcallerpc()
            sp := getcallersp()
            gp := getg()
            var docrash bool
        
            systemstack(func() {
                if startpanic_m() && msgs != nil {
                    ...
                    printpanics(msgs)
                }
        
                docrash = dopanic_m(gp, pc, sp)
            })
        
            systemstack(func() {
                exit(2)
            })
        
            *(*int)(nil) = 0
        }
    我们看到在异常处理的最后会执行该方法，似乎它承担了所有收尾工作。
    实际呢，它是在最后对程序执行 exit 指令来达到中止运行的作用，
    但在结束前它会通过 printpanics 递归输出所有的异常消息及参数。代码如下：
        func printpanics(p *_panic) {
            if p.link != nil {
                printpanics(p.link)
                print("\t")
            }
            print("panic: ")
            printany(p.arg)
            
            if p.recovered {
                print(" [recovered]")
            }
            print("\n")
         }
    所以不要以为所有的异常都能够被 recover 到，实际上像 fatal error 和 runtime.throw 都是无法被 recover 到的，
    甚至是 oom 也是直接中止程序的，也有反手就给你来个 exit(2) 教做人。因此在写代码时你应该要相对注意些，
    “恐慌” 是存在无法恢复的场景的.
    
    [恢复 recover]
        func main() {
            defer func() {
                if err := recover(); err != nil {
                    log.Printf("recover: %v", err)
                }
            }()
        
            panic("EDDYCJY.")
        }
    
    输出结果：
    $ go run main.go 
    2019/05/11 23:39:47 recover: EDDYCJY.

    和预期一致，成功捕获到了异常。但是 recover 是怎么恢复 panic 的呢？再看看汇编代码，如下：

    $ go tool compile -S main.go
    "".main STEXT size=110 args=0x0 locals=0x18
    0x0000 00000 (main.go:5)	TEXT	"".main(SB), ABIInternal, $24-0
    ...
    0x0024 00036 (main.go:6)	LEAQ	"".main.func1·f(SB), AX
    0x002b 00043 (main.go:6)	PCDATA	$2, $0
    0x002b 00043 (main.go:6)	MOVQ	AX, 8(SP)
    0x0030 00048 (main.go:6)	CALL	runtime.deferproc(SB)
    ...
    0x0050 00080 (main.go:12)	CALL	runtime.gopanic(SB)
    0x0055 00085 (main.go:12)	UNDEF
    0x0057 00087 (main.go:6)	XCHGL	AX, AX
    0x0058 00088 (main.go:6)	CALL	runtime.deferreturn(SB)
    ...
    0x0022 00034 (main.go:7)	MOVQ	AX, (SP)
    0x0026 00038 (main.go:7)	CALL	runtime.gorecover(SB)
    0x002b 00043 (main.go:7)	PCDATA	$2, $1
    0x002b 00043 (main.go:7)	MOVQ	16(SP), AX
    0x0030 00048 (main.go:7)	MOVQ	8(SP), CX
    ...
    0x0056 00086 (main.go:8)	LEAQ	go.string."recover: %v"(SB), AX
    ...
    0x0086 00134 (main.go:8)	CALL	log.Printf(SB)
    ...
    
    通过分析底层调用，可得知主要是如下几个方法：
    runtime.deferproc
    runtime.gopanic
    runtime.deferreturn
    runtime.gorecover

    在上小节中，我们讲述了简单的流程，gopanic 方法会调用当前 Goroutine 下的 defer 链表，
    若 reflectcall 执行中遇到 recover 就会调用 gorecover 进行处理，该方法代码如下：
    
        func gorecover(argp uintptr) interface{} {
            gp := getg()
            p := gp._panic
            if p != nil && !p.recovered && argp == uintptr(p.argp) {
                // 核心就是修改 recovered 字段
                p.recovered = true
                return p.arg
            }
            return nil
        }
    这代码，看上去挺简单的，核心就是修改 recovered 字段。该字段是用于标识当前 panic 是否已经被 recover 处理。
    但是这和我们想象的并不一样啊，程序是怎么从 panic 流转回去的呢？
    是不是在核心方法里处理了呢？我们再看看 gopanic 的代码，如下：

        func gopanic(e interface{}) {
                ...
                for {
                    // defer...
                    ...
                    pc := d.pc
                    sp := unsafe.Pointer(d.sp) // must be pointer so it gets adjusted during stack copy
                    freedefer(d)
                    
                    // recover... 判断当前 _panic 中的 recover 是否已标注为处理
                    if p.recovered {

                        atomic.Xadd(&runningPanicDefers, -1)
            
                        gp._panic = p.link
                        for gp._panic != nil && gp._panic.aborted {
                            gp._panic = gp._panic.link
                        }
                        if gp._panic == nil { 
                            gp.sig = 0
                        }
                        // 将相关需要恢复的栈帧信息传递给 recovery 方法的 gp 参数
                        // （每个栈帧对应着一个未运行完的函数。栈帧中保存了该函数的返回地址和局部变量）
                        gp.sigcode0 = uintptr(sp)
                        gp.sigcode1 = pc
                        //执行 recovery 进行恢复动作
                        mcall(recovery)
                        throw("recovery failed") 
                    }
                }
                ...
         }
    我们回到 gopanic 方法中再仔细看看，发现实际上是包含对 recover 流转的处理代码的。恢复流程如下：
    1、判断当前 _panic 中的 recover 是否已标注为处理
    2、从 _panic 链表中删除已标注中止的 panic 事件，也就是删除已经被恢复的 panic 事件
    3、将相关需要恢复的栈帧信息传递给 recovery 方法的 gp 参数
        （每个栈帧对应着一个未运行完的函数。栈帧中保存了该函数的返回地址和局部变量）
    4、执行 recovery 进行恢复动作

    从流程来看，最核心的是 recovery 方法。它承担了异常流转控制的职责。代码如下：
        func recovery(gp *g) {
            sp := gp.sigcode0
            pc := gp.sigcode1
        
            if sp != 0 && (sp < gp.stack.lo || gp.stack.hi < sp) {
                print("recover: ", hex(sp), " not in [", hex(gp.stack.lo), ", ", hex(gp.stack.hi), "]\n")
                throw("bad recovery")
            }
            // 实际上设置的是编译器中伪寄存器的值，常常被用于维护上下文等。
            gp.sched.sp = sp
            gp.sched.pc = pc
            gp.sched.lr = 0
            gp.sched.ret = 1
            
            gogo(&gp.sched)
        }
    粗略一看，似乎就是很简单的设置了一些值？但实际上设置的是编译器中伪寄存器的值，常常被用于维护上下文等。
    在这里我们需要结合 gopanic 方法一同观察 recovery 方法。

    它所使用的栈指针 sp 和程序计数器 pc 是由当前 defer 在调用流程中的 deferproc 传递下来的，
    因此实际上最后是通过 gogo 方法跳回了 deferproc 方法。另外我们注意到：
            gp.sched.ret = 1
    
    在底层中程序将 gp.sched.ret 设置为了 1，"也就是没有实际调用 deferproc 方法"，直接修改了其返回值。
    意味着默认它已经处理完成。直接转移到 deferproc 方法的下一条指令去。
    至此为止，异常状态的流转控制就已经结束了。接下来就是继续走 defer 的流程了。
    
    为了验证这个想法，我们可以看一下核心的跳转方法 gogo ，代码如下：

    // void gogo(Gobuf*)
    // restore state from Gobuf; longjmp
    TEXT runtime·gogo(SB),NOSPLIT,$8-4
    MOVW	buf+0(FP), R1
    MOVW	gobuf_g(R1), R0
    BL	setg<>(SB)
    
        MOVW	gobuf_sp(R1), R13	// restore SP==R13
        MOVW	gobuf_lr(R1), LR
        MOVW	gobuf_ret(R1), R0
        MOVW	gobuf_ctxt(R1), R7
        MOVW	$0, R11
        MOVW	R11, gobuf_sp(R1)	// clear to help garbage collector
        MOVW	R11, gobuf_ret(R1)
        MOVW	R11, gobuf_lr(R1)
        MOVW	R11, gobuf_ctxt(R1)
        MOVW	gobuf_pc(R1), R11
        CMP	R11, R11 // set condition codes for == test, needed by stack split
        B	(R11)

    通过查看代码可得知其主要作用是从 Gobuf 恢复状态。
    简单来讲就是--》将寄存器的值修改为对应 Goroutine(g) 的值，而在文中讲了很多次的 Gobuf，如下：
        type gobuf struct {
            sp   uintptr
            pc   uintptr
            g    guintptr
            ctxt unsafe.Pointer
            ret  sys.Uintreg
            lr   uintptr
            bp   uintptr
        }
    其实它存储的就是 Goroutine 切换上下文时所需要的一些东西
    【拓展】
    const(
	OPANIC       // panic(Left)
	ORECOVER     // recover()
	...
    )
    ...
    func walkexpr(n *Node, init *Nodes) *Node {
    ...
    switch n.Op {
    default:
    Dump("walk", n)
    Fatalf("walkexpr: switch 1 unknown op %+S", n)
    
        case ONONAME, OINDREGSP, OEMPTY, OGETG:
        case OTYPE, ONAME, OLITERAL:
            ...
        case OPANIC:
            n = mkcall("gopanic", nil, init, n.Left)
    
        case ORECOVER:
            n = mkcall("gorecover", n.Type, init, nod(OADDR, nodfp, nil))
        ...
    }
    实际上在调用 panic 和 recover 关键字时，是在编译阶段先转换为相应的 OPCODE 后，
    再由编译器转换为对应的运行时方法。并不是你所想像那样一步到位，有兴趣的小伙伴可以研究一下
    
    
【设计思想】
在本文中我们不能仅限于源码，需要更深挖，Go 设计者他的思想是什么，为什么就是不支持？
在 Go issues 中《proposal: spec: allow fatal panic handler[5]》、
《No way to catch errors from goroutines automatically[6] 》分别的针对性探讨过上述问题。

Go 团队的大当家 @Russ Cox 给出了明确的答复：
    Go 语言的设计立场是错误恢复应该在本地完成，或者完全在一个单独的进程中完成。

这就是为什么 Go 语言不能跨 goroutines 从 panic 中恢复，也不能从 throw 中恢复的根本原因，
是语言设计层面的思想所决定。

    在源码剖析时，你所看到的整套 GMP+defer+panic+recover 的机制机制，就是跟随着这个设计思想去编写和发展的。
    设计思想决定源码实现。

【建议方式】

从 Go 语言层面去动摇这个设计思想，目前来看可能性很低。至少 2021 年的现在没有看到改观。
整体上会建议提供公共的 Go 方法去规避这种情况。参考 issues 所提供的范例如下：
    
        recovery.SafeGo(logger, func() {
            method(all parameters)
        })
    
    func SafeGo(logger logging.ILogger, f func()) {
        go func() {
        defer func() {
        if panicMessage := recover(); panicMessage != nil {
            ...
        }
        }()
        
        f()
        }()
    }
是不是感觉似曾相识？
每家公司的内部库都应该有这么一个工具方法，规避偶尔忘记的 goroutine recover 所引发的奇奇怪怪问题。
也可以参照建议，利用一个单独的进程（Go 语言中是 goroutine）去统一处理这些 panic，不过这比较麻烦，较少见。

【未来会如何】
Go 社区对 Go 语言未来的错误处理机制非常关心，因为 Go1 已经米已成炊，希望在 Go2 上解决错误处理机制的问题。
期望 Go2 核心要处理的包含如下几点（#40432）：
    1、对于 Go2，我们希望使错误检查更加轻便，减少专门用于错误检查的 Go 程序代码的数量。我们还想让写错误处理更方便，减少程序员花时间写错误处理的可能性。
    2、错误检查和错误处理都必须保持明确，即在程序文本中可见。我们不希望重复异常处理的陷阱。
    3、现有的代码必须继续工作，并保持与现在一样的有效性。任何改变都必须与现有的代码相互配合。

为此，许多人提过不少新的提案...很可惜，截止 2021.08 月底为止，有许多人试图改变语言层面以达到这些目标，但没有一个新的提案被接受。
现在也有许多变更并入 Go2 提案，主要是 error-handling 方面的优化。
大家有兴趣可以看看我之前写的：《先睹为快，Go2 Error 的挣扎之路》，相信能给你带来不少新知识

    https://mp.weixin.qq.com/s?__biz=MzUxMDI4MDc1NA==&mid=2247484512&idx=1&sn=05fc8da62f29397b9931b1cdc68b952c&scene=21#wechat_redirect
