package main

import (
	"io"
	"log"
	"os"
)

func main() {
	// 打开原始文件
	originalFile, err := os.Open("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer originalFile.Close()

	// 创建新的文件作为目标文件
	newFile, err := os.Create("test_copy.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	// 从源中复制字节到目标文件
	bytesWritten, err := io.Copy(newFile, originalFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Copied %d bytes.", bytesWritten)

	// 将内存中的文件内容--》flush--》到硬盘中s
	err = newFile.Sync()
	if err != nil {
		log.Fatal(err)
	}
}
