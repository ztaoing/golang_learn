package main

import (
	"github.com/mojocn/base64Captcha"
)

// 验证码
// DefaultMemStore 是一个通过New function创建的用于验证码的共享存储
var store = base64Captcha.DefaultMemStore
