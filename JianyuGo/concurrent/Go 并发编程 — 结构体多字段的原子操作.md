【多字段更新？ 】

---
并发编程中，原子更新多个字段是常见的需求。
举个例子，有一个 struct Person 的结构体，里面有两个字段。
我们先更新 Person.name，再更新 Person.age ，这是两个步骤，但我们必须保证原子性。

有童鞋可能奇怪了，为什么要保证原子性？

我们以一个示例程序开端，公用内存简化成一个全局变量，开 10 个并发协程去更新。你猜最后的结果是啥？

main.go

打印结果是啥？你能猜到吗？
可能是这样的：
    
    p.name=nobody:2
    p.age=3

也可能是：

    p.name=nobody:8
    p.age=7
按照排列组合来算，一共有 10*10 种结果。

那我们想要什么结果？我们想要 name 和 age 一定要是匹配的，不能牛头不对马嘴。
换句话说，name 和 age 的更新一定要原子操作，不能出现未定义的状态。

我们想要的是 （ nobody:i，i ），正确的结果只能在以下预定的 10 种结果出现：

    （ nobody:0, 0 ）
    （ nobody:1, 1 ）
    （ nobody:2, 2 ）
    （ nobody:3, 3 ）
    ...
    （ nobody:9, 9 ）

这仅仅是一个简单的示例，童鞋们思考下自己现实的需求，应该是非常常见的。

现在有两个问题：
[第一个问题：这个 demo 观察下运行时间，用 time 来观察，时间大概是 200 ms 左右，为什么？]
    
    root@ubuntu:~/code/gopher/src/atomic_test# time ./atomic_test 
    p.name=nobody:8
    p.age=7
    
    real 0m0.203s
    user 0m0.000s
    sys 0m0.000s

如上就是 203 毫秒。划重点：这个时间大家请先记住了，对我们分析下面的例子有帮助。

这个 200 毫秒是因为奇伢在 update 函数中故意加入了一点点时延，这样可以让程序估计跑慢一点。
每个协程跑 update 的时候至少需要 200 毫秒，10 个协程并发跑，没有任何互斥，时间重叠，所以整个程序的时间也是差不都 200 毫秒左右。

[第二个问题：怎么解决这个正确性的问题。]

大概两个办法：
* 锁互斥
* 原子操作
下面详细分析下异同和优劣。
  
---

[锁实现]

在并发的上下文，用锁来互斥，这是最常见的思路。 
锁能形成一个临界区，锁内的一系列操作任何时刻都只会有一个人更新，如此就能确保更新不会混乱，从而保证多步操作的原子性。

首先配合变量，对应一把互斥锁：

    // 全局变量（简单处理）
    var p Person
    // 互斥锁，保护变量更新
    var mu sync.Mutex

更新的逻辑在锁内：

    func update(name string, age int) {
        // 更新：加锁，逻辑串行化
        mu.Lock()
        defer mu.Unlock()
        
            // 以下逻辑不变
    }

大家按照上面的把程序改了之后，逻辑是不是就正确了。一定是 （ nobody:i，i ）配套更新的。
但你注意到另一个可怕的问题吗？

程序运行变的好慢！！！！

同样用 time 命令统计下程序运行时间，竟然耗费 2 秒！！！，10 倍的时延增长，每次都是这样。

    root@ubuntu:~/code/gopher/src/atomic_test# time ./atomic_test
    p.name=nobody:8
    p.age=8
    
    real 0m2.017s
    user 0m0.000s
    sys 0m0.000s

不禁要问自己，为啥？
还记得上面我提到过，一个 update 固定要 200 毫秒。

加锁之后的 update 函数逻辑全部在锁内，10 个协程并发跑 update 函数，
但由于锁的互斥性，抢锁不到就阻塞等待，保证 update 内部逻辑的串行化。
第 1 个协程加上锁了，后面 9 个都要等待，依次类推。最长的等待时间应该是 1.8 秒。
换句话说，程序串行执行了 10 次 update 函数，时间是累加的。程序 2 秒的运行时延就这样来的。

加锁不怕，抢锁等待才可怕。在大量并发的时候，由于锁的互斥特性，这里的性能可能堪忧。

[还有就是抢锁失败的话，是要把调度权让出去的，直到下一次被唤醒。
这里还增加了协程调度的开销，一来一回可能性能就更慢了下来。]

思考：用锁之后正确性是保证了，某些场景性能可能堪忧。那咋吧？
在本次的例子，下一步的进化就是：原子化操作。

温馨提示：
怕童鞋误会，声明一下：锁不是不能用，是要区分场景，不分场景的性能优化措施是没有意义的哈。
大部分的场景，用锁没啥问题。且锁是可以细化的，比如读锁和写锁，更新加写锁，只读操作加读锁。
这样确实能带来较大的性能提升，特别是在写少读多的时候。

