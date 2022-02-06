package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

/**
ioutil提供了两个函数: TempDir() 和 TempFile()。
使用完毕后，调用者负责删除这些临时文件和文件夹。

有一点好处就是当你传递一个空字符串作为文件夹名的时候，它会在操作系统的临时文件夹中创建这些项目（/tmp on Linux）。
os.TempDir()返回当前操作系统的临时文件夹。
*/
func main() {
	// 在系统临时文件夹中创建一个临时文件夹
	tempDirPath, err := ioutil.TempDir("", "myTempDir")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Temp dir created:", tempDirPath)

	// 在临时文件夹中创建临时文件
	tempFile, err := ioutil.TempFile(tempDirPath, "myTempFile.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Temp file created:", tempFile.Name())

	// ... 做一些操作 ...

	// 关闭文件
	err = tempFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	// 删除我们创建的资源
	// 先删除临时文件
	err = os.Remove(tempFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	// 再删除临时目录
	err = os.Remove(tempDirPath)
	if err != nil {
		log.Fatal(err)
	}
}
