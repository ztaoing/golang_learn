package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "woai慕课网"
	//X 可以打出具体的字节
	fmt.Printf("%X\n", []byte(s))
	//每个字节的存放
	for _, v := range []byte(s) {
		fmt.Printf("%X ", v)
	}
	//utf-8的编码
	//77 6F 61 69 E6 85 95 E8 AF BE E7 BD 91
	fmt.Println()

	for k, v := range s { //v 是一个rune
		fmt.Printf("(%d %X)", k, v)
	}
	fmt.Println()
	//unicode编码
	//(0 77)(1 6F)(2 61)(3 69)(4 6155)(7 8BFE)(10 7F51)

	fmt.Println("rune count:", utf8.RuneCountInString(s))

	//获取s的字节
	bytes := []byte(s)
	//返回字符和大小
	for len(bytes) > 0 {
		cha, size := utf8.DecodeRune(bytes)
		//获取到每一个字符
		bytes = bytes[size:]
		//%c 打印字符
		fmt.Printf("%c ", cha)
	}

	fmt.Println()

	for k, v := range []rune(s) { //v 是一个rune
		fmt.Printf("(%d %c)", k, v)
	}
	/*
		strings.Fields()
		strings.Split()
		strings.Join()

		strings.Contains()
		strings.Index()
		strings.ToLower()
		strings.ToUpper()

		strings.Trim()
		strings.TrimRight()
		strings.TrimLeft()
	*/
}
