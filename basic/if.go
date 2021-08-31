/**
* @Author:zhoutao
* @Date:2021/8/13 10:08 上午
* @Desc:
 */

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	const filename = "abc.text"
	if contents, err := ioutil.ReadFile("golang_learn/gop/" + filename); err != nil {
		fmt.Println("file not exist !")
	} else {
		fmt.Printf("%s\n", contents)
	}

	fmt.Println(grade(59))
	fmt.Println(grade(61))
	//将十进制 转二进制
	fmt.Println(converToBin(13), converToBin(0))

	//逐行读取文件的内容

	printFile(filename)
}

func grade(score int) string {
	q := ""
	switch {
	case score < 0 || score > 100:
		panic("out of range ")
	case score < 60:
		q = "F"
	case score < 80:
		q = "C"
	case score < 90:
		q = "B"
	case score < 100:
		q = "A"
	}
	return q
}

//将数字抓换为二进制
func converToBin(n int) string {
	result := ""
	if n > 0 {
		for ; n > 0; n /= 2 {
			lsb := n % 2
			result = strconv.Itoa(lsb) + result
		}
	} else {
		result = "0000"
	}

	return result
}

//读取文件的内容
func printFile(filename string) {
	file, err := os.Open("golang_learn/gop/" + filename)
	if err != nil {
		panic(err)
	}
	//一行一行的读取文件
	scaner := bufio.NewScanner(file)
	for scaner.Scan() {
		fmt.Println(scaner.Text())
	}
}
