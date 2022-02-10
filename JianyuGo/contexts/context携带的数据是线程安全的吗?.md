[context携带value是线程安全的吗？]

这道题其实就是考察面试者对context实现原理的理解，如果不知道context的实现原理，很容易答错这道题

[先说答案]，context本身就是线程安全的，所以context携带value也是线程安全的

[为什么线程安全？]

context包提供两种创建根context的方式：

* context.Backgroud()
* context.TODO()

又提供了四个函数基于父Context衍生，
其中使用WithValue函数来衍生context并携带数据，
每次调用WithValue函数都会基于当前context衍生一个新的子context，
WithValue内部主要就是调用valueCtx类：

    func WithValue(parent Context, key, val interface{}) Context {
        if parent == nil {
            panic("cannot create context from nil parent")
        }
        if key == nil {
            panic("nil key")
        }
        if !reflectlite.TypeOf(key).Comparable() {
            panic("key is not comparable")
        }
        return &valueCtx{parent, key, val}
    }

valueCtx结构如下:

    // valueCtx继承父Context，这种是采用匿名接口的继承实现方式，key,val用来存储携带的键值对。
    type valueCtx struct {
        Context
        key, 
        val interface{}
    }
通过上面的代码分析，可以看到添加键值对不是在原context结构体上直接添加，
而是以此context作为父节点，重新创建一个新的valueCtx子节点，将键值对添加在子节点上，由此形成一条context链。

获取键值过程也是层层向上调用直到最终的根节点，中间要是找到了key就会返回，
否会就会找到最终的emptyCtx返回nil。画个图表示一下：

总结：context添加的键值对一个链式的，会不断衍生新的context，所以context本身是不可变的，因此是线程安全的。