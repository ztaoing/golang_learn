package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

/**
Scanner是bufio包下的类型,在处理文件中以分隔符分隔的文本时很有用。通常我们使用换行符作为分隔符将文件内容分成多行。
在CSV文件中，逗号一般作为分隔符。os.File文件可以被包装成bufio.Scanner，它就像一个缓存reader。
我们会调用Scan()方法去读取下一个分隔符，使用Text()或者Bytes()获取读取的数据。

分隔符可以不是一个简单的字节或者字符，有一个特殊的方法可以实现分隔符的功能，以及将指针移动多少，
返回什么数据。如果没有定制的SplitFunc提供，缺省的ScanLines会使用newline字符作为分隔符，
其它的分隔函数还包括ScanRunes和ScanWords,皆在bufio包中。


// To define your own split function, match this fingerprint
type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)

// Returning (0, nil, nil) will tell the scanner
// to scan again, but with algorithm bigger buffer because
// it wasn't enough data to reach the delimiter
*/

func main() {
	file, err := os.Open("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)

	// 默认的分隔函数是bufio.ScanLines,我们这里使用ScanWords。
	// 也可以定制一个SplitFunc类型的分隔函数
	scanner.Split(bufio.ScanWords)

	// scan下一个token.
	success := scanner.Scan()
	if success == false {
		// 出现错误或者EOF是返回Error
		err = scanner.Err()
		if err == nil {
			log.Println("Scan completed and reached EOF")
		} else {
			log.Fatal(err)
		}
	}

	// 得到数据，Bytes() 或者 Text()
	fmt.Println("First word found:", scanner.Text())

	// 再次调用scanner.Scan()发现下一个token
}
