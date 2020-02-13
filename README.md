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
    分布式锁：
    etcd
    redis redlock
    consul
* [高可用分布式存储 etcd 的实现原理](https://draveness.me/etcd-introduction)
* [源码分析golang consul分布式锁lock delay问题](http://xiaorui.cc/2019/05/19/%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90golang-consul%E5%88%86%E5%B8%83%E5%BC%8F%E9%94%81lock-delay%E9%97%AE%E9%A2%98/)

-------
    内存分配
* [图解Go语言内存分配](https://zhuanlan.zhihu.com/p/59125443)

-------
    内存泄漏
* [一起 goroutine 泄漏问题的排查](https://zhuanlan.zhihu.com/p/100740270)
* [实战Go内存泄露](http://lessisbetter.site/2019/05/18/go-goroutine-leak/)
* [slice类型内存泄漏的逻辑(曹大)](https://xargin.com/logic-of-slice-memory-leak/)
* [分析golang time.After引起内存暴增OOM问题](http://xiaorui.cc/?p=5745)

-------
    channel
* [深度解密Go语言之channel ](https://zhuanlan.zhihu.com/p/74613114)
* [Golang并发：再也不愁选channel还是选锁](http://lessisbetter.site/2019/01/14/golang-channel-and-mutex/)
* [channel-draveness](https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-channel/)

-------
    锁🔐
* [golang多场景下RwMutex和mutex锁性能对比](http://xiaorui.cc/?p=5611)
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
* [完全兼容golang定时器的高性能时间轮实现(go-timewheel)](http://xiaorui.cc/2019/09/27/%e5%85%bc%e5%ae%b9golang-time%e5%ae%9a%e6%97%b6%e5%99%a8%e7%9a%84%e6%97%b6%e9%97%b4%e8%bd%ae%e5%ae%9e%e7%8e%b0/)
* [源码分析go time.timer和ticker的stop问题](http://xiaorui.cc/2019/09/09/%e6%ba%90%e7%a0%81%e5%88%86%e6%9e%90go-time-timer%e5%92%8cticker%e7%9a%84stop%e9%97%ae%e9%a2%98/)
* [分析golang time.After引起内存暴增OOM问题](http://xiaorui.cc/?p=5745)

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
* [runtime.stack加锁引起高时延及阻塞](http://xiaorui.cc/2020/01/03/go-runtime-stack%e5%8a%a0%e9%94%81%e5%bc%95%e8%b5%b7%e9%ab%98%e6%97%b6%e5%bb%b6%e5%8f%8a%e9%98%bb%e5%a1%9e/)
* [万字长文深入浅出 Golang Runtime](https://zhuanlan.zhihu.com/p/95056679)

-------
    grpc
    
* [golang grpc网关使用连接池提吞吐量](http://xiaorui.cc/2019/08/13/golang-grpc%e7%bd%91%e5%85%b3%e7%94%a8%e8%bf%9e%e6%8e%a5%e6%b1%a0%e6%8f%90%e9%ab%98%e5%90%9e%e5%90%90%e9%87%8f/)

-------
    goroutine池
* [使用golang协程池控制并发请求](http://xiaorui.cc/2019/05/24/%e4%bd%bf%e7%94%a8golang%e5%8d%8f%e7%a8%8b%e6%b1%a0%e6%8e%a7%e5%88%b6%e5%b9%b6%e5%8f%91%e8%af%b7%e6%b1%82/)

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
    crontab
* [开源golang兼容crontab的定时任务管理器](http://xiaorui.cc/?p=5625)