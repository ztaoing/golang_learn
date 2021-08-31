package main

import (
	"encoding/json"
	"fmt"
)

type Peopless struct {
	name string `json:"name"`
}

func main() {
	js := `{
        "name":"11"
    }`
	var p Peopless
	err := json.Unmarshal([]byte(js), &p)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Println("peopless: ", p)
}

/**
考点:结构体访问控制
这道题坑很大，很多同学一看就以为是p的初始化问题，实际上是因为name首字母是小写，导致其他包不能访问，所以输出为空结构体。
输出结果：
peopless:  {}

*/
