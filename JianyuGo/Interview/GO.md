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
    

[数据结构]