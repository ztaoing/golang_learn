
Go语言扩展包提供了一个带权重的信号量库Semaphore，使用信号量我们可以实现一个"工作池"控制一定数量的goroutine并发工作。
因为对源码抱有好奇的态度，所以在周末仔细看了一下这个库并进行了解析，在这里记录一下。

[何为信号量]

    信号量(Semaphore)，有时被称为信号灯，是[多线程环境下使用的一种设施，是可以用来保证两个或多个关键代码段不被并发调用。
    在进入一个关键代码段之前，线程必须获取一个信号量；一旦该关键代码段完成了，那么该线程必须释放信号量。
    其它想进入该关键代码段的线程必须等待直到第一个线程释放信号量。为了完成这个过程，需要创建一个信号量VI，
    然后将Acquire Semaphore VI以及Release Semaphore VI分别放置在每个关键代码段的首末端。
    确认这些信号量VI引用的是初始创建的信号量。

通过这段解释我们可以得知什么是信号量，其实信号量就是一种变量或者抽象数据类型，
用于控制并发系统中多个进程对公共资源的访问，访问具有原子性。

信号量主要分为两类：

* 二值信号量：顾名思义，其值只有两种0或者1，相当于互斥量，当值为1时资源可用，当值为0时，资源被锁住，进程阻塞无法继续执行。
  
* 计数信号量：信号量是一个任意的整数，起始时，如果计数器的计数值为0，那么创建出来的信号量就是不可获得的状态，
  如果计数器的计数值大于0，那么创建出来的信号量就是可获得的状态，并且总共获取的次数等于计数器的值。
  
---

【信号量工作原理】

信号量是由操作系统来维护的，信号量只能进行两种操作：等待、发送信号，操作总结来说，核心就是PV操作：
* P原语：P是荷兰语Proberen(测试)的首字母。为阻塞原语，负责把当前进程由运行状态转换为阻塞状态，直到另外一个进程唤醒它。
    操作为：申请一个空闲资源(把信号量减1)，若成功，则退出；若失败，则该进程被阻塞；
* V原语：V是荷兰语Verhogen(增加)的首字母。为唤醒原语，负责把一个被阻塞的进程唤醒，
    它有一个参数表，存放着等待被唤醒的进程信息。
  操作为：释放一个被占用的资源(把信号量加1)，如果发现有被阻塞的进程，则选择一个唤醒之。

在信号量进行PV操作时都为原子操作，并且在PV原语执行期间不允许有中断的发生。

PV原语对信号量的操作可以分为三种情况：
* 把信号量视为某种类型的"共享资源的剩余个数"，实现对一类共享资源的访问。（剩余可用的资源个数）
* 把信号量用作进程间的同步
* 视信号量为一个加锁标志，实现对一个共享变量的访问

具体在什么场景使用本文就不在继续分析，
接下来我们重点来看一下Go语言提供的扩展包Semaphore，看看它是怎样实现的。

【官方扩展包Semaphore】
我们之前在分析Go语言源码时总会看到这几个函数：
    
    func runtime_Semacquire(s *uint32)
    func runtime_SemacquireMutex(s *uint32, lifo bool, skipframes int)
    func runtime_Semrelease(s *uint32, handoff bool, skipframes int)

这几个函数就是信号量的PV操作，不过他们都是给Go内部使用的，

如果想使用信号量，那就可以使用官方的扩展包：Semaphore，这是一个带权重的信号量，接下来我们就重点分析一下这个库。

    安装方法：go get -u golang.org/x/sync

【数据结构】
semaphore库核心结构就是Weighted

    type Weighted struct {
        size    int64 // 这个代表的是最大权值，在创建Weighted对象指定
        cur     int64 // 相当于一个游标，标识当前已被使用的资源数
        mu      sync.Mutex // 并发情况下做临界区保护
        waiters list.List // 阻塞等待的调用者列表，使用链表数据结构保证先进先出的顺序，存储的数据是waiter对象，waiter数据结构如下：
    }

    type waiter struct {
        n     int64 // 这个就是等待调用者的权重值
        ready chan<- struct{} // 这就是一个channel，利用channel的close机制实现唤醒
    }

semaphore还提供了一个创建Weighted对象的方法，在初始化时需要给定最大权值：
    
    // NewWeighted为并发访问创建一个新的加权信号量，该信号量具有给定的最大权值。
    func NewWeighted(n int64) *Weighted {
        w := &Weighted{size: n}
        return w
    }

【阻塞获取权值的方法 - Acquire】

        func (s *Weighted) Acquire(ctx context.Context, n int64) error {
            s.mu.Lock() // 加锁保护临界区
            
            // 有资源可用并且没有等待获取权值的goroutine
            if s.size-s.cur >= n && s.waiters.Len() == 0 {
                s.cur += n // 加权
                s.mu.Unlock() // 释放锁
                return nil
            }

            // 要获取的权值n大于最大的权值了
            if n > s.size {
                // 先释放锁，确保其他goroutine调用Acquire的地方不被阻塞
                s.mu.Unlock()
                // 阻塞等待context的返回
                <-ctx.Done()
                return ctx.Err()
            }

            // 走到这里就说明现在没有资源可用了，需要等待，所以需要加入到等待队列
            // 创建一个channel用来做通知唤醒
            ready := make(chan struct{})

            // 创建waiter对象
            w := waiter{n: n, ready: ready}

            // waiter按顺序入队
            elem := s.waiters.PushBack(w)
            // 释放锁，等待唤醒，别阻塞其他goroutine
            s.mu.Unlock()
            
            // 阻塞等待唤醒
            select {
            // context关闭
            case <-ctx.Done():
                err := ctx.Err() // 先获取context的错误信息
                s.mu.Lock()
                select {
                    case <-ready:
                        // 在context被关闭后被唤醒了，那么试图修复队列，假装我们没有取消
                        err = nil
                    default:
                        // 判断是否是第一个元素
                        isFront := s.waiters.Front() == elem
                        // 移除第一个元素
                        s.waiters.Remove(elem)
                        // 如果是第一个元素且有资源可用通知其他waiter
                        if isFront && s.size > s.cur {
                            s.notifyWaiters()
                        }
                }
                s.mu.Unlock()
                return err

            // 被唤醒了
            case <-ready:
            return nil
            }
        }

注释已经加到代码中了，总结一下这个方法主要有三个流程：

* 流程一：有资源可用时并且没有等待权值的goroutine，走正常加权流程；

* 流程二：想要获取的权值n大于初始化时设置最大的权值了，这个goroutine永远不会获取到信号量，所以阻塞等待context的关闭；

* 流程三：前两步都没问题的话，就说明现在系统没有资源可用了，这时就需要阻塞等待唤醒，在阻塞等待唤醒这里有特殊逻辑；
  1、特殊逻辑二：context关闭后，则根据是否有可用资源决定通知后面等待唤醒的调用者，
    这样做的目的其实是为了避免当不同的context控制不同的goroutine时，未关闭的goroutine不会被阻塞住，
    依然执行，来看这样一个例子（因为goroutine的抢占式调度，所以这个例子也会具有偶然性）：
  
  2、特殊逻辑一：如果在context被关闭后被唤醒了，那么就先忽略掉这个cancel，试图修复队列。

[不阻塞获取权值的方法 - TryAcquire]
这个方法就简单很多了，不阻塞地获取权重为n的信号量，成功时返回true，失败时返回false并保持信号量不变。

    func (s *Weighted) TryAcquire(n int64) bool {
        s.mu.Lock() // 加锁

        // 有资源可用并且没有等待获取资源的goroutine
        success := s.size-s.cur >= n && s.waiters.Len() == 0

        if success {
            s.cur += n
        }
        s.mu.Unlock()
        return success
    }

[释放权重]
这里就是很常规的操作，主要就是资源释放，同时进行安全性判断，如果释放资源大于持有的资源，则会发生panic。
    func (s *Weighted) Release(n int64) {
        s.mu.Lock()

        // 释放资源, 利用负值判断释放的资源是否大于持有的资源
        s.cur -= n
        // 释放资源大于持有的资源，则会发生panic
        if s.cur < 0 {
            s.mu.Unlock()
            panic("semaphore: released more than held")
        }
        // 通知其他等待的调用者
        s.notifyWaiters()
        s.mu.Unlock()
    }

【唤醒waiter】
在Acquire和Release方法中都调用了notifyWaiters，我们来分析一下这个方法：
这里只需要注意一个点：唤醒waiter采用先进先出的原则，避免需要资源数比较大的waiter被饿死。

        func (s *Weighted) notifyWaiters() {
            for {
                // 获取等待调用者队列中的队员
                next := s.waiters.Front()
                // 没有要通知的调用者了
                if next == nil {
                    break // No more waiters blocked.
                }
        
                // 断言出waiter信息
                w := next.Value.(waiter)
                if s.size-s.cur < w.n {
                    // 没有足够资源为下一个调用者使用时，继续阻塞该调用者，遵循先进先出的原则，
                    // 避免需要资源数比较大的waiter被饿死
                    //
                    // 考虑一个场景，使用信号量作为读写锁，现有N个令牌，N个reader和一个writer
                    // 每个reader都可以通过Acquire（1）获取读锁，writer写入可以通过Acquire（N）获得写锁定
                    // 但不包括所有的reader，如果我们允许reader在队列中前进，writer将会饿死-总是有一个令牌可供每个reader
                    break
                }
                // 获取资源
                s.cur += w.n

                // 从waiter列表中移除
                s.waiters.Remove(next)

                // 使用channel的close机制唤醒waiter
                close(w.ready)
            }
        }

【何时使用Semaphore】
到这里我们就把Semaphore的源代码看了一篇，代码行数不多，封装的也很巧妙，那么我们该什么时候选择使用它呢？

目前能想到一个场景就是Semaphore配合上errgroup实现一个"工作池"，

使用Semaphore限制goroutine的数量，配合上errgroup做并发控制，示例如下：use

本文我们主要赏析了Go官方扩展库Semaphore的实现，他的设计思路简单，仅仅用几十行就完成了完美的封装，
值得我们借鉴学习。

不过在实际业务场景中，我们使用信号量的场景并不多，大多数场景我们都可以使用channel来替代，
但是有些场景使用Semaphore来实现会更好，比如上篇文章【[警惕] 请勿滥用goroutine】我们使用channel+sync来控制goroutine数量，
这种实现方式并不好，因为实际已经起来了多个goroutine，只不过控制了工作的goroutine数量，
如果改用semaphore实现才是真正的控制了goroutine数量。


