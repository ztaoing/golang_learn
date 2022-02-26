package main

import (
	"fmt"
	"strconv"
)

func main() {
	// int 转 字符串
	/*	var logic int64 = 56
		fmt.Printf("%T\n", strconv.Itoa(int(logic)))*/

	// 字符串 转 int
	// 他的返回值有int, error
	/*_, err := strconv.Atoi("12")
	if err != nil {
		// 转换出错
		return
	}*/

	//1、 parse 用于将 字符串 转换为 给定的类型
	b, _ := strconv.ParseBool("flase")
	fmt.Println(b)
	// 是float64 还是float32
	f, _ := strconv.ParseFloat("3.1415", 64)
	fmt.Println(f)
	fmt.Printf("%T\n", f)
	//base参数表示以什么进制的方式去解析给定的字符串
	i, _ := strconv.ParseInt("012", 0, 64)
	fmt.Println(i)
	fmt.Printf("%T\n", i)

	u, _ := strconv.ParseUint("42", 10, 64)
	fmt.Println(u)

	//2、 将 给定的类型 格式化 为string类型
	fb := strconv.FormatBool(true)
	fmt.Println(fb)

}
