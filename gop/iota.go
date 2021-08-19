/**
* @Author:zhoutao
* @Date:2021/8/13 10:01 上午
* @Desc:
 */

package main

import "fmt"

func consts() {
	const (
		b = 1 << (10 * iota)
		kb
		mb
		gb
		tb
		pb
	)
	fmt.Println(b, kb, mb, gb, tb, pb)
}

func main() {
	consts()

}
