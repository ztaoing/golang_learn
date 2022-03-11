/**
* @Author:zhoutao
* @Date:2022/3/11 08:09
* @Desc:
 */

package main

import (
	"fmt"
	"net"
)

func main() {
	// step 1 监听服务器
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 9090,
	})
	if err != nil {
		fmt.Printf("listen failed,err %v\n", err)
		return
	}

	// step 2 循环读取消息
	for {
		var data [1024]byte
		// 不是使用套接字读取，而是直接使用listen这个句柄去读取的
		_, addr, err := listen.ReadFromUDP(data[:])
		if err != nil {
			fmt.Printf("read failed from addr:%v,err:%v\n", addr, err)
			break
		}
		go func() {
			// step 3 回复数据

			_, err := listen.WriteToUDP([]byte("received success!"), addr)
			if err != nil {
				fmt.Printf("write faield,err%v\n", err)
			}

		}()
	}
}
