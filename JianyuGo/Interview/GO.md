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
* go是传值还是传引用--》传值，传递的总是参数的副本（go没有引用传递）
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


[调度模型]


[数据结构]