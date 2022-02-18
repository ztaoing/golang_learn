【以下内容来自煎鱼的微信公众号】

[Go: A Documentary 发布！] https://mp.weixin.qq.com/s/5MtBE8vecKPOmRUYu2E-lg


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

[Go 缓冲系列之-free-cache] https://mp.weixin.qq.com/s/VmPIW6HhVrDyeADiRmkC_Q （也是使用减小锁的粒度、go 1.5版本之后，如果你使用的map的key和value中都不包含指针，那么GC会忽略这个map。）
    

[runtime]
* 什么是go runtime.KeepAlive
go 官方文档: https://pkg.go.dev/runtime#KeepAlive
文档: https://medium.com/a-journey-with-go/go-keeping-a-variable-alive-c28e3633673a
  
[编程思考：对象生命周期的问题] https://mp.weixin.qq.com/s/Hoy9cqHe9RZqEA5T9Dys5w
  


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



[多路复用] channel--》multiplex
[atomic.Value为什么不加锁也能保证数据线程安全]

[go中的零值，它有什么作用？] 官方：https://golang.org/ref/spec#the_zero_value
布尔型为false；数字型为0；字符串型为""；指针、函数、接口、切片、通道和映射都为nil

[go是如何实现启动参数的加载的？]compile


[select机制] select




    

[在go容器里设置gomaxprocs的正确姿势：]dockers

---

[unsafe包]unsafe
[详解 Go 团队不建议用的 unsafe.Pointer] https://mp.weixin.qq.com/s/8qtHdw2JiRQ1cXlzbJ0ANA

---


[go网络编程和tcp抓包实操] network-》getTCPPackage
[go中如何强制关闭tcp连接] network-》getTCPPackage

[网友很强大，发现了Go并发下载的Bug] https://mp.weixin.qq.com/s/Kd4np83CKEOLQ3K0eWxlKg

[连接一个ip不存在的主机时，握手过程是怎样的？]network-》ConnIP
    连接一个ip不存在的主机时，握手过程是怎样的？
    连接一个IP地址存在但是端口不存在的主机时，握手过程是怎样的？
[context使用不当引发的一个bug]


---

[从CPU角度理解go中的结构体内存对齐]memory-》align

[详解 Go 内存对齐]memory-》align

[Go程序内存分配过多？] https://mp.weixin.qq.com/s/zBHPYJWnGf67Ex8i4cV8Eg

[Go 编程怎么也有踩内存？]  https://mp.weixin.qq.com/s/tXAP8_U63QLNj1h0ZMvXPw

[Go 内存泄露之痛，这篇把 Go timer.After 问题根因讲透了！] https://mp.weixin.qq.com/s/KSBdPkkvonSES9Z9iggElg

[为什么 Go 占用那么多的虚拟内存？] https://mp.weixin.qq.com/s/eVHK_ey8SgqEtl8v_Nurxg


---


[go五种原子性操作的用法详解] memory-》atomic-》cas +atomicMutex
原子性：外界不会看到只执行到一半的状态！
CPU执行一些列操作时不可能不发生中断，但是如果我们在执行多个操作时，能让他们的中间状态对外不可见，
那我们就可以说他拥有了"不可分割"的原子性
Go语言通过内置包sync/atomic提供了对原子操作的支持，其提供的原子操作有以下几大类：
1、增减，操作的方法名方式为AddXXXType，保证对操作数进行原子的增减，支持的类型为int32、int64、uint32、uint64、uintptr，使用时以实际类型替换前面我说的XXXType就是对应的操作方法。
2、载入，保证了读取到操作数前没有其他任务对它进行变更，操作方法的命名方式为LoadXXXType，支持的类型除了基础类型外还支持Pointer，也就是支持载入任何类型的指针。
3、存储，有载入了就必然有存储操作，这类操作的方法名以Store开头，支持的类型跟载入操作支持的那些一样。
4、比较并交换，也就是CAS （Compare And Swap），像Go的很多并发原语实现就是依赖的CAS操作，同样是支持上面列的那些类型。
5、交换，这个简单粗暴一些，不比较直接交换，这个操作很少会用。

互斥锁和院系操作的区别：
1、使用目的：互斥锁是用来保护一段逻辑，原子操作用于对一个变量的更新保护。
2、底层实现：Mutex由操作系统的调度器实现，而atomic包中的原子操作则由底层硬件指令直接提供支持，这些指令在执行的过程中是不允许中断的，
因此原子操作可以在lock-free的情况下保证并发安全，并且它的性能也能做到随CPU个数的增多而线性扩展。

