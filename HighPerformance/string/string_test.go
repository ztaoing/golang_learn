/**
* @Author:zhoutao
* @Date:2021/1/1 上午10:59
* @Desc: 字符串的拼接
 */

package string

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

//在go语言中，string是不可以变化的，拼接字符串实际上是创建新的字符串对象

const letterByte = "qwertyuiopasdfghjklzxcvbnm"

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterByte[rand.Intn(len(letterByte))]
	}
	return string(b)
}

// 使用 +
func plusConcat(n int, str string) string {
	s := ""
	for i := 0; i < n; i++ {
		s += str
	}
	return s
}

// 使用fmt.printf
func sprintfConcat(n int, str string) string {
	s := ""
	for i := 0; i < n; i++ {
		s = fmt.Sprintf("%s%s", s, str)
	}
	return s
}

// 使用strings.Builder
func builderConcat(n int, str string) string {
	var s strings.Builder
	for i := 0; i < n; i++ {
		s.WriteString(str)
	}
	return s.String()
}

// 使用strings.Builder
func preBuilderConcat(n int, str string) string {
	var s strings.Builder
	s.Grow(n * len(str))
	for i := 0; i < n; i++ {
		s.WriteString(str)
	}
	//strings.Builder 直接将底层的[]byte转换成了字符串类型返回
	return s.String()
}

// 使用bytes.Buffer
func bufferConcat(n int, str string) string {
	buf := new(bytes.Buffer)
	for i := 0; i < n; i++ {
		buf.WriteString(str)
	}
	//bytes.Buffer 转换为字符串时，重新申请了一块内存空间，用来存放新生成的字符串变量
	return buf.String()
}

//使用[]byte
func byteConcat(n int, str string) string {
	buf := make([]byte, 0)
	for i := 0; i < n; i++ {
		buf = append(buf, str...)
	}
	return string(buf)
}

//如果长度是可以预知的，那么创建[]byte时，我们还可以预分配切片的容量
func preByteConcat(n int, str string) string {
	//预分配cap大小的内存
	buf := make([]byte, 0, n*len(str))
	for i := 0; i < n; i++ {
		buf = append(buf, str...)
	}
	return string(buf)
}

func benchmark(b *testing.B, f func(int, string) string) {
	var str = randomString(10)
	for i := 0; i < b.N; i++ {
		f(10000, str)
	}
}

func BenchmarkPlusConcat(b *testing.B) {
	benchmark(b, plusConcat)
}

func BenchmarkSprintfConcat(b *testing.B) {
	benchmark(b, sprintfConcat)
}

func BenchmarkPreBuilderConcat(b *testing.B) {
	benchmark(b, preBuilderConcat)
}

func BenchmarkBuilderConcat(b *testing.B) {
	benchmark(b, builderConcat)
}

func BenchmarkBufferConcat(b *testing.B) {
	benchmark(b, bufferConcat)
}

func BenchmarkByteConcat(b *testing.B) {
	benchmark(b, byteConcat)
}

func BenchmarkPreByteConcat(b *testing.B) {
	benchmark(b, preByteConcat)
}

/**
goos: darwin
goarch: amd64
BenchmarkPlusConcat-4         	       9	 114650379 ns/op	530996580 B/op	   10046 allocs/op
BenchmarkSprintfConcat-4      	       6	 296908030 ns/op	833263162 B/op	   37418 allocs/op

BenchmarkPreBuilderConcat-4   	    9890	    120975 ns/op	  106496 B/op 内存占用	       1 allocs/op 内存分配1次
BenchmarkBuilderConcat-4      	    6146	    204457 ns/op	  522224 B/op	      23 allocs/op

BenchmarkBufferConcat-4       	    4376	    272747 ns/op	  423536 B/op	      13 allocs/op
BenchmarkByteConcat-4         	    3108	    343294 ns/op	  628721 B/op	      24 allocs/op
BenchmarkPreByteConcat-4      	   13729	     85385 ns/op	  212992 B/op	       2 allocs/op
PASS
ok  	command-20line-arguments	10.749s


*/
