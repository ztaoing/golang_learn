package main

import (
	"io/ioutil"
	"log"
)

/**
ioutil包有一个非常有用的方法WriteFile()可以处理创建或者打开文件、写入字节切片和关闭文件一系列的操作。
如果你需要简洁快速地写字节切片到文件中，你可以使用它。
*/

func main() {
	err := ioutil.WriteFile("test.txt", []byte("Hi\n"), 0666)
	if err != nil {
		log.Fatal(err)
	}
}
