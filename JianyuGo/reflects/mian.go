/**
* @Author:zhoutao
* @Date:2022/2/12 09:28
* @Desc:
 */

package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Name   string
	Age    uint64
	Gender bool // true：男 false: 女
}

func main() {
	u := User{
		Name:   "asong",
		Age:    18,
		Gender: false,
	}
	// 返回的type是一个接口
	getType := reflect.TypeOf(u)
	for i := 0; i < getType.NumField(); i++ {
		fieldType := getType.Field(i)
		// 输出成员名
		fmt.Printf("name: %v \n", fieldType.Name)
	}
}

/*
运行结果
name: Name
name: Age
name: Gender */
