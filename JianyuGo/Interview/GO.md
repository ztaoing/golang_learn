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
    
        
    

[调度模型]


[数据结构]