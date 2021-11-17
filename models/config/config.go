package config

type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}

type ServerConfig struct {
	JWTInfo JWTConfig `mapstructure:"jwt"`
}

// 这是一个全局变量
var SigningKey string = "header.payload.signature"
