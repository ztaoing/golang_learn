package main

import (
	"log"
	"net"
	"time"
)

func main() {
	// 1、openClose algorithm tcp session to server建立到8000端口的连接
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatalf("error to openClose tcp connection:%s", err)
	}
	defer conn.Close()

	// 2、向server写入数据
	log.Printf("tcp session openClose")
	b := []byte("hi gopher?")
	_, err = conn.Write(b)
	if err != nil {
		log.Fatalf("error writing tcp session:%s", err)
	}

	// 3、创建一个goroutine在10秒钟后关闭tcp  session
	go func() {
		<-time.After(time.Duration(10) * time.Second)
		defer conn.Close()
	}()

	// 4、接收response，直到发生错误
	for {
		d := make([]byte, 100)
		_, err = conn.Read(d)
		if err != nil {
			log.Fatalf("error reading session:%s", err)
		}
		log.Printf("reading data from server:%s", string(d))
	}

}
