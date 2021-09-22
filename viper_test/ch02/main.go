package main

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type MysqlConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type ServerConfig struct {
	Name      string      `mapstructure:"name"`
	MysqlInfo MysqlConfig `mapstructure:"mysql"`
}

// 获取环境变量
func GetEnvInfo(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
	// 新设置的环境变量想要生效，必须重启goland编程软件。
}

func main() {
	// 开发环境和测试环境的配置文件的隔离
	data := GetEnvInfo("Debug")
	var configFileName string
	configFileNamePrefix := "config"
	if data == "true" {
		configFileName = fmt.Sprintf("viper_test/%s-debug.yaml", configFileNamePrefix)
	} else {
		configFileName = fmt.Sprintf("viper_test/%s-pro.yaml", configFileNamePrefix)
	}

	serverConfig := ServerConfig{}

	fmt.Println(data)

	v := viper.New()
	v.SetConfigFile(configFileName)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	fmt.Println(serverConfig)

	go func() {
		// 监控配置的变化
		// watcher.Add(configDir)
		v.WatchConfig()

		// 定义当发生变化，需要执行的操作
		v.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name)
			_ = v.ReadInConfig() // 重新读取配置数据
			_ = v.Unmarshal(&serverConfig)
			fmt.Println(serverConfig)
		})
	}()

	time.Sleep(time.Second * 3000)

}