[Go 并发编程 — 深入浅出 sync.Pool] https://mp.weixin.qq.com/s/1hLgu2XBBJkLzvI6pOR80g

[一口气搞懂 Go sync.map 所有知识点] https://mp.weixin.qq.com/s/8aufz1IzElaYR43ccuwMyA
---


[i/o timeout net/http的坑] network->timeout
[go的io库如何选择]network->io库
[Go语言常用文件操作汇总]network->func

[怎么使用 direct io？]io-》io.md
[浅析 Go IO 的知识框架]io
[深入理解 Linux 的 epoll 机制] https://mp.weixin.qq.com/s/GEoG23wz2JfQQQ9MgoM8tg

---
[go精妙的互斥锁设计]lock
[golang的位运算]lock


---

[go程序错误处理的一些建议]errors
[对go错误处理的4个误解]errors
[go的panic的三种诞生方式]errors
[go的panic的秘密都在这里]errors
[Go 错误处理：用 panic 取代 err != nil 的模式] https://mp.weixin.qq.com/s/p77V3_LkREuXPVLdebmmmQ

[你考虑过defer close的风险吗？]errors

[说好 defer 在 return 之后执行，为什么结果却不是？] https://mp.weixin.qq.com/s/oP90maykSzMXjdnudOKdSw

[使用 Go defer 要小心这 2 个雷区！] https://mp.weixin.qq.com/s/ZEsWa4xUb0a7tWemVZMXVw

[Go 群友提问：学习 defer 时很懵逼，这道不会做！] https://mp.weixin.qq.com/s/lELMqKho003h0gfKkZxhHQ

[Go 中的 error 居然可以这样封装]  https://mp.weixin.qq.com/s/-X8MKIQFRXmENsdwyRXcCg

[如何保存go程序崩溃的现场] errors 

[详解并发编程包之 Go errgroup] https://mp.weixin.qq.com/s/0_bV3DyrIqx5sph4sjNuUA

[一文带你由浅入深地解读 Go Zap 的高性能] https://mp.weixin.qq.com/s/zqYNu2uTJe1UXiWvm98dOw
---


[go什么时候会触发gc]gc

---

[go语言中的零值，他有什么用？]zero
[两个nil比较结果是什么？]zero
[true != true？Go 面试官，你坑人！！！]zero  https://mp.weixin.qq.com/s/Xb0ZUUeOw-IgHwGVsCaycA

[问个 Go 问题，字符串 len == 0 和 字符串== "" ，有啥区别？] https://mp.weixin.qq.com/s/rMygbfaLAF5NF206uEUJKA

[小技巧分享：在 Go 如何实现枚举？] https://mp.weixin.qq.com/s/4jvhRQvKlMiYweSOG6xCrA

---
[Goroutine 一泄露就看到他，这是个什么？]memory->gopark
[go切片导致内存泄漏、slice的data字段、边界取值] memory->slice
[go map的初始化、访问、赋值、扩容、缩容]memory->map
[go 的负载因子为什么是6.5]memory->map
[聊一聊内存逃逸]memory->
[透过内存看slice和array的异同]

[Go 数组比切片好在哪？] https://mp.weixin.qq.com/s/zp1vdhGukEYKpzAdPt--Mw

[灵魂拷问 Go 语言：这个变量到底分配到哪里了？] https://mp.weixin.qq.com/s/mFfza7DayFqsiS93Ep15BA

[搞 Go 要了解的 2 个 Header，你知道吗？] https://mp.weixin.qq.com/s/rGqM1wMlqQEoJSgyrgZNLg

[用 Go map 要注意这 1 个细节，避免依赖他！] https://mp.weixin.qq.com/s/MzAktbjNyZD0xRVTPRKHpw

[Go1.16 新特性：详解内存管理机制的变更，你需要了解] https://mp.weixin.qq.com/s/l4oEJdskbWpff1E3tTNUxQ


---
[一文吃透 Go 语言解密之上下文 context] https://mp.weixin.qq.com/s/A03G3_kCvVFN3TxB-92GVw

[面试官：context携带的数据是线程安全的吗？]contexts

[为什么 Go map 和 slice 是非线性安全的？] https://mp.weixin.qq.com/s/TzHvDdtfp0FZ9y1ndqeCRw

