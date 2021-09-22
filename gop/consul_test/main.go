package main

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

func main() {
	_ = Register("192.168.1.102", 8021, "user-web", []string{"shop", "bod"}, "user-web")

}

// 服务的注册
func Register(address string, port int, name string, tags []string, id string) error {
	// 定义配置
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.1.103:8500"
	// 根据配置生成client
	client, err := api.NewClient(cfg)
	// 定义对应的检查对象
	check := &api.AgentServiceCheck{
		HTTP:                           "http://192.168.1.102:8021/health",
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}
	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	/*
	   registration := api.AgentServiceRegistration{}
	*/
	// 定义注册对象
	registration.Address = address
	registration.Name = name
	registration.Port = port
	registration.Tags = tags
	registration.ID = id

	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	return nil
}

// 服务的发现
func AllService() {
	// 定义配置
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.1.103:8500"
	// 根据配置生成client
	client, err := api.NewClient(cfg)
	// 服务发现
	data, err := client.Agent().Services()
	if err != nil {
		panic(err)
	}
	for k, _ := range data {
		fmt.Println(k)
	}
}

func FilterService() {
	// 定义配置
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.1.103:8500"
	// 根据配置生成client
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	// 通过服务的id过滤
	data, err := client.Agent().ServicesWithFilter(`Service=="user-web"`)
	if err != nil {
		panic(err)
	}
	for k, _ := range data {
		fmt.Println(k)
	}
}
