/**
* @Author:zhoutao
* @Date:2021/5/14 下午1:16
* @Desc:
 */

package main

import (
	"bufio"
	"fmt"
	"golang_learn/golang_learn/socket"
	"io"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:3030")
	if err != nil {
		fmt.Println("listen failed,err:", err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed,err:", err)
			continue
		}
		go process(conn)
	}
}

func process(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		msg, err := socket.Decode(reader)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("decode msg failed,err:", err)
			return
		}
		fmt.Println("收到client发来的数据:", msg)
	}
}
