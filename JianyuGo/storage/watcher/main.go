/**
* @Author:zhoutao
* @Date:2022/2/14 10:17
* @Desc:
 */

package main

import (
	"fmt"

	"runtime"

	"github.com/howeyc/fsnotify"
)

/**
注意这里使用的是：
github.com/howeyc/fsnotify

而不是
github.com/fsnotify/fsnotify
*/
var exit chan bool

func main() {
	//1、初始化监控对象watcher
	watcher, err := fsnotify.NewWatcher()

	if err != nil {

		fmt.Printf("Fail to create new Watcher[ %s ]\n", err)
	}

	//3、启动监听文件对象事件协程
	go func() {
		fmt.Println("开始监听文件变化")
		for {
			select {
			case e := <-watcher.Event:
				// 这里添加根据文件变化的业务逻辑
				fmt.Printf("监听到文件 - %s变化\n", e.Name)
				if e.IsCreate() {
					fmt.Println("监听到文件创建事件")
				}
				if e.IsDelete() {
					fmt.Println("监听到文件删除事件")
				}
				if e.IsModify() {
					fmt.Println("监听到文件修改事件")
				}
				if e.IsRename() {
					fmt.Println("监听到文件重命名事件")
				}
				if e.IsAttrib() {
					fmt.Println("监听到文件属性修改事件")
				}

				fmt.Println("根据文件变化开始执行业务逻辑")

			case err := <-watcher.Error:

				fmt.Printf(" %s\n", err.Error())
			}
		}
	}()

	// 2、将需要监听的文件加入到watcher的监听队列中
	paths := []string{"config.yml"}

	for _, path := range paths {

		err = watcher.Watch(path) //将文件加入监听

		if err != nil {

			fmt.Sprintf("Fail to watch directory[ %s ]\n", err)
		}
	}

	<-exit
	runtime.Goexit()
}
