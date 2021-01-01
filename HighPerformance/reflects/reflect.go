/**
* @Author:zhoutao
* @Date:2021/1/1 下午1:20
* @Desc:
 */

package reflects

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Config struct {
	Name    string `json:"server-name"`
	IP      string `json:"server-ip"`
	URL     string `json:"server-url"`
	Timeout string `json:"timeout"`
}

func readConfig() *Config {

	config := Config{}
	typ := reflect.TypeOf(config)
	fmt.Println("reflect.TypeOf:", typ)

	value := reflect.Indirect(reflect.ValueOf(&config))
	fmt.Println("reflect.ValueOf", reflect.ValueOf(&config))
	fmt.Println("reflect.Indirect", value)

	fmt.Println("typ.NumField", typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		fmt.Println("typ.Field", f)

		if v, ok := f.Tag.Lookup("json"); ok {
			fmt.Println("f.Tag.Lookup(\"json\")", v, "\n")
			// the  upper key
			key := fmt.Sprintf("CONFIG_%s", strings.ReplaceAll(strings.ToUpper(v), "-", "_"))
			if env, exist := os.LookupEnv(key); exist {
				value.FieldByName(f.Name).Set(reflect.ValueOf(env))
			}
		}
	}
	return &config
}

func mains() {
	os.Setenv("CONFIG_SERVER_NAME", "global_server")
	os.Setenv("CONFIG_SERVER_IP", "10.0.0.1")
	os.Setenv("CONFIG_SERVER_URL", "go1234.cn")

	c := readConfig()
	fmt.Printf("%+v", c)
}

/**
reflect.TypeOf: main.Config
reflect.ValueOf &{   }
reflect.Indirect {   }
typ.NumField 4
typ.Field {Name  string json:"server-name" 0 [0] false}
f.Tag.Lookup("json") server-name

typ.Field {IP  string json:"server-ip" 16 [1] false}
f.Tag.Lookup("json") server-ip

typ.Field {URL  string json:"server-url" 32 [2] false}
f.Tag.Lookup("json") server-url

typ.Field {Timeout  string json:"timeout" 48 [3] false}
f.Tag.Lookup("json") timeout

&{Name:global_server IP:10.0.0.1 URL:go1234.cn Timeout:}

*/
