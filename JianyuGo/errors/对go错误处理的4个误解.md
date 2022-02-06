熟悉C语言的朋友知道在C语言中默认情况下不初始化局部变量。
    
    未初始化的变量可以包含任何值，其使用会导致未定义的行为；
    如果我们未初始局部变量，在编译时就会报警告 C4700，这个警告指示一个Bug，
    这个Bug可能导致程序中出现不可预测的结果或故障。而在Go语言就不会有这样的问题，
    Go语言的设计者吸取了在设计C语言时的一些经验，所以Go语言的零值规范如下：

以下内容来自官方blog：https://golang.org/ref/spec#The_zero_value
当通过声明或 new 调用为变量分配存储空间时，或通过复合文字或 make 调用创建新值时，且未提供显式初始化，
则给出变量或值一个默认值。
此类变量或值的每个元素都为其类型设置为零值：

    布尔型为 false，
    数字类型为 0，
    字符串为 ""，
    指针、函数、接口、切片、通道和映射为 nil。
此初始化是递归完成的，例如，如果未指定任何值，则结构体数组的每个元素的字段都将被清零（覆盖）。
例如这两个简单的声明是等价的：

    var i int 
    var i int = 0
在或者这个结构体的声明：

    type T struct { i int; f float64; next *T }
    t := new(T)
这个结构体t中成员字段零值如下：

    t.i == 0
    t.f == 0.0
    t.next == nil

Go语言中这种始终将值设置为已知默认值的特性对于程序的安全性和正确性起到了很重要的作用，这样也使整个Go程序更简单、更紧凑。

[零值有什么用]
* 通过零值来提供默认值
  
  我们在看一些Go语言库的时候，都会看到在初始化对象时采用-->"动态初始化"的模式，
  其实就是在创建对象时判断如果是零值就使用默认值，
  比如我们在分析hystrix-go这个库时，在配置Command时就是使用的这种方式：


      func ConfigureCommand(name string, config CommandConfig) {
          settingsMutex.Lock()
          defer settingsMutex.Unlock()
    
        timeout := DefaultTimeout
        // 通过零值判断进行默认值赋值，增强了go程序的健壮性
        if config.Timeout != 0 {
            timeout = config.Timeout
        }
        
        max := DefaultMaxConcurrent
        if config.MaxConcurrentRequests != 0 {
            max = config.MaxConcurrentRequests
        }
        
        volume := DefaultVolumeThreshold
        if config.RequestVolumeThreshold != 0 {
            volume = config.RequestVolumeThreshold
        }
        
        sleep := DefaultSleepWindow
        if config.SleepWindow != 0 {
            sleep = config.SleepWindow
        }
        
        errorPercent := DefaultErrorPercentThreshold
        if config.ErrorPercentThreshold != 0 {
            errorPercent = config.ErrorPercentThreshold
        }
        
        circuitSettings[name] = &Settings{
        Timeout:                time.Duration(timeout) * time.Millisecond,
        MaxConcurrentRequests:  max,
        RequestVolumeThreshold: uint64(volume),
        SleepWindow:            time.Duration(sleep) * time.Millisecond,
        ErrorPercentThreshold:  errorPercent,
        }
        }

[开箱即用]
为什么叫开箱即用呢？因为Go语言的零值让程序变得更简单了，
    有些场景我们不需要显示初始化就可以直接用，举几个例子：

* 切片，他的零值是nil，即使不用make进行初始化也是可以直接使用的，例如
  
        package main
    
        import (
        "fmt"
        "strings"
        )
        
        func main() {
        var s []string
        
            s = append(s, "asong")
            s = append(s, "真帅")
            fmt.Println(strings.Join(s, " "))
        }
但是零值也并不是万能的，零值切片不能直接进行赋值操作：这样的程序就报错了。
    
    var s []string
    s[0] = "asong真帅"

* 方法接收者的归纳
  
  利用零值可用的特性，我们配合空结构体的方法接收者特性，可以将方法组合起来，在业务代码中便于后续扩展和维护：
  我在一些开源项目中看到很多地方都这样使用了，这样的代码最结构化～。

    // 不包含任何属性，用于将方法组合起来
    type T struct{}

    func (t *T) Run() {
        fmt.Println("we run")
    }
    
    func main() {
        var t T
        t.Run()
    }
  
* 标准库无需显示初始化
  我们经常使用sync包中的mutex、once、waitgroup都是无需显示初始化即可使用，
  拿mutex包来举例说明，我们看到mutex的结构如下：
  

      type Mutex struct {
          state int32
          sema  uint32
      }

    这两个字段在未显示初始化时默认零值都是0，所以我们就看到上锁代码就针对这个特性来写的：
    func (m *Mutex) Lock() {

    // Fast path: grab unlocked mutex.
    // 上锁代码就针对这个特性来写的
    // 原子操作交换时使用的old值就是0，这种设计让mutex调用者无需考虑对mutex的初始化则可以直接使用。
    if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
        if race.Enabled {
            race.Acquire(unsafe.Pointer(m))
        }
         return
    }

    // Slow path (outlined so that the fast path can be inlined)
    m.lockSlow()
    }
    还有一些其他标准库也使用零值可用的特性，使用方法都一样，就不在举例了。

[零值并不是万能]
Go语言零值的设计大大便利了开发者，但是零值并不是万能的，
有些场景下零值是不可以直接使用的：

* 未显示初始化的切片、map，他们可以直接操作，但是不能写入数据，否则会引发程序panic：
    

    var s []string
    s[0] = "append"

    var m map[string]bool
    m["s"] = true
    
    这两种写法都是错误的使用!

* 零值的指针

  零值的指针就是指向nil的指针，无法直接进行运算，
  因为是没有无内容的地址：

    
    var p *uint32
    *p++ // panic: panic: runtime error: invalid memory address or nil pointer dereference
    这样才可以：
    func main() {
        var p *uint64
    
        a := uint64(0)
        p = &a
    
        *p++
        fmt.Println(*p) // 1
    }
* 零值的error类型
  error内置接口类型是表示错误条件的常规接口，nil值表示没有错误，
  所以调用Error方法时类型error不能是零值，否则会引发panic：
    
      func main() {
      rs := res()
        // rs是nil，所以调用Error()时会panic
      fmt.Println(rs.Error())
      }
    
        func res() error {
        return nil
        }
  
        panic: runtime error: invalid memory address or nil pointer dereference
        [signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x10a6f27]
    
* 闭包中的nil函数
  在日常开发中我们会使用到闭包，但是这其中隐藏一个问题，如果我们函数忘记初始化了，那么就会引发panic：


    var f func(a,b,c int)

    func main(){
        f(1,2,3) // panic: runtime error: invalid memory address or nil pointer dereference
    }

* 零值channels
我们都知道channels的默认值是nil，给定一个nil channel c:


    1、 <-c 从c 接收将永远阻塞
    2、 c <- v 发送值到c 会永远阻塞
    3、 close(c) 关闭c 引发panic
【总结】
* Go语言中所有变量或者值都有默认值，对程序的安全性和正确性起到了很重要的作用
* Go语言中的一些标准库利用零值特性来实现，简化操作
* 可以利用"零值可用"的特性可以提升代码的结构化、使代码更简单、更紧凑
* 零值也不是万能的，有一些场景下零值是不可用的，开发时要注意