以下内容来自【煎鱼的微信公众号】【网管叨bi叨】【奇伢云存储】【Golang技术分享】

---

* [Go: A Documentary 发布！]( https://mp.weixin.qq.com/s/5MtBE8vecKPOmRUYu2E-lg)

---


* [6 万 Star！ Go 语言资源大全（上）](https://mp.weixin.qq.com/s/gL3p0pCVlZzrLCwYk7gTvw)

* [6 万 Star！ Go 语言资源大全（中）]( https://mp.weixin.qq.com/s/DR39kTPz9xLCwNVKV6K4Xw)

* [6 万 Star！ Go 语言资源大全（下）]( https://mp.weixin.qq.com/s/KPb4rxv-BuzCpzYv9DWyiQ)

---
* [如何让Gitlab私有仓库支持Go Get]( https://mp.weixin.qq.com/s/nMg4HB4sJkgrEC9iyfT4_A)

 ---

* [Golang 数据结构到底是怎么回事？gdb调一调？](https://mp.weixin.qq.com/s/qtQoZaX_SJi6_TD-uGUPWA) “ 不仅限于语法，使用gdb，dlv工具更深层的剖析golang的数据结构”

slice，map，channel 这三种类型必须使用make来创建，就是这个道理。因为如果仅仅定义了类型变量，那仅仅是代表了分配了这个变量本身的内存空间，并且初始化是nil，一旦你直接用，那么就会导致非法地址引用的问题

---



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
  2、 减少GC：BigCache如何加速并发访问以及避免高额的GC开销： https://pengrl.com/p/35302/

[Go 缓冲系列之-free-cache]( https://mp.weixin.qq.com/s/VmPIW6HhVrDyeADiRmkC_Q) （也是使用减小锁的粒度、go 1.5版本之后，如果你使用的map的key和value中都不包含指针，那么GC会忽略这个map。）


[runtime]
* 什么是go runtime.KeepAlive
  go 官方文档: https://pkg.go.dev/runtime#KeepAlive
  文档: https://medium.com/a-journey-with-go/go-keeping-a-variable-alive-c28e3633673a

* [编程思考：对象生命周期的问题](https://mp.weixin.qq.com/s/Hoy9cqHe9RZqEA5T9Dys5w)



[类型的比较]( golang.org/ref/spec#comparison_operators)
* 可比较类型和不可不叫类型。对于不可比较类型，如何比较他们包含的值是否相等呢？使用reflect.DeepEqual


* [json.unmarshal](pkg.go.dev/encoding/json#unmarshal)
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



* [多路复用] channel--》multiplex
  *[atomic.Value为什么不加锁也能保证数据线程安全]

* [go中的零值，它有什么作用？](https://golang.org/ref/spec#the_zero_value)
  布尔型为false；数字型为0；字符串型为""；指针、函数、接口、切片、通道和映射都为nil

* [go是如何实现启动参数的加载的？]compile


* [select机制] select


* [在go容器里设置gomaxprocs的正确姿势：][dockers](  https://mp.weixin.qq.com/s/kbZsJncgVZv30_TwVrLyLQ)

---

* [unsafe包]( https://mp.weixin.qq.com/s/dulgHWM-mjrYIdD9nHZyYg)

* [详解 Go 团队不建议用的 unsafe.Pointer]( https://mp.weixin.qq.com/s/8qtHdw2JiRQ1cXlzbJ0ANA)

---

* [http 请求怎么确定边界？]( https://mp.weixin.qq.com/s/1SzIWYxgAV6Ourb9HSrQZQ )，HTTP 是基于TCP协议的应用层协议，而 TCP 是面向数据流的协议，是没有边界的。HTTP 作为应用层协议需要自己明确定义数据边界。


* [Go原生网络轮询器（netpoller）剖析](https://mp.weixin.qq.com/s/oDLYJqkwF2Em_hcRNLZ9qg) net.Listen；l.Accept；conn.Read

* [Go udp 的高性能优化](  https://mp.weixin.qq.com/s/ZfjXhgoFP0InA18uWlQByw)  golang udp 的锁竞争问题

* [go网络编程和tcp抓包实操 : network-》getTCPPackage](https://mp.weixin.qq.com/s/Ou7YSLR1seHfS27rAgdbQQ)

* [go中如何强制关闭tcp连接 : network-》getTCPPackage](   https://mp.weixin.qq.com/s/ySvp47waWKjy4YK7NZIauQ)

* [网友很强大，发现了Go并发下载的Bug]( https://mp.weixin.qq.com/s/Kd4np83CKEOLQ3K0eWxlKg)

* [连接一个ip不存在的主机时，握手过程是怎样的？: network-》ConnIP  ](https://mp.weixin.qq.com/s/Czv0CxDKqr2QNItO4aNZMA)
  连接一个ip不存在的主机时，握手过程是怎样的？
  连接一个IP地址存在但是端口不存在的主机时，握手过程是怎样的？

* [context使用不当引发的一个bug]( https://mp.weixin.qq.com/s/lJxjlDg5SkQyNLZBpOPP5A)

* [解决golang开发socket服务时粘包半包bug]( http://xiaorui.cc/?p=2888)

---

* [从CPU角度理解go中的结构体内存对齐 memory-》align ](https://mp.weixin.qq.com/s/TDIM1tspIEWpQCH_SNGnog)

* [详解 Go 内存对齐 memory-》align](https://mp.weixin.qq.com/s/ApJCdMSTydxN5zgxhzj21w)

* [Go程序内存分配过多？]( https://mp.weixin.qq.com/s/zBHPYJWnGf67Ex8i4cV8Eg) (如何优化内存)

* [Go 编程怎么也有踩内存？](  https://mp.weixin.qq.com/s/tXAP8_U63QLNj1h0ZMvXPw) (由小结构 向大的结构转换，导致内存占用变大，变大后的结构占用了后边结构的内存，导致后边结构的前边的内存的内容被覆盖了)

* [Go 内存泄露之痛，这篇把 Go timer.After 问题根因讲透了！]( https://mp.weixin.qq.com/s/KSBdPkkvonSES9Z9iggElg)

* [为什么 Go 占用那么多的虚拟内存？]( https://mp.weixin.qq.com/s/eVHK_ey8SgqEtl8v_Nurxg) （需要多次阅读）


---


* [go五种原子性操作的用法详解] memory-》atomic-》cas +atomicMutex
  原子性：外界不会看到只执行到一半的状态！
  CPU执行一些列操作时不可能不发生中断，但是如果我们在执行多个操作时，能让他们的中间状态对外不可见，
  那我们就可以说他拥有了"不可分割"的原子性
  Go语言通过内置包sync/atomic提供了对原子操作的支持，其提供的原子操作有以下几大类：
  1、增减，操作的方法名方式为AddXXXType，保证对操作数进行原子的增减，支持的类型为int32、int64、uint32、uint64、uintptr，使用时以实际类型替换前面我说的XXXType就是对应的操作方法。
  2、载入，保证了读取到操作数前没有其他任务对它进行变更，操作方法的命名方式为LoadXXXType，支持的类型除了基础类型外还支持Pointer，也就是支持载入任何类型的指针。
  3、存储，有载入了就必然有存储操作，这类操作的方法名以Store开头，支持的类型跟载入操作支持的那些一样。
  4、比较并交换，也就是CAS （Compare And Swap），像Go的很多并发原语实现就是依赖的CAS操作，同样是支持上面列的那些类型。
  5、交换，这个简单粗暴一些，不比较直接交换，这个操作很少会用。

* *互斥锁和院系操作的区别：
  1、使用目的：互斥锁是用来保护一段逻辑，原子操作用于对一个变量的更新保护。
  2、底层实现：Mutex由操作系统的调度器实现，而atomic包中的原子操作则由底层硬件指令直接提供支持，这些指令在执行的过程中是不允许中断的，
  因此原子操作可以在lock-free的情况下保证并发安全，并且它的性能也能做到随CPU个数的增多而线性扩展。

* [Go 并发编程 — 深入浅出 sync.Pool]( https://mp.weixin.qq.com/s/1hLgu2XBBJkLzvI6pOR80g) (解释了pool的每个特点)

* [一口气搞懂 Go sync.map 所有知识点]( https://mp.weixin.qq.com/s/8aufz1IzElaYR43ccuwMyA )(原生 map + 互斥锁或读写锁 mutex ;
  标准库 sync.Map（Go1.9及以后）。sync.Map 的读操作性能如此之高的原因，就在于存在 read 这一巧妙的设计，其作为一个缓存层，提供了快路径（fast path）的查找。
  同时其结合 amended 属性，配套解决了每次读取都涉及锁的问题，实现了读这一个使用场景的高性能。read缓存层的设计却使写入变慢了。)

---


* [i/o timeout net/http的坑] network->timeout

* [go的io库如何选择 network->io库]( https://mp.weixin.qq.com/s/TtN6NZ8hQ2AIf0C8wVzkjA)

* [Go语言常用文件操作汇总 network->func]( https://mp.weixin.qq.com/s/dQUEq0lJekEUH4CHEMwANw)

* [怎么使用 direct io？:io-》io.md]( https://mp.weixin.qq.com/s/gW_3JD52rtRdEqXvyg-lJQ)

* [浅析 Go IO 的知识框架 io ](https://mp.weixin.qq.com/s/JniBMBw__WLbYtigj3eiXQ)



---
* [go精妙的互斥锁设计:lock](  https://mp.weixin.qq.com/s/j0NCgrU6M8ps0zIOkOT3FQ)

* [golang的位运算:lock](  https://mp.weixin.qq.com/s/8wubPDKO6-CLLhqjGsJ7xw)


---

errors

* [go程序错误处理的一些建议]( https://mp.weixin.qq.com/s/HuZn9hnHUBx3J4bAGYBYpw)

* [对go错误处理的4个误解](  https://mp.weixin.qq.com/s/vrcn2N4ddcAHiZl5UoqTZg)

* [go的panic的三种诞生方式]( https://mp.weixin.qq.com/s/sGdTVSRxqxIezdlEASB39A)

* [go的panic的秘密都在这里]( https://mp.weixin.qq.com/s/pxWf762ODDkcYO-xCGMm2g)

* [Go 错误处理：用 panic 取代 err != nil 的模式]( https://mp.weixin.qq.com/s/p77V3_LkREuXPVLdebmmmQ) （与panic的设计理念相违背）

* [你考虑过defer close的风险吗？]errors

* [说好 defer 在 return 之后执行，为什么结果却不是？]( https://mp.weixin.qq.com/s/oP90maykSzMXjdnudOKdSw) （需要再次阅读）

* [使用 Go defer 要小心这 2 个雷区！]( https://mp.weixin.qq.com/s/ZEsWa4xUb0a7tWemVZMXVw) (问题就是针对在 for 循环里搞 defer 关键字，是否会造成什么性能影响？)

* [Go 群友提问：学习 defer 时很懵逼，这道不会做！]( https://mp.weixin.qq.com/s/lELMqKho003h0gfKkZxhHQ)

* [Go 中的 error 居然可以这样封装](  https://mp.weixin.qq.com/s/-X8MKIQFRXmENsdwyRXcCg) (封装的目的都是为了添加更多的注解信息)

* [如何保存go程序崩溃的现场](   https://mp.weixin.qq.com/s/RktnMydDtOZFwEFLLYzlCA)



* [一文带你由浅入深地解读 Go Zap 的高性能]( https://mp.weixin.qq.com/s/zqYNu2uTJe1UXiWvm98dOw )（介绍了代码结构，没有介绍高性能的根本原因）

---


* [go什么时候会触发gc]( https://mp.weixin.qq.com/s/e2-NXWCS0bd2BPWzdeiS_A)

---

* [go语言中的零值，他有什么用？](https://mp.weixin.qq.com/s/pVLs9mCOevEpQtbJVnWPbQ)

* [两个nil比较结果是什么？]( https://mp.weixin.qq.com/s/T-qmiyzlOx5T5S6Ca7X9aQ)

* [true != true？Go 面试官，你坑人！！！](  https://mp.weixin.qq.com/s/Xb0ZUUeOw-IgHwGVsCaycA)

* [问个 Go 问题，字符串 len == 0 和 字符串== "" ，有啥区别？]( https://mp.weixin.qq.com/s/rMygbfaLAF5NF206uEUJKA)

* [小技巧分享：在 Go 如何实现枚举？]( https://mp.weixin.qq.com/s/4jvhRQvKlMiYweSOG6xCrA) (go的iota实现是不完全的enum)

---
* [Goroutine 一泄露就看到他，这是个什么？](https://mp.weixin.qq.com/s/x6Kzn7VA1wUz7g8txcBX7A)

* [go切片导致内存泄漏、slice的data字段、边界取值] memory->slice

* [go map的初始化、访问、赋值、扩容、缩容]memory->map

* [go 的负载因子为什么是6.5](https://mp.weixin.qq.com/s/vxf7VxRcWL27ST2_VDHhOg)

* [聊一聊内存逃逸](https://mp.weixin.qq.com/s/J-WjYpZ40ehGLoJ0dDTWDw)

* [透过内存看slice和array的异同]

* [Go 数组比切片好在哪？]( https://mp.weixin.qq.com/s/zp1vdhGukEYKpzAdPt--Mw )（定长，可控的内存）


* [灵魂拷问 Go 语言：这个变量到底分配到哪里了？]( https://mp.weixin.qq.com/s/mFfza7DayFqsiS93Ep15BA)
  go build -gcflags '-m -l' main.go ;
  go tool compile -S main.go

* [搞 Go 要了解的 2 个 Header，你知道吗？]( https://mp.weixin.qq.com/s/rGqM1wMlqQEoJSgyrgZNLg) StringHeader 和 SliceHeader。

* [用 Go map 要注意这 1 个细节，避免依赖他！]( https://mp.weixin.qq.com/s/MzAktbjNyZD0xRVTPRKHpw) 输出是乱序的 rand随机

* [Go1.16 新特性：详解内存管理机制的变更，你需要了解]( https://mp.weixin.qq.com/s/l4oEJdskbWpff1E3tTNUxQ) madvise释放内存的策略


---
* [一文吃透 Go 语言解密之上下文 context]( https://mp.weixin.qq.com/s/A03G3_kCvVFN3TxB-92GVw)

* [面试官：context携带的数据是线程安全的吗？](https://mp.weixin.qq.com/s/7L2H3ulyTk4hXQFbFfa79A)

* [为什么 Go map 和 slice 是非线性安全的？]( https://mp.weixin.qq.com/s/TzHvDdtfp0FZ9y1ndqeCRw)  Go Slice 主要还是索引位覆写问题

* [Context 是怎么在 Go 语言中发挥关键作用的]( https://mp.weixin.qq.com/s/NNYyBLOO949ElFriLVRWiA)

* [一起聊聊 Go Context 的正确使用姿势]( https://mp.weixin.qq.com/s/5JDSqNIimNrgm5__Z4FNjw)

* [一文搞懂如何实现 Go 超时控制]( https://mp.weixin.qq.com/s/S4d9CJYmViJT8EbhyNCIMg)

---

* [读者提问：反射是如何获取结构体成员信息的？](https://mp.weixin.qq.com/s/BYVYhpP70gX4Vp1W9ckkMQ)

* [解密 Go 语言之反射 reflect]( https://mp.weixin.qq.com/s/onl3sBCSNs8l42uihi_p4A)  反射在工程实践中，目的一就是可以获取到值和类型，其二就是要能够修改他的值。；Elem 方法来获取指针所指向的源变量；反射本质上与 Interface 存在直接关系

---
    compile

* [内联函数和编译器对go代码的优化]( https://mp.weixin.qq.com/s/Or4FmVAf9nvMQzPct87Ecw)

* [终于识破这个go编译器把戏]( https://mp.weixin.qq.com/s/rbIIT26rFQzjVcfFnwso5Q)

* [翻译了一篇关于Go编译器的文章]( https://mp.weixin.qq.com/s/G7sVQNbgXmyeAT9ZI2q2OA)

* [迷惑了，Go len() 是怎么计算出来的？]( https://mp.weixin.qq.com/s/VId32MuVA3ZRvxAHBKHXJA)

* [一道关于 len 函数的诡异 Go 面试题解析](  https://mp.weixin.qq.com/s/FUNE8dI-yFArJF2KCNFCgw)

* [面试官：小松子知道什么是内联函数吗？]( https://mp.weixin.qq.com/s/TaiRDMt0yaG89meT0eaghw )
  将函数调用展开,避免了频繁调用函数对栈内存重复开辟所带来的消耗
  --gcflags=-m参数可以查看编译器的优化策略，传入--gcflags="-m -m"会查看更完整的优化策略！

Go在内部维持了一份内联函数的映射关系，会生成一个内联树，我们可以通过-gcflags="-d pctab=pctoinline"参数查看

---

* [单元测试] unitTest


---

* [文件存储] stroge

* [Go 存储基础 — “文件”被偷偷修改？来，给它装个监控！storage-->fsnofify](  https://mp.weixin.qq.com/s/Czv0CxDKqr2QNItO4aNZMA)

* [浅析gowatch监听文件变动实现原理 storage-->gowatch ](https://mp.weixin.qq.com/s/tsavVgnxFb6anLvcjvQwlA)

* [Go 存储基础 — 内存结构体怎么写入文件？storage->file ](https://mp.weixin.qq.com/s/mfNz7r76vZOOgiMSmuVeJA)

* [深入理解 Linux 的 epoll 机制]( https://mp.weixin.qq.com/s/GEoG23wz2JfQQQ9MgoM8tg) （IO 多路复用就是 1 个线程处理 多个 fd 的模式）

* [Linux fd 系列 — eventfd 是什么？]( https://mp.weixin.qq.com/s/9S1kYiDs6aVQXVtPY_fTBg)

* [自制文件系统 — 来看看文件系统的样子](https://mp.weixin.qq.com/s/7qq3AugMKqjlwx26PT20sw)

* [自制文件系统 —— Go实战：hello world 的文件系统]( https://mp.weixin.qq.com/s/oaxYWrlXaeu5mil4lNVbvg)

---

* [详解 Go 程序的启动流程，你知道 g0，m0 是什么吗？（Go 程序是怎么引导起来的）]( https://mp.weixin.qq.com/s/YK-TD3bZGEgqC0j-8U6VkQ)
  go build GOFLAGS="-ldflags=-compressdwarf=false"
  在命令中指定了 GOFLAGS 参数，这是因为在 Go1.11 起，为了减少二进制文件大小，调试信息会被压缩。
  导致在 MacOS 上使用 gdb 时无法理解压缩的 DWARF 的含义是什么

* [会诱发goroutine挂起的27个原因]( https://mp.weixin.qq.com/s/h1zrFLWoryA7P5I19kwkpg)

* [再见 Go 面试官：单核 CPU，开两个 Goroutine，其中一个死循环，会怎么样？]( https://mp.weixin.qq.com/s/h27GXmfGYVLHRG3Mu_8axw)

* [嗯，你觉得 Go 在什么时候会抢占 P？]( https://mp.weixin.qq.com/s/WAPogwLJ2BZvrquoKTQXzg)

* [跟读者聊 Goroutine 泄露的 N 种方法，真刺激！]( https://mp.weixin.qq.com/s/ql01K1nOnEZpdbp--6EDYw)  一直不能释放goroutine


* [技巧分享：多 Goroutine 如何优雅处理错误？]( https://mp.weixin.qq.com/s/NX6kVJP-RdUzcCmG2MF31w) sync/errgroup

* [详解并发编程包之 Go errgroup]( https://mp.weixin.qq.com/s/0_bV3DyrIqx5sph4sjNuUA)

* [回答我，停止 Goroutine 有几种方法？]( https://mp.weixin.qq.com/s/tN8Q1GRmphZyAuaHrkYFEg)

* [Go 群友提问：Goroutine 数量控制在多少合适，会影响 GC 和调度？]( https://mp.weixin.qq.com/s/uWP2X6iFu7BtwjIv5H55vw)  还是得看 Goroutine 里面跑的是什么东西。

* [go官方信号量库]Semaphore

---




---
    concurrent


* [Go 并发编程 — 结构体多字段的原子操作]( https://mp.weixin.qq.com/s/u5NKKqAtrJt-sgTM1iQ0gA)


---
    limiter

* [go官方限流器的详解]( https://mp.weixin.qq.com/s/S3_YEyhLcaAUuaSabXdimw)

* [常用限流算法的应用场景和实现原理]( https://mp.weixin.qq.com/s/krrUFEHVBw4c-47ziXOK2w)


* [go-monitor：服务质量统计分析警告工具]( https://mp.weixin.qq.com/s/WF_-XrzI3lS3tqmrzxMjdg)

---

* [Go 的相对路径问题 path]( https://mp.weixin.qq.com/s/QOA3Mk20M4rRM9oXOGB9HA)

* [面试官：你能聊聊string和[]byte的转换吗？bytes](  https://mp.weixin.qq.com/s/6vBreVLyPQc-WRBh_s90oA)

---
    debugs

* [编写和优化Go代码]( https://github.com/dgryski/go-perfbook/blob/master/performance-zh.md)

* [学会使用 GDB 调试 Go 代码](  https://mp.weixin.qq.com/s/O9Ngzgg9xfHMs5RSK5wHQQ)

* [一个 Demo 学会使用 Go Delve 调试]( https://mp.weixin.qq.com/s/Yz_p0S5N4ubf8wxLm5wbmQ)

* [Go 程序崩了？煎鱼教你用 PProf 工具来救火！]( https://mp.weixin.qq.com/s/9yLd7kkYzmbCriolhbvK_g)

* [Go 工程师必学：Go 大杀器之跟踪剖析 trace]( https://mp.weixin.qq.com/s/7DY0hDwidgx0zezP1ml3Ig)  (有时候单单使用 pprof 还不一定足够完整观查并解决问题，因为在真实的程序中还包含许多的隐藏动作。
  Goroutine 在执行时会做哪些操作？
  Goroutine 执行/阻塞了多长时间？
  Syscall 在什么时候被阻止？在哪里被阻止的？
  谁又锁/解锁了 Goroutine ？
  GC 是怎么影响到 Goroutine 的执行的？
  这些东西用 pprof 是很难分析出来的，但如果你又想知道上述的答案的话，你可以用本章节的主角 go tool trace 来打开新世界的大门)



* [必须要学的 Go 进程诊断工具 gops]( https://mp.weixin.qq.com/s/iS7R0NTZcTlonUw8bq0jKQ)

* [生产环境Go程序内存泄露，用pprof如何快速定位]( https://mp.weixin.qq.com/s/8UG7qJabqHi6yWARKkZsgA)


* [Golang Profiling: 关于 pprof]( https://mp.weixin.qq.com/s/YpUUj4xqlaZ9paEJe7VPYg)

* [Go 应用的性能优化](  https://xargin.com/go-perf-optimization/)

* [Go 语言中的一些非常规优化]( https://xargin.com/unusual-opt-in-go/)


* [注释竟然还有特殊用途？一文解惑 //go:linkname 指令]( https://mp.weixin.qq.com/s/_d1Q0Sx_KPrzEd4psPccMg)

* [我无语了，Go 中 +-*/ 四个运算符竟然可以连着用]( https://mp.weixin.qq.com/s/8GRq6At23fMho3BKkylcGw)



* [想要4个9？本文告诉你监控告警如何做]( https://mp.weixin.qq.com/s/qaNWBlDGgE2hNnu6SV4EBg)

* [我要提高 Go 程序健壮性，Fuzzing 来了！]( https://mp.weixin.qq.com/s/zdsrmlwVR0bP1Q_Xg_VlpQ) (Go 在 dev.fuzz 分支上提供了该功能的 Beta 测试 https://github.com/dvyukov/go-fuzz)


---

* [助力你成为优秀 Gopher 的 7 个Go库]dev

* [Go项目实战：从零构建一个并发文件下载器]( https://mp.weixin.qq.com/s/28CjAeINvlvNqLXP0g2oMw)

* [用 Go 来了解一下 Redis 通讯协议](https://mp.weixin.qq.com/s/pLwRiG1H_EAANadzz3VaLg ) （redis协议的组成）


* [一道 Go 闭包题，面试官说原来自己答错了：面别人也涨知识]( https://mp.weixin.qq.com/s/OLgsdhXGEMltmjcpTW2ICw) 闭包通过一个结构体来实现，它存储一个函数和一个关联的上下文环境。

* [Go函数闭包底层实现]( https://mp.weixin.qq.com/s/JsnuIyLy3XhQQuuxFIMzrA )变量逃逸


* [我这样升级 Go 版本，你呢？]( https://mp.weixin.qq.com/s/bGS5D0UYVp6BxSLjuZy0pg) (go的多版本)

* [又吵起来了，Go 是传值还是传引用？]( https://mp.weixin.qq.com/s/qsxvfiyZfRCtgTymO9LBZQ) （传递的是副本，值的副本，指针的副本，原指针和指针副本指向同一个数据地址;map 和 slice 的行为类似于指针，它们是包含指向底层 map 或 slice 数据的指针的描述符”）
* func makemap(t *maptype, hint int, h *hmap) *hmap {} 返回的是一个指针

* [Go 面试官问我如何实现面向对象？]( https://mp.weixin.qq.com/s/2x4Sajv7HkAjWFPe4oD96g) (封装、继承、多态：在 Go 语言中，多态是通过接口来实现的)

* [Go 面试官：什么是协程，协程和线程的区别和联系？]( https://mp.weixin.qq.com/s/vW5n_JWa3I-Qopbx4TmIgQ)

* [手撕 Go 面试官：Go 结构体是否可以比较，为什么？](https://mp.weixin.qq.com/s/HScH6nm3xf4POXVk774jUA)

* [用 Go struct 不能犯的一个低级错误！]( https://mp.weixin.qq.com/s/K5B2ItkzOb4eCFLxZI5Wvw) (空结构体，分配在栈(刻意优化)和堆(zerobase)上的不同处理方式)

* [详解 Go 空结构体strcut的 3 种使用场景]( https://mp.weixin.qq.com/s/zbYIdB0HlYwYSQRXFFpqSg) (Go 编译器在内存分配时做的优化)

* [你知道 Go 结构体和结构体指针调用有什么区别吗？]( https://mp.weixin.qq.com/s/g-D_eVh-8JaIoRne09bJ3Q)



* [Go 群友提问：进程、线程都有 ID，为什么 Goroutine 没有 ID？](https://mp.weixin.qq.com/s/qFAtgpbAsHSPVLuo3PYIhg)



* [生产环境遇到一个 Go 问题，整组人都懵逼了...]( https://mp.weixin.qq.com/s/F9II4xc4yimOCSTeKBDWqw) interface{}与nil的比较

* [Go 面试题：Go interface 的一个 “坑” 及原理分析]( https://mp.weixin.qq.com/s/vNACbdSDxC9S0LOAr7ngLQ)  interface包括类型和值

* [Go 面试题： new 和 make 是什么，差异在哪？]( https://mp.weixin.qq.com/s/tZg3zmESlLmefAWdTR96Tg) 主要用途都是用于分配相应类型的内存空间。 调用 make 函数去初始化切片（slice）的类型时，会带有零值，需要明确是否需要。

---

* [一文吃透 Go 语言解密之接口 interface]( https://mp.weixin.qq.com/s/vSgV_9bfoifnh2LEX0Y7cQ)

* [一文带你解密 Go 语言之通道 channel](https://mp.weixin.qq.com/s/ZXYpfLNGyej0df2zXqfnHQ) 当缓冲区满了后，发送者就会阻塞并等待。而当缓冲区为空时，接受者就会阻塞并等待，直至有新的数据：

---


* [项目实战：使用 Fiber + Gorm 构建 REST API]( https://mp.weixin.qq.com/s/TKphSzgM443DuO9KgZlgKw)

---

* [漫谈 MQ：要消息队列（MQ）有什么用？]( https://mp.weixin.qq.com/s/aN4VKhzmiqMF7a2GKI2ADQ)  解耦 削峰 异步

* [《漫谈 MQ》设计 MQ 的 3 个难点]( https://mp.weixin.qq.com/s/_QZ1mOtSFECab7TkvPePvQ) 高可用(水平扩展+配套服务：服务注册、发现机制、负载均衡) 高并发（队列划分，起到分而治之的作用） 高可靠（主要是针对消息发送、存储消息、处理消息这三块进行展开，和 MySQL 数据库的存储模式是有一定的神似之处）

---

* [上帝视角看 “Go 项目标准布局” 之争]( https://mp.weixin.qq.com/s/KnsB9cTGnM0X7hNR9VDzxg)  golang-standards/project-layout

---

* [干货满满的 Go Modules 知识分享] (https://mp.weixin.qq.com/s/uUNTH06_s6yzy5urtjPMsg)

* [最新提案：维持 GOPATH 的传统使用方式（Go1.17 移除 GOPATH）] (https://mp.weixin.qq.com/s/AzfKHfs6AOolxutdVpZibw)

* [Go1.16 新特性：Go mod 的后悔药，仅需这一招] (https://mp.weixin.qq.com/s/0g89yj9sc1oIz9kS9ZIAEA) retract



---

* [万字长文 | 从实践到原理，带你参透 gRPC] (https://mp.weixin.qq.com/s/o-K7G9ywCdmW7et6Q4WMeA) gRPC 基于 HTTP/2 标准设计，带来诸如双向流、流控、头部压缩、单 TCP 连接上的多复用请求等特性。这些特性使得其在移动设备上表现更好，更省电和节省空间占用。
  grpc.NewServer()；grpc.DialContext()

---

* [使用golang进行证书签发和双向认证](https://mp.weixin.qq.com/s/JtIWAyOPNgc08JSvqoFBmA)

* [这 Go 的边界检查，简直让人抓狂~](https://mp.weixin.qq.com/s/397sL-TCaZrOGR2-s1NFLw) 是 Go 语言中防止数组、切片越界而导致内存不安全的检查手段。 go build -gcflags="-d=ssa/check_bce/debug=1" main.go

* [边界检查消除] (https://gfw.go101.org/article/bounds-check-elimination.html)

* [一个活跃在众多 Go 项目中的编程模式](  https://mp.weixin.qq.com/s/dWY1ZzOl1TwpmM-rrF0m4Q)  函数式选项模式( Functional Options)。该模式解决的问题是，如何更动态灵活地为对象配置参数。



* [超全总结：Go 读文件的 10 种方法](https://mp.weixin.qq.com/s/ww27OPuD_Pse_KDNQWyjzA )

* [选择合适的 Go 字符串拼接方式]( https://mp.weixin.qq.com/s/BnJlP7co44__ZCl2lnSENw) 在Go语言中就提供了6种方式进行字符串拼接，那这几种拼接方式该如何选择呢？ 无论什么情况下使用strings.builder进行字符串拼接都是最高效的，不过要主要使用方法，记得调用grow进行容量分配，才会高效



* [在实现小工具的过程中学会 Go 反射]( https://mp.weixin.qq.com/s/6_zhqUB3aQr-s_ftTQTR_g)

* [Go 如何实现启动参数的加载]( https://mp.weixin.qq.com/s/NYlAXYdfA0g8JpSdpksPGg) os.Args 函数，获取命令行参数； runtime.argslice； flag 包

Go 汇编语言对 CPU 的重新抽象。Go汇编为了简化汇编代码的编写，引入了 PC、FP、SP、SB 四个伪寄存器。
四个伪寄存器加上其它的通用寄存器就是 Go 汇编语言对 CPU 的重新抽象。

* [写 Go 时如何优雅地查文档]( https://mp.weixin.qq.com/s/cCLKCPWEminsC1BJcaguSQ)

* [Go 的结构体标签]( https://mp.weixin.qq.com/s/4FmxImNLcU0-up5aVZLMzw)  
  由空格分隔;

  type User struct {
  Name string `json:"name" xml:"name"`
  }
  键，通常表示后面跟的“值”是被哪个包使用的，例如json这个键会被encoding/json包处理使用。如果要在“键”对应的“值”中传递多个信息，通常通过用逗号（'，'）分隔来指定，;

  Name string `json:"name,omitempty"`

按照惯例，如果一个字段的结构体标签里某个键的“值”被设置成了的破折号 ('-')，那么就意味着告诉处理该结构体标签键值的进程排除该字段。

    Name string `json:"-"`

* [线上实战:大内存 Go 服务性能优化]( https://mp.weixin.qq.com/s/SHcBZNO_t9dNOiWug3weSw)  good

* [应该如何去选择 Go router？]( https://mp.weixin.qq.com/s/OoZRkIVVK9Yz63NMYJ34tw)

* [如何保留 Go 程序崩溃现场]( https://mp.weixin.qq.com/s/RktnMydDtOZFwEFLLYzlCA) core dump 文件是操作系统提供给我们的一把利器，它是程序意外终止时产生的内存快照

* [如何有效控制 Go 线程数？]( https://mp.weixin.qq.com/s/HYcHfKScBlYCD0IUd0t4jA) 如果真的存在线程数暴涨的问题，那么你应该思考代码逻辑是否合理（为什么你能允许短时间内如此多的系统同步调用），是否可以做一些例如限流之类的处理。





* [含有CGO代码的项目如何实现跨平台编译]( https://mp.weixin.qq.com/s/Xd-YuN-v2OWIFO2wrpruCA)

* [Go 如何利用 Linux 内核的负载均衡能力](  https://mp.weixin.qq.com/s/lnOTaraGKINxaqbrMHPP5Q) socket五元组 ;linux 内核自 3.9 提供的 SO_REUSEPORT 选项，可以让多进程监听同一个端口。

* [SO_REUSEPORT学习笔记](  http://www.blogjava.net/yongboy/archive/2015/02/12/422893.html )

---

* [golang 垃圾回收 （一）概述篇](https://mp.weixin.qq.com/s/GYYLLlVWMoI-ls8IgrzndA)

* [golang 垃圾回收（二）屏障技术](https://mp.weixin.qq.com/s/z0Pt0gUUsHfJGAhMVg4vuQ) 写屏障确保在 GC 运行时正确跟踪新的写入（这样它们就不会被意外释放或保留）。

* [golang 垃圾回收 - 删除写屏障]( https://mp.weixin.qq.com/s/T8HvENFlkKuEm2U7rbZTzg)

* [通过 eBPF 深入探究 Go GC]( https://mp.weixin.qq.com/s/gBhxNwLmdQjmB87y6qOvBg  )

---





netFD、poll.FD、pollDesc（这三个数据结构可以理解为对操作系统接口调用的层层封装）。


* [几个秒杀 Go 官方库的第三方开源库](https://mp.weixin.qq.com/s/JRsstunuD2UClWb237kPTQ) fasthttp；jsoniter；gogo/protobuf；valyala/quicktemplate （它们的重点都是优化对应官方库的性能问题）

* [fasthttp 快在哪里](  https://xargin.com/why-fasthttp-is-fast-and-the-cost-of-it/)

---





* [学会这几招让 Go 程序自己监控自己 ]( https://mp.weixin.qq.com/s/H-eCNw7s4e3oz2ReI6Hu_A) （在宿主机、虚拟机、容器获取性能指标 https://github.com/shirou/gopsutil）

* [如何让 Go 程序自动采样](  https://mp.weixin.qq.com/s/0KL9r4osbFwRQTKcscARDg) 判断采样时间点的规则

* [无人值守的自动 dump（一）]( https://mp.weixin.qq.com/s/2nbyWSZMT1HzvYAoaeWK_A)

* [无人值守的自动 dump（二）]( https://mp.weixin.qq.com/s/wKpTiyc1VkZQy0-J8x519g)

---

* [go-swagger源码解析]( https://zhuanlan.zhihu.com/p/294069197)


---
    wbsocket

*[Golang Websocket 实践](  https://mp.weixin.qq.com/s/wZVkWLswzN3YtSdZMXF1jg)

---
    redis
* [Golang使用redigo实现redis的分布式锁](http://xiaorui.cc/?p=3028)
* [Golang使用redis protocol实现pubsub通信](http://xiaorui.cc/?p=4847)
* [golang基于redis lua封装的优先级去重队列](http://xiaorui.cc/?p=4828)
* [Golang基于redis实现的分布式信号量(semaphore)](http://xiaorui.cc/?p=4822)
* [golang redigo lua解决性能问题]( http://xiaorui.cc/?p=4737)

  monkey补丁
* [使用monkey补丁替换golang的标准库]( http://xiaorui.cc/?p=5128)
* [通过火焰图排查golang json的性能问题](http://xiaorui.cc/?p=5108)

  crontab
* [开源golang兼容crontab的定时任务管理器](http://xiaorui.cc/?p=5625)


    log
* [Golang logrus的高级配置(hook, logrotate)]( http://xiaorui.cc/?p=4963)
* [使用golang log库包实现日志文件输出](http://xiaorui.cc/?p=2972)


    etcd:

    分布式锁：
    etcd
    redis redlock
    consul
* [高可用分布式存储 etcd 的实现原理](https://draveness.me/etcd-introduction)
* [源码分析golang consul分布式锁lock delay问题](http://xiaorui.cc/2019/05/19/%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90golang-consul%E5%88%86%E5%B8%83%E5%BC%8F%E9%94%81lock-delay%E9%97%AE%E9%A2%98/)
* [分布式一致性raft算法实现原理](http://xiaorui.cc/2016/07/08/%e6%8a%80%e6%9c%af%e5%88%86%e4%ba%ab-%e3%80%8a%e5%88%86%e5%b8%83%e5%bc%8f%e4%b8%80%e8%87%b4%e6%80%a7%e7%ae%97%e6%b3%95%e5%ae%9e%e7%8e%b0%e5%8e%9f%e7%90%86%e3%80%8b/)
* [Golang使用redigo实现redis的分布式锁](http://xiaorui.cc/?p=3028)
* [etcd技术内幕]-百里~
-------
    内存分配
* [图解Go语言内存分配](https://zhuanlan.zhihu.com/p/59125443)
* [strace分析追踪malloc申请内存过程](http://xiaorui.cc/?p=5334)
* [go内存分配那些事，就这么简单](https://www.cnblogs.com/shijingxiang/articles/11466957.html)
* go内存分配器可视化指南(go语言中文网)
* [TCMalloc:Thread-Caching Malloc](http://goog-perftools.sourceforge.net/doc/tcmalloc.html)
* [Golang源码探索(三) GC的实现原理](https://www.cnblogs.com/zkweb/p/7880099.html)
* [第九章 虚拟内存（深入理解计算机系统）]
* [内存分配器](https://draveness.me/golang/)
* [栈空间管理](https://draveness.me/golang/)

-------
    内存泄漏
* [一起 goroutine 泄漏问题的排查](https://zhuanlan.zhihu.com/p/100740270)
* [实战Go内存泄露](http://lessisbetter.site/2019/05/18/go-goroutine-leak/)
* [slice类型内存泄漏的逻辑(曹大)](https://xargin.com/logic-of-slice-memory-leak/)
* [分析golang time.After引起内存暴增OOM问题](http://xiaorui.cc/?p=5745)
* [探究golang的channel和map内存释放问题](http://xiaorui.cc/?p=5450)

-------
    连接池
* [使用golang协程池控制并发请求](http://xiaorui.cc/2019/05/24/%e4%bd%bf%e7%94%a8golang%e5%8d%8f%e7%a8%8b%e6%b1%a0%e6%8e%a7%e5%88%b6%e5%b9%b6%e5%8f%91%e8%af%b7%e6%b1%82/)
* [golang通用自定义连接池的实现](http://xiaorui.cc/?p=5434)
* [解决golang redis连接池的io异常BUG?](http://xiaorui.cc/?p=5513)
* [深入研究golang net/http连接池可用性](http://xiaorui.cc/?p=5056)
* [golang grpc网关使用连接池提吞吐量](http://xiaorui.cc/2019/08/13/golang-grpc%e7%bd%91%e5%85%b3%e7%94%a8%e8%bf%9e%e6%8e%a5%e6%b1%a0%e6%8f%90%e9%ab%98%e5%90%9e%e5%90%90%e9%87%8f/)

-------
    channel
* [深度解密Go语言之channel ](https://zhuanlan.zhihu.com/p/74613114)
* [Golang并发：再也不愁选channel还是选锁](http://lessisbetter.site/2019/01/14/golang-channel-and-mutex/)
* [channel-draveness](https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-channel/)
* [golang channel提前close丢失数据?](http://xiaorui.cc/?p=5007)

-------
    锁🔐
* [golang多场景下RwMutex和mutex锁性能对比](http://xiaorui.cc/?p=5611)
* [golang log日志里为什么需要加锁?](http://xiaorui.cc/?p=5195)
* [通过golang goroutine stack分析死锁问题](http://xiaorui.cc/?p=5160)
* [扩展golang的sync mutex的trylock及islocked](http://xiaorui.cc/?p=5084)
* [golang新版如何优化sync.pool锁竞争消耗？](http://xiaorui.cc/?p=5878﻿)

-------
    context
* [深度解密Go语言之context](https://zhuanlan.zhihu.com/p/68792989)
* [上下文context](https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-context/)
* [源码分析context的超时及关闭实现](http://xiaorui.cc/?p=5604)

-------
    map
* [深度解密Go语言之 map](https://zhuanlan.zhihu.com/p/66676224)
* [map并发崩溃一例(非线程安全)(曹大)](https://xargin.com/map-concurrent-throw/)

-------
    scheduler
* [深度解密Go语言之 scheduler](https://zhuanlan.zhihu.com/p/80853548)
* [调度器-draveness](https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-goroutine/)
* [golang密集场景下协程调度饥饿问题](http://xiaorui.cc/?p=5251)

-------
    error
* [Golang error 的突围](https://zhuanlan.zhihu.com/p/82985617)

-------
    pprof
* [深度解密Go语言之 pprof](https://zhuanlan.zhihu.com/p/91241270)
* [golang pprof分析net/http的性能瓶颈](http://xiaorui.cc/?p=5577)
* [通过火焰图排查golang json的性能问题](http://xiaorui.cc/?p=5108)
* [Golang使用pprof监控性能及GC调优](http://xiaorui.cc/?p=3000)

-------
    内存重排
* [曹大谈内存重排](https://zhuanlan.zhihu.com/p/69414216)

-------
    unsafe
* [深度解密Go语言之unsafe](https://zhuanlan.zhihu.com/p/67852800)

-------
    reflect
* [深度解密Go语言之反射](https://zhuanlan.zhihu.com/p/64884660)
* [反射-dravness](https://draveness.me/golang/docs/part2-foundation/ch04-basic/golang-reflect/)

-------
    slice
* [深度解密Go语言之Slice](https://zhuanlan.zhihu.com/p/61121325)


-------
    逃逸分析：栈与堆
* [Golang之变量去哪儿？](https://zhuanlan.zhihu.com/p/58065429)

-------
    defer
* [Golang之轻松化解defer的温柔陷阱](https://zhuanlan.zhihu.com/p/56557423)
* [defer-draveness](https://draveness.me/golang/docs/part2-foundation/ch05-keyword/golang-defer/)

-------
    interface
* [深度解密Go语言之关于 interface 的10个问题](https://zhuanlan.zhihu.com/p/63649977)
* [接口-draveness](https://draveness.me/golang/docs/part2-foundation/ch04-basic/golang-interface/)

-------
    编译-链接-运行
* [Go 程序是怎样跑起来的](https://zhuanlan.zhihu.com/p/71993748)

-------

    sync
* [Golang并发的次优选择：sync包](http://lessisbetter.site/2019/01/04/golang-pkg-sync/)
* [golang新版如何优化sync.pool锁竞争消耗？](http://xiaorui.cc/?p=5878﻿)
* [go sync.pool []byte导致grpc解包异常](http://xiaorui.cc/?p=5969)
* [扩展go sync.map的length和delete方法](http://xiaorui.cc/?p=4972)

-------
    select
* [Golang并发模型：轻松入门select](http://lessisbetter.site/2018/12/13/golang-slect/)
* [Golang并发模型：select进阶](http://lessisbetter.site/2018/12/17/golang-selete-advance/)
* [select](https://draveness.me/golang/docs/part2-foundation/ch05-keyword/golang-select/)
* [Golang利用select实现goroutine的超时控制](http://xiaorui.cc/?p=2997)

-------
    make  new
* [make 和 new](https://draveness.me/golang/docs/part2-foundation/ch05-keyword/golang-make-and-new/)

------
    time
* [定时器-draveness](https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-timer/)
* [完全兼容golang定时器的高性能时间轮实现(go-timewheel)](http://xiaorui.cc/2019/09/27/%e5%85%bc%e5%ae%b9golang-time%e5%ae%9a%e6%97%b6%e5%99%a8%e7%9a%84%e6%97%b6%e9%97%b4%e8%bd%ae%e5%ae%9e%e7%8e%b0/)
* [源码分析go time.timer和ticker的stop问题](http://xiaorui.cc/2019/09/09/%e6%ba%90%e7%a0%81%e5%88%86%e6%9e%90go-time-timer%e5%92%8cticker%e7%9a%84stop%e9%97%ae%e9%a2%98/)
* [分析golang time.After引起内存暴增OOM问题](http://xiaorui.cc/?p=5745)
* [分析golang定时器CPU使用率高的现象](http://xiaorui.cc/?p=5117)
* [golang随机time.sleep的Duration问题](http://xiaorui.cc/?p=3034)

-------
    函数
* [函数调用-draveness](https://draveness.me/golang/docs/part2-foundation/ch04-basic/golang-function-call/)

-------
* [关于golang的panic recover异常错误处理](http://xiaorui.cc/?p=2909)

-------
    array
* [数组-draveness](https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-array/)

-------
    slice
* [slice-draveness](https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-array-and-slice/)

-------
    struct
* [一个空struct的坑（曹大）](https://xargin.com/addr-of-empty-struct-may-not-eq/)

-------
    MPG
* [为什么 Go 模块在下游服务抖动恢复后，CPU 占用无法恢复(曹大)](https://xargin.com/cpu-idle-cannot-recover-after-peak-load/)
* [disk io引起golang线程数暴涨的问题](http://xiaorui.cc/?p=5171)

-------
    gomaxprocs
* [golang gomaxprocs调高引起调度性能损耗](http://xiaorui.cc/2020/01/11/golang-gomaxprocs%e8%b0%83%e9%ab%98%e5%bc%95%e8%b5%b7%e8%b0%83%e5%ba%a6%e6%80%a7%e8%83%bd%e6%8d%9f%e8%80%97/)

-------
    runtime
* [runtime.stack加锁引起高时延及阻塞](http://xiaorui.cc/2020/01/03/go-runtime-stack%e5%8a%a0%e9%94%81%e5%bc%95%e8%b5%b7%e9%ab%98%e6%97%b6%e5%bb%b6%e5%8f%8a%e9%98%bb%e5%a1%9e/)
* [万字长文深入浅出 Golang Runtime](https://zhuanlan.zhihu.com/p/95056679)
* [系统监控](https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-sysmon/)
-------
    grpc

* [golang grpc网关使用连接池提吞吐量](http://xiaorui.cc/2019/08/13/golang-grpc%e7%bd%91%e5%85%b3%e7%94%a8%e8%bf%9e%e6%8e%a5%e6%b1%a0%e6%8f%90%e9%ab%98%e5%90%9e%e5%90%90%e9%87%8f/)


-------
    http
* [源码解析http shutdown优雅退出的原理](http://xiaorui.cc/?p=5803)

-------
    Grpc protobuf
* [Grpc protobuf的动态加载及类型反射实战](http://xiaorui.cc/2019/04/01/grpc-protobuf%e7%9a%84%e5%8a%a8%e6%80%81%e5%8a%a0%e8%bd%bd%e5%8f%8a%e7%b1%bb%e5%9e%8b%e5%8f%8d%e5%b0%84%e5%ae%9e%e6%88%98/)
* [深入 ProtoBuf - 简介](https://www.jianshu.com/p/a24c88c0526a)

-------
    database/sql
* [分析golang sql连接池大量的time wait问题](http://xiaorui.cc/?p=5771)

-------
    udp
* [关于golang udp的高性能优化](http://xiaorui.cc/?p=5684)

-------

[Go 使用场景和用 Go 的顶级公司]( https://mp.weixin.qq.com/s/Y1Dr3UykTvWuzjNKr-NMTg)







