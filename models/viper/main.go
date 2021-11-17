package main

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

/**
viper ：获取yaml文件中的内容，映射为struct
*/

type ServerConfig struct {
	ServerName string `mapstructure:"name"` //将yaml中的name映射到ServerName中
	Port       int    `mapstructure:"port"`
}

/**
config.yaml
name:'user-web'
port:8201
*/
func main() {
	v := viper.New()
	// 设置yaml文件的路径
	v.SetConfigFile("./config.yaml")
	// 读取配置
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	//将yaml中内容映射到ServerConfig中
	serverConfig := ServerConfig{}
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}

	// 获取name
	fmt.Printf("%V", v.Get("name"))

	//二： 通过环境变量获取不同的配置
	envCfg := GetEnvInfo("DEBUG")

	var configFileName string
	configFIleNamePrefix := "config"
	if envCfg == "true" {
		configFileName = fmt.Sprintf("viper_test/%s-debug.yaml", configFIleNamePrefix)
	} else {
		configFileName = fmt.Sprintf("viper_test/%s-pro.yaml", configFIleNamePrefix)
	}

	// 然后读取yaml中的内容映射到struct中
	serverConfig = ServerConfig{}

	v = viper.New()
	//设置配置文件
	v.SetConfigFile(configFileName)
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}

	// 记录到日志中
	zap.S().Infof("配置信息:%v", serverConfig)
	// 监控文件的变更
	go func() {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			// 重新读取文件
			_ = v.ReadInConfig()
			//重新映射 global.serverConfig
			v.Unmarshal(&serverConfig)
		})
	}()

	time.Sleep(time.Second * 3000)
}

/**
如何将线上和线下的配置文件隔离？生产使用生产的配置文件，测试环境使用测试的配置文件
不用修改任何代码而且线上和线下的配置文件能隔离开来。解决：环境变量（）
*/

func GetEnvInfo(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
	// 注意：新设置的环境变量想要生效，必须重启goland编程软件
}
