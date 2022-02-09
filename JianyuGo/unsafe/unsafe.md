go标准库中大量使用了unsafe.pointer,想要更好的理解源码实现，就要知道unsafe.pointer到底是什么？

1、什么是unsafe
众所周知，Go语言被设计成一门强类型的静态语言，那么他的类型就不能改变了，静态也是意味着类型检查在运行前就做了。
所以在Go语言中是不允许两个指针类型进行转换的，使用过C语言的朋友应该知道这在C语言中是可以实现的，
Go中不允许这么使用是处于安全考虑，毕竟强制转型会引起各种各样的麻烦，有时这些麻烦很容易被察觉，
有时他们却又隐藏极深，难以察觉。大多数读者可能不明白为什么类型转换是不安全的，这里用C语言举一个简单的例子：

        int main(){
            double pi = 3.1415926;
            double *pv = &pi;
            void *temp = pd;
            int *p = temp;
        }
    在标准C语言中，任何非void类型的指针都可以和void类型的指针相互指派，也可以通过void类型指针作为中介，
    实现不同类型的指针间接相互转换。上面示例中，指针pv指向的空间本是一个双精度数据，占8个字节，
    但是经过转换后，p指向的是一个4字节的int类型。这种发生内存截断的设计缺陷会在转换后进行内存访问是存在安全隐患。
    我想这就是Go语言被设计成强类型语言的原因之一吧。
    
    虽然类型转换是不安全的，但是在一些特殊场景下，使用了它，可以打破Go的类型和内存安全机制，
    可以绕过低效的类型系统，提高运行效率。所以Go标准库中提供了一个unsafe包，之所以叫这个名字，
    就是不推荐大家使用，但是不是不能用，如果你掌握的特别好，还是可以实践的。

    unsafe实现原理：unsafe包中只提供了3种方法：
    func Sizeof(x ArbitraryType) uintptr
    func Offsetof(x ArbitraryType) uintptr
    func Alignof(x ArbitraryType) uintptr

    Sizeof(x ArbitrayType)方法主要作用是用返回类型x所占据的字节数，但并不包含x所指向的内容的大小，
    与C语言标准库中的Sizeof()方法功能一样，比如在32位机器上，一个指针返回大小就是4字节。

    Offsetof(x ArbitraryType)方法主要作用是返回结构体成员在内存中的位置离结构体起始处(结构体的第一个字段的偏移量都是0)的字节数，
    即偏移量，我们在注释中看一看到其入参必须是一个结构体，其返回值是一个常量。
    
    Alignof(x ArbitratyType)的主要作用是返回一个类型的对齐值，也可以叫做对齐系数或者对齐倍数。
    对齐值是一个和内存对齐有关的值，合理的内存对齐可以提高内存读写的性能。一般对齐值是2^n，最大不会超过8(受内存对齐影响).
    获取对齐值还可以使用反射包的函数，也就是说：unsafe.Alignof(x)等价于reflect.TypeOf(x).Align()。
    对于任意类型的变量x，unsafe.Alignof(x)至少为1。对于struct结构体类型的变量x，计算x每一个字段f的unsafe.Alignof(x，f)，
    unsafe.Alignof(x)等于其中的最大值。对于array数组类型的变量x，unsafe.Alignof(x)等于构成数组的元素类型的对齐倍数。
    没有任何字段的空struct{}和没有任何元素的array占据的内存空间大小为0，不同大小为0的变量可能指向同一块地址。

    细心的朋友会发发现这三个方法返回的都是uintptr类型，这个目的就是可以和unsafe.poniter类型相互转换，
    因为*T是不能计算偏移量的，也不能进行计算，但是uintptr是可以的，所以可以使用uintptr类型进行计算，
    这样就可以可以访问特定的内存了，达到对不同的内存读写的目的。三个方法的入参都是ArbitraryType类型，
    代表着任意类型的意思，同时还提供了一个Pointer指针类型，即像void *一样的通用型指针。

    type ArbitraryType int
    type Pointer *ArbitraryType
    // uintptr 是一个整数类型，它足够大，可以存储
    type uintptr uintptr

    *T：普通类型指针类型，用于传递对象地址，不能进行指针运算。
    unsafe.poniter：通用指针类型，用于转换不同类型的指针，不能进行指针运算，不能读取内存存储的值(需转换到某一类型的普通指针)
    uintptr：用于指针运算，GC不把uintptr当指针，uintptr无法持有对象。uintptr类型的目标会被回收。

    三者关系就是：unsafe.Pointer是桥梁，可以让任意类型的指针实现相互转换，
    也可以将任意类型的指针转换为uintptr进行指针运算，也就说uintptr是用来与unsafe.Pointer打配合，用于指针运算。画个图表示一下：
    
    unsafe.Pointer基本使用:
    我们在上一篇分析atomic.Value源码时，看到atomic/value.go中定义了一个ifaceWords结构，
    其中typ和data字段类型就是unsafe.Poniter，这里使用unsafe.Poniter类型的原因是传入的值就是interface{}类型，
    使用unsafe.Pointer强转成ifaceWords类型，这样可以把类型和值都保存了下来，方便后面的写入类型检查。
    截取部分代码如下：

    // ifaceWords is interface{} internal representation.
        type ifaceWords struct {
            typ  unsafe.Pointer
            data unsafe.Pointer
        }
        // Load returns the value set by the most recent Store.
        // It returns nil if there has been no call to Store for this Value.
        func (v *Value) Load() (x interface{}) {
        vp := (*ifaceWords)(unsafe.Pointer(v))
        for {
            typ := LoadPointer(&vp.typ) // 读取已经存在值的类型
            /**
            ..... 中间省略
            **/
            // First store completed. Check type and overwrite data.
            if typ != xp.typ { //当前类型与要存入的类型做对比
                panic("sync/atomic: store of inconsistently typed value into Value")
            }
        }
    
        上面就是源码中使用unsafe.Pointer的一个例子，有一天当你准备读源码时，unsafe.pointer的使用到处可见。
        好啦，接下来我们写一个简单的例子，看看unsafe.Pointer是如何使用的。
        
        func main()  {
            number := 5
            pointer := &number
            fmt.Printf("number:addr:%p, value:%d\n",pointer,*pointer)
            
            float32Number := (*float32)(unsafe.Pointer(pointer))
            *float32Number = *float32Number + 3
            
            fmt.Printf("float64:addr:%p, value:%f\n",float32Number,*float32Number)
        }
        
        运行结果：
        number:addr:0xc000018090, value:5
        float64:addr:0xc000018090, value:3.000000
        由运行可知使用unsafe.Pointer强制类型转换后指针指向的地址是没有改变，只是类型发生了改变。
        这个例子本身没什么意义，正常项目中也不会这样使用。

        注：总结一下基本使用：先把*T类型转换成unsafe.Pointer类型，然后在进行强制转换转成你需要的指针类型即可。
        
        Sizeof、Alignof、Offsetof三个函数的基本使用：
        type User struct {
            Name string
            Age uint32
            Gender bool // 男:true 女：false 就是举个例子别吐槽我这么用。。。。
        }
        
        func func_example()  {
            // sizeof
            fmt.Println(unsafe.Sizeof(true))
            fmt.Println(unsafe.Sizeof(int8(0)))
            fmt.Println(unsafe.Sizeof(int16(10)))
            fmt.Println(unsafe.Sizeof(int(10)))
            fmt.Println(unsafe.Sizeof(int32(190)))
            fmt.Println(unsafe.Sizeof("asong"))
            fmt.Println(unsafe.Sizeof([]int{1,3,4}))

            // Offsetof
            user := User{Name: "Asong", Age: 23,Gender: true}
            userNamePointer := unsafe.Pointer(&user)
            
            nNamePointer := (*string)(unsafe.Pointer(userNamePointer))
            *nNamePointer = "Golang梦工厂"
            
            nAgePointer := (*uint32)(unsafe.Pointer(uintptr(userNamePointer) + unsafe.Offsetof(user.Age)))
            *nAgePointer = 25
            
            nGender := (*bool)(unsafe.Pointer(uintptr(userNamePointer)+unsafe.Offsetof(user.Gender)))
            *nGender = false
            
            fmt.Printf("u.Name: %s, u.Age: %d,  u.Gender: %v\n", user.Name, user.Age,user.Gender)
            // Alignof
            var b bool
            var i8 int8
            var i16 int16
            var i64 int64
            var f32 float32
            var s string
            var m map[string]string
            var p *int32
            
            fmt.Println(unsafe.Alignof(b))
            fmt.Println(unsafe.Alignof(i8))
            fmt.Println(unsafe.Alignof(i16))
            fmt.Println(unsafe.Alignof(i64))
            fmt.Println(unsafe.Alignof(f32))
            fmt.Println(unsafe.Alignof(s))
            fmt.Println(unsafe.Alignof(m))
            fmt.Println(unsafe.Alignof(p))
        }
