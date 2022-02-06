package main

import (
	"log"
	"os"
)

func main() {
	// 这个例子测试写权限，如果没有写权限则返回error。

	file, err := os.OpenFile("test.txt", os.O_WRONLY, 0666)
	if err != nil {
		// 注意文件不存在也会返回error，需要检查error的信息来获取到底是哪个错误导致。
		if os.IsPermission(err) {
			log.Println("Error: Write permission denied.")
		}
	}
	file.Close()

	// 测试读权限
	file, err = os.OpenFile("test.txt", os.O_RDONLY, 0666)
	if err != nil {
		// 注意文件不存在也会返回error，需要检查error的信息来获取到底是哪个错误导致。
		if os.IsPermission(err) {
			log.Println("Error: Read permission denied.")
		}
	}
	file.Close()

	// os.O_RDONLY // 只读
	// os.O_WRONLY // 只写
	// os.O_RDWR // 读写
	// os.O_APPEND // 往文件中添建（Append）
	// os.O_CREATE // 如果文件不存在则先创建
	// os.O_TRUNC // 文件打开时裁剪文件
	// os.O_EXCL // 和O_CREATE一起使用，文件不能存在
	// os.O_SYNC // 以同步I/O的方式打开
}
