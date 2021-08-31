package main

import "time"

//实现两个go轮流输出：A1B2C3.....Z26
func main() {

}

//方法一：有缓冲chan
func ChannelFunc() {
	zimu := make(chan int, 1)
	suzi := make(chan int, 1)
	zimu <- 0
	// zimu
	go func() {
		for i := 65; i <= 90; i++ {
			<-zimu
			fmt.Printf("%v", string(rune(i)))
			suzi <- i
		}
		return
	}()

	go func() {
		for i := 1; i <= 26; i++ {
			<-suzi
			fmt.Printf("%v", i)
			zimu <- i
		}
		return
	}()

	time.Sleep(1 * time.Second)
	fmt.Println()
}

//方法二：无缓冲chan
func Channel1Func() {
	zimu := make(chan int)
	suzi := make(chan int)

	// zimu
	go func() {
		for i := 65; i <= 90; i++ {
			fmt.Printf("%v", string(rune(i)))
			zimu <- i
			<-suzi
		}
		return
	}()

	go func() {
		for i := 1; i <= 26; i++ {
			<-zimu
			fmt.Printf("%v", i)
			suzi <- i
		}
		return
	}()

	time.Sleep(10 * time.Second)
	fmt.Println()
}

//方法三：使用锁
