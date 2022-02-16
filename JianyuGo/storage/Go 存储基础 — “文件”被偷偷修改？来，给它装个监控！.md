我们总有这样的担忧：总有刁民想害朕，总有人偷偷在目录下删改文件，高危操作想第一时间了解，怎么办？
而且通常我们还有这样的需求：

* 监听一个目录中所有文件，文件大小到一定阀值，则处理；
* 监控某个目录，当有文件新增，立马处理；
* 监控某个目录或文件，当有文件被修改或者删除，立马能感知，进行处理；

怎么做到这个事情呢？最常见的通常有三个办法：
1、第一种：当事人主动通知你，这是侵入式的，需要当事人修改这部分代码来支持，依赖于当事人的自觉；
2、第二种：轮询观察，这个是无侵入式的，你可以自己写个轮询程序，每隔一段时间唤醒一次，对文件和目录做各种判断，从而得到这个目录的变化；
3、第三种：操作系统支持，以事件的方式通知到订阅这个事件的用户，达到及时处理的目的；

很明显，第三种最好：
纯旁路的逻辑，对线上程序无侵入；
操作系统直接支持，以事件的形式通知，性能也最好，100% 准确率（比较自己轮询判断要好）；

怎么做到这个事情呢？
既然是操作系统的支持，那么就涉及到系统调用。
系统调用直接使用略微复杂了些，Go 里面有个库 fsnotify ，就是封装了系统调用，用来监控文件事件的。
当指定目录或者文件，发生了创建，删除，修改，重命名的事件，里面就能得到通知。

【Go 的 fsnotify 的使用】
使用方法非常简单：
1、先用 fsnotify 创建一个监听器；
2、然后放到一个单独的 Goroutine 监听事件即可，通过 channel 的方式传递；
    
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

我们测试一下（有惊喜哦）。先把上述程序编译，然后跑起来：
再打开一个终端，准备进行你的操作：
先 touch 一个新文件 hello.txt
    
    touch hello.txt
使用 vim 打开这个文件，写入一行数据，然后关闭退出：
    
    vim hello.txt
    oot@ubuntu:~/code/gopher/src/notify# ./notify 
    # 触发事件：创建的时候
    2021/08/20 17:02:52 event: "./hello.txt": CREATE
    2021/08/20 17:02:52 event: "./hello.txt": CHMOD
    # 触发事件：vim 打开初始化的时候（创建 swp 文件）
    2021/08/20 17:17:08 event: "./.hello.txt.swp": CREATE
    2021/08/20 17:17:08 event: "./.hello.txt.swx": REMOVE
    2021/08/20 17:17:08 event: "./.hello.txt.swp": REMOVE
    2021/08/20 17:17:08 event: "./.hello.txt.swp": CREATE
    2021/08/20 17:17:08 event: "./.hello.txt.swp": WRITE
    2021/08/20 17:17:08 event: "./.hello.txt.swp": CHMOD
    # 触发事件：:w 写入保存的时候
    2021/08/20 17:17:53 event: "./4913": REMOVE
    2021/08/20 17:17:53 event: "./hello.txt": RENAME
    2021/08/20 17:17:53 event: "./hello.txt~": CREATE
    2021/08/20 17:17:53 event: "./hello.txt": CREATE
    2021/08/20 17:17:53 event: "./hello.txt": WRITE
    2021/08/20 17:17:53 event: "./hello.txt": CHMOD
    2021/08/20 17:17:53 event: "./hello.txt": CHMOD
    2021/08/20 17:17:53 event: "./hello.txt~": REMOVE
    # 触发事件：:q 的退出时候
    2021/08/20 17:17:57 event: "./.hello.txt.swp": WRITE
    2021/08/20 17:18:11 event: "./.hello.txt.swp": REMOVE


惊喜就是，这里能和之前 Linux 编辑器之神 vim 的 IO 存储原理 篇能结合上：
看到了 ~ 镜像文件，还看到了 swp 文件，竟然还看到了 一个 4913 的文件（这个文件也是个临时文件，感兴趣的可以了解一下）；

太神奇了，这样你就有一个新的手段监控你的文件发生的任何事情了。这是什么原理呢？

[深层原理:]
fsnotify 是跨平台的实现，奇伢这里只讲 Linux 平台的实现机制。
fsnotify 本质上就是对系统能力的一个浅层封装，主要封装了操作系统提供的两个机制：

* inotify 机制；
* epoll 机制；
  旁白：真的是何处都有 epoll 呀。如果还有对 epoll 不明白的赶紧复习下 Linux fd 系列，深度 epoll 剖析。

环境声明：Linux 内核版本 4.19
* inotify 机制
  什么是 inotify 机制？
  这是一个内核用于通知用户空间程序文件系统变化的机制。
  划重点：其实 inotify 机制的诞生源于一个通用的需求，由于IO/硬件管理都在内核，
        但用户态是有获悉内核事件的强烈需求，比如磁盘的热插拔，文件的增删改。
        这里就诞生了三个异曲同工的机制：hotplug 机制、udev 管理机制、inotify 机制。
  
  inotify 的三个接口:
  操作系统提供了三个接口来支撑，非常简洁：


    // fs/notify/inotify/inotify_user.c
    
    // 创建 notify fd
    inotify_init1
    
    // 添加监控路径
    inotify_add_watch
    
    // 删除一个监控
    inotify_rm_watch

用法非常简单，分别对应 inotify fd 的创建，监控的添加和删除。
    