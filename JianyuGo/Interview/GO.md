* for select时，如果通道已经关闭会怎样？如果select中只有一个case会怎么？
* 对已经关闭的chan进行读写会怎样？
* 对未初始化的chan进行读写会怎样？

[slice问题]
# slice的数据结构为：
    type SliceHeader struct{
        Data uintptr // 引用数组指针地址
        Len int // 切片的目前使用的长度
        Cap int // 切片的容量
}
  * nil的切片和空切片指向的地址是不一样的  -->nilSlice.go

[基本理解]

1、 go是传值还是传引用--》传值，传递的总是参数的副本（go没有引用传递）
    1、 向一个函数传递一个int值，就会得到int的副本，传一个一个指针就会得到指针的副本
    2、map和slice的行为类似于指针：他们包含指向底层map或slice数据的指针的描述符
        * 复制一个map或slice值并不会复制它指向的数据
        * 复制一个接口值会复制存储在接口值中的东西
        * 如果接口值持有一个结构，复制接口值就会复制该结构。如果接口值持有一个
            指针，复制接口值会复制该指针，但同样不会复制它指向的数据
        

    传值：也叫值传递，在调用函数时将实际参数 复制一份 传递到函数中。这样在函数中如果修改了这个参数副本，将不会影响到原来的参数
        值传递本质上不能认为是一个东西，指向的不是一个内存地址
    证明：
  
        func main(){
        s:="这是一本书"
        fmt.Println("main 内存地址%p",s)  // main中的内存地址: 0xc000116220
        hello(&s)
        }
        func hello(s *string){
            fmt.Println("hello 内存地址%p",s) // hello中的内存地址: 0xc000132020
        }
  
    从上边的，可以看到在main中的变量s，在经过hello函数的参数副本传递后，内部的输出地址发生了改变。


      func main() {
      s := "脑子进煎鱼了"
      fmt.Printf("main 内存地址：%p\n", &s)
      hello(&s)
      fmt.Println(s)
      }
    
        func hello(s *string) {
        fmt.Printf("hello 内存地址：%p\n", &s)
        *s = "煎鱼进脑子了"
        }
    输出结果：
    main 内存地址：0xc000010240
    hello 内存地址：0xc00000e030
    煎鱼进脑子了
    
    上边的main证明了什么呢？ 传递过去的值是指向内存空间的地址，那么是可以对这块内存空间做修改的
    也就是说，这两个内存地址，其实是指针的指针，其根源都指向着同一个指针，也就是指向着变量s。

    传引用：也叫引用传递，指在调用函数时将实际参数的地址直接传递到函数中，那么在函数中对参数所做的修改，将影响到实际参数
    在go语言中，官方已经明确了没有传引用！

    争议最大的map和slice
    go语言中的map和slice类型，能直接修改，难道不是同一个内存地址，不是引用了？
    注意：map和slice的行为类似于指针，他们是包含指向底层map和slice数据的 指针的描述符。

    map：
    func main(){
     m:=make(map[string]string)
     m["index"] = "这是第一次"
     fmt.Printf("main 内存地址:%p\n",&m)  //main中的m内存地址：0x0000e028
    
     hello(m)
     fmt.Printf("%v",m)
    }
    
    func hello(p map[string]string){
        fmt.Printf("hello 内存地址：%p\n",&p) // hello中的m内存地址:0x0000e038
        p["index"] = "这是第二次"
    
    地址发生了变化，确实是传值，
    修改之后：map["index"] = "这是第二次",修改成功
    为什么值传递，又能做到类似传引用的效果呢，能修改到原值呢？
    
    创建map：
        func makemap(t *maptype ,hint int, h *hmap)*hmap{}
    这是创建map类型的底层runtime方法，注意其返回的*hmap类型，是一个指针。 也就是说在调用hello方法时，相当于
    是传入一个指针参数hello(*hmap)。这种情况称其为"引用类型"，但是"引用类型"不等于就是传引用，还是有明确的区别的
    
    在go语言中与map类型类似的还有chan类型:
    func makechan(t *chantype,size int)*hchan{}
    
    
slice 
     
        func main() {
            s := []string{"烤鱼", "咸鱼", "摸鱼"}
            fmt.Printf("main 内存地址：%p\n", s)
            hello(s)
            fmt.Println(s)
        }
    
        func hello(s []string) {
            fmt.Printf("hello 内存地址：%p\n", s)
            s[0] = "煎鱼"
        }
        输出结果：
        main 内存地址：0xc000098180
        hello 内存地址：0xc000098180
        [煎鱼 咸鱼 摸鱼]

        从结果来看，两者的内存地址一样，也成功的变更到了变量 s 的值。这难道不是引用传递吗?
        关注两个细节：
            没有用 & 来取地址。
            可以直接用 %p 来打印。
        之所以可以同时做到上面这两件事，是因为标准库 fmt 针对在这一块做了优化：
        func (p *pp) fmtPointer(value reflect.Value, verb rune) {
            var u uintptr

            switch value.Kind() {
            case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
                         u = value.Pointer()
            default:
                         p.badVerb(verb)
            return
        }
        留意到代码 value.Pointer，标准库进行了特殊处理，直接对应的值的指针地址，当然就不需要取地址符了。
        标准库 fmt 能够输出 slice 类型对应的值的原因也在此：
            func (v Value) Pointer() uintptr {
                ...
                case Slice:
                return (*SliceHeader)(v.ptr).Data
                }
            }
            
            type SliceHeader struct {
                Data uintptr
                Len  int
                Cap  int
            }
        其在内部转换的 Data 属性，正是 Go 语言中 slice 类型的运行时表现 SliceHeader。我们在调用 %p 输出时，是在输出 slice 的底层存储数组元素的地址。
        
        下一个问题是：为什么 slice 类型可以直接修改源数据的值呢。
        其实和输出的原理是一样的，在 Go 语言运行时，传递的也是相应 slice 类型的底层数组的指针，但需要注意，其使用的是指针的副本。严格意义是引用类型，依旧是值传递。
        妙不妙？
2、 Go面试官问我如何实现面向对象（go语言如何实现面向对象？进一步就是go语言如何实现面向对象特性中的继承）
  

    什么是面向对象？OOP
    面向对象编程（OOP）是一种基于"对象"概念的编程范式，它可以包含数据和代码：数据以字段的形式存在（通常成为属性），
        代码以程序的形式存在（通常成为方法）
    对象的程序可以访问、修改自己的数据字段
    对象经常被定义为类的一个实例
    对象利用属性和方法的私有、受保护、公共可见，对象的内部状态收到保护，不受外接影响（被封装）
    面向对象的三大基本特征：封装、继承、多态
    
    go语言是一门面向对象的语言吗？是也不是！：
    1、go有类型和方法，并且运行面向对象的编程风格，但没有类型层次
    2、go中的"接口"概念提供了一种不同的实现方法，我们认为这种方法玉玉使用，而且在某些方面更加通用，还有一些方法
        可以将类型嵌入到其他的类型中，以提供类似的东西，但不等同于子类！
    3、go中的方法比c++或java中的方法更通用：他们可以为任何类型的数据定义，甚至是内置类型，
    4、go由于没有类型层次，所以，go中的对象比c++或Java等语言更轻巧。
    
    go实现面向对象编程
    封装：面向对象中的封装值得是可以隐藏对象内部属性、和实现细节，仅对外提供公开接口调用，这样用户就不需要关注对象的内部是怎么实现的
    
    go语言中的属性访问权，通过首字母大小写来控制
    
    继承：面向对象中的继承指的是子类继承父类的特征和行为，是的子类对象具有父类的实例和方法。
    在go语言中没有extends关键字的这种继承方式，在语言设计上采取的是组合的方式：
    type Animal struct{ Name string}
    type Cat struct{ 
        Animal
        FeatureA string
    }
    type Dog strcut{
        Animal 
        FeatureB string
    }
    
    多态：面向对象中的多态指的是 同一个行为具有多种不同表现形式或形态的能力，具体是指一个类实例（对象）的相同方法
        在不同情形有不同表现形式
    多态使得不同内部结构的对象可以共享相同的外部接口，也就是都是一套外部模板，内部实际是什么，只要符合规格就可以
    
    type AnimalSounder interface { MakeDNA()}
    
    // 传入的是一个接口
    func MakeSomeDNA(animalSounder AnimalSounder){
        animalSounder.MakeDNA()
    }
    func (c *Cat)MakeDNA(){
        fmt.Println("猫")
    }
    func (d *Dog)MakeDNA(){
        fmt.Println("狗")
    }
    func main(){
        MakeSomeDNA(&Cat{})
        MakeSomeDNA(&Dog{})
    }
    
    面向对象的三大特征，五大基本原则：
    三大特征：封装、继承、多态
    五大原则：单一职责原则、开放封闭原则、里式替换原则、依赖倒置原则、接口隔离原则

3、go结构体和结构体指针调用有什么区别吗？在什么情况下用哪种？有没有什么注意事项：

    1、在使用上，方法是否需要修改接收器？如果需要，接收器必须是一个指针
    2、在效率上：如果接收器很大，比如：一个大的结构体，使用指针占用的内容比使用结构体占用的内存更小时，则使用指针
    3、在一致性的考虑，如果类型的某些方法必须有指针接收器，按摩其余的方法应该有指针接收器，所以无论类型如何使用，方法集都是一致的
    
4、new和make是什么，差异在哪？
        
    new和make，主要用途是用于分配相应类型的内存空间
    make：内置函数make支持slice、map、channel三种数据类型的内存创建，
         注意：其返回值是所创建类型的本身，而不是新的指针引用
    func make(t Type,size ...integerType)Type
    
    func main(){
        slice1:=make([]int,1,5)
        map1:=make(map[int]bool,5)
        chan1:=make(chan int,1)
    fmt.Println(slice1,map1,chan1)
    )
    // 注意：使用make初始化切片的类型时，会带有零值

    new:可以对类型进行内存创建和初始化
        注意：其返回值是所创建类型的指针引用，所以与make在存在区别
    func new(Type)*Type
    
    其实new函数在日常工程代码中是比较少见的，因为他可被替代
    一般会直接使用快捷的T{}来进行初始化，因为常规的结构体都会带有结构体的字面属性：这种初始化方式更方便
    func NewT()*T{
        return &T{Name:"煎鱼"}
    }

    new函数也能初始化make的三种类型：
    v1:=new(chan bool)
    v2:=new(map[string]struct{})

    那new和make的区别，优势是什么呢？
    本质上在于make函数在初始化时，会初始化slice、chan、map类型的内部数据结构，new函数并不会
    例如：在map类型中，合理的长度len和容量cap可以提高效率和减少开销。
    
    更进一步的区别：
    make：
        1、能够 分配并初始化 类型所需的内存空间和结构，返回引用类型的本身。
        2、仅支持channel、map、slice是那种类型。
        3、make函数会对三种类型的内部数据结构（长度、容量等）赋值。
    new:
        1、能够分配类型所需的内存空间，返回指针引用（指向内存的指针）。
        2、可被替代，能够通过字面量快速初始化。
    
5、什么是协程，协程和线程的区别和联系？ 需要将进程、线程、协程都介绍一遍，关键是将协程和线程的区别和联系介绍清除
        
    进程：进程是操作系统对一个正在运行的程序的一种抽象，进程是资源分配的最小单位。
    为什么有进程：为了合理压榨CPU的性能和分配运行的时间片，不能闲着！
    
    单个CPU一次只能运行一个任务。如果一个进程跑着，就把这个CPU霸占了，那是非常不合理的。
    为什么要压榨CPU的性能？因为CPU是在是太快了，寄存器仅仅能够住的上他的步伐，ram和其他挂在各总线上的设备则更是望尘莫及
    
    多进程的缘由：
    进程上的任务，可能是计算型的任务，也有可能是网络io调用。
    
    线程：为什么有了进程还要有线程呢？
    1、进程间的信息难以共享，进程间的通信信息交换，性能开销比较大。
    2、创建进程的性能开销比较大
    
    一个进程可以由多个成为线程的执行单元组成，每个线程都运行在进程的上下文中，共享着同样的代码和全局数据。
    多线程之间比多进程之间更容易共享数据，在上下文切换中线程比进程更高效。
    
    1、线程之间可以非常方便、快速的共享数据
        只需要将数据复制到进程汇总的共享区域就可以了，但是需要注意避免多个线程修改同一份内存
        
    2、创建线程比创建进程要快10倍，甚至更多
        同一个进程中的所有线程，像内存页、页表就需要了

    协程：是用户态的线程。通常创建协程时，会从进程的堆中分配一段内存作为协程的栈
    线程的栈有8MB，而协程栈的大小通常只有KB，而go语言的协程更夸张，只有2-4KB，非常轻巧。
    
    
    

[调度模型]  《scalable go scheduler design doc》go调度器

1、 GMP模型，为甚要有P？
    为什么不是G和M直接绑定就完了，还另外需要一个P？是要解决什么问题？
   

    GM、GMP模型的变迁原因：
    GM: 
    在go1.1之前的go的调度模型是GM，是没有P的：也就是只有一个全局队列s
    // 停止正在运行的goroutine，运行另一个可运行的goroutine,在用户态模仿了时间片用完后的切换操作
    schedule(G *gp){
        ...
        // 调用schedlock方法获取全局锁，此版本只有M和G
       
        schedlock();
        
        if(gp != nil) {
        ...
        // 获取全局锁成功
        switch(gp->status){
        case Grunnable:
        case Gdead:
        // Shouldn't have been running!
        runtime·throw("bad gp->status in sched");
        
        case Grunning:
        // 获取全局锁成功后，将当前的goroutine状态从running（正在被调度）修改为runable（可被调度）
        gp->status = Grunnable;
        // 调用gput方法保存当前goroutine的运行状态等信息，以便后续的使用
        gput(gp);
        break;
        }
        
        // 调用nextgandunlock，寻找下一个可运行的goroutine，并释放全局锁
        gp = nextgandunlock();
        gp->readyonstop = 0;
        // 获取到下一个待运行的goroutine后，将其运行状态修改为running
        gp->status = Grunning;
        
        m->curg = gp;
        gp->m = m;
        ...
        // 调用runtime.gogo，将刚刚获取到的下一个待执行的goroutine运行起来， 进入下一轮调度
        runtime·gogo(&gp->sched, 0);
        }
            
        思考：schedule方法，在正常刘晨下，是不会返回的，也即是不会结束主流程。
        GM模型的缺点：
        1、因为只有一个全局队列，所以使用的是全局锁
            mutex需要保护所有与goroutine相关的操作（创建、完成、重排等）导致锁竞争严重
        2、goroutine传递的问题：
            goroutine（G）交接(G.nextg):工作线程(M)之间会经常交接可运行的goroutine。
            上述操作可能会导致延迟增加和额外的开销。每个M必须能够执行任何可运行的G，特别是刚刚创建的G的M
        3、每个M都需要做内存缓存（M.mcache）
            会导致资源消耗过大(每个mcache可以吸纳到2M的内存缓存和其他缓存)，数据的局部性差
        4、频繁的线程阻塞/解阻塞：
            在存在syscall的情况下，线程经常被阻塞和解阻塞。这增加了额很多额外的性能开销。

    GMP: 为了解决GM模型的以上问题，在Go1.1时，增加了P。并实现了work stealing算法来解决一些新产生的问题。
        

        加入P之后，带来了什么改变呢？
        1、每个p有自己的本地队列，大幅度减轻了对全局队列的直接依赖，减少了锁竞争。而GM模型的性能开发大头就是锁竞争
        2、在GMP模型中实现了work stealing算法，是每个p相对平衡。如果p的本地队列空了，他会去全局队列或其他p的队列中
          窃取可运行的G，减少空转即没有任务执行，提高了利用率。
    
    为什么要有P：
        如果想实现本地队列，为什么不直接在M上加呢？M上照样可以实现类似的功能，为什么再多加一个P呢？：
            M是系统线程
            1、一般来讲，M的数量都会多于P，在go中，M的数量最大限制时10000，p的默认数量是CPU的核数。
               另外由于M的属性，也就是如果存在系统阻塞调用，阻塞了M，又不够用的情况下，M会不断增加
            2、M不断增加的，如果本地队列挂载在M上，那就意味着本地队列也会随之增加。这显示是不合理的
            3、M被系统调用阻塞后，我们是希望把它已有的未执行的任务分配给其他继续运行，而不是一阻塞就导致全部停止
        因此在M上使用本地队列是不合理的。所以引入了新组件P
        
2、单机goroutine数量控制在多少合适，会影响GC和调度？
    goroutine太多了会影响gc和调度吧，主要是怎么预算这个数时合理的呢？
    
    go调度的流程：
    1、当执行 go func()时，实际上就是创建一个全新的goroutine，即G
    2、创建的G会被放入P的本地队列或全局队列中。需要注意一点，这里的P指的是创建G的P
    3、唤醒或创建M以便执行G
    4、不断进行时间循环
    5、寻找 可用状态的G，然后执行任务
    6、清除后，重新进入事件循环

    本地队列数量有限制，不允许超过256个。
    在创建G时，会优先选择P的本地队列，如果本地队列满了，则将P的本地队列的一半G移动到全局队列
    
    在协程的运行过程中，真正干活的GPM又分别被什么约束？

    M的限制：
    首先要知道在协程的执行中，真正干活的是GPM中的哪一个？M
    G是用户态上的东西，最终执行都需要映射，绑定到M上去执行。
    M有没有限制？
    在go语言中，M的默认数量限制时10000，如果超出则会报错：GO: runtime: program exceeds 10000-thread limit
    通常只在goroutine出现阻塞的情况下，才会遇到这种情况，这可能也预示着你的程序有问题。
    如果要设置固定数量的M：debug.SetMaxThreads 进行设置

    G的限制：
    没有限制，但是理论上会受到内存你的影响。 假设一个goroutine创建需要4K
        4k*80000 = 320000k 0.3G内存
        4k*1000000 = 400 0000k 4G内存
    依次就可以大概计算出一台单机能够创建goroutine的大概数量级别。
    goroutine创建所需申请的2-4K是需要连续的内存块.

    P的限制：p的数量是否有限制？受什么影响？
    有限制，P的数量受环境变量GOMAXPROCS的直接影响

    在go语言中，通过设置GOMAXPROCS可以调整调度器中P的数量
    另一点，与P相关的M，是需要绑定P才能进行具体的任务，因此p的多少会影响到Go程序的运行表现
    
    p的数量基本是受本机的核数影响。
    p的数量是否影响goroutine的数量创建呢？
    不影响，
    
    单机的goroutine数量只要控制在限额以下，就是合理的。
    真实场景得看里面具体跑的是什么，跑的如果是资源怪兽，只运行几个goroutine都可以跑死
    
    因此想定义预算，就得看跑的是什么？
    
3、go结构体是否可以比较？为什么？
    
    在go中，结构体有时并不能直接比较，当其中包含：slice、map、function时，是不能比较的。若强行比较会导致出现报错
    而指针引用，其虽然都是new(string)，从表面看是一个东西，单其具体返回的地址是不一样的。
        package main

            import "fmt"
            
            type Vertex struct {
            Name1 string
            Name2 string
            }
            
            func main() {
                v := Vertex{"脑子进了", "煎鱼"}
                v.Name2 = "蒸鱼"
                fmt.Println(v.Name2)
            }
        输出结果： 蒸鱼

            type Value struct {
            Name   string
            Gender string
        }
        
        func main() {
            v1 := Value{Name: "煎鱼", Gender: "男"}
            v2 := Value{Name: "煎鱼", Gender: "男"}
            if v1 == v2 {
            fmt.Println("脑子进煎鱼了")
            return
        }
        
            fmt.Println("脑子没进煎鱼")
        }
    输出结果：脑子进煎鱼了

            type Value struct {
            Name   string
            Gender *string
        }
        
        func main() {
            v1 := Value{Name: "煎鱼", Gender: new(string)}
            v2 := Value{Name: "煎鱼", Gender: new(string)}
            if v1 == v2 {
            fmt.Println("脑子进煎鱼了")
            return
        }
        
            fmt.Println("脑子没进煎鱼")
        }
        答案是：脑子没进煎鱼。

                type Value struct {
            Name   string
            GoodAt []string
        }
        
        func main() {
            v1 := Value{Name: "煎鱼", GoodAt: []string{"炸", "煎", "蒸"}}
            v2 := Value{Name: "煎鱼", GoodAt: []string{"炸", "煎", "蒸"}}
            if v1 == v2 {
            fmt.Println("脑子进煎鱼了")
            return
        }
        
            fmt.Println("脑子没进煎鱼")
        }
        答案：./main.go:15:8: invalid operation: v1 == v2 (struct containing []string cannot be compared)
        type Value1 struct {
        Name string
        }
        
        type Value2 struct {
        Name string
        }
        
        func main() {
            v1 := Value1{Name: "煎鱼"}
            v2 := Value2{Name: "煎鱼"}
            if v1 == v2 {
            fmt.Println("脑子进煎鱼了")
            return
        }
        
            fmt.Println("脑子没进煎鱼")
        }
        报错：./main.go:18:8: invalid operation: v1 == v2 (mismatched types Value1 and Value2)

    那是不是就完全没法比较了呢？并不，我们可以借助强制转换来实现：
             if v1 == Value1(v2) {
                fmt.Println("脑子进煎鱼了")
                return
        }

    如果必须要比较，可以使用反射，reflect.DeepEqual()，其用来判定两个值是否深度一致。
            func main() {
            v1 := Value{Name: "煎鱼", GoodAt: []string{"炸", "煎", "蒸"}}
            v2 := Value{Name: "煎鱼", GoodAt: []string{"炸", "煎", "蒸"}}
            if reflect.DeepEqual(v1, v2) {
                fmt.Println("脑子进煎鱼了")
                return
            }
    
              fmt.Println("脑子没进煎鱼")
            }

    1、相同的类型的值是深度相等的，不同类型的值永远不会深度相等
    2、当数组值的对应元素深度相等时，数组值是深度相等的
    3、当结构体值的对应字段都是深度相等的，则值是深度相等
    4、当接口值的深度相等，则深度相等
    golang.org/pkg/reflect/#DeepEqual 
    该方法对go语言中的各种类型都进行了兼容处理和判定。

4、单核CPU，开两个goroutine，其中一个死循环，会怎样？
    
    第一点：计算机只有一个单核CPU对go程序产生什么影响？
    从单核CPU来看，最大的影响就是GMP模型中的P，因为P的数量默认是与CPU核数（GOMAXPROCS）保持一致
    M必须与P绑定，然后不断在M上循环查找可运行的G来执行相应的任务。
    
    第二点：goroutine受限
    goroutine的数量和运行模式都是受限的。有两个goroutine，一个在死循环，另外一个在正常运行
    可以理解为 main goroutine+ 一个新的goroutine跑死循环。
    需要注意的是，goroutine里跑着死循环，也就是时时刻刻在运行着业务逻辑，
    这块需要与单核CPU关联起来，考虑是否会一直阻塞主，把整个go进程运行给hang住了
    
    注意：面试时，可以先举出这个场景，解释清楚后，再补充提问面试官是否是这类场景

    第三点：go版本的问题
    go的调度是会变动的，在不同的go版本中，结果可能会不一样

    实战演练：
    func main(){
        // 模拟单核CPU
        runtime.GOMAXPROCES(1)
    
        // 模拟死循环
        go func(){
            for {}
        }
        time.sleep(time.Millisecond)
        fmt.Println("end")
    }
    
    答案：
    在go1.14之前，不会输出任何结果：
        这个死循环的goroutine是无法被抢占的，这个死循环中没有设计主动放弃执行权或被动放弃的行为，所以会一直执行
        那为什么main goroutine会无法运行呢？因为会优先调用休眠，但由于单核CPU，只有一个P。唯一的P又一直在运行，无法停止
        导致main goroutine没有机会被调度，所以这个程序一直阻塞在了死循环中。
    在go1.14及之后，能够正常输出结果
        主要是因为在1.14实现了基于信号的抢占调度，
    1、抢占阻塞在系统调用上的P
    2、抢占运行时间过长的G（类似时间片？）
    该方法会检测符合场景的P，当满足上述两个场景之一时，就会发送信号给M，M收到信号后将会休眠正在阻塞的goroutine
    调用绑定的信号方法，并进行重新调度。

    注意：在go语言中，sysmon会用于检测抢占。sysmon是go的runtime的系统检测器，
        sysmon可进行forcegc、netpoll、retake等一系列操作。
    
5、详解go程序的启动流程，g0，m0是什么？


6、goroutine泄露的N种方法
    
     泄露的原因大多集中在：
    1、goroutine内正在进行channel/mutex等读写操作，单由于逻辑问题，某些情况下回被一直阻塞。
    2、goroutine内的业务逻辑进入死循环，资源一直无法释放
    3、goroutine内的业务逻辑进入长时间等待，又不断有新增的goroutine进入等待
    
    1、channel使用不当：channel读写操作时的逻辑问题
        a、发送不接收：
        func main (){
            for i:=0;i<4;i++{
                queryAll()
                fmt.printf("goroutines:%d\n,runtime.NumGoroutine())
            }
        }
        func queryAll()int{
            // 无缓冲
            ch:=make(chan int)
            // 发送三次
            for i:=0;i<3;i++{
                go func (){
                    ch<-query()
                }()
            }
            // 只接收一次
          return <-ch
        }

        func query()int{
            n:=rand.Intn(100)
            time.Sleep(time.Duration(n)*time.millisecond)
            return n
        }
    
        输出结果： 每次向channel发送3次，但是接收端只接收一次,有两个没有接收
         gorotines:3
         gorotines:5
         gorotines:7
         gorotines:9
            
        b、接收不发送：没有发送端
        func main(){
            defer func(){
                fmt.println("goroutines:",runtime.NumGoroutine())
            }()
            var ch chan struct{}
            go func(){
                ch<-struct{}{}
            }()
            time.Sleep(time.Second)
        }
        输出：
            goroutines:2 

        c、nil channel
        func main(){
            defer func(){
                fmt.println("goroutines:",runtime.NumGoroutine())
            }
            var ch chan int
            go func(){
                <-ch
            }
         time.sleep(time.second)
        }
        输出结果： goroutines:2
        注意：channel如果忘记初始化，无论读，还是写操作，都会造成阻塞
        正确的姿势：
        // 使用make函数进行初始化
        ch:= make(ch chan,int)
        go func(){
            <-ch
        }()
        ch<-0

        d、奇怪的慢等待:经典的事故场景
        func main() {
            for {
                go func() {
                    _, err := http.Get("https://www.xxx.com/")
                    if err != nil {
                        fmt.Printf("http.Get err: %v\n", err)
                    }
                    // do something...
            }()

            time.Sleep(time.Second * 1)
            fmt.Println("goroutines: ", runtime.NumGoroutine())
            
            输出结果：
                goroutines:  5
                goroutines:  9
                goroutines:  13
                goroutines:  17
                goroutines:  21
                goroutines:  25
                ...
                在这个例子中，展示了一个 Go 语言中经典的事故场景。也就是一般我们会在应用程序中去调用第三方服务的接口。
                但是第三方接口，有时候会很慢，久久不返回响应结果。恰好，Go 语言中默认的 http.Client 是没有设置超时时间的。
                因此就会导致一直阻塞，一直阻塞就一直爽，Goroutine 自然也就持续暴涨，不断泄露，最终占满资源，导致事故。
            在 Go 工程中，我们一般建议至少对 http.Client 设置超时时间：
                    httpClient := http.Client{
                        Timeout: time.Second * 15,
                    }
            并且要做限流、熔断等措施，以防突发流量造成依赖崩塌，依然吃 P0。
        
            e、互斥锁忘记解锁
                    func main() {
                            total := 0
                            defer func() {
                                time.Sleep(time.Second)
                                fmt.Println("total: ", total)
                                fmt.Println("goroutines: ", runtime.NumGoroutine())
                        }()
                        
                            var mutex sync.Mutex
                            for i := 0; i < 10; i++ {
                                go func() {
                                    mutex.Lock()
                                    total += 1
                                }()
                            }
                        }
                    }
             输出结果：
             total:  1
             goroutines:  10  
    第一个互斥锁 sync.Mutex 加锁了，但是他可能在处理业务逻辑，又或是忘记 Unlock 了。
    因此导致后面的所有 sync.Mutex 想加锁，却因未释放又都阻塞住了

    我们建议如下写法：
    var mutex sync.Mutex
    for i := 0; i < 10; i++ {
        go func() {
            mutex.Lock()
            defer mutex.Unlock()
            total += 1
    }()
    }

    f、同步锁使用不当
        func handle(v int) {
            var wg sync.WaitGroup
            wg.Add(5)
            for i := 0; i < v; i++ {
                fmt.Println("脑子进煎鱼了")
                wg.Done()
            }
            wg.Wait()
    }
    
    func main() {
        defer func() {
        fmt.Println("goroutines: ", runtime.NumGoroutine())
        }()
        
            go handle(3)
            time.Sleep(time.Second)
    }
    由于 wg.Add 的数量与 wg.Done 数量并不匹配，因此在调用 wg.Wait 方法后一直阻塞等待。
    
    建议如下写法：
      var wg sync.WaitGroup
    for i := 0; i < v; i++ {
        wg.Add(1)
        defer wg.Done()
        fmt.Println("脑子进煎鱼了")
    }
    wg.Wait()
    
    g、排查方法
    我们可以调用 runtime.NumGoroutine 方法来获取 Goroutine 的运行数量，进行前后一比较，就能知道有没有泄露了。
    但在业务服务的运行场景中，Goroutine 内导致的泄露，大多数处于生产、测试环境，因此更多的是使用 PProf：
    import (
    "net/http"
     _ "net/http/pprof"
    )
    
    http.ListenAndServe("localhost:6060", nil))
    只要我们调用 http://localhost:6060/debug/pprof/goroutine?debug=1，PProf 会返回所有带有堆栈跟踪的 Goroutine 列表。
    也可以利用 PProf 的其他特性进行综合查看和分析，这块参考我之前写的《Go 大杀器之性能剖析 PProf》，基本是全村最全的教程了。

7、go在什么时候会抢占P？
    
    goroutine早期是没有设计成抢占式的，早期的goroutine只有在读写、主动让出、锁等操作时才会触发调度切换
    这样有一个严重的问题，就是在垃圾收集器进行STW时，如果有一个goroutine一直都在阻塞调用，垃圾回收器就会一直等待他...
    这种情况就需要抢占式调度来解决问题。
    
    如果一个goroutine运行时间过久，就需要进行抢占来解决。
    
    go语言在1.12起开始实现抢占式调度器，并不断完善：
    g0.x：基于单线程的调度器
    g1.0：基于多线程的调度器
    g1.1：基于任务窃取的调度器
    g1.2-1.13：基于协作式的抢占调度器
    g1.14：基于信号的抢占调度器
    
    调度器的新体感：非均匀存储器访问调度，但是非常复杂：NUMA-aware scheduler for GO
    
    为什么要抢占P？
    答：如果不抢，就没有机会运行，会hang死，又或是资源分配不均了。出现这种情况，显然是不合理的。
    
    举个栗子：
        func main(){
            runtime.GOMAXPROCS(1)
            go func(){
                for{}
            }
            time.Sleep(time.Millisecond)
            fmt.Println("end")
        }
    在老版本中，因为只有一个p，子goroutine在不断执行，并无法停止，导致main goroutine没有执行的机会。
    
    如果goroutine从阻塞状态恢复，该怎么继续运行呢？没了p怎么办？
    该goroutine需要先检查所在M是否仍然绑定着P：
    1、若有P，则可以调整状态，继续运行
    2、没有p，M可以重新抢占P，再绑定P

    抢占P，本身就是一个双向行为，你抢了我的P，我也可以抢占别人的P来运行。
    
    怎么抢占P？
    具体的处理在runtime.retake方法：处理以下两种场景
    1、抢占阻塞在系统调用上的P  
    2、抢占运行时间过长的G
    
    以下主要这对抢占P的场景：
        func retake(now int64) uint32 {
                n := 0
                // 防止发生变更，对所有 P 加锁，allpLock的所有P的锁
                // 其会保护allp、idleMask、timerpMask的无P读取和大小变化，以及对allp的所有写入操作
                lock(&allpLock)
                
                // 走入主逻辑，对所有 P 开始循环处理
                for i := 0; i < len(allp); i++ {
                _p_ := allp[i]
                pd := &_p_.sysmontick
                s := _p_.status
                sysretake := false
                ...
                if s == _Psyscall {
    
                // 进入主逻辑：
                // 场景1：会使用万能的for循环对所有的P进行一个个的处理。
                // 判断是否超过 1 个 sysmon tick 周期(20us)的任务，则会从系统调用中抢占P，否则跳过
                // 
                    t := int64(_p_.syscalltick)
                    if !sysretake && int64(pd.syscalltick) != t {
                    pd.syscalltick = uint32(t)
                    pd.syscallwhen = now
                    continue
                }
                
                ...
                }
                }
                unlock(&allpLock)
                return uint32(n)
        }
    
        场景2：
                func retake(now int64) uint32 {
                        for i := 0; i < len(allp); i++ {
                        ...
                        if s == _Psyscall {
                        // 从此处开始分析
                        // runqempty(_p_) == true 方法会判断任务队列 P 是否为空，以此来检测有没有其他任务需要执行。
                        // atomic.Load(&sched.nmspinning)+atomic.Load(&sched.npidle) > 0 会判断是否存在空闲 P 和正在进行调度窃取 G 的 P。
                        // pd.syscallwhen+10*1000*1000 > now 会判断系统调用时间是否超过了 10ms。
                    
                        // 这里奇怪的是 runqempty 方法明明已经判断了没有其他任务，这就代表了没有任务需要执行，是不需要抢夺 P 的。
                        // 但实际情况是，由于可能会阻止 sysmon 线程的深度睡眠，最终还是希望继续占有 P。
                            if runqempty(_p_) &&
                            atomic.Load(&sched.nmspinning)+atomic.Load(&sched.npidle) > 0 &&
                            pd.syscallwhen+10*1000*1000 > now {
                            continue
                        }
                        ...
                        }
                        }
                        unlock(&allpLock)
                        return uint32(n)
                }

            抢夺p的阶段：
                    func retake(now int64) uint32 {
                            for i := 0; i < len(allp); i++ {
                            ...
                            if s == _Psyscall {

                            // 承接上半部分
                            // 解锁相关属性：需要调用 unlock 方法解锁 allpLock，从而实现获取 sched.lock，以便继续下一步。
                            // 减少闲置 M：需要在原子操作（CAS）之前减少闲置 M 的数量（假设有一个正在运行）。否则在发生抢夺 M 时可能会退出系统调用，递增 nmidle 并报告死锁事件。
                           //  修改 P 状态：调用 atomic.Cas 方法将所抢夺的 P 状态设为 idle，以便于交于其他 M 使用。
                           //  抢夺 P 和调控 M：调用 handoffp 方法从系统调用或锁定的 M 中抢夺 P，会由新的 M 接管这个 P。
                            unlock(&allpLock)

                            incidlelocked(-1)
                            if atomic.Cas(&_p_.status, s, _Pidle) {
                                if trace.enabled {
                                traceGoSysBlock(_p_)
                                traceProcStop(_p_)
                            }
                                n++
                                _p_.syscalltick++
                                handoffp(_p_)
                            }
                                incidlelocked(1)
                                lock(&allpLock)
                            }
                            }
                            unlock(&allpLock)
                            return uint32(n)
            }

8、诱发goroutine挂起的27个原因？
    goroutine一泄漏就看到他，这个是什么？runtime.gopark(),这个函数到底是什么？作用是？
    
    runtime.gopark是什么？
    最快的办法就是看源码，其实现细节在src/runtime/proc.go文件中。
    
    func gopark(unlockf func(*g, unsafe.Pointer) bool, lock unsafe.Pointer, reason waitReason, traceEv byte, traceskip int) {
       // 调用acquirem，
            1、湖区当前goroutine锁绑定的m，设置各类所需数据；
            2、调用releasem()将当前goroutine和其m的绑定关系解除
        mp := acquirem()
        gp := mp.curg
        status := readgstatus(gp)
        mp.waitlock = lock
        mp.waitunlockf = unlockf
        gp.waitreason = reason
        mp.waittraceev = traceEv
        mp.waittraceskip = traceskip
        releasem(mp)
        // 调用park_m函数
            1、将当前goroutine的状态从_Grunning 切换为 _Gwaiting等待状态，
            2、删除m和当前goroutine m->curg（简称gp）之间的关联
        // 调用mcall()，仅仅会在需要进行goroutine切换时会被调用：
            1、切换当前线程的堆栈，从g的堆栈切换到g0的堆栈并调用fn(g)函数
            2、将g的当前PC/SP保存在g-sched中，以便后续调用goready()时可以恢复运行现场
        mcall(park_m)
    }
    注：该函数的关键作用是将当前的goroutine让如等待状态，这意味着goroutine被暂时搁置，也就是被运行时调度器暂停了。
    
    缘由：
    之所以goroutine泄漏，你会看到大量的runtime.gopark()，也就是大量的goroutine被暂停了，进入到休眠状态
    直至满足条件后再被runtime.goready函数唤醒，该函数沪江已经转呗就绪的goroutine切换状态，再加入运行队列。
    等待调度器的新一轮调度。

    goroutine一共有9中状态：
    _Gidle：刚刚被分配，还没有进行初始化。
    _Grunnable：已经在运行队列中，还没有执行用户代码。
    _Grunning：不在运行队列里中，已经可以执行用户代码，此时已经分配了 M 和 P。
    _Gsyscall：正在执行系统调用，此时分配了 M。
    _Gwaiting：在运行时被阻止，没有执行用户代码，也不在运行队列中，此时它正在某处阻塞等待中。
    _Gmoribund_unused：尚未使用，但是在 gdb 中进行了硬编码。
    _Gdead：尚未使用，这个状态可能是刚退出或是刚被初始化，此时它并没有执行用户代码，有可能有也有可能没有分配堆栈。
    _Genqueue_unused：尚未使用。
    _Gcopystack：正在复制堆栈，并没有执行用户代码，也不在运行队列中。

9、gopark 的 27 个诱发原因  
    
    第一类：
    waitReasonZero：无正式解释，从使用情况来看。主要在 sleep 和 lock 的 2 个场景中使用。
    waitReasonGCAssistMarking：GC 辅助标记阶段会使得阻塞等待。
    waitReasonIOWait：IO 阻塞等待时，例如：网络请求等。

    第二部分：
    waitReasonChanReceiveNilChan：对未初始化的 channel 进行读
    waitReasonChanSendNilChan：对未初始化的 channel 进行写
    
    第三部分
    waitReasonDumpingHeap：对 Go Heap 堆 dump 时，这个的使用场景仅在 runtime.debug 时，也就是常见的 pprof 这一类采集时阻塞。
    waitReasonGarbageCollection：在垃圾回收时，主要场景是 GC 标记终止（GC Mark Termination）阶段时触发。
    waitReasonGarbageCollectionScan：在垃圾回收扫描时，主要场景是 GC 标记（GC Mark）扫描 Root 阶段时触发。

    第四部分
    waitReasonPanicWait：在 main goroutine 发生 panic 时，会触发。
    waitReasonSelect：在调用关键字 select 时会触发。
    waitReasonSelectNoCases：在调用关键字 select 时，若一个 case 都没有，会直接触发。

    第五部分
    waitReasonGCAssistWait：GC 辅助标记阶段中的结束行为，会触发。
    waitReasonGCSweepWait：GC 清扫阶段中的结束行为，会触发。
    waitReasonGCScavengeWait：GC scavenge 阶段的结束行为，会触发。GC Scavenge 主要是新空间的垃圾回收，是一种经常运行、快速的 GC，负责从新空间中清理较小的对象。

    第六部分
    waitReasonChanReceive：在 channel 进行读操作，会触发。
    waitReasonChanSend：在 channel 进行写操作，会触发。
    waitReasonFinalizerWait：在 finalizer 结束的阶段，会触发。在 Go 程序中，可以通过调用 runtime.SetFinalizer 函数来为一个对象设置一个终结者函数。这个行为对应着结束阶段造成的回收。

    第七部分
    waitReasonForceGCIdle：强制 GC（空闲时间）结束时，会触发。
    waitReasonSemacquire：信号量处理结束时，会触发。
    waitReasonSleep：经典的 sleep 行为，会触发。

    第八部分
    waitReasonSyncCondWait：结合 sync.Cond 用法能知道，是在调用 sync.Wait 方法时所触发。
    waitReasonTimerGoroutineIdle：与 Timer 相关，在没有定时器需要执行任务时，会触发。
    waitReasonTraceReaderBlocked：与 Trace 相关，ReadTrace会返回二进制跟踪数据，将会阻塞直到数据可用。

    第九部分
    waitReasonWaitForGCCycle：等待 GC 周期，会休眠造成阻塞。
    waitReasonGCWorkerIdle：GC Worker 空闲时，会休眠造成阻塞。
    waitReasonPreempted：发生循环调用抢占时，会会休眠等待调度。
    waitReasonDebugCall：调用 GODEBUG 时，会触发。
[数据结构]
1、go interface的一个坑及原理分析：本文是go比较有名的一个坑。因为在线上真实出现过这个坑。写给不了解的嗯在使用
    if err !=nil 的时候提高警惕
    
    go语言的interface{}在使用过程中有一个特别坑的特性，当你比较一个interface{}类型是否是nil的时候，这是需要特别注意的问题
    
    例子一：
        func main(){
            var v interface{}
            v = (*int)(nil)
            fmt.println(v==nil)
        }
        结果：false ,为什么不是true，命名已经强行设置为nil了。

    例子二：
        func main (){
            var data *byte
            var in interface{}
            
            fmt.println(data,data==nil) // nil true
            fmt.println(in,in==nil)     // nil true 

            in = data
            fmt.println(in,in==nil)     // nil false
        }
      为什么刚刚声明的data和in变量，输出结果是nil，判断结果也是true
      为什么把data变量赋值给变量in之后，输出结果依然是nil，但是判断确实false？
        
    原因：
        interface的判断与想象中不一样的根本原因是，interface不是一个指针类型，虽然看一起来像，但是确实不是！
        
        interface有两类数据结构：
                            1、runtime.iface ：表示包含方法的接口
                            2、runtime.eface ：表示不包含任何方法的空接口，也成为empty interface。
        // 为什么iface能藏住那么多的方法集呢？
        type iface struct{
            tab *itab  //类型
            data unsafe.Pointer //值信息指针
        }
        type itab struct {
            inter *interfacetype  // 接口的类型信息
            _type *_type          // 具体类型
            hash  uint32          //注： _type.hash的副本，用于目标类型和接口变量的类型对比判断
            _     [4]byte       
            fun   [1]uintptr      //注： 底层数组，存储接口的方法集的具体实现地址，其中包含一组函数指针，实现了接口方法的动态分配，且每次在接口发生变更时都会更新
                                    // 问：长度为1的uintptr是如何存储定义的所有方法的？
        }

        接下来进一步展开interfacetype:
        type nameOff int32
        type typeOff int32
        
        type imethod struct {
        name nameOff
        ityp typeOff
        }
        
        type interfacetype struct {
        typ     _type  // 接口的具体类型信息
        pkgpath name   // 接口的package名信息
        mhdr    []imethod // 接口所定义的函数列表
        }
        除了interfacetype，还有各种类型的type：
        例如：maptype、arraytype、chantype、slicetype 等，都是针对具体的类型做的具体类型定义：
        type arraytype struct {
            typ   _type
            elem  *_type
            slice *_type
            len   uintptr
            }
            
            type chantype struct {
            typ  _type
            elem *_type
            dir  uintptr
            }
            ...
        
    
        type eface struct{
            _type *_type  //类型
            data unsafe.Pointer //值信息指针
        }
        
        //类型信息所需要的信息都会存储在这里
        type _type struct {
            size       uintptr  //类型的大小
            ptrdata    uintptr  //包含所有指针的内存前缀的大小
            hash       uint32   //类型的hash值。此处提前计算好，可以避免在哈希表中计算
            tflag      tflag    //额外的类型信息标志，此处为类型的flag标志，用于反射
            align      uint8    //对应变量与该类型的内存对齐大小
            fieldAlign uint8    //对应类型的结构体的内存对齐大小
            kind       uint8    //类型的枚举值。包含go语言中的所有类型：kindBool、kindInt、kindIn8等
            equal func(unsafe.Pointer, unsafe.Pointer) bool  // 用于比较此对象的回调函数
            gcdata    *byte  // 存储垃圾收集器的GC类型数据
            str       nameOff
            ptrToThis typeOff
        }
        注：必须类型和值同时为nil的情况下，interface的nil判断才会是true

    解决方法：
        在不改变类型的情况下，方法之一就是利用反射
        
        func main(){
            var data *byte
            var in interface{}
            
            in = data
            fmt.Println(IsNil(in))
        
        }
    // 利用反射来做nil判断，在反射中会有针对interface类型的特殊处理
        func IsNil(in interface{})bool{
            vi:=reflect.ValueOf(i)
            if vi.Kind()==reflect.Ptr{
                return vi.IsNil()
            }
            return false
        }
        
        其他方法：改变原有的程序逻辑，例如：
        1、对值进行nil判断的，再返回给interface设置
        2、返回具体的值类型，而不是返回interface
    
2、类型的断言
    
    var i interface{} = "煎鱼"
    
    //进行变量断言，若不判断容易出现panic
        s:=i.(string)
    // 进行安全断言
        s,ok:=i.(string)

    在switch case中，还有另外一种写法：
    var i interface{} = "煎鱼"
    
    //进行switch断言
    switch i.(type){
        case string:
            // do something
        case int :
            // do something
        case float64:
            // do something
    }
    采取的是 变量.(type)的调用方式。
    在go语言的背后，类型断言其实是在编译器翻译后，根据iface和eface分别对应了下述方法：
    // 主要是根据接口的类型信息进行新一轮判断和识别：主要核心在于getitab方法
    func assertI2I2(inter *interfacetype, i iface) (r iface, b bool) {
            tab := i.tab
            if tab == nil {
            return
            }
            if tab.inter != inter {
            // 主要核心在于getitab方法
            // getitab 方法的主要作用是获取 itab 元素，若不存在则新增
            tab = getitab(inter, tab._type, true)
            if tab == nil {
            return
            }
            }
            r.tab = tab
            r.data = i.data
            b = true
            return
    }
    func assertI2I(inter *interfacetype, i iface) (r iface)
    
    func assertE2I2(inter *interfacetype, e eface) (r iface, b bool)
    func assertE2I(inter *interfacetype, e eface) (r iface)
        
    类型转换：
    func main() {
        x := "煎鱼"
        var v interface{} = x
        fmt.Println(v)
    }
    // 查看汇编代码：
    0x0021 00033 (main.go:9) LEAQ go.string."煎鱼"(SB), AX
    0x0028 00040 (main.go:9) MOVQ AX, (SP)
    0x002c 00044 (main.go:9) MOVQ $6, 8(SP)
    0x0035 00053 (main.go:9) PCDATA $1, $0
    // 主要对应了 runtime.convTstring 方法。同时很显然其是根据类型来区分来方法：
    0x0035 00053 (main.go:9) CALL runtime.convTstring(SB)
    0x003a 00058 (main.go:9) MOVQ 16(SP), AX
    0x003f 00063 (main.go:10) XORPS X0, X0
    
    func convTstring(val string) (x unsafe.Pointer) {
        if val == "" {
             x = unsafe.Pointer(&zeroVal[0])
        } else {
             x = mallocgc(unsafe.Sizeof(val), stringType, true)
             *(*string)(x) = val
        }
        return
    }
    
    func convT16(val uint16) (x unsafe.Pointer)
    func convT32(val uint32) (x unsafe.Pointer)
    func convT64(val uint64) (x unsafe.Pointer)
    func convTstring(val string) (x unsafe.Pointer)
    func convTslice(val []byte) (x unsafe.Pointer)
    func convT2Enoptr(t *_type, elem unsafe.Pointer) (e eface)
    func convT2I(tab *itab, elem unsafe.Pointer) (i iface)

    动态分派:
    前面有提到接口中的 fun [1]uintptr 属性会可以存储接口的方法集，但不知道为什么。
    接下来我们将进行具体的分析，演示代码：

    type Human interface {
        Say(s string) error
        Eat(s string) error
        Walk(s string) error
    }
    
    type TestA string
    
    func (t TestA) Say(s string) error {
        fmt.Printf("煎鱼：%s\n", s)
        return nil
    }
    func (t TestA) Eat(s string) error {
        fmt.Printf("煎鱼：%s\n", s)
        return nil
    }
    
    func (t TestA) Walk(s string) error {
    fmt.Printf("煎鱼：%s\n", s)
    return nil
    }
    
    func main() {
        var h Human
        var t TestA
        h = t
        _ = h.Eat("烤羊排")
        _ = h.Say("炸鸡翅")
        _ = h.Walk("去炸鸡翅")
    }
    存储方式:
    执行: go build -gcflags '-l' -o awesomeProject .
    编译后，再次执行 :go tool objdump -s "main" awesomeProject。
        // 首个方法的地址
        LEAQ go.itab.main.TestA,main.Human(SB), AX 
        TESTB AL, 0(AX)     
        MOVQ 0x10(SP), AX    
        MOVQ AX, 0x28(SP)    
        // 32
        MOVQ go.itab.main.TestA,main.Human+32(SB), CX
        MOVQ AX, 0(SP)     
        LEAQ go.string.*+3048(SB), DX   
        MOVQ DX, 0x8(SP)    
        MOVQ $0x9, 0x10(SP)    
        CALL CX      
        //24
        MOVQ go.itab.main.TestA,main.Human+24(SB), AX
        MOVQ 0x28(SP), CX    
        MOVQ CX, 0(SP)     
        LEAQ go.string.*+3057(SB), DX   
        MOVQ DX, 0x8(SP)    
        MOVQ $0x9, 0x10(SP)    
        CALL AX      
        //40
        MOVQ go.itab.main.TestA,main.Human+40(SB), AX
        MOVQ 0x28(SP), CX    
        MOVQ CX, 0(SP)     
        LEAQ go.string.*+4973(SB), CX   
        MOVQ CX, 0x8(SP)    
        MOVQ $0xc, 0x10(SP)    
        CALL AX   
        
        结合来看，虽然 fun 属性的类型是 [1]uintptr，只有一个元素，
        但其实就是存放了接口方法集的 "首个方法的地址信息",
        接着根据顺序往后计算并获取就好了。也就是说其实存在一定规律的。在存入方法时就决定了，所以获取也能明确。

        我们进一步展开，看看 itab hash table 是如何获取和新增的。
        getitab 方法的主要作用是获取 itab 元素，若不存在则新增。源码如下：
        func getitab(inter *interfacetype, typ *_type, canfail bool) *itab {
            // 省略一些边界、异常处理
            var m *itab
            //调用 atomic.Loadp 方法加载并查找现有的 itab hash table，看看是否是否可以找到所需的 itab 元素。
            // 若没有找到，则调用 lock 方法对 itabLock 上锁，并进行重试（再一次查找）。
            // 若找到，则跳到 finish 标识的收尾步骤。
            // 若没有找到，则新生成一个 itab 元素，并调用 itabAdd 方法新增到全局的 hash table 中。
            t := (*itabTableType)(atomic.Loadp(unsafe.Pointer(&itabTable)))
            if m = t.find(inter, typ); m != nil {
                goto finish
            }
            
            lock(&itabLock)
            if m = itabTable.find(inter, typ); m != nil {
            unlock(&itabLock)
            goto finish
            }
            
            m = (*itab)(persistentalloc(unsafe.Sizeof(itab{})+uintptr(len(inter.mhdr)-1)*sys.PtrSize, 0, &memstats.other_sys))
            m.inter = inter
            m._type = typ
            m.hash = 0
            m.init()
            itabAdd(m)
            unlock(&itabLock)
            finish:
            if m.fun[0] != 0 {
            // 返回 fun 属性的首位地址，继续后续业务逻辑。
            return m
        }
            
            panic(&TypeAssertionError{concrete: typ, asserted: &inter.typ, missingMethod: m.init()})
            }

    新增 itab 元素:
        itabAdd 方法的主要作用是将所生成好的 itab 元素新增到 itab hash table 中。源码如下：

        func itabAdd(m *itab) {
            // 省略一些边界、异常处理
            t := itabTable
          检查 itab hash table 的容量情况，查看容量情况是否已经满足大于或等于 75%。
            if t.count >= 3*(t.size/4) { // 75% load factor
                // 若满足扩容策略，则调用 mallocgc 方法申请内存，按既有 size 大小扩容双倍容量。
                t2 := (*itabTableType)(mallocgc((2+2*t.size)*sys.PtrSize, nil, true))
                t2.size = t.size * 2
                iterate_itabs(t2.add)
                if t2.count != t.count {
                throw("mismatched count during itab table copy")
            }
            
            atomicstorep(unsafe.Pointer(&itabTable), unsafe.Pointer(t2))
            t = itabTable
            }
            // 若不满足扩容策略，则直接新增 itab 元素到 hash table 中。
            t.add(m)
        }
      
      
3、 go defer万恶的闭包问题

    func main() {
        var whatever [6]struct{}
        for i := range whatever {
            defer func() {
                fmt.Println(i)
            }()
        }
    }
    请自己先想一下输出的结果答案是什么。
    这段程序的输出结果是：
                    5
                    5
                    5
                    5
                    5
                    5
    为什么全是 5，为什么不是 0, 1, 2, 3, 4, 5 这样的输出结果呢？
    其根本原因是闭包所导致的，有两点原因：
    1、在 for 循环结束后，局部变量 i 的值已经是 5 了，并且 defer的闭包是直接引用变量的 i
    2、结合defer 关键字的特性，可得知会在 main 方法主体结束后再执行。
    
    既然了解了为什么，我们再变形一下。再看看另外一种情况，代码如下：
    func main() {
        var whatever [6]struct{}
        for i := range whatever {
            defer func(i int) {
                fmt.Println(i)
            }(i)
        }
    }
    与第一个案例不同，我们这回把变量 i 传了进去。那么他的输出结果是什么呢？
    这段程序的输出结果是：
                    5
                    4
                    3
                    2
                    1
                    0
    为什么是 5, 4, 3, 2, 1, 0 呢，为什么不是 0, 1, 2, 3, 4, 5？
    其根本原因在于两点：
    1、在 for 循环时，局部变量 i 已经传入进 defer func 中 ，属于值传递。其值在 defer 语句声明时的时候就已经确定下来了。
    2、结合 defer 关键字的特性，是按先进后出的顺序来执行的
    
    下一个疑问：
    func f1() (r int) {
        defer func() {
            r++
        }()
        return 0
    }
        
        func f2() (r int) {
        t := 5
        defer func() {
            t = t + 5
        }()
        return t
    }
        
        func f3() (r int) {
            defer func(r int) {
                r = r + 5
            }(r)
        return 1
    }
    
    主函数：
    func main() {
        println(f1())
        println(f2())
        println(f3())
    }
    这段程序的输出结果是：
                    1
                    5
                    1
    为什么是 1, 5, 1 呢，而不是 0, 10, 5，又或是其他答案？
    f1引用了返回值
    f2引用了局部变量，跟返回值没关系
    f3值传递，跟返回值没关系

4、为什么go map和slice是非线性安全的？
    
    slice： 我们使用多个 goroutine 对类型为 slice 的变量进行操作，看看结果会变的怎么样。
    
    func main() {
        var s []string
        for i := 0; i < 9999; i++ {
            go func() {
                 s = append(s, "脑子进煎鱼了")
            }()
        }
        
        fmt.Printf("进了 %d 只煎鱼", len(s))
    }
    输出结果：
            // 第一次执行
            进了 5790 只煎鱼
            // 第二次执行
            进了 7370 只煎鱼
            // 第三次执行
            进了 6792 只煎鱼
    你会发现无论你执行多少次，每次输出的值大概率都不会一样。也就是追加进 slice 的值，出现了覆盖的情况。
    因此在循环中所追加的数量，与最终的值并不相等。且这种情况，是不会报错的，是一个出现率不算高的隐式问题。
    注：这个产生的主要原因是程序逻辑本身就有问题，同时读取到"相同索引位"，自然也就会产生覆盖的写入了
        
    map：同样针对 map 也如法炮制一下。重复针对类型为 map 的变量进行写入。
    func main() {
        s := make(map[string]string)
        for i := 0; i < 99; i++ {
            go func() {
            s["煎鱼"] = "吸鱼"
            }()
        }
        
        fmt.Printf("进了 %d 只煎鱼", len(s))
    }
    输出结果： 在go1.6起会进行原生map的并发检测
    fatal error: concurrent map writes

    goroutine 18 [running]:
    runtime.throw(0x10cb861, 0x15)
    /usr/local/Cellar/go/1.16.2/libexec/src/runtime/panic.go:1117 +0x72 fp=0xc00002e738 sp=0xc00002e708 pc=0x1032472
    runtime.mapassign_faststr(0x10b3360, 0xc0000a2180, 0x10c91da, 0x6, 0x0)
    /usr/local/Cellar/go/1.16.2/libexec/src/runtime/map_faststr.go:211 +0x3f1 fp=0xc00002e7a0 sp=0xc00002e738 pc=0x1011a71
    main.main.func1(0xc0000a2180)
    /Users/eddycjy/go-application/awesomeProject/main.go:9 +0x4c fp=0xc00002e7d8 sp=0xc00002e7a0 pc=0x10a474c
    runtime.goexit()
    /usr/local/Cellar/go/1.16.2/libexec/src/runtime/asm_amd64.s:1371 +0x1 fp=0xc00002e7e0 sp=0xc00002e7d8 pc=0x1063fe1
    created by main.main
    /Users/eddycjy/go-application/awesomeProject/main.go:8 +0x55
    
    程序运行会直接报错。并且是 Go 源码调用 throw 方法所导致的致命错误，也就是说 Go 进程会中断。
    不得不说，这个并发写 map 导致的 fatal error: concurrent map writes 错误提示。我有一个朋友，已经看过少说几十次了，不同组，不同人...
    
    如何支持并发读写：
    对 map 上锁：实际上我们仍然存在并发读写 map 的诉求（程序逻辑决定）
    像是一般写爬虫任务时，基本会用到多个 goroutine，获取到数据后再写入到 map 或者 slice 中去。
    Go 官方在 Go maps in action 中提供了一种简单又便利的方式来实现：
    
    var counter = struct{
        sync.RWMutex
        m map[string]int
    }{m: make(map[string]int)}
    这条语句声明了一个变量，它是一个匿名结构（struct）体，包含一个原生和一个嵌入读写锁 sync.RWMutex。
    
    要想从变量中中读出数据，则调用读锁：
    counter.RLock()
    n := counter.m["煎鱼"]
    counter.RUnlock()
    fmt.Println("煎鱼:", n)
    
    要往变量中写数据，则调用写锁：
    counter.Lock()
    counter.m["煎鱼"]++
    counter.Unlock()
    这就是一个最常见的 Map 支持并发读写的方式了。

    sync.Map:其是专门为 append-only 场景设计的，也就是适合读多写少的场景
    
    虽然有了 Map+Mutex 的极简方案，但是也仍然存在一定问题。那就是在 map 的数据量非常大时，
    只有一把锁（Mutex）就非常可怕了，一把锁会导致大量的争夺锁，导致各种冲突和性能低下。
    常见的解决方案是分片化，将一个大 map 分成多个区间，各区间使用多个锁，这样子锁的粒度就大大降低了。不过该方案实现起来很复杂，很容易出错。
    因此 Go 团队到此时为止暂无推荐，而是采取了其他方案。
    该方案就是在 Go1.9 起支持的 sync.Map，其支持并发读写 map，起到一个补充的作用。

    Go 语言的 sync.Map 支持并发读写 map，采取了 “空间换时间” 的机制，
    冗余了两个数据结构，分别是：read 和 dirty，减少加锁对性能的影响：
    type Map struct {
        mu Mutex
        read atomic.Value // readOnly
        dirty map[interface{}]*entry
        misses int
    }

    注：其是专门为 append-only 场景设计的，也就是适合读多写少的场景。这是他的优点之一。
    若出现写多/并发多的场景，会导致 read map 缓存失效，需要加锁，冲突变多，性能急剧下降。这是他的重大缺点。
    
    提供了以下常用方法：
    //Delete：删除某一个键的值。
    func (m *Map) Delete(key interface{})
    //Load：返回存储在 map 中的键的值，如果没有值，则返回 nil。ok 结果表示是否在 map 中找到了值。
    func (m *Map) Load(key interface{}) (value interface{}, ok bool)
    //LoadAndDelete：删除一个键的值，如果有的话返回之前的值。
    func (m *Map) LoadAndDelete(key interface{}) (value interface{}, loaded bool)
    //LoadOrStore：如果存在的话，则返回键的现有值。否则，它存储并返回给定的值。如果值被加载，加载的结果为 true，如果被存储，则为 false。
    func (m *Map) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool)
    //Range：递归调用，对 map 中存在的每个键和值依次调用闭包函数 f。如果 f 返回 false 就停止迭代。
    func (m *Map) Range(f func(key, value interface{}) bool)
    //Store：存储并设置一个键的值。
    func (m *Map) Store(key, value interface{})
    
        var m sync.Map

        func main() {
        //写入
        data := []string{"煎鱼", "咸鱼", "烤鱼", "蒸鱼"}
        for i := 0; i < 4; i++ {
            go func(i int) {
                m.Store(i, data[i])
            }(i)
        }
        time.Sleep(time.Second)
        
        //读取
        v, ok := m.Load(0)
        fmt.Printf("Load: %v, %v\n", v, ok) //Load: 煎鱼, true
        
        //删除
        m.Delete(1)
        
        //读或写
        v, ok = m.LoadOrStore(1, "吸鱼")
        fmt.Printf("LoadOrStore: %v, %v\n", v, ok) //LoadOrStore: 吸鱼, false
        
        //遍历
        m.Range(func(key, value interface{}) bool {
            fmt.Printf("Range: %v, %v\n", key, value)
            return true
        })
       // Range: 0, 煎鱼
       // Range: 1, 吸鱼
       // Range: 3, 蒸鱼
       // Range: 2, 烤鱼
        }
    
    Go Slice 的话，主要还是索引位覆写问题，这个就不需要纠结了，势必是程序逻辑在编写上有明显缺陷，自行改之就好。
    但 Go map 就不大一样了，很多人以为是默认支持的，一个不小心就翻车，这么的常见。那凭什么 Go 官方还不支持，难不成太复杂了，性能太差了，到底是为什么？
    
    原因如下:
    1、典型使用场景：map 的典型使用场景是不需要从多个 goroutine 中进行安全访问。
    2、非典型场景（需要原子操作）：map 可能是一些更大的数据结构或已经同步的计算的一部分。
    3、性能场景考虑：若是只是为少数程序增加安全性，导致 map 所有的操作都要处理 mutex，将会降低大多数程序的性能。
    
    汇总来讲：就是 Go 官方在经过了长时间的讨论后，认为 Go map 更应适配典型使用场景，而不是为了小部分情况，
    导致大部分程序付出代价（性能），决定了不支持。

5、go sync.map和原生map谁的性能更高，为什么？
    
    map 的两种目前在业界使用的最多的并发支持的模式。分别是：
    1、原生 map + 互斥锁或读写锁 mutex。
    2、标准库 sync.Map（Go1.9及以后）。
    
    这两种到底怎么选，谁的性能更加的好？
    我们先会了解清楚什么场景下，Go map 的多种类型怎么用，谁的性能最好！
    
    在 Go 官方文档中明确指出 Map 类型的一些建议：
    1、多个 goroutine 的并发使用是安全的，不需要额外的锁定或协调控制
    2、大多数代码应该使用原生的 map，而不是单独的锁定或协调控制，以获得更好的类型安全性和维护性。

    同时 Map 类型，还针对以下场景进行了性能优化：
    1、当一个给定的键的条目只被写入一次但被多次读取时。例如在仅会增长的缓存中，就会有这种业务场景。
    2、当多个 goroutines 读取、写入和覆盖不相干的键集合的条目时。
    
    这两种情况与 Go map 搭配单独的 Mutex 或 RWMutex 相比较，使用 Map 类型可以大大减少锁的争夺。
    性能测试：
    
    // 代表互斥锁
    type FooMap struct {
        sync.Mutex
        data map[int]int
    }
    
    // 代表读写锁
    type BarRwMap struct {
        sync.RWMutex
        data map[int]int
    }
    
    var fooMap *FooMap
    var barRwMap *BarRwMap
    var syncMap *sync.Map
    
    // 初始化基本数据结构
    func init() {
        fooMap = &FooMap{data: make(map[int]int, 100)}
        barRwMap = &BarRwMap{data: make(map[int]int, 100)}
        syncMap = &sync.Map{}
    }

    在配套方法上，常见的增删改查动作我们都编写了相应的方法。用于后续的压测（只展示部分代码）：
    func builtinRwMapStore(k, v int) {
        barRwMap.Lock()
        defer barRwMap.Unlock()
        barRwMap.data[k] = v
    }
    
    func builtinRwMapLookup(k int) int {
        barRwMap.RLock()
        defer barRwMap.RUnlock()
        if v, ok := barRwMap.data[k]; !ok {
             return -1
        } else {
            return v
        }
    }
    
    func builtinRwMapDelete(k int) {
        barRwMap.Lock()
        defer barRwMap.Unlock()
        if _, ok := barRwMap.data[k]; !ok {
            return
        } else {
            delete(barRwMap.data, k)
        }
    }

    压测方法基本代码如下：
    这块主要就是增删改查的代码和压测方法的准备，压测代码直接复用的是大白大佬的 go19-examples/benchmark-for-map 项目。
    func BenchmarkBuiltinRwMapDeleteParalell(b *testing.B) {
        b.RunParallel(func(pb *testing.PB) {
        r := rand.New(rand.NewSource(time.Now().Unix()))
        for pb.Next() {
            k := r.Intn(100000000)
            builtinRwMapDelete(k)
        }
        })
    }
    压测结果：
    在写入元素上：最慢的是 sync.map 类型，其次是原生 map+互斥锁（Mutex），最快的是原生 map+读写锁（RwMutex）。
    总体的排序（从慢到快）为：SyncMapStore < MapStore < RwMapStore。
    
    在查找元素上：最慢的是原生 map+互斥锁，其次是原生 map+读写锁。最快的是 sync.map 类型。
    总体的排序为：MapLookup < RwMapLookup < SyncMapLookup。

    在删除元素上，最慢的是原生 map+读写锁，其次是原生 map+互斥锁，最快的是 sync.map 类型。
    总体的排序为：RwMapDelete < MapDelete < SyncMapDelete。

    根据上述的压测结果，我们可以得出 sync.Map 类型：
    1、在读和删场景上的性能是最佳的，领先一倍有多。
    2、在写入场景上的性能非常差，落后原生 map+锁整整有一倍之多。
    
    因此在实际的业务场景中。假设是读多写少的场景，会更建议使用 sync.Map 类型。
    但若是那种写多的场景，例如多 goroutine 批量的循环写入，那就建议另辟途径了，性能不忍直视（无性能要求另当别论）。

    sync.Map 剖析：
    为什么 sync.Map 类型的测试结果这么的 “偏科”，为什么读操作性能这么高，写操作性能低的可怕，他是怎么设计的？
    sync.Map 类型的底层数据结构如下：
    type Map struct {
        mu Mutex //互斥锁，用于保护 read 和 dirty。
        read atomic.Value // readOnly 只读数据，支持并发读取（atomic.Value 类型）。如果涉及到更新操作，则只需要加锁来保证数据安全。
                          //read 实际存储的是 readOnly 结构体，内部也是一个原生 map，amended 属性用于标记 read 和 dirty 的数据是否一致。
        dirty map[interface{}]*entry //dirty：读写数据，是一个原生 map，也就是非线程安全。操作 dirty 需要加锁来保证数据安全。
        misses int  //misses：统计有多少次读取 read 没有命中。每次 read 中读取失败后，misses 的计数值都会加 1。
    }
    
    // Map.read 属性实际存储的是 readOnly。
    type readOnly struct {
        m       map[interface{}]*entry
        amended bool  //amended 属性用于标记 read 和 dirty 的数据是否一致。
    }
    
    在 read 和 dirty 中，都有涉及到的结构体：
    type entry struct {
        p unsafe.Pointer // *interface{} 其包含一个指针 p, 用于指向用户存储的元素（key）所指向的 value 值。
    }
    在此建议你必须搞懂 read、dirty、entry，再往下看，食用效果会更佳，后续会围绕着这几个概念流转。

    查找过程:
    划重点，sync.Map 类型本质上是有两个 “map”。一个叫 read、一个叫 dirty，长的也差不多
    当我们从 sync.Map 类型中读取数据时，其会先查看 read 中是否包含所需的元素：
    1、若有，则通过 atomic 原子操作读取数据并返回。
    2、若无，则会判断 read.readOnly 中的 amended 属性，他会告诉程序 dirty 是否包含 read.readOnly.m 中没有的数据；
            因此若存在，也就是 amended 为 true，将会进一步到 dirty 中查找数据。
    注：sync.Map 的读操作性能如此之高的原因，就在于存在 read 这一巧妙的设计，其"作为一个缓存层"，提供了快路径（fast path）的查找。
    同时其结合 amended 属性，配套解决了每次读取都涉及锁的问题，实现了读这一个使用场景的高性能。

    写入过程：
    我们直接关注 sync.Map 类型的 Store 方法，该方法的作用是新增或更新一个元素
    func (m *Map) Store(key, value interface{}) {
        read, _ := m.read.Load().(readOnly)
        if e, ok := read.m[key]; ok && e.tryStore(&value) {
            return
        }
        ...
    }
    调用 Load 方法检查 m.read 中是否存在这个元素。若存在，且没有被标记为删除状态，则尝试存储。
    若该元素不存在或已经被标记为删除状态，则继续走到下面流程：
    func (m *Map) Store(key, value interface{}) {
    ...
    m.mu.Lock()
        read, _ = m.read.Load().(readOnly)
        if e, ok := read.m[key]; ok {
            if e.unexpungeLocked() {
            m.dirty[key] = e
        }
        e.storeLocked(&value)
        } else if e, ok := m.dirty[key]; ok {
             e.storeLocked(&value)
        } else {
        if !read.amended {
            m.dirtyLocked()
            m.read.Store(readOnly{m: read.m, amended: true})
        }
        m.dirty[key] = newEntry(value)
        }
    m.mu.Unlock()
    }
    
    由于已经走到了 dirty 的流程，因此开头就直接调用了 Lock 方法上互斥锁，保证数据安全，也是凸显性能变差的第一幕。
    其分为以下三个处理分支：
    
    1、若发现 read 中存在该元素，但已经被标记为已删除（expunged），则说明 dirty 不等于 nil（dirty 中肯定不存在该元素）。其将会执行如下操作。
        将元素状态从已删除（expunged）更改为 nil。
        将元素插入 dirty 中。
    2、若发现 read 中不存在该元素，但 dirty 中存在该元素，则直接写入更新 entry 的指向。
    3、若发现 read 和 dirty 都不存在该元素，则从 read 中复制未被标记删除的数据，并向 dirty 中插入该元素，赋予元素值 entry 的指向。
    
    写入过程的整体流程就是：
    1、查 read，read 上没有，或者已标记删除状态。
    2、上互斥锁（Mutex）。
    3、操作 dirty，根据各种数据情况和状态进行处理。
    4、回到最初的话题，为什么他写入性能差那么多。究其原因：
    5、写入一定要会经过 read，无论如何都比别人多一层，后续还要查数据情况和状态，性能开销相较更大。
    （第三个处理分支）当初始化或者 dirty 被提升后，会从 read 中复制全量的数据，若 read 中数据量大，则会影响性能。
    可得知 sync.Map 类型不适合写多的场景，读多写少是比较好的
    若有大数据量的场景，则需要考虑 read 复制数据时的偶然性能抖动是否能够接受。

    删除过程：只是标记为删除
    写入过程，理论上和删除不会差太远。怎么 sync.Map 类型的删除的性能似乎还行，这里面有什么猫腻？
    func (m *Map) LoadAndDelete(key interface{}) (value interface{}, loaded bool) {
        read, _ := m.read.Load().(readOnly)
        e, ok := read.m[key]
        ...
        if ok {
            return e.delete()
        }
    }
    删除是标准的开场，依然先到 read 检查该元素是否存在。：
    1、若存在，则调用 delete 标记为 expunged（删除状态），非常高效。可以明确在 read 中的元素，被删除，性能是非常好的。
    2、若不存在，也就是走到 dirty 流程中：
        func (m *Map) LoadAndDelete(key interface{}) (value interface{}, loaded bool) {
            ...
            if !ok && read.amended {
                m.mu.Lock()
                read, _ = m.read.Load().(readOnly)
                e, ok = read.m[key]
    
                if !ok && read.amended {
                e, ok = m.dirty[key]
                delete(m.dirty, key)
                m.missLocked()
                }
                m.mu.Unlock()
            }
            ...
            return nil, false
        }
        若 read 中不存在该元素，dirty 不为空，read 与 dirty 不一致（利用 amended 判别），则表明要操作 dirty，上互斥锁。
        再重复进行双重检查，若 read 仍然不存在该元素。则调用 delete 方法从 dirty 中标记该元素的删除。
        需要注意，出现频率较高的 delete 方法：
        func (e *entry) delete() (value interface{}, ok bool) {
            for {
                p := atomic.LoadPointer(&e.p)
                if p == nil || p == expunged {
                return nil, false
            }
            if atomic.CompareAndSwapPointer(&e.p, p, nil) {
                 return *(*interface{})(p), true
            }
            }
        }
    该方法都是将 entry.p 置为 nil，并且标记为 expunged（删除状态），而不是真真正正的删除。

    注：不要误用 sync.Map，前段时间从字节大佬分享的案例来看，他们将一个连接作为 key 放了进去，
        于是和这个连接相关的，例如：buffer 的内存就永远无法释放了...

    我们针对 sync.Map 的性能差异，进行了深入的源码剖析，了解到了其背后快、慢的原因，实现了知其然知其所以然。

5、为什么go map的负载因子是6.5？
    
    map中最中重要的一个基本单本是hmap：
        type hmap struct {
        count     int    // map 的大小，也就是 len() 的值，代指 map 中的键值对个数。
        flags     uint8  // 状态标识，主要是 goroutine 写入和扩容机制的相关状态控制。并发读写的判断条件之一就是该值。
        B         uint8  // 桶，最大可容纳的元素数量，值为 负载因子（默认 6.5） * 2 ^ B，是 2 的指数。其值就是 6.5
        noverflow uint16 // 溢出桶的数量。
        hash0     uint32 // 哈希因子。
        buckets    unsafe.Pointer // 保存当前桶数据的指针地址（指向一段连续的内存地址，主要存储键值对数据）。
        oldbuckets unsafe.Pointer // 保存旧桶的指针地址。
        nevacuate  uintptr        // 迁移进度。
        extra *mapextra // 原有 buckets 满载后，会发生扩容动作，在 Go 的机制中使用了增量扩容，如下为细项：
        }
        
        type mapextra struct {
        overflow    *[]*bmap  // 为 hmap.buckets （当前）溢出桶的指针地址。
        oldoverflow *[]*bmap  // 为 hmap.oldbuckets （旧）溢出桶的指针地址。
        nextOverflow *bmap    // 为空闲溢出桶的指针地址。
        }
    我们关注到 hmap 的 B 字段，其值就是 6.5，他就是我们在苦苦寻找的 6.5，但他又是什么呢？
        什么是负载因子?
        B 值，这里就涉及到一个概念：负载因子（load factor），用于衡量当前哈希表中"空间占用率"的核心指标，
        也就是每个 bucket 桶存储的"平均元素个数"。
    
        另外负载因子与扩容、迁移等重新散列（rehash）行为有直接关系：
        1、在程序运行时，会不断地进行插入、删除等，会导致 bucket 不均，内存利用率低，需要迁移。
        2、在程序运行时，出现负载因子过大，需要做扩容，解决 bucket 过大的问题
        
        负载因子是哈希表中的一个重要指标，在各种版本的哈希表实现中都有类似的东西，
        主要目的是为了"平衡 buckets 的存储空间大小和查找元素时的性能高低"。
        
        在接触各种哈希表时都可以关注一下，做不同的对比，看看各家的考量。
        
        为什么是 6.5：为什么 Go 语言中哈希表的负载因子是 6.5，为什么不是 8 ，也不是 1。这里面有可靠的数据支撑吗？
        实际上这是 Go 官方的经过认真的测试得出的数字，一起来看看官方的这份测试报告。
        
        原因：Go 官方发现：负载因子太大了，会有很多溢出的桶。太小了，就会浪费很多空间（too large and we have lots of overflow buckets, too small and we waste a lot of space）。
        根据这份测试结果和讨论，Go 官方把 Go 中的 map 的负载因子硬编码为 6.5，这就是 6.5 的选择缘由。
        
        这意味着在 Go 语言中，当 B（bucket）平均每个存储的元素大于或等于 6.5 时，就会触发扩容行为，
        这是作为我们用户对这个数值最近的接触。

[怎么避免内存逃逸?]
    在 runtime/stub.go:133 有个noescape函数。noescape可以在逃逸分析中隐藏一个指针。
    让这个指针在逃逸分析中不会被检查为逃逸。

    // noescape hides a pointer from escape analysis.  noescape is
    // the identity function but escape analysis doesn't think the
    // output depends on the input.  noescape is inlined and currently
    // compiles down to zero instructions.
    // USE CAREFULLY!
    //go:nosplit
    
    func noescape(p unsafe.Pointer) unsafe.Pointer {
         x := uintptr(p)
    return unsafe.Pointer(x ^ 0)
    }
    
    举例：
    使用 go build -gcflags=-m查看逃逸分析情况：
    
    package main
    
    import (
    "unsafe"
    )
    
    type A struct {
    S *string
    }
    
    func (f *A) String() string {
    return *f.S
    }
    
    type ATrick struct {
    S unsafe.Pointer
    }
    
    func (f *ATrick) String() string {
    return *(*string)(f.S)
    }
    
    func NewA(s string) A {
    return A{S: &s}
    }
    
    func NewATrick(s string) ATrick {
    return ATrick{S: noescape(unsafe.Pointer(&s))}
    }
    
    func noescape(p unsafe.Pointer) unsafe.Pointer {
    x := uintptr(p)
    return unsafe.Pointer(x ^ 0)
    }
    
    func main() {
    s := "hello"
    f1 := NewA(s)
    f2 := NewATrick(s)
    s1 := f1.String()
    s2 := f2.String()
    _ = s1 + s2
    }

    执行go build -gcflags=-m main.go：
    
    $go build -gcflags=-m main.go
    # command-line-arguments
    ./main.go:11:6: can inline (*A).String
    ./main.go:19:6: can inline (*ATrick).String
    ./main.go:23:6: can inline NewA
    ./main.go:31:6: can inline noescape
    ./main.go:27:6: can inline NewATrick
    ./main.go:28:29: inlining call to noescape
    ./main.go:36:6: can inline main
    ./main.go:38:14: inlining call to NewA
    ./main.go:39:19: inlining call to NewATrick
    ./main.go:39:19: inlining call to noescape
    ./main.go:40:17: inlining call to (*A).String
    ./main.go:41:17: inlining call to (*ATrick).String
    /var/folders/45/qx9lfw2s2zzgvhzg3mtzkwzc0000gn/T/go-build763863171/b001/_gomod_.go:6:6: can inline init.0
    ./main.go:11:7: leaking param: f to result ~r0 level=2
    ./main.go:19:7: leaking param: f to result ~r0 level=2
    ./main.go:24:16: &s escapes to heap  // 这个是NewA中的逃逸
    ./main.go:23:13: moved to heap: s
    ./main.go:27:18: NewATrick s does not escape //NewATrick里的s却没有逃逸
    ./main.go:28:45: NewATrick &s does not escape
    ./main.go:31:15: noescape p does not escape
    ./main.go:38:14: main &s does not escape
    ./main.go:39:19: main &s does not escape
    ./main.go:40:10: main f1 does not escape
    ./main.go:41:10: main f2 does not escape
    ./main.go:42:9: main s1 + s2 does not escape

分析：

    1、上段代码对A和ATrick同样的功能有两种实现：他们包含一个 string ，
        然后用 String() 方法返回这个字符串。但是从逃逸分析看ATrick 版本没有逃逸。
    2、noescape() 函数的作用是遮蔽输入和输出的依赖关系。
        使编译器不认为 p 会通过 x 逃逸， 因为 uintptr() 产生的引用是编译器无法理解的。
    3、内置的 uintptr 类型是一个真正的指针类型，但是在编译器层面，
        它只是一个存储一个 指针地址 的 int 类型。代码的最后一行返回 unsafe.Pointer 也是一个 int。
    4、noescape() 在 runtime 包中使用 unsafe.Pointer 的地方被大量使用。
        如果作者清楚被 unsafe.Pointer 引用的数据肯定不会被逃逸，但编译器却不知道的情况下，这是很有用的。
    5、面试中秀一秀是可以的，如果在实际项目中如果使用这种unsafe包大概率会被同事打死。
        不建议使用！  毕竟包的名字就叫做 unsafe, 而且源码中的注释也写明了 USE CAREFULLY!。
