/**
* @Author:zhoutao
* @Date:2021/5/14 下午12:54
* @Desc:粘包
 */

/**
为什么会出现粘包
主要原因就是tcp数据传递模式是流模式，在保持长连接的时候可以进行多次的收和发。
“粘包”可发生在发送端也可发生在接收端：
1. 由Nagle算法造成的发送端的粘包：Nagle算法是一种改善网络传输效率的算法。
简单来说就是当我们提交一段数据给TCP发送时，TCP并不立刻发送此段数据，
而是等待一小段时间看看在等待期间是否还有要发送的数据，若有则会一次把这两段数据发送出去。
2. 接收端接收不及时造成的接收端粘包：TCP会把接收到的数据存在自己的缓冲区中，然后通知应用层取数据。
当应用层由于某些原因不能及时的把TCP的数据取出来，就会造成TCP缓冲区中存放了几段数据。

解决办法
出现”粘包”的关键在于接收方不确定将要传输的数据包的大小，因此我们可以对数据包进行封包和拆包的操作。
封包：封包就是给一段数据加上包头，这样一来数据包就分为包头和包体两部分内容了(过滤非法包时封包会加入”包尾”内容)。
包头部分的长度是固定的，并且它存储了包体的长度，根据包头长度固定以及包头中含有包体长度的变量就能正确的拆分出一个完整的数据包。
我们可以自己定义一个协议，比如数据包的前4个字节为包头，里面存储的是发送的数据的长度。
*/
package main

import (
	"fmt"
	"golang_learn/golang_learn/socket"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:3030")
	if err != nil {
		fmt.Println("dial failed,err:", err)
		return
	}
	defer conn.Close()
	for i := 0; i < 20; i++ {
		msg := "hello, hello,how are you?"
		data, err := socket.Encode(msg)
		if err != nil {
			fmt.Println("encode msg failed,err:", err)
			return
		}
		conn.Write(data)
	}
}
