文章来源于[网管叨bi叨]

---
Context 是 Go 语言独有的设计，在其他编程语言中很少见到类似的概念，
用一句话解释 Context 在 Go 语言中的作用就是：

Context 为同一任务的多个 goroutine 之间提供了 退出信号通知 和 元数据传递的功能。

那么如果不用 Context，就不能在 Go 语言里实现多个 goroutine  间的信号通知和元数据传递了吗？
答案是：简单场景下可以，在多层级 goroutine 的控制中就行不通了。

我们举一个例子来理解上面那段话:
假如主协程中有多个任务，主协程对这些任务有超时控制；
而其中任务1又有多个子任务，任务1对这些子任务也有自己的超时控制，那么这些子任务既要感知主协程的取消信号，也需要感知任务1的取消信号。

任务的 goroutine 层级越深，想要自己做退出信号感知和元数据共享就越难。
所以我们需要一种优雅的方案来实现这样一种机制：
* 上层任务取消后，所有的下层任务都会被取消；
* 中间某一层的任务取消后，只会将当前任务的下层任务取消，而不会影响上层的任务以及同级任务；
* 可以线程安全地在 goroutine 之间共享元数据(子context是在父context的基础上生成了一个新的context，所以是线程安全的)；

为此 Go 官方在1.7 版本就引入了 Context 来实现上面阐述的机制。

Go 的 context 包提供了对Context的接口定义和类型实现，我通过一张类图给大家描述下 context 提供的接口和类型。
通过上面的类图我们能获取到这些信息:
* 除了Context 接口外还定义了一个叫做 canceler 的接口，实现了它的类型即为带取消功能的 Context。
* emptyCtx 什么属性也没有，啥也不能干。
* valueCtx 只能携带一个键值对，且自身要依附在上一级 Context 上。
* timerCtx 继承自 cancelCtx 他们都是带取消功能的 Context。
* 除了emptyCtx，其他类型的 Context 都依附在上级 Context 上
  
  看完这个类图，你可能会问 Context 是怎么实现元数据在任务间传递的呢？
  毕竟一个valueCtx只携带一个键值对。
  其实原理也很简单，它实现的 Value 方法能够在整个Context链路上查找指定键的值，直到回源到根 Context。

通过查找Context 携带的键值对的示意图我们能看到Context链路的根节点是一个 emptyCtx，
这也就是emptyCtx 任何功能也不提供的原因，它是用来作为根节点而存在的。

每次要在Context链路上增加要携带的键值对时，都要在上级Context的基础上新建一个 valueCtx 存储键值对，
且只能增加不能修改，读取 Context 上的键值又是一个幂等的操作，
所以 Context 就这样实现了线程安全的数据共享机制，且全程无锁，不会影响性能。

---
那么 “上层任务取消后，所有的下层任务都会被取消”，“中间某一层的任务取消后，只会将当前任务的下层任务取消，而不会影响上层的任务以及同级任务” 
这两个取消信号同步的关键点， Context 又是怎么实现的呢？

    下文把cancelCtx，timerCtx统称带取消功能的Context，且示意图中也会用 cancelCtx 这个标记代表他们。

首先在 创建 带取消功能的Context的时候还是要在父级Context节点的基础上创建，保持整个Context链路的连续性。
除此之外还会在Context链路中找到上一个带取消功能的 Context，把自己加入到它的 children 列表里。
这样在整个Context链路里，除了父子Context之间有直接关联外，
可取消的Context还会通过维护自身携带的children 属性建立自己与下级可取消Context之间的关联。

    看源码的同学，重点关注用 WithCancel、WithDeadline这些方法里对propagateCancel的调用
    propagateCancel中会在祖先Context节点中找到可取消的Context，把自己维护到祖先的children属性里

经过这个结构设计，如果要在整个任务链路上取消某个cancelCtx时，就能做到既取消自己，
也能通知下级 cancelCtx 进行取消，同时还不会影响到上级和同级的其他节点。

现在让我们再回到开头那个例子，有了 Context 之后，我们的任务会变成什么样呢？

我们让每个 goroutine 都携带了 Context ，那些做子任务的 goroutine 只要监听了这些子 cancelCtx 也就能收到信号，
结束自己的运行，即通过Context 完成上级 goroutine 对下级 goroutine 的取消控制。

面对不同层级的 goroutine 的取消条件不同的情况，
代码里只需要监听传递到 goroutine 里的 Context 就能做到，免除了监听多个信号的繁琐。

---
针对Context的使用建议，Go官方提到了下面几点：
* 不要将 Context 塞到结构体里。直接将 Context 类型作为函数的第一参数，而且一般都命名为 ctx。
* 不要向函数传入一个 nil 的 context，如果你实在不知道传什么，标准库的TODO方法给你准备好了一个 emptyCtx。
* 不要把本应该作为函数参数的类型塞到 context 中，context 存储的应该是一些在 goroutine 共享的数据，比如Server的信息等等。

第一点的前半部分我觉得说的很牵强，比如在官方的 net/http 包里就把 Context 放到了 Request的结构体里，其他几点确实是需要注意的地方。
好了，今天是不卷源码的一天，
我用通俗的语言和几张图示向你展示了Context的设计理念和它在Go语言里起到的重要作用，
如果你能喜欢这种形式，请不要吝啬你的点赞和在看，感谢你的支持。

