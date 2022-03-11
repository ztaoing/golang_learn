/**
* @Author:zhoutao
* @Date:2022/3/11 08:32
* @Desc:
 */

package main

import (
	"fmt"
	"net"
)

func main() {
	// step 1 监听端口
	listen, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		fmt.Printf("listen fail,err:%v\n", err)
		return
	}

	// step2 建立套接字连接
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept fail,err:%v\n", err)
			continue
		}
		// step 3 创建处理goroutine
		go process(conn)
	}
}

func process(conn net.Conn) {
	// 如果不关闭会怎样
	defer conn.Close()

	for {
		var buf [128]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Printf("read from connect failed,err:%v\n", err)
			return
		}
		str := string(buf[:n])
		fmt.Printf("recevie from client,data:%v\n", str)

	}
}
