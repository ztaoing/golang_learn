文章来源于[Go学堂]

---

刚开始接触go时，发现go程序和php程序的其中一个不同是php是解释性语言，go是编译型语言，
即每次在有程序改动后，需要重新运行 go run或go build进行重新编译，更改才能生效，实则不便。

于是乎在网络上搜索发现了gowatch这个包，该包可通过监听当前目录下相关文件的变动，对go文件实时编译，提高研发效率。
那gowatch又是如何做到监听文件变化的呢?

通过阅读源码我们发现，在linux内核中，有一种用于通知用户空间程序文件系统变化的机制—Inotify。

它监控文件系统，并且及时向专门的应用程序发出相关的事件警告，比如删除、读、写和卸载操作等。
您还可以跟踪活动的源头和目标等细节。

---

Golang的标准库syscall实现了该机制。为进一步扩展，实现了fsnotify包实现了一个基于通道的、跨平台的实时监听接口。如下图：

                        gowatch             fsnotify package        syscall package             kernel/INotify接口
1、创建watcher对象       NewWatcher()------->fsnotify.NewWatcher()--->syscall.NotifyInit()------>fd=inotify_init()
                     获得watcher对象<-------将fd封装成watcher结构<------返回fd<------------------返回fd

2、将目录加入监听阶段    watcher.Watch(path)-->addWatch(path,event)----->InotifyAddWatch(path,event)-->inotify_add_watch(fd,path,mask)

3、读取监听事件阶段      e:=<-watcher.Event<---watcher.readEvent()<-----read(fd,buffer)<------------read(fd,buf,BUF_LEN)

4、应用逻辑执行阶段      根据事件类型执行应用逻辑

---

根据上图可知，监听文件的变化主要依赖于linux内核的INotify接口机制。Go的标准库中对其做了实现。
而fsnotify package的主要作用就是将进一步封装成watcher结构体和事件类型结构体的封装，从而实现事件的判断以及目录的监听。
下面看下 fsnotify package中对watcher的封装。


    type Watcher struct {
    
        mu sync.Mutex // Map access
    
        fd int // File descriptor (as returned by the inotify_init() syscall)
    
        watches map[string]*watch // Map of inotify watches (key: path)
    
        fsnFlags map[string]uint32 // Map of watched files to flags used for filter
    
        fsnmut sync.Mutex // Protects access to fsnFlags.
    
        paths map[int]string // Map of watched paths (key: watch descriptor)
    
        Error chan error // Errors are sent on this channel
    
        internalEvent chan *FileEvent // Events are queued on this channel
    
        Event chan *FileEvent // Events are returned on this channel
    
        done chan bool // Channel for sending a "quit message" to the reader goroutine
    
        isClosed bool // Set to true when Close() is first called
    
    }


---
[linux内核Inotify接口简介]

inotify中主要涉及3个接口。分别是inotify_init, inotify_add_watch,read。具体如下：

* int fd = inotify_init()	创建inotify实例，返回对应的文件描述符
* inotify_add_watch (fd, path, mask)	注册被监视目录或文件的事件
* read (fd, buf, BUF_LEN)	读取监听到的文件事件 


Inotify可以监听的文件系统事件列表：
        事件名称	                 事件说明
        IN_ACCESS                文件被访问
        IN_MODIFY                文件被 write
        IN_CLOSE_WRITE           可写文件被 close
        IN_OPEN                  文件被 open
        IN_MOVED_TO              文件被移来，如 mv、cp
        IN_CREATE                创建新文件
        IN_DELETE                文件被删除，如 rm 
        IN_DELETE_SELF           自删除，即一个可执行文件在执行时删除自己
        IN_MOVE_SELF	         自移动，即一个可执行文件在执行时移动自己
        IN_ATTRIB	             文件属性被修改，如 chmod、chown、touch 等
        IN_CLOSE_NOWRITE	     不可写文件被 close
        IN_MOVED_FROM	         文件被移走,如 mv
        IN_UNMOUNT	             宿主文件系统被 umount
        IN_CLOSE	             文件被关闭，等同于(IN_CLOSE_WRITE | IN_CLOSE_NOWRITE)
        IN_MOVE	                 文件被移动，等同于(IN_MOVED_FROM | IN_MOVED_TO)

---
[示例应用]
接下来是一个简易的示例应用，具体的应用实例可参考github.com/silenceper/gowatch包源代码 。

主要逻辑如下：
1、初始化watcher对象
2、将文件或目录加入到watcher监控对象的队列
3、启动监听协程，实时获取文件对象事件
