文章来源于[Golang梦工厂]

---
我们在之前的文章—— 源码剖析sync.WaitGroup(文末思考题你能解释一下吗?)，
从源码层面分析了sync.WaitGroup的实现，使用waitGroup可以实现一个goroutine等待一组goroutine干活结束，更好的实现了任务同步，

但是waitGroup却无法返回错误，当一组Goroutine中的某个goroutine出错时，我们是无法感知到的，
所以errGroup对waitGroup进行了一层封装，封装代码仅仅不到50行，下面我们就来看一看他是如何封装的？

---

[errGroup如何使用]

以下来自官方文档的例子：errorGroup

上面这个例子来自官方文档，代码量有点多，但是核心主要是在Google这个闭包中，

首先我们使用errgroup.WithContext创建一个errGroup对象和ctx对象，
然后我们直接调用errGroup对象的Go方法就可以启动一个协程了，
Go方法中已经封装了waitGroup的控制操作，不需要我们手动添加了，
最后我们调用Wait方法，其实就是调用了waitGroup方法。

这个包不仅减少了我们的代码量，而且还增加了错误处理，对于一些业务可以更好的进行并发处理。


---
[赏析errGroup]

数据结构:
我们先看一下Group的数据结构：

    type Group struct {
        cancel func() // 这个存的是context的cancel方法
        
        wg sync.WaitGroup // 封装sync.WaitGroup
        
        errOnce sync.Once // 保证只接受一次错误
        err     error // 保存第一个返回的错误
    }

方法解析:

    func WithContext(ctx context.Context) (*Group, context.Context)
    func (g *Group) Go(f func() error)
    func (g *Group) Wait() error

errGroup总共只有三个方法：

WithContext方法:
    
    func WithContext(ctx context.Context) (*Group, context.Context) {
        // 使用context的WithCancel()方法创建一个可取消的Context
        ctx, cancel := context.WithCancel(ctx)
        // 创建cancel()方法赋值给Group对象
        return &Group{cancel: cancel}, ctx
    }

    func (g *Group) Go(f func() error) {
        // 执行Add()方法增加一个计数器
        g.wg.Add(1)
        // 开启一个协程，运行我们传入的函数f，使用waitGroup的Done()方法控制是否结束
        go func() {
            defer g.wg.Done()

            //如果有一个函数f运行出错了，我们把它保存起来，如果有cancel()方法，则执行cancel()取消其他goroutine
            if err := f(); err != nil {
                // 这里大家应该会好奇为什么使用errOnce，也就是sync.Once，这里的目的就是保证获取到第一个出错的信息，避免被后面的Goroutine的错误覆盖。
                g.errOnce.Do(func() {
                        g.err = err
                        // 如果有cancel()方法，则执行cancel()取消其他goroutine
                        if g.cancel != nil {
                            g.cancel()
                        }
                })
            }
        }()
    }


[wait方法]

    func (g *Group) Wait() error {
        //调用waitGroup的Wait()等待一组Goroutine的运行结束
        g.wg.Wait()
        // 这里为了保证代码的健壮性，如果前面赋值了cancel，要执行cancel()方法
        if g.cancel != nil {
            g.cancel()
        }
        // 返回错误信息，如果有goroutine出现了错误才会有值
        return g.err
    }

---

* 我们可以使用withContext方法创建一个可取消的Group，也可以直接使用一个零值的Group或new一个Group，
    不过直接使用零值的Group和new出来的Group出现错误之后就不能取消其他Goroutine了。
* 如果多个Goroutine出现错误，我们只会获取到第一个出错的Goroutine的错误信息，
    晚于第一个出错的Goroutine的错误信息将不会被感知到。
* errGroup中没有做panic处理，我们在Go方法中传入func() error方法时要保证程序的健壮性

---

[踩坑日记] errorGroupBug
使用errGroup也并不是一番风顺的，我之前在项目中使用errGroup就出现了一个BUG，把它分享出来，避免踩坑。

    代码没啥问题吧，但是日志一直没有写入，排查了好久，终于找到问题原因。原因就是这个ctx

因为这个ctx是WithContext方法返回的一个带取消的ctx，我们把这个ctx当作父context传入WriteChangeLog方法中了，
如果errGroup取消了，也会导致上下文的context都取消了，所以WriteChangelog方法就一直执行不到。
