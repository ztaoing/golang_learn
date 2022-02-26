package main

import (
	"fmt"
	"runtime"
	"unsafe"
)

/**
从CPU角度理解go中的结构体内存对齐

*/

type t2 struct {
	f1 int8
	f2 int32
	f3 int64
}
type t1 struct {
	f1 int8
	f3 int64
	f2 int32
}

func main() {
	fmt.Println(runtime.GOARCH)

	t1 := t1{}
	fmt.Println(unsafe.Sizeof(t1))

	t2 := t2{}
	fmt.Println(unsafe.Sizeof(t2))
}

/**
64位的电脑，指的是CPU一次可以从内存中读取64 位的数据，即8字节。这个长度也成为CPU的字长。
虽然CPU可以一次拿取8字节，但是也是想从哪里拿就可以从哪里拿。因为内存也会以8字节（1个字）为单位分割，也就是每次最多可以拿一个字。
如果所需读取的数据正好跨越两个字的内存，就得需要两个CPU周期的时间来读取！
*/

/**
3、struct内存对齐
T1结构体在内存中如果紧密排列的话是什么样的。在T1结构体中各字段的顺序是int8、int64、int32，
在第一个字中的顺序是：int8+int64（前部分）
在第二个字中的顺序是：int64（后部分）+int32

这样排列会有什么问题呢？如果程序想要读取t1.f2字段的数据，那CPU就得需要两个时钟周期把f2字段完整的从内存中读取出来。
所以，为了能让CPU可以更快的存取各个字段，go编译器会把struct结构体的数据做对齐。

数据对齐：是指内存地址是存储数据大小（按字节单位）的整数倍，以便CPU可以一次将该数据从内存中读取出来。
编译器通过在T1结构体的各个字段之间填充一些空白来达到对齐的目的！

重新排列后，内存的布局是这样的：
f1：int8+填充7个字节
f2：int64（正好8个字节）
f3：int32+填充4字节

*/

/**
如何减少struct的填充：
虽然通过填充可以提高CPU的读写效率，但是这些填充的内存中是没有任何数据的。以T1为例，实际数据只有13字节，但实际却申请了24字节的内存
那有没有好的办法可以做到内存对齐之后既可以提高CPU的读取效率，又能较少内存浪费呢？
答案：可以调整struct字段的顺序：
f1:int8+f3:int32(正好一个字)
f2:int64(正好一个字)
把f1和f3保存在一起，就可以节省8字节的空间
*/

/**
在go中，go会按照结构体中字段的顺序在内存中布局，所以需要程序员将字段f2和f3的位置交换，定义的顺序变为int8、int32、int64，
这样go的编译器才会按照最节省内存的方式来排列
*/

/** 5、在同一个字中的内存布局
上面都是跨字长的存储实例，如果是有n个小于一个字长的类型在同一个字长中是否可以连续分配呢？

var x strcut{
	logic bool
	b int16
}
在内存中是如下排列的：
logic+一个空白+b
为什么会这样呢？
答案是从内存对齐的定义中推到出来的。

注：内存对齐是指-》数据存放的地址是数据大小的整数倍。也就是说会有（数据存放的其实地址）%（数据的大小）=0
验证：
假设结构体的其实地址为0，那么a从0爱是占用1个字节。b字段如果放在地址1处，套用上边的公式1%2=1，就不满足对齐的要求。
所以在地址为2处存放b字段。

*/

/**6、什么时候该关注结构体字段顺序
由此可知，对结构体字段的合理排序可以节省内存，但是我们需要这么做吗？
type Student struct {
    id int8 //学号
    name string //姓名
    classID int8 //班级
    phone [10]byte //联系电话
    address string // 地址
    grade int32 //成绩
}

会有很多的填充，共浪费了16字节

type Student struct {
    name string //姓名
    address string // 地址
    grade int32 //成绩
    phone [10]byte //联系电话
    id int8 //学号
    classID int8 //班级
}

我们通过调整Student结构体的字段顺序来进行下优化，可以看到从开始的64字节，可以优化到48字节，共省下了25%的空间。

我们看到，通过调整结构体中的字段顺序确实节省了内存空间，那我们真的有必要这样节省空间吗？

以student结构体为例，经过重新排列后，节省了16字节的空间，假设我们在程序中需要排列全校同学的成绩，
需要定义一个长度为10万的Student类型的数组，那剩下的内存也不过16MB的空间，跟现在个人电脑的8G到16G的内存比起来微不足道。
而且在字段重新排列后，可读性也变得很差了。像Student原本是以学号，姓名，班级...这样依次排列的，而重新调整后变成了姓名，地址，成绩...，一直到最后才是学号跟班级，
不符合人们的思维习惯。

所以，我的建议是对于结构体的字段排列不需要过早的进行优化，除非一开始就知道你的程序瓶颈就卡在这里。否则，就按照正常的习惯编写Go程序即可。

*/
