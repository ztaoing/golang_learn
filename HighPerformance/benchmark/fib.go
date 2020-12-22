/**
* @Author:zhoutao
* @Date:2020/12/22 下午12:28
* @Desc:
 */

package benchmark

import "time"

func fib(n int) int {
	time.Sleep(time.Second * 3) //模拟耗时的准备工作
	if n == 0 || n == 1 {
		return n
	}
	return fib(n-2) + fib(n-1)
}