[Context 是怎么在 Go 语言中发挥关键作用的]contexts

[一起聊聊 Go Context 的正确使用姿势] https://mp.weixin.qq.com/s/5JDSqNIimNrgm5__Z4FNjw

[一文搞懂如何实现 Go 超时控制] https://mp.weixin.qq.com/s/S4d9CJYmViJT8EbhyNCIMg
---

[读者提问：反射是如何获取结构体成员信息的？]reflects

[解密 Go 语言之反射 reflect] https://mp.weixin.qq.com/s/onl3sBCSNs8l42uihi_p4A

---


[内联函数和编译器对go代码的优化]compile
[终于识破这个go编译器把戏]compile
[翻译了一篇关于Go编译器的文章]compile
[迷惑了，Go len() 是怎么计算出来的？]compile->slice
[一道关于 len 函数的诡异 Go 面试题解析] compile

[面试官：小松子知道什么是内联函数吗？] https://mp.weixin.qq.com/s/TaiRDMt0yaG89meT0eaghw

---

[单元测试] unitTest


---

[文件存储] stroge
[Go 存储基础 — “文件”被偷偷修改？来，给它装个监控！] storage-->fsnofify
[浅析gowatch监听文件变动实现原理]storage-->gowatch
[Go 存储基础 — 内存结构体怎么写入文件？]storage->file

[Linux fd 系列 — eventfd 是什么？] https://mp.weixin.qq.com/s/9S1kYiDs6aVQXVtPY_fTBg

[自制文件系统 — 来看看文件系统的样子]https://mp.weixin.qq.com/s/7qq3AugMKqjlwx26PT20sw

[自制文件系统 —— Go实战：hello world 的文件系统] https://mp.weixin.qq.com/s/oaxYWrlXaeu5mil4lNVbvg

---


[会诱发goroutine挂起的27个原因]
[嗯，你觉得 Go 在什么时候会抢占 P？] https://mp.weixin.qq.com/s/WAPogwLJ2BZvrquoKTQXzg

[跟读者聊 Goroutine 泄露的 N 种方法，真刺激！] https://mp.weixin.qq.com/s/ql01K1nOnEZpdbp--6EDYw

[详解 Go 程序的启动流程，你知道 g0，m0 是什么吗？] https://mp.weixin.qq.com/s/YK-TD3bZGEgqC0j-8U6VkQ

[再见 Go 面试官：单核 CPU，开两个 Goroutine，其中一个死循环，会怎么样？] https://mp.weixin.qq.com/s/h27GXmfGYVLHRG3Mu_8axw

[技巧分享：多 Goroutine 如何优雅处理错误？] https://mp.weixin.qq.com/s/NX6kVJP-RdUzcCmG2MF31w

[回答我，停止 Goroutine 有几种方法？] https://mp.weixin.qq.com/s/tN8Q1GRmphZyAuaHrkYFEg

[Go 群友提问：Goroutine 数量控制在多少合适，会影响 GC 和调度？] https://mp.weixin.qq.com/s/uWP2X6iFu7BtwjIv5H55vw

[go官方信号量库]Semaphore

---

[生产环境遇到一个 Go 问题，整组人都懵逼了...] https://mp.weixin.qq.com/s/F9II4xc4yimOCSTeKBDWqw


---



[Go 并发编程 — 结构体多字段的原子操作]concurrent


---

[go官方限流器的详解]limiter
[常用限流算法的应用场景和实现原理]limiter


[go-monitor：服务质量统计分析警告工具]monitor

---
[Go 的相对路径问题]path

[面试官：你能聊聊string和[]byte的转换吗？] bytes

---
[学会使用 GDB 调试 Go 代码] debugs
[一个 Demo 学会使用 Go Delve 调试]debugs
[Go 工程师必学：Go 大杀器之跟踪剖析 trace] https://mp.weixin.qq.com/s/7DY0hDwidgx0zezP1ml3Ig

[Go 程序崩了？煎鱼教你用 PProf 工具来救火！] https://mp.weixin.qq.com/s/9yLd7kkYzmbCriolhbvK_g

[必须要学的 Go 进程诊断工具 gops] https://mp.weixin.qq.com/s/iS7R0NTZcTlonUw8bq0jKQ

[生产环境Go程序内存泄露，用pprof如何快速定位] https://mp.weixin.qq.com/s/8UG7qJabqHi6yWARKkZsgA


