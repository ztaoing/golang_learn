/**
* @Author:zhoutao
* @Date:2022/2/16 10:21
* @Desc:
 */

package main

type MyErr struct {
}

func (MyErr) Error() string {
	return "MyErr"
}

func main() {
	var e error = GetErr()
	println(e == nil)

}

func GetErr() *MyErr {
	return nil
}
