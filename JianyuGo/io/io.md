
go存储编程怎么使用O_DIRECT模式？
    操作系统的IO经过文件系统的时候，默认是会使用到page cache，并且采用的是write back方式，
即系统异步刷盘。由于是异步的，如果在数据还未刷盘前掉电，就会导致数据丢失。
    如果要确保数据写到磁盘的话，有两种方式：
    1、每次写完主动调用sync
    2、使用direct io，等待此次io数据都写到磁盘后才返回

那，如何使用O_DIRECT模式呢？是在open文件的时候flag用O_DIRECT吗？有这么简单吗？
    1、O_DIRECT这个定义在go标准库的哪个文件？
    2、direct io需要io大小和偏移扇区对齐，且还要满足内存buffer地址对齐，这个怎么做到？

[O_DIRECT知识点]
    O_DIRECT是在open的时候通过flag来指定O_DIRECT参数，之后的数据的write/read都是绕过page cache的，
直接和磁盘操作，从而避免了掉电数据丢失数据的情况，同时也让应用层可以自己决定内存的使用（避免不必要的cache消耗）

    direct io 一般解决两个问题：
    1、数据的落盘，确保掉电不丢失数据
    2、减少内核page cache的内存使用，业务层自己控制内存更加灵活
    
    direct io模式需要用户保证对齐规则，否则io会报错，有3个需要对齐的规则：
    1、io的大小必须与扇区大小（512字节）对齐
    2、io偏移按照扇区大小对齐
    3、内存buffer的地址也必须是扇区对齐
    
    direct io 模式不对齐会怎样？
    读写会报错，会抛出"无效参数"的错误

    为什么go的O_DIRECT知识点值得提一提？
    1、O_DIRECT平台不兼容
    注：go标准库os中是没有O_DIRECT这个参数的。为什么？
    go的os包中实现的是各个操作系统兼容的实现，direct io这个在不同的操作系统下实现形态不一样。O_DIRECT这个open的flag参数
    值存在于Linux系统中。

    以下才是各个平台兼容的open参数（os/file.go）:其中没有O_DIRECT！
    const (
        // Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
        O_RDONLY int = syscall.O_RDONLY // open the file read-only.
        O_WRONLY int = syscall.O_WRONLY // open the file write-only.
        O_RDWR   int = syscall.O_RDWR   // open the file read-write.
        // The remaining values may be or'ed in to control behavior.
        O_APPEND int = syscall.O_APPEND // append data to the file when writing.
        O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
        O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist.
        O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O.
        O_TRUNC  int = syscall.O_TRUNC  // truncate regular writable file when opened.
    )
    O_DIRECT是一个和系统平台相关的参数。

    那O_DIRECT被定义在哪里了呢？
    根操作系统强相关的自然定义在syscall库中：
    // syscall/zerrors_linux_amd64.go
    const (
        O_DIRECT = 0x4000
    )

    如何打开文件呢？
    // 指明在Linux平台系统下编译:+build linux
        fp:=os.OpenFile(name,syscall.O_DIRECT|flag,perm)
    
    2、go无法精确控制内存分配地址
    标准库或者内置函数没有提供，自行分配对齐内存的函数！
    
    direct io必须要满足3种对齐规则：io偏移扇区对齐，长度扇区偏移对齐、内存buffer地址扇区对齐
    前两个还比好好满足，但是分配的内存地址作为一个程序员无法精确控制。
    
    在c语言中，libc库是调用posix——memaalign直接分配出符合要求的内存块。go是怎么做的呢？
    先问个问题：go里面怎么分配buffer内存的呢？
    io的buffer其实就算是字节数组，最常见自然是make来分配，如下：
    buffer:=make([]byte,4096)
    那这个地址是对齐的吗？
    答：不确认！
    
    怎么才能获取到对齐的地址呢？
    注重点：方法很简单，就是先分配一个比预期要大的内存块，然后再这个内存块里找对齐位置。(这是任何语言皆通用的方法，在go里面也是可用的)
    
    举个例子：我们需要一个4096大小的内存块，要求地址按照512对齐，可以这样做：
    1、先分配4096+512大小的内存块，假设得到的地址是p1
    2、然后再[p1,p1+512]这个地址范围内找，一定能找到512对齐的地址，这个地址是p2
    3、返回p2这个地址给用户使用，用户能正常使用[p2,p2+4096]这个范围内的内存块而不越界。
    以上就是节本原理！
    
    const(
        AlignSize = 512
    )
    // 在block这个字节数组首地址，往后找，找到符合AlignSize对齐的地址，并返回
    // 这里用到位操作，速度很快
    func alignment(block []byte,AlignSize int)int{
        return  int(uintptr(unsafe.Pointer(&block[0])) & uintptr(AlignSize-1))
    }

    // 分配BlockSize大小的内存块
    // 地址按512对齐
    func AlignedBlock(BlockSize int)[]byte{
       // 分配一个字节数组，分配的大小比实际需要的稍大
        block:=make([]byte,BlockSize+AlignSize)
        
        // 计算这个block内存块，向后偏移多少才能对齐512
        a:=alignment(block,AlignSize)
        offset:=0
        if a!=0{
            offset = AlignSize-a
        }
        // 偏移到指定位置，生成一个新的block，这个block将满足地址对齐512
        block = block[offset:offset+BlockSize]
        if BlockSize!=0{
            // 最后做一次校验
            a = alignment(block,AlignSize)
            if a!=0{
                log.Fatal("Failed to align block")
            }
        }
        return block
    }
    通过以上AlignedBlock函数分配的内存一定是512地址对齐的，但是会浪费空间！实际需要4k，但是却分配了4k+512
    
    3、有开源库吗？
    github.com/ncw/directio 内部实现机器简单，就是和上面的一样。
    封装的关键在于O_DIRECT是从syscall库获取的
    如何使用：
    1、 fp,err:=directio.OpenFile(file,os.O_RDONLY,0666)
    2、读取数据
            // 创建地址按照4k对齐的内存
            buffer:=directio,AlignedBlock(directio.BlockSize)
            // 把文件数据读取到内存块中
            _,err:=io.ReadFull(fp,buffer)
    关键在于：buffer必须是定制的[]byte数组，不能使用make([]byte,512)来创建出内存对齐的内存块！
            

    
    
    
    
    