package main

import "fmt"

//适配器模式
//应用最多的是在接口升级，而又需要保证老接口的兼容性，为了让老接口能够继续工作，所以提供了一个中间层，让老接口对外的接口不变
//但是实际上调用了新的代码
func main() {
	printName("", "")
}

func printName(familyName, name string) {
	if isChinese(familyName) {
		printChinesehName(name, familyName)
	} else if isEnglish(familyName) {
		printEnglishName(familyName, name)
	} else {
		fmt.Println("暂不支持此类型 Not supported yet")
	}
}

func isChinese(string) bool {
	return true
}

func isEnglish(string) bool {
	return true
}

func printEnglishName(firstName, secondName string) {
	fmt.Printf("firstName is %s,secondName is %s\n", firstName, secondName)
}

func printChinesehName(familyName, name string) {
	fmt.Printf("姓：%s,名：%s", familyName, name)
}
