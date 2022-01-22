/**
* @Author:zhoutao
* @Date:2021/1/1 上午9:19
* @Desc: slice.md and array
* @Desc: https://github.com/golang/go/wiki/SliceTricks  slice的操作
 */

package main

import (
	"fmt"
	"math/rand"
	"runtime"
	testing2 "testing"
	"time"
)

func square(arr *[3]int) {
	for i, num := range *arr {
		(*arr)[i] = num * num
	}
}

func TestArrayPointer(t *testing2.T) {
	// ... 会自动计算数组长度
	a := [...]int{1, 2, 3}
	square(&a)
	fmt.Println(a)

	if a[1] != 2 && a[2] != 9 {
		t.Fatal("failed")
	}
}

func printfLenCap(nums []int) {
	fmt.Printf("len: %d , cap: %d %v\n", len(nums), cap(nums), nums)
}

func TestSliceLenCap(t *testing2.T) {
	//切片的本质是一个数组片段的描述,因此切片操作并不复制切片指向的元素，
	//创建一个新的切片会复用原来切片数组的底层数组，因此切片操作是非常高效的
	nums := []int{1}
	printfLenCap(nums)

	nums = append(nums, 2)
	printfLenCap(nums)

	nums = append(nums, 3)
	printfLenCap(nums)

	nums = append(nums, 4)
	printfLenCap(nums)

	nums = append(nums, 5)
	printfLenCap(nums)

}

/*切片是对底层数组的引用，所以底层数组在被引用的情况下是不会释放的*/

//直接在原切片的基础上进行切片
func lastNumsBySlice(origin []int) []int {
	//引用了底层的数组，底层数组得不到释放
	return origin[len(origin)-2:]
}

//创建了一个新切片，将origin的最后两个元素拷贝到新切片上，然后返回新切片
func lastNumsByCopy(origin []int) []int {
	result := make([]int, 2)
	//通过copy指向了一个新的底层数组，origin不再被引用后，内存会被垃圾回收
	copy(result, origin[len(origin)-2:])
	return result
}

// 生成n个随机int值 64位机器上，一个int占用8个字节，128*1024个整数切好占用1Mb的空间
func generateWithCap(n int) []int {
	rand.Seed(time.Now().UnixNano())

	nums := make([]int, 0, n)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

//打印程序运行时占用内存的大小
func printMem(t *testing2.T) {
	t.Helper()
	var rumtimeMemory runtime.MemStats

	runtime.ReadMemStats(&rumtimeMemory)

	t.Logf("%.2f MB", float64(rumtimeMemory.Alloc)/1024./1024.)
}

func testLastChars(t *testing2.T, f func([]int) []int) {
	t.Helper()

	ans := make([][]int, 0)
	//执行100次， 也就是将100M的数据放入ans中
	for k := 0; k < 100; k++ {
		origin := generateWithCap(128 * 1024) // 1M
		ans = append(ans, f(origin))
	}
	printMem(t)

	// ?
	_ = ans
}

func TestLastCharBySlice(t *testing2.T) {
	//执行100次，使用了100M的内存，由于引用了底层数组，所以底层数组占用的100M内存得不到释放
	testLastChars(t, lastNumsBySlice)
}

func TestLastCharByCopy(t *testing2.T) {
	testLastChars(t, lastNumsByCopy)
}

// -bench regexp 执行相应的 benchmarks，例如 -bench=.；
// -cover 开启测试覆盖率；
// -run regexp 只运行 regexp 匹配的函数，例如 -run=Array 那么就执行包含有 Array 开头的函数；
// -v 显示测试的详细命令
// go test slice_test.go -run=^TestLastChar -v

/*
=== RUN   TestLastCharBySlice
slice_test.go:108: 100.10 MB
--- PASS: TestLastCharBySlice (0.37s)
=== RUN   TestLastCharByCopy
slice_test.go:112: 3.10 MB
--- PASS: TestLastCharByCopy (0.32s)
PASS

lastNumsBySlice 耗费了100.10 MB内存
lastNumsByCopy  耗费了3.10 MB 内存
*/
