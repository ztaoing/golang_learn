每个case 上的操作例如方法的结果给channel， 每次循环，所有的方法都会执行，但是只会选择其中一个case，其他的case的操作就被丢失了

这个问题很常见：最多的就是time.After导致内存泄漏问题，网上有很多的文章解释原因，如何避免，其实最根本原因是select这个机制导致的

以下代码会内存泄漏：
func main(){
ch:= make(chan int ,10)
go func(){
var i = 1
for {
i++
ch<-i
}
}()
for {
select{
case x:=<-ch:
println(x)
case <-time.After(30*time.Second):
println(time.now().Unix())
}
}

    }
    为什么会内存泄漏？
    答： 每次循环都会执行time.After(30*time.Second)，导致堆内存不断升高，最后泄漏。哪怕没有选择这个case
        time.After(30*time.Second)也会执行。