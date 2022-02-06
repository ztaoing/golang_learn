package main

import (
	"log"
	"os"
)

/**
可以使用os包写入一个打开的文件。
因为Go可执行包是静态链接的可执行文件，
你import的每一个包都会增加你的可执行文件的大小。
其它的包如io、｀ioutil｀、｀bufio｀提供了一些方法，但是它们不是必须的。
*/

func main() {
	// 可写方式打开文件
	file, err := os.OpenFile(
		"test.txt",
		os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 写字节到文件中
	byteSlice := []byte("write Bytes!\n")
	// 以覆盖的方式写入
	bytesWritten, err := file.Write(byteSlice)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes.\n", bytesWritten)

	// os.O_RDONLY // 只读
	// os.O_WRONLY // 只写
	// os.O_RDWR // 读写
	// os.O_APPEND // 往文件中添建（Append）
	// os.O_CREATE // 如果文件不存在则先创建
	// os.O_TRUNC // 文件打开时裁剪文件
	// os.O_EXCL // 和O_CREATE一起使用，文件不能存在
	// os.O_SYNC // 以同步I/O的方式打开
}
