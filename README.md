# golang_learn
# 只争朝夕，不负韶华!

此处整理的知识点来源于以下书籍📚:
* 《Go语言编程入门与实战技巧》-黄靖钧
* 《Go语言核心编程》-李文塔
* 《Go编发编程实战》-郝林
* 《Go语言编程》-许式伟
* 《Go程序设计语言》-译本
* 《Go语言学习笔记》-雨痕

-------
    etcd
* [高可用分布式存储 etcd 的实现原理](https://draveness.me/etcd-introduction)

-------
    channel
* [深度解密Go语言之channel ](https://zhuanlan.zhihu.com/p/74613114)
* [Golang并发：再也不愁选channel还是选锁](http://lessisbetter.site/2019/01/14/golang-channel-and-mutex/)
* [channel-draveness](https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-channel/)
-------
    context
* [深度解密Go语言之context](https://zhuanlan.zhihu.com/p/68792989)
* [上下文context](https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-context/)

-------
    map
* [深度解密Go语言之 map](https://zhuanlan.zhihu.com/p/66676224)
* [map并发崩溃一例(非线程安全)(曹大)](https://xargin.com/map-concurrent-throw/)

-------
    scheduler
* [深度解密Go语言之 scheduler](https://zhuanlan.zhihu.com/p/80853548)
* [调度器-draveness](https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-goroutine/)

-------
    error
* [Golang error 的突围](https://zhuanlan.zhihu.com/p/82985617)

-------
    pprof
* [深度解密Go语言之 pprof](https://zhuanlan.zhihu.com/p/91241270)

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
    内存分配
* [图解Go语言内存分配](https://zhuanlan.zhihu.com/p/59125443)

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
    内存泄漏
* [一起 goroutine 泄漏问题的排查](https://zhuanlan.zhihu.com/p/100740270)
* [实战Go内存泄露](http://lessisbetter.site/2019/05/18/go-goroutine-leak/)
* [slice类型内存泄漏的逻辑(曹大)](https://xargin.com/logic-of-slice-memory-leak/)
-------
    runtime
* [万字长文深入浅出 Golang Runtime](https://zhuanlan.zhihu.com/p/95056679)

-------
    sync
* [Golang并发的次优选择：sync包](http://lessisbetter.site/2019/01/04/golang-pkg-sync/)

-------
    select
* [Golang并发模型：轻松入门select](http://lessisbetter.site/2018/12/13/golang-slect/)
* [Golang并发模型：select进阶](http://lessisbetter.site/2018/12/17/golang-selete-advance/)
* [select](https://draveness.me/golang/docs/part2-foundation/ch05-keyword/golang-select/)

-------
    make  new
* [make 和 new](https://draveness.me/golang/docs/part2-foundation/ch05-keyword/golang-make-and-new/)

------
    time
* [定时器-draveness](https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-timer/)

-------
    函数
* [函数调用-draveness](https://draveness.me/golang/docs/part2-foundation/ch04-basic/golang-function-call/)

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

-------
    gomaxprocs
* [golang gomaxprocs调高引起调度性能损耗](http://xiaorui.cc/2020/01/11/golang-gomaxprocs%e8%b0%83%e9%ab%98%e5%bc%95%e8%b5%b7%e8%b0%83%e5%ba%a6%e6%80%a7%e8%83%bd%e6%8d%9f%e8%80%97/)

-------
    runtime
* [go runtime.stack加锁引起高时延及阻塞](http://xiaorui.cc/2020/01/03/go-runtime-stack%e5%8a%a0%e9%94%81%e5%bc%95%e8%b5%b7%e9%ab%98%e6%97%b6%e5%bb%b6%e5%8f%8a%e9%98%bb%e5%a1%9e/)