---
原子操作

其实我们再深究下，这里本质上是想要保证更新 name 和 age 的原子性，要保证他们配套。
其实可以先在局部环境设置好 Person 结构体，然后一把原子赋值给全局变量即可。
Go 提供了 atomic.Value 这个类型。

怎么改造？
首先把并发更新的目标设置为 atomic.Value 类型：

    // 全局变量（简单处理）
    var p atomic.Value

然后 update 函数改造成先局部构造，再原子赋值的方式：

    func update(name string, age int) {
        // lp是一个局部变量
        lp := &Person{}
        // 更新第一个字段
        lp.name = name
        // 加点随机性
        time.Sleep(time.Millisecond * 200)
        // 更新第二个字段
        lp.age = age
        // 原子设置到全局变量
        p.Store(lp)
    }

最后 main 函数读取全局变量打印的地方，需要使用原子 Load 方式：
    
     // 结果是啥？你能猜到吗？
    _p := p.Load().(*Person)
    fmt.Printf("p.name=%s\np.age=%v\n", _p.name, _p.age)

这样就解决并发更新的正确性问题啦。感兴趣的童鞋可以运行下，结果都是正确的 （ nobody:i，i ）。

下面再看一下程序的运行时间：

    root@ubuntu:~/code/gopher/src/atomic_test# time ./atomic_test
    p.name=nobody:7
    p.age=7
    
    real 0m0.202s
    user 0m0.000s
    sys 0m0.000s

竟然是 200 毫秒作用，比锁的实现时延少 10 倍，并且保证了正确性。
    

---

[为什么会这样？]

因为这 10 个协程还是并发的，没有类似于锁阻塞等待的操作，只有最后 p.Store(lp) 调用内才有做状态的同步，
而这个时间微乎其微，所以 10 个协程的运行时间是重叠起来的，自然整个程序就只有 200 毫秒左右。
锁和原子变量都能保证正确的逻辑。在我们这个简要的场景里，我相信你已经感受到性能的差距了。

    当然了，还是那句话，具体用哪个，要看具体场景，不能一概而论。
    而且，锁有自己无可替代的作用，它能保证多个步骤的原子性，而不仅仅是字段的赋值。

相信你已经非常好奇 atomic.Value 了，下面简要的分析下原理，是否真的很神秘呢？
原理可能要大跌眼镜。


趁现在我们还不懂内部原理，先思考个问题（不然待会一下子看懂了就没意思了）？
Value.Store  和 Value.Load 是用来赋值和取值的。

我的问题是，这两个函数里面有没有用户数据拷贝？
Store 和 Load 是否是保证了多字段拷贝的原子性？

提前透露下：并非如此。

---
atomic.Value 原理

[atomic.Value 结构体]

atomic.Value  定义于文件 src/sync/atomic/value.go  ，结构本身非常简单，就是一个空接口：

    type Value struct {
        v interface{}
    }

在之前文章中，奇伢有分享过 Go 的空接口类型（ interface {} ）
在 Go 内部实现是一个叫做 eface 的结构体（ src/runtime/iface.go ）：

    type eface struct {
        _type *_type
        data  unsafe.Pointer
    }

[重要]

interface {} 是给程序猿用的，eface  是 Go 内部自己用的， 位于不同层面的同一个东西，这个请先记住了，
因为 atomic.Value 就利用了这个特性，
在 value.go 定义了一个 ifaceWords 的结构体。

[划重点]：interface {} ，eface ，ifaceWords 这三个结构体内存布局完全一致，
        只是用的地方不同而已，本质无差别。这给类型的强制转化创造了前提。

[Value.Store 方法]

看一下简要的代码，这是一个简单的 for 循环：

        func (v *Value) Store(x interface{}) {
            // 强制转化类型，转变成 ifaceWords （三种类型，相同的内存布局，这是前提）
            vp := (*ifaceWords)(unsafe.Pointer(v))
            xp := (*ifaceWords)(unsafe.Pointer(&x))

            for {
                // 获取数据类型
                typ := LoadPointer(&vp.typ)
                // 第一个判断：atomic.Value 初始的时候是 nil 值，那么就是走这里进去的；
                if typ == nil {
                    runtime_procPin()
                    // 通过 CompareAndSwapPointer 来确保 ^uintptr(0)  只能被一个执行体抢到，其他没抢到的走 continue ，再循环一次；
                    if !CompareAndSwapPointer(&vp.typ, nil, unsafe.Pointer(^uintptr(0))) {
                        runtime_procUnpin()
                        continue
                    }
                // 初始赋值
                    StorePointer(&vp.data, xp.data) // 保存数据
                    StorePointer(&vp.typ, xp.typ)   // 保存类型
                    runtime_procUnpin()
                    return
                }

                // 第二个判断：这个也是初始的时候，这是一个中间状态； 这个中间状态具体指什么？
                // ^uintptr(0) 是第一次存取的标志位，也就是如果是第一次存取就continue
                if uintptr(typ) == ^uintptr(0) {
                    continue
                }
                
                // 第三个判断：类型校验，通过这里就能看出来，Value 里面的类型不能变，否则会 panic；
                if typ != xp.typ {
                    panic("sync/atomic: store of inconsistently typed value into Value")
                }
    
                // 划重点啦：只要过了初始化赋值阶段，基本上就是直接跑到这行代码啦
                StorePointer(&vp.data, xp.data)
                return
            }
        }

