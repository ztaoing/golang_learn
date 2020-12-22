/**
* @Author:zhoutao
* @Date:2020/12/22 下午1:00
* @Desc: -benchmem 可以度量内存分配的次数,内存分配的次数和性能息息相关
 */

package benchmark

import (
	"math/rand"
	"testing"
	"time"
)

func generateWithCap(n int) []int {
	rand.Seed(time.Now().UnixNano())
	// allocate memory all in one
	nums := make([]int, 0, n)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

func generate(n int) []int {
	rand.Seed(time.Now().UnixNano())

	//allocate memory when needed
	nums := make([]int, 0)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

func BenchmarkGenerateWithCap(b *testing.B) {
	for n := 0; n < b.N; n++ {

		generateWithCap(1000000)

	}
}

func BenchmarkGenerate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		generate(1000000)
	}
}

/**
BenchmarkGenerateWithCap-4   	      46	  23157814 ns/op
BenchmarkGenerate-4          	      34	  29628389 ns/op
*/
