/**
* @Author:zhoutao
* @Date:2022/3/11 08:19
* @Desc:
 */

package main

import (
	"fmt"
	"net"
)

func main() {
	// step 1 连接服务器
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 9090,
	})
	if err != nil {
		fmt.Printf("connet failed,err %v\n", err)
		return
	}
	// step 2 发送数据
	for i := 0; i < 100; i++ {
		_, err := conn.Write([]byte("hello server!"))
		if err != nil {
			fmt.Printf("send data failed,err:%v\n", err)
			return
		}
		// step 3 结束数据
		result := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFromUDP(result)
		if err != nil {
			fmt.Printf("receive data failed,err % v\n", err)
			return
		}
		fmt.Printf("receive from addr:%v data:%v\n", remoteAddr, string(result[:n]))
	}
}