有几个点稍微解释下：
1、atomic.Value 使用 ^uintptr(0) 作为第一次存取的标志位，这个标识位是设置在 type 字段里，这是一个中间状态；

2、通过 CompareAndSwapPointer 来确保 ^uintptr(0)  只能被一个执行体抢到，其他没抢到的走 continue ，再循环一次；

3、atomic.Value 第一次写入数据时，将当前协程设置为不可抢占，当存储完毕后，即可解除不可抢占；

4、真正的赋值，无论是第一次，还是后续的 data 赋值，在 Store 内，只涉及到指针的原子操作，"不涉及到数据拷贝"；


这里有没有大跌眼镜？

Store 内部并不是保证多字段的原子拷贝！！！！Store  里面处理的是个结构体指针

只通过了 StorePointer 保证了指针的原子赋值操作。
我的天？是这样的吗？那何来的原子操作。

[核心在于]：Value.Store()  的参数必须是个局部变量（或者说是一块全新的内存）。

---
这里就回答了上面的问题：Store，Load 是否有数据拷贝？
【划重点】：没有！没动数据

原来你是这样子的 atomic.Value ！

回忆一下我上面的 update 函数，真的是局部变量，全新的内存块：

    func update(name string, age int) {
        // 注意哦，局部变量哦
        lp := &Person{}
        // 更新字段 。。。。
     
        // 设置的是全新的内存地址给全局的 atomic.Value 变量
        p.Store(lp)
    }

---

又有个问题，你可能会想了，如果 p.Store( /* */ ) 传入的不是指针，而是一个结构体呢？

事情会是这样的：
1、编译器识别到这种情况，编译期间就会多生成一段代码，用 runtime.convT2E  函数把结构体赋值转化成 eface （注意，这里会涉及到结构体数据的拷贝）；
2、然后再调用 Value.Store 方法，所以就 Store 方法而言，行为还是不变；

---

再思考一个问题：既然是指针的操作，为什么还要有个 for 循环，还要有个  CompareAndSwapPointer  ？

这是因为 ifaceWords 是两个字段的结构体，初始赋值的时候，要赋值类型和数据指针两部分。
atomic.Value 是服务所有类型，此类需求的，通用封装。

----
【Value.Load 方法】

有写就有读嘛，看一下读的简要的实现：

        func (v *Value) Load() (x interface{}) {
            vp := (*ifaceWords)(unsafe.Pointer(v))
            typ := LoadPointer(&vp.typ)
            // 初始赋值还未完成
            if typ == nil || uintptr(typ) == ^uintptr(0) {
                return nil
            }
            // 划重点啦：只要过了初始化赋值阶段，原子读的时候基本上就直接跑到这行代码啦；
            data := LoadPointer(&vp.data)
            xp := (*ifaceWords)(unsafe.Pointer(&x))
            // 赋值类型，和数据结构体的地址
            xp.typ = typ
            xp.data = data
            return
        }

哇，太简单了。处理做了一下初始赋值的判断（返回 nil ），后续基本就只靠 LoadPointer 函数来个原子读指针值而已。

---

【总结】

1、interface {} ，eface ，ifaceWords  本质是一个东西，同一种内存的三种类型解释，用在不同层面和场景。它们可以通过强制类型转化进行切换；
2、atomic.Value 使用 cas 操作只在初始赋值的时候，一旦赋值过，后续赋值的原子操作更简单，依赖于 StorePointer ，指针值得原子赋值；
3、atomic.Value 的 Store 和 Load 方法都不涉及到数据拷贝，只涉及到指针操作；
4、atomic.Value 的神奇的核心在于：每次 Store 的时候用的是全新的内存块 ！！！ 且 Load 和 Store 都是"以完整结构体的地址"进行操作，所以才有原子操作的效果。
5、atomic.Value 实现多字段原子赋值的原理，千万不要以为是并发操作同一块多字段内存，还能保证原子性；