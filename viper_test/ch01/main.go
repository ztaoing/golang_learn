package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// ServerConfig tag将viper获取的yaml内容，映射为struct
// 底层使用了mapstructure库
type ServerConfig struct {
	ServiceName string `mapstructure:"name"` // 将yaml中的name映射到ServiceName
	Port        int    `mapstructure:"port"` // 将yaml中的name映射到Port
}

func main() {
	// *Viper
	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile("viper_test/ch01/config-debug.yaml")
	// 读取配置
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	serverConfig := ServerConfig{}
	// 将yaml文件中的内容映射成struct
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	fmt.Println(serverConfig)
	// v.Get() 获取
	fmt.Printf("%V", v.Get("name"))
}
