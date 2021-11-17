package main

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	password "github.com/anaskhan96/go-password-encoder"
)

func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}

func main() {
	// 加盐
	// 1、通过生成随机数 和 MD5 生成字符串进行组合
	// 2、数据库同时存储MD5值 和 salt值,以及加密的算法，验证的时候使用salt进行MD5即可
	fmt.Println(genMd5("123456"))

	// 使用sha512.New加密算法是比较安全的
	option := &password.Options{10, 100, 50, sha512.New}
	// 返回生成的salt 和 已经加密的key值
	salt, encodedPwd := password.Encode("a password", option)

	check := password.Verify("a password", salt, encodedPwd, option)
	fmt.Println(check)

	//密码的存储问题：
	// 如果把salt保存到user的struct中，就会有侵入性，一般不这样做！：可以把密码、salt值、使用的加密算法都保存起来！
	// 这是很多大型的框架内置的做法！

	newPasswrod := fmt.Sprintf("$ps-sha512$%s$%s", salt, encodedPwd)
	fmt.Println(newPasswrod)

	// 从newPasswrod中解析出密码
	passwordInfo := strings.Split(newPasswrod, "$")
	fmt.Println(passwordInfo)

	//验证密码
	check2 := password.Verify("a password", passwordInfo[2], passwordInfo[3], option)
	fmt.Println("验证密码", check2)

}
