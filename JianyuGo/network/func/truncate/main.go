package main

import (
	"log"
	"os"
)

func main() {
	/**
	裁剪一个文件到100字节：

	如果文件本来就少于100字节，则文件中原始内容得以保留，剩余的字节以null字节填充。
	如果文件超过100字节，则超过的字节会被丢弃。

	这样我们总是得到精确的100字节的文件。

	传入0，则会清空文件。
	*/
	err := os.Truncate("test.txt", 100)
	if err != nil {
		log.Fatal(err)
	}
}
