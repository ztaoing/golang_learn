
string类型本质也是一个结构体

    type stringStruct struct {
        str unsafe.Pointer
        len int
    }

stringStruct和slice还是很相似的，str指针指向的是某个数组的首地址，len代表的就是数组长度。
怎么和slice这么相似，底层指向的也是数组，是什么数组呢？我们看看他在实例化时调用的方法：

    //go:nosplit
    func gostringnocopy(str *byte) string {
        ss := stringStruct{str: unsafe.Pointer(str), len: findnull(str)}
        s := *(*string)(unsafe.Pointer(&ss))
        return s
    }

入参是一个byte类型的指针，从这我们可以看出string类型底层是一个byte类型的数组.

string类型本质上就是一个byte类型的数组，在Go语言中string类型被设计为不可变的，

---

