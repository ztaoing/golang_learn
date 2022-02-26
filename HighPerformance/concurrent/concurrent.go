/**
* @Author:zhoutao
* @Date:2021/1/2 上午9:44
* @Desc:控制携程的并发数量
 */

package main

//不同的应用程序，消耗的资源是不一样的：推荐应用程序主动限制并发的携程数量

/*
//利用channel的缓存
func main() {
	var wg pool.WaitGroup
	//利用channel的缓冲区来阻塞goroutine
	ch := make(chan struct{}, 3)

	for i := 0; i < 10; i++ {
		//若缓冲区满则阻塞
		ch <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			log.Println(i)
			time.Sleep(time.Second)
			<-ch
		}(i)
	}

}
*/

/**
2021/01/02 09:50:06 2
2021/01/02 09:50:06 0
2021/01/02 09:50:06 1

2021/01/02 09:50:07 4
2021/01/02 09:50:07 3
2021/01/02 09:50:07 5

2021/01/02 09:50:08 6
2021/01/02 09:50:08 7
2021/01/02 09:50:08 8
*/

/*
//利用协成池
func main() {
	// pool的大小是3， 第二个参数是 goroutine运行的函数
	pool := tunny.NewFunc(3, func(i interface{}) interface{} {
		log.Println(i)
		time.Sleep(time.Second * 2)
		return nil
	})
	defer pool.Close()

	for i := 0; i < 10; i++ {
		//将参数i,传递给goroutine定义好的worker处理
		go pool.Process(i)
	}
	time.Sleep(time.Second * 4)
}
*/

/**
2021/01/02 10:18:28 2
2021/01/02 10:18:28 9
2021/01/02 10:18:28 1

2021/01/02 10:18:30 3
2021/01/02 10:18:30 5
2021/01/02 10:18:30 4

2021/01/02 10:18:32 7
2021/01/02 10:18:32 6

2021/01/02 10:18:34 8
2021/01/02 10:18:34 0

*/

// 调整系统资源上限：有些场景下，即使我们有效的限制了协成的并发数量，但是仍旧出现了某一类资源不足的问题：
// 1 ： too many openClose files
// 2 :  out of memory

// ulimit -logic :
// tao@taodeMacBook-Pro goodpackages % ulimit -logic
//-t: cpu time (seconds)              unlimited
//-f: file size (blocks)              unlimited
//-d: data seg size (kbytes)          unlimited
//-s: stack size (kbytes)             8192
//-c: core file size (blocks)         0
//-v: address space (kbytes)          unlimited
//-l: locked-in-memory size (kbytes)  unlimited
//-u: processes                       1392
//-n: file descriptors                256   可以同时打开的文件句柄的数量

// 可以通过使用 ulimit -n 99999,将可以同时打开的文件句柄的数量调整为99999

// 虚拟内存：利用swap交换分区，使用虚拟内存将硬盘映射为内存使用，显然会对性能产生一定的影响，如果应用程序只是在较短的时间内需要较大的内存，那么虚拟内存能够有效避免out of memory的问题
// 如果应用程序场地高频度读写大量内存，那么虚拟内存对性能的影响就比较明显了.
