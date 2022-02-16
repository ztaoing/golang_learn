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

[多路复用] channel--》multiplex
[atomic.Value为什么不加锁也能保证数据线程安全]

[go中的零值，它有什么作用？] 官方：https://golang.org/ref/spec#the_zero_value
布尔型为false；数字型为0；字符串型为""；指针、函数、接口、切片、通道和映射都为nil

[go是如何实现启动参数的加载的？]compile


[select机制] select


[如何保存go程序崩溃的现场] errors

    

[在go容器里设置gomaxprocs的正确姿势：]dockers


[unsafe包]unsafe

[go网络编程和tcp抓包实操] network-》getTCPPackage
[go中如何强制关闭tcp连接] network-》getTCPPackage

[连接一个ip不存在的主机时，握手过程是怎样的？]network-》ConnIP
    连接一个ip不存在的主机时，握手过程是怎样的？
    连接一个IP地址存在但是端口不存在的主机时，握手过程是怎样的？
[context使用不当引发的一个bug]

[怎么使用 direct io？]io-》io.md

---

[从CPU角度理解go中的结构体内存对齐]memory-》align
[详解 Go 内存对齐]memory-》align

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


---
[透过内存看slice和array的异同]

[i/o timeout net/http的坑] network->timeout
[go的io库如何选择]network->io库
[Go语言常用文件操作汇总]network->func

---
[go精妙的互斥锁设计]lock
[golang的位运算]lock


---

[go程序错误处理的一些建议]errors
[对go错误处理的4个误解]errors
[go的panic的三种诞生方式]errors
[go的panic的秘密都在这里]errors
[你考虑过defer close的风险吗？]errors

---


[go什么时候会触发gc]gc

---

[go语言中的零值，他有什么用？]zero
[两个nil比较结果是什么？]zero
[true != true？Go 面试官，你坑人！！！]zero

---
[Goroutine 一泄露就看到他，这是个什么？]memory->gopark
[go切片导致内存泄漏、slice的data字段、边界取值] memory->slice
[go map的初始化、访问、赋值、扩容、缩容]memory->map
[go 的负载因子为什么是6.5]memory->map
[聊一聊内存逃逸]memory->


---


[面试官：context携带的数据是线程安全的吗？]contexts
[Context 是怎么在 Go 语言中发挥关键作用的]contexts

---

[读者提问：反射是如何获取结构体成员信息的？]reflects

---


[内联函数和编译器对go代码的优化]compile
[终于识破这个go编译器把戏]compile
[翻译了一篇关于Go编译器的文章]compile
[迷惑了，Go len() 是怎么计算出来的？]compile->slice

---

[单元测试] unitTest


---

[文件存储] stroge
[Go 存储基础 — “文件”被偷偷修改？来，给它装个监控！] storage-->fsnofify
[浅析gowatch监听文件变动实现原理]storage-->gowatch
[Go 存储基础 — 内存结构体怎么写入文件？]storage->file

---


[会诱发goroutine挂起的27个原因]

[go官方信号量库]Semaphore

---

[详解并发编程包之 Go errgroup]
[Go 并发编程 — 结构体多字段的原子操作]concurrent


---

[go官方限流器的详解]limiter
[常用限流算法的应用场景和实现原理]limiter


[go-monitor：服务质量统计分析警告工具]monitor

---
[Go 的相对路径问题]path

[面试官：你能聊聊string和[]byte的转换吗？] bytes








