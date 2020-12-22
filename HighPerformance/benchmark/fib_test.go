/**
* @Author:zhoutao
* @Date:2020/12/22 下午12:54
* @Desc:
 */

package benchmark

import "testing"

//基准测试
func BenchmarkFib(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(30) //run fib(30) b.N times
	}
}
