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
  * nil的切片和空切片指向的地址是不一样的