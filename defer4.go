/**
* @Author:zhoutao
* @Date:2021/3/31 下午12:50
* @Desc:
 */

package main

import "fmt"

func main() {

	deferCall()
}
func deferCall() {
	defer func() { fmt.Println("಑ڹܦ") }()
	defer func() { fmt.Println("಑ܦӾ") }()
	defer func() { fmt.Println("಑ݸܦ") }()
	panic("恐慌")
}
