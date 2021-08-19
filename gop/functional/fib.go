package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func fibonacci() genInt {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

type genInt func() int

// 函数实现reader接口
// 如果p 即[]byte 给的比较小怎么办？缓冲区太小。解决：做一个struct把函数和reader缓存起来
func (g genInt) Read(p []byte) (n int, err error) {
	// 取得下一个元素
	next := g()
	if next > 10000 {
		return 0, io.EOF
	}
	// 把数字转换为字符串
	s := fmt.Sprintf("%d\n", next)
	// 然后把s写进p中
	return strings.NewReader(s).Read(p)
}
func printFileContent(reader io.Reader) {
	// 从reader返回一个scanner
	scanner := bufio.NewScanner(reader)
	// &Scanner{
	//		r:            r,
	//		split:        ScanLines,
	//		maxTokenSize: MaxScanTokenSize,
	//	}
	// 无限的循环scan，即无限的调用scanner的read方法，即genInt.Read()
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
func main() {
	// 返回的是一个函数接口，并且实现了reader接口
	f := fibonacci()
	printFileContent(f)
}
