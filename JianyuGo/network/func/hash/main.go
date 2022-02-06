package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	// 得到文件内容
	data, err := ioutil.ReadFile("test.txt")
	if err != nil {
		log.Fatal(err)
	}

	// 计算Hash
	// 复制整个文件内容到内存中，传递给hash函数。
	fmt.Printf("Md5: %x\n\n", md5.Sum(data))
	fmt.Printf("Sha1: %x\n\n", sha1.Sum(data))
	fmt.Printf("Sha256: %x\n\n", sha256.Sum256(data))
	fmt.Printf("Sha512: %x\n\n", sha512.Sum512(data))

	//另一个方式是创建一个hash writer,使用Write、WriteString、Copy将数据传给它。
	//下面的例子使用 md5 hash,但你可以使用其它的Writer。
	file, err := os.Open("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//创建一个新的hasher,满足writer接口
	hasher := md5.New()

	_, err = io.Copy(hasher, file)
	if err != nil {
		log.Fatal(err)
	}

	// 计算hash并打印结果。
	// 传递 nil 作为参数，因为我们不通过参数传递数据，而是通过writer接口。
	sum := hasher.Sum(nil)
	fmt.Printf("Md5 checksum: %x\n", sum)
}
