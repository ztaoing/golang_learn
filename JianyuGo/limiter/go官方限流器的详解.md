限流器是提升服务稳定性的非常重要的组件，可以用来限制请求速率，保护服务，以免服务过载。
限流器的实现方法有很多种，常见的限流算法有固定窗口、滑动窗口、漏桶、令牌桶，
我在前面的文章 「常用限流算法的应用场景和实现原理」 中给大家讲解了这几种限流方法自身的特点和应用场景，
其中令牌桶在限流的同时还可以应对一定的突发流量，与互联网应用容易因为热点事件出现突发流量高峰的特点更契合。


令牌桶就是想象有一个固定大小的桶，系统会以恒定速率向桶中放 Token，桶满则暂时不放。
在请求比较的少的时候桶可以先"攒"一些Token，应对突发的流量，如果桶中有剩余 Token 就可以一直取。
如果没有剩余 Token，则需要等到桶中被放置了 Token 才行。

有的同学在看明白令牌桶的原理后就非常想去自己实现一个限流器应用到自己的项目里，em... 怎么说呢，
造个轮子确实有利于自己水平提高，不过要是应用到商用项目里的话其实大可不必自己去造轮子，
Golang官方已经替我们造好轮子啦 ......~！

Golang 官方提供的扩展库里就自带了限流算法的实现，即 golang.org/x/time/rate。
该限流器也是基于 Token Bucket(令牌桶) 实现的。


[限流器的内部结构]
time/rate包的Limiter类型对限流器进行了定义，所有限流功能都是通过基于Limiter类型实现的，其内部结构如下：

    type Limiter struct {
        mu     sync.Mutex
        limit  Limit    // 表示往桶里放Token的速率，它的类型是Limit，是int64的类型别名。
                        // 设置limit时既可以用数字指定每秒向桶中放多少个Token，也可以指定向桶中放Token的时间间隔，
                        // 其实指定了每秒放Token的个数后就能计算出放每个Token的时间间隔了。
        burst  int // 令牌桶的大小
        tokens float64 // 桶中的令牌。
        last time.Time // 上次往桶中放 Token 的时间。
        lastEvent time.Time // 上次发生限速器事件的时间（通过或者限制都是限速器事件）
    }
    可以看到在 timer/rate 的限流器实现中，并没有单独维护一个 Timer 和队列去真的每隔一段时间向桶中放令牌，
    而是仅仅通过计数的方式表示桶中剩余的令牌。每次消费取 Token 之前会先根据上次更新令牌数的时间差更新桶中Token数。

大概了解了time/rate限流器的内部实现后，下面的内容我们会集中介绍下该组件的具体使用方法：

[构造限流器]
    
    我们可以使用以下方法构造一个限流器对象：每秒向桶中产生10个token；这个桶的容量时100
    limiter := rate.NewLimiter(10, 100);

    这里有两个参数：
1、第一个参数是 r Limit，设置的是限流器Limiter的limit字段，
    代表每秒可以向 Token 桶中产生多少 token。
    Limit 实际上是 float64 的别名。
2、第二个参数是 b int，b 代表 Token 桶的容量大小，也就是设置的限流器 Limiter 的burst字段。

那么，对于以上例子来说，其构造出的限流器的令牌桶大小为 100, 以每秒 10 个 Token 的速率向桶中放置 Token。

除了给r Limit参数直接指定每秒产生的 Token 个数外，还可以用 Every 方法来指定向桶中放置 Token 的间隔，例如：
    
    limit := rate.Every(100 * time.Millisecond); // 每 100ms 往桶中放一个 Token。本质上也是一秒钟往桶里放 10 个
    limiter := rate.NewLimiter(limit, 100);

[使用限流器]
Limiter 提供了三类方法供程序消费 Token，可以每次消费一个 Token，也可以一次性消费多个 Token。
每种方法代表了当 Token 不足时，各自不同的对应手段:
1、可以阻塞等待桶中Token补充，
2、也可以直接返回取Token失败。

【Wait/WaitN】
    
    func (lim *Limiter) Wait(ctx context.Context) (err error)
    func (lim *Limiter) WaitN(ctx context.Context, n int) (err error)

Wait 实际上就是 WaitN(ctx,1)。
当使用 Wait 方法消费 Token 时，如果此时桶内 Token 数组不足 (小于 N)，
那么 Wait 方法将会阻塞一段时间，直至 Token 满足条件。如果充足则直接返回。

这里可以看到，Wait 方法有一个 context 参数。我们可以设置 context 的 Deadline 或者 Timeout，来决定此次 Wait 的最长时间。

    // 一直等到获取到桶中的令牌
    err := limiter.Wait(context.Background())
    if err != nil {
    fmt.Println("Error: ", err)
    }
    
    // 设置一秒的等待超时时间
    ctx, _ := context.WithTimeout(context.Background(), time.Second * 1)
    err := limiter.Wait(ctx)
    if err != nil {
    fmt.Println("Error: ", err)
    }

【Allow/AllowN】
        
    func (lim *Limiter) Allow() bool
    func (lim *Limiter) AllowN(now time.Time, n int) bool

Allow 实际上就是对 AllowN(time.Now(),1) 进行简化的函数。
AllowN 方法表示，截止到某一时刻，目前桶中数目是否至少为 n 个，
满足则返回 true，同时从桶中消费 n 个 token。
反之不消费桶中的Token，返回false。

【对应线上的使用场景是，如果请求速率超过限制，就直接丢弃超频后的请求。】


    if limiter.AllowN(time.Now(), 2) {
        fmt.Println("event allowed")
    } else {
       fmt.Println("event not allowed")
    }

【Reserve/ReserveN】
    
    func (lim *Limiter) Reserve() *Reservation
    func (lim *Limiter) ReserveN(now time.Time, n int) *Reservation

Reserve 相当于 ReserveN(time.Now(), 1)。
ReserveN 的用法就相对来说复杂一些，当调用完成后，
无论 Token 是否充足，都会返回一个 *Reservation 对象。

你可以调用该对象的Delay()方法，该方法返回的参数类型为time.Duration，反映了需要等待的时间，
必须等到等待时间之后，才能进行接下来的工作。
如果不想等待，可以调用Cancel()方法，该方法会将 Token 归还。

【举一个简单的例子，我们可以这么使用 Reserve 方法。】

    r := limiter.Reserve()
    if !r.OK() {
    // Not allowed to act! Did you remember to set lim.burst to be > 0 ?
        return
    }
    time.Sleep(r.Delay())
    Act() // 执行相关逻辑

[动态调整速率和桶大小]
Limiter 支持创建后动态调整速率和桶大小：
1、SetLimit(Limit) 改变放入 Token 的速率
2、SetBurst(int) 改变 Token 桶大小

有了这两个方法，可以根据现有环境和条件以及我们的需求，动态地改变 Token 桶大小和速率。

除了Golang官方提供的限流器实现，Uber公司开源的限流器uber-go/ratelimit也是一个很好的选择，
与Golang官方限流器不同的是Uber的限流器是通过漏桶算法实现的，
不过对传统的漏桶算法进行了改良，有兴趣的同学可以自行去体验一下。



