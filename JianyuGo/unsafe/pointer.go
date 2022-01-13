package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type User struct {
	Name   string
	Age    uint32
	Gender bool // 男:true 女：false 就是举个例子别吐槽我这么用。。。。
}

func main() {
	// sizeof 各个类型所占的字节大小
	fmt.Println(unsafe.Sizeof(true))
	fmt.Println(unsafe.Sizeof(int8(0)))
	fmt.Println(unsafe.Sizeof(int16(10)))

	fmt.Println(unsafe.Sizeof(int(10))) // int类型的具体大小是跟机器的CPU位数相关的，如果是CPU是32的，那么int占4个字节；如果CPU是64位的，那么int占8字节

	fmt.Println(unsafe.Sizeof(int32(190)))
	fmt.Println(unsafe.Sizeof("asong"))
	fmt.Println(unsafe.Sizeof([]int{1, 3, 4}))

	// Offsetof
	user := User{Name: "Asong", Age: 23, Gender: true}
	userNamePointer := unsafe.Pointer(&user)

	// 结构体中第一个成员变量的地址是不需要进行偏移量计算的，直接取出后转换为unsafe.pointer类型，在强制转换成字符串类型的指针即可
	nNamePointer := (*string)(unsafe.Pointer(userNamePointer))
	*nNamePointer = "Golang梦工厂"

	// Offsetof返回成员变量在结构体中的偏移量
	// 注意uintptr的使用：不可以用一个临时变量存储uintptr类型，在用于指针运算时，gc不把uintptr当指针，uintptr无法持有对象。
	// uintptr类型的目标会被回收，所以你不知道他什么时候会被gc掉，那样接下来的内存操作会发生什么样的错误，咱也不知道
	/**
	切记不要这样做
		p1:=uintptr(userNamePointer)

	*/
	nAgePointer := (*uint32)(unsafe.Pointer(uintptr(userNamePointer) + unsafe.Offsetof(user.Age)))
	*nAgePointer = 25

	nGender := (*bool)(unsafe.Pointer(uintptr(userNamePointer) + unsafe.Offsetof(user.Gender)))
	*nGender = false

	fmt.Printf("u.Name: %s, u.Age: %d,  u.Gender: %v\n", user.Name, user.Age, user.Gender)
	// Alignof
	var b bool
	var i8 int8
	var i16 int16
	var i64 int64
	var f32 float32
	var s string
	var m map[string]string
	var p *int32

	// Alignof主要是获取变量的对齐值，除了int、uintptr这些依赖CPU位数的类型，基本类型的对齐值都是固定的，
	// 结构体中对齐值取他的成员对齐值的最大值，结构体的对齐，涉及到内存对齐
	fmt.Println(unsafe.Alignof(b))
	fmt.Println(unsafe.Alignof(i8))
	fmt.Println(unsafe.Alignof(i16))
	fmt.Println(unsafe.Alignof(i64))
	fmt.Println(unsafe.Alignof(f32))
	fmt.Println(unsafe.Alignof(s))
	fmt.Println(unsafe.Alignof(m))
	fmt.Println(unsafe.Alignof(p))

	/**
	经典应用：string与[]byte相互转换
	使用这种方式进行转换都会设计底层数值的拷贝，所以想要实现零拷贝，可以使用unsafe.pointer来实现，通过强转换直接完成指针的指向，从而使string和[]byte指向同一个底层数据

	在reflect包中有：string和slice的对应的结构体：
		string运行时的表现形式
		type StringHeader struct{
			Data uintptr
			Len int
		}

		slice运行时的表现形式
		type SliceHeader struct{
			Data uintptr
			Len int
			Cap int
		}
	string和slice运行时只有一个cap字段不同，所以他们的内存布局是对齐的，所以可以通过unsafe.pointer进行转换，
	因为可以写出stringToBytes(s string)[]byte 和bytesToString(b []byte)string方法
	*/

	// string to []byte
	str1 := "golang梦工厂"
	by := []byte(str1)
	fmt.Println(by)

	// []byte to stirng
	str2 := string(by)
	fmt.Println(str2)

}

