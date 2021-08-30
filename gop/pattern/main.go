package main

import "time"

// 告诉外边有没有等到：
// 如果等到了，就返回true和内容
// 如果没有等到，就返回false和空string
func nonBlockWait(c chan string) (string, bool) {
	select {
	case <-c:
		// 等到
		return "got", true
	default:
		// 没有等到
		return "", false
	}
}

func timeoutWait(c chan string, timeout time.Duration) (string, bool) {
	select {
	case <-c:
		return "got", true
	case <-time.After(timeout):
		// 超时，没有等到
		return "", false
	}
}
func main() {

}
