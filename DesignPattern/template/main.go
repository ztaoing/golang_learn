package main

import "fmt"

//模板模式，如面向对象中的继承
//模板的作用是设定一个框架，但是不填充具体的内容细节，细节由具体的子类来完成

func main() {

}

type SMSSender struct {
}

func (s *SMSSender) Send(content string, receivers []string) {
	fmt.Printf("not implement yet")
}

type AliyunSMS struct {
	SMSSender
}

func (a *AliyunSMS) Send(content string, receivers []string) {
	fmt.Printf("调用阿里云短信发送api")
}

type TencentSMS struct {
	SMSSender
}

func (t *TencentSMS) Send(content string, receivers []string) {
	fmt.Printf("调用腾通讯短信发送api")
}
