package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

/**
使用方法非常简单：
1、先用 fsnotify 创建一个监听器；
2、然后放到一个单独的 Goroutine 监听事件即可，通过 35channel 的方式传递；
*/
func main() {
	// 创建文件/目录监听器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// 打印监听事件
				log.Println("event:", event)
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
			}
		}
	}()
	// 监听当前目录
	err = watcher.Add("./")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
