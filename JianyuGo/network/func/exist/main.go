package main

import (
	"log"
	"os"
)

var (
	fileInfo *os.FileInfo
	err      error
)

func main() {
	// 文件不存在则返回error
	fileInfo, err := os.Stat("test.txt")
	if err != nil {
		// 如果是文件不存在的错误
		if os.IsNotExist(err) {
			log.Fatal("File does not exist.")
		}
	}
	log.Println("File does exist. File information:")
	log.Println(fileInfo)
}
