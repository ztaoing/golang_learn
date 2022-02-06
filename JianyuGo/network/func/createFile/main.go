package main

import (
	"log"
	"os"
)

// 创建文件

var (
	newFile *os.File
	err     error
)

func main() {
	newFile, err = os.Create("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(newFile)
	newFile.Close()
}
