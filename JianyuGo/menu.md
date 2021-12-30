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
    