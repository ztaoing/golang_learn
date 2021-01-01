/**
* @Author:zhoutao
* @Date:2021/1/1 下午1:53
* @Desc: 反射的性能
 */

package reflects

import (
	"reflect"
	"testing"
)

type Config struct {
	Name    string `json:"server-name"`
	IP      string `json:"server-ip"`
	URL     string `json:"server-url"`
	Timeout string `json:"timeout"`
}

func BenchmarkNew(b *testing.B) {
	var config *Config
	for i := 0; i < b.N; i++ {
		config = new(Config)
	}
	_ = config
}

func BenchmarkReflectNew(b *testing.B) {
	var config *Config
	typ := reflect.TypeOf(Config{})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		config, _ = reflect.New(typ).Interface().(*Config)
	}
	_ = config
}

/**
通过反射创建对象耗时约为new的1.5倍
tao@taodeMacBook-Pro reflects % go test reflect_test.go -bench="^Benchmark"
goos: darwin
goarch: amd64
BenchmarkNew-4          	21698186	        49.1 ns/op
BenchmarkReflectNew-4   	16383436	        68.5 ns/op
PASS
ok  	command-line-arguments	2.332s

*/

// 通过反射获取结构体的字段有两种方式：FieldByName 和 Field

func BenchmarkSet(b *testing.B) {
	config := new(Config)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.Name = "name"
		config.IP = "ip"
		config.Timeout = "timeout"
		config.URL = "url"
	}
}

// find by index O(1)
func BenchmarkReflect_FieldSet(b *testing.B) {
	typ := reflect.TypeOf(Config{})

	ins := reflect.New(typ).Elem()
	//fmt.Println("reflect.New: ", reflect.New(typ))
	//fmt.Println("reflect.New(typ).Elem(): ", ins)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ins.Field(0).SetString("name")
		ins.Field(1).SetString("ip")
		ins.Field(2).SetString("url")
		ins.Field(3).SetString("timeout")
	}

}

// find by name O(N)
func BenchmarkReflect_FieldByName(b *testing.B) {
	typ := reflect.TypeOf(Config{})
	ins := reflect.New(typ).Elem()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ins.FieldByName("Name").SetString("name")
		ins.FieldByName("IP").SetString("ip")
		ins.FieldByName("URL").SetString("url")
		ins.FieldByName("Timeout").SetString("timeout")
	}
}

/**
tao@taodeMacBook-Pro reflects % go test reflect_test.go -bench="^BenchmarkReflect_Field"
goos: darwin
goarch: amd64
BenchmarkReflect_FieldSet-4      	32578636	        35.2 ns/op
BenchmarkReflect_FieldByName-4   	 3432375	       345 ns/op
PASS
ok  	command-line-arguments	2.738s

*/

//GO 语言自带的json的marshal和unmarshal方法，在序列化和反序列化的时候使用了反射
