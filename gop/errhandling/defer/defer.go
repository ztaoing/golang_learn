package main

import (
	"bufio"
	"fmt"
	"os"
)

func writeFile(filename string) {

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//如果直接写文件会比较慢，可以使用buffer，先写缓存，后写磁盘
	bWrite := bufio.NewWriter(file)
	//将缓存中的内容，存储到磁盘
	defer bWrite.Flush()
	f := fibonacci()
	for i := 0; i < 10; i++ {
		//将f（）的结果写入bWrite
		fmt.Fprintln(bWrite, f())
	}
}

func fibonacci() genInt {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

type genInt func() int

//参数在defer语句时计算
func tryDefer() {
	for i := 0; i < 100; i++ {
		defer fmt.Println(i)
		if i == 30 {
			panic("too many")
		}
	}
}
func main() {
	tryDefer()
	writeFile("fib.text")
}
