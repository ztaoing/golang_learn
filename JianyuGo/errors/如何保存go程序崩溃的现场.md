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