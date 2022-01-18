package main

import (
	"log"
	"net"
	"time"
)

func main() {
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("error listener returned:%s", err)
	}
	defer l.Close()

	for {
		// c实现了net.conn接口
		c, err := l.Accept()
		if err != nil {
			log.Fatalf("error to accept connection:%s", err)
		}

		// 负责接收和发送的goroutine
		go func() {
			log.Println("tcp session open")
			defer c.Close()

			for {
				d := make([]byte, 100)
				_, err = c.Read(d)
				if err != nil {
					log.Printf("error reading session:%s", err)
					break
				}
				log.Printf("reading data from client:%s\n", string(d))

				_, err := c.Write(d)
				if err != nil {
					log.Printf("error writing tcp session:%s", err)
					break
				}
			}
		}()
		// 负责在10秒后关闭tcp连接的goroutine
		go func() {
			// TCPConn中的conn实现了conn接口
			err := c.(*net.TCPConn).SetLinger(0)
			if err != nil {
				log.Printf("error when setting linger:%s", err)
			}
			<-time.After(time.Duration(10) * time.Second)
			defer c.Close()
		}()
		/**
		11:02:42.872011 IP6 (flowlabel 0x20900, hlim 64, next-header TCP (6) payload length: 20)
		::1.8000 > ::1.51239: Flags [R.], cksum 0x001c (incorrect -> 0x69ba), seq 2212035947, ack 3397515779, win 6371, length 0
		[R]代表了rst包，用于重置连接。这个过程不需要被关闭方的回复，就可以关闭连接。
		*/

	}
}