// 通过重构slice header和string header完成类型转换
func stringToBytes(s string) []byte {
	header := (*reflect.StringHeader)(unsafe.Pointer(&s))
	newHeader := reflect.SliceHeader{
		Data: header.Data,
		Len:  header.Len,
		Cap:  header.Len,
	}
	// 返回的是值，而不是指针
	return *(*[]byte)(unsafe.Pointer(&newHeader))
}
func bytesToString(b []byte) string {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	newHeader := reflect.StringHeader{
		Data: header.Data,
		Len:  header.Len,
	}
	return *(*string)(unsafe.Pointer(&newHeader))
}

// []byte转换成string可以直接强转，因为string的底层也是[]byte，强转会自动构造
// 最燃这种更高效，但是不推荐使用，主要是不安全，使用不当会造成极大的隐患，设置连recover也不能捕获
func bytesToString2(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

/**
内存对齐

对齐的作用和原因：
1、CPU访问内存时，并不是逐个字节访问，认识以字节单位访问。比如32位的CPU，字长为4字节，那么CPU访问内存的单位就是4字节。这样可以减少CPU访问内存你的次数，
增大CPU访问内存的吞吐量。加入我们需要读取8字节的数据，一次读取4个字节那么只需要读取2次就可以了。
2、内存对齐对实现变量的原子性操作也是有好处的，每次内存访问都是原子的，如果变量的大小不超过字长，那么内存对齐后，对该变量的访问就是原子的，而不需要对分开的数据单独操作，
	这个特性在并发场景下非常重要。

	从结果可以看出，字段防止不同的顺序，占用内存也不一样，这就是因为内存对齐，影响了struct的大小，所以有时候合理的字段排序可以减少内存的开销。

	下面我们就一起来分析一下内存对齐：
	首先要明白什么是内存对齐的规则，c语言的对齐规则与go语言一样，所以c语言的对齐规则对go同样使用：
		1、对于结构体的各个成员，第一个成员位于偏移为0的位置，结构体第一个成员的偏移（offset）为0，"以后每个成员对于结构体首地址的offset都是该成员大小与有效对齐值中 较小 哪个的整数倍"
			如有需要，编译器会在成员之间加上填充字节
		2、除了结构成员需要对齐，结构本身也需要对齐，结构的长度必须是编译器默认的对齐长度和成员中 最长 类型中最小的数据代销的倍数对齐。

		好啦，知道规则了，我们现在来分析一下上面的例子，

		// 64位平台，对齐参数是8
		type User1 struct {
		 A int32 // 4
		  B []int32 // 24
		  C string // 16
		  D bool // 1
		}

		根据我的mac使用的64位CPU,对齐参数是8来分析，
		int32、[]int32、string、bool对齐值分别是4、8、8、1，
		占用内存大小分别是4、24、16、1，我们先根据第一条对齐规则分析User1：

		第一个字段类型是int32，对齐值是4，大小为4，所以放在内存布局中的第一位.：A占4个字节，补4个字节，凑成8字节

		第二个字段类型是[]int32，对齐值是8，大小为24，所以他的内存偏移值必须是8的倍数，
		所以在当前user1中，就不能从第4位开始了，必须从第5位开始，也就偏移量为8。
		第4,5,6,7位由编译器进行填充，一般为0值，也称之为空洞。第9位到第32位为第二个字段B.：B占24字节，第9字节-第32字节

		第三个字段类型是string，对齐值是8，大小为16，所以他的内存偏移值必须是8的倍数，
		因为user1前两个字段就已经排到了第32位，所以下一位的偏移量正好是32，正好是字段C的对齐值的倍数，
		不用填充，可以直接排列第三个字段，也就是从第32位到48位第三个字段C. ：C占16字节，第32字节-第48字节

		第四个字段类型是bool，对齐值是1，大小为1，所以他的内存偏移值必须是1的倍数，
		因为user1前三个字段就已经排到了第48位，所以下一位的偏移量正好是48。正好是字段D的对齐值的倍数，
		不用填充，可以直接排列到第四个字段，也就是从48到第49位是第四个字段D. D：占1字节，即第48字节-第49字节

		好了现在第一条内存对齐规则后，内存长度已经为49字节，我们开始使用内存的第2条规则进行对齐。

		根据第二条规则，默认对齐值是8，字段中最大类型程度是24，取最小的那一个，所以求出结构体的对齐值是8，
		我们目前的内存长度是49字节，不是8的倍数，所以需要补齐，所以最终的结果就是56，补了7字节。

*/
