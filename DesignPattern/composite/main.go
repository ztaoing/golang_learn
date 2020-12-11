package main

import "fmt"

//composite模式
// 对于多个对象，由于我们只需要其中一部分共同的操作，因此我们可以通过定义一个父类，来规定我们所需要的操作，却并不管具体每个子类
// 到底是什么样的。
// 这里说明几个问题：
// 1、我们会把一组独享当做同样地类型，也就是说，我们并不在乎他是什么类的实例，我们只在乎有什么操作
// 2、通常会用树状结构来表示，对应到编程语言，其实就是使用集成的方式
func main() {

}

//composite 的核心并不是一定要用树状模式（也就是对应编程语言的继承）来表示，而是我们只关心是否实现了接口，并不关系它具体是什么，这不就是go里面接口的用法吗！
//这就是composite模式在go语言里的应用。
//composite模式可以用于递归的表示某些东西，比如：文件系统、窗口系统等大量共同属性、操作的情况
type Sender interface {
	Send(user, message string)
}

type AliyunSender struct {
	Sender
}

func (a *AliyunSender) Send(user, message string) {
	fmt.Printf("使用aliyun向%s发送信息%s\n", user, message)
}

type TencentSender struct {
	Sender
}

func (a *TencentSender) Send(user, message string) {
	fmt.Printf("使用tencent向%s发送信息%s", user, message)
}

var (
	_ Sender = &AliyunSender{}
	_ Sender = &TencentSender{}
)