[注释竟然还有特殊用途？一文解惑 //go:linkname 指令] https://mp.weixin.qq.com/s/_d1Q0Sx_KPrzEd4psPccMg

[我无语了，Go 中 +-*/ 四个运算符竟然可以连着用] https://mp.weixin.qq.com/s/8GRq6At23fMho3BKkylcGw

[go程序自己监控自己]
[想要4个9？本文告诉你监控告警如何做] https://mp.weixin.qq.com/s/qaNWBlDGgE2hNnu6SV4EBg


---

[助力你成为优秀 Gopher 的 7 个Go库]dev
[Go项目实战：从零构建一个并发文件下载器] https://mp.weixin.qq.com/s/28CjAeINvlvNqLXP0g2oMw

[用 Go 来了解一下 Redis 通讯协议] https://mp.weixin.qq.com/s/pLwRiG1H_EAANadzz3VaLg

[我要提高 Go 程序健壮性，Fuzzing 来了！] https://mp.weixin.qq.com/s/zdsrmlwVR0bP1Q_Xg_VlpQ

[一道 Go 闭包题，面试官说原来自己答错了：面别人也涨知识] https://mp.weixin.qq.com/s/OLgsdhXGEMltmjcpTW2ICw

[我这样升级 Go 版本，你呢？] https://mp.weixin.qq.com/s/bGS5D0UYVp6BxSLjuZy0pg

[又吵起来了，Go 是传值还是传引用？] https://mp.weixin.qq.com/s/qsxvfiyZfRCtgTymO9LBZQ

[Go 面试官问我如何实现面向对象？] https://mp.weixin.qq.com/s/2x4Sajv7HkAjWFPe4oD96g

[Go 面试官：什么是协程，协程和线程的区别和联系？] https://mp.weixin.qq.com/s/vW5n_JWa3I-Qopbx4TmIgQ

[用 Go struct 不能犯的一个低级错误！] https://mp.weixin.qq.com/s/K5B2ItkzOb4eCFLxZI5Wvw

[详解 Go 空结构体strcut的 3 种使用场景] https://mp.weixin.qq.com/s/zbYIdB0HlYwYSQRXFFpqSg

[你知道 Go 结构体和结构体指针调用有什么区别吗？] https://mp.weixin.qq.com/s/g-D_eVh-8JaIoRne09bJ3Q

[手撕 Go 面试官：Go 结构体是否可以比较，为什么？] https://mp.weixin.qq.com/s/HScH6nm3xf4POXVk774jUA

[Go 群友提问：进程、线程都有 ID，为什么 Goroutine 没有 ID？] https://mp.weixin.qq.com/s/qFAtgpbAsHSPVLuo3PYIhg

[Go 面试题：Go interface 的一个 “坑” 及原理分析] https://mp.weixin.qq.com/s/vNACbdSDxC9S0LOAr7ngLQ

[一文吃透 Go 语言解密之接口 interface] https://mp.weixin.qq.com/s/vSgV_9bfoifnh2LEX0Y7cQ

[Go 面试题： new 和 make 是什么，差异在哪？] https://mp.weixin.qq.com/s/tZg3zmESlLmefAWdTR96Tg

[一文带你解密 Go 语言之通道 channel] https://mp.weixin.qq.com/s/ZXYpfLNGyej0df2zXqfnHQ

[项目实战：使用 Fiber + Gorm 构建 REST API] https://mp.weixin.qq.com/s/TKphSzgM443DuO9KgZlgKw
---
[《漫谈 MQ》设计 MQ 的 3 个难点] https://mp.weixin.qq.com/s/_QZ1mOtSFECab7TkvPePvQ

---
[上帝视角看 “Go 项目标准布局” 之争] https://mp.weixin.qq.com/s/KnsB9cTGnM0X7hNR9VDzxg

[最新提案：维持 GOPATH 的传统使用方式（Go1.17 移除 GOPATH）] https://mp.weixin.qq.com/s/AzfKHfs6AOolxutdVpZibw

[Go1.16 新特性：Go mod 的后悔药，仅需这一招] https://mp.weixin.qq.com/s/0g89yj9sc1oIz9kS9ZIAEA

[干货满满的 Go Modules 知识分享] https://mp.weixin.qq.com/s/uUNTH06_s6yzy5urtjPMsg

---
[万字长文 | 从实践到原理，带你参透 gRPC] https://mp.weixin.qq.com/s/o-K7G9ywCdmW7et6Q4WMeA






