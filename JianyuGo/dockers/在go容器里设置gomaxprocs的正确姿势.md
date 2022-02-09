gomaxprocs是go提供的非常重要的一个环境变量。通过设置gomaxprocs，用户可以调整调度器中processor（即P）
的数量，由于每个系统线程必须要绑定P，P才能把G交给M执行。
所以，p的数量会很大程度上影响go runtime的并发表现。
gomaxprocs在go1.5之后默认值是机器的CPU核数（runtime.NumCPU）。通过下面的代码可以获取当前机器的核心数和给
gomaxprocs设置值。

    func getGOMAXPROCS()int{
       _:= runtime.NumCPU() //获取机器的CPU核心数
        return runtime.GOMAXPROCS(0) //参数为0时用于获取给gomaxprocs设置的值
    }
    func main(){
        fmt.Printf("GOMAXPROCS:%d\n",getGOMAXPROCS())
    }
    
    但是，以docker为代表的容器虚拟化技术，会通过cgroup等技术对CPU资源进行隔离。
    以k8s为代表的基于容器虚拟化实现的资源管理系统，也支持这样的特性，比如在podTemplate的容器定义里：
    limits.cpu=1000m就嗲表给这个容器分配1个核心的使用时间。

    这类隔离技术，导致runtime.NumCPU()无法正确的获取到容器被分配CPU资源数。
    runtime.NumCPU()获取的是物理机的实际核心数。

    设置gomaxprocs高于真正可以使用的核心数会导致go调度器不停的进行os线程切换，
    从而给调度器增加很多不必要的工作。
    
    目前go官方没有好的方式来规避容器里获取到实际使用的核心数，而Uber的uber-go/automaxprocs这个包，
    可以在运行时根据cgroup为容器分配的CPU资源限制来修改稿gomaxprocs

    import _ "go.uber.org/automaxprocs"
    func main(){
        //逻辑
    }