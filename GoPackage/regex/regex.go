/**
* @Author:zhoutao
* @Date:2021/1/2 下午4:00
* @Desc:
 */

package main

import (
	"errors"
	"fmt"
	"regexp"
)

func matchString(str string) (match bool, err error) {
	if str == "" {
		return false, errors.New("str can not be empty")
	} else {
		match, err = regexp.MatchString(`^Golang`, str)
	}
	return match, err
}

func main() {
	var str string
	str = "Golang regular expressions example"

	match, err := matchString(str)
	fmt.Println("match: ", match, " error: ", err)

	//创建一个编译好的正则表达式对象
	//Compile()方法返回error，而MustCompile()编译非法正则表达式时不会返回error，而是返回panic
	// 如果需要好的性能，不要在使用的时候才调用Compile临时进行编译，而是预先调用Compile编译好正则表达式对象
	regex1, err := regexp.Compile("<code>regexp</code>")
	regexp2 := regexp.MustCompile("<code>regexp</code>")
	fmt.Println(regex1, regexp2)

	//返回第一个匹配的结果，如果没有匹配字符串，那么返回一个空的字符串，当然如果你的正则表达式就是要匹配空字符串的话，它也会返回空字符串,使用FindStringIndex,FindStringSubmatch可以区分这两种情况
	regexp, _ := regexp.Compile("Gola([algorithm-z]+)g")
	fmt.Println(regexp.FindString(str))

	//得到匹配字符串的索引
	//regex2, err := regexp.Compile("<code>regexp</code>")

	//返回匹配到的字符串以外的部分
	//FindStringSubmatch()
}
