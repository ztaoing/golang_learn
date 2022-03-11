/**
* @Author:zhoutao
* @Date:2022/3/11 08:32
* @Desc:
 */

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	// step 1 连接服务器
	conn, err := net.Dial("tcp", "localhost:9090")
	// 如果不关闭会怎样?
	defer conn.Close()
	if err != nil {
		fmt.Printf("connect failed,err:%v\n", err)
		return
	}

	// step 2 读取命令行输入
	inputReader := bufio.NewReader(os.Stdin)
	for {
		// step 3 一直读取到\n
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Printf("read from console failed, err:%v\n", err)
			break
		}
		// step 4 读取Q时停止
		trimmedInput := strings.TrimSpace(input)
		if trimmedInput == "Q" {
			break
		}
		// step 5 回复服务器信息
		_, err = conn.Write([]byte(trimmedInput))
		if err != nil {
			fmt.Printf("write faield,err:%v\n", err)
			break
		}

	}

	fmt.Println("client close")
	time.Sleep(time.Second * 100)
}

// netstat -Aaln | grep 9090
