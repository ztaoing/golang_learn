作为一名 Gopher，我们很容易形成一个编程惯例：每当有一个实现了 io.Closer 接口的对象 x 时，
在得到对象并检查错误之后，会立即使用 defer x.Close() 以保证函数返回时 x 对象的关闭 。
以下给出两个惯用写法例子:

* HTTP 请求
  

    resp, err := http.Get("https://golang.google.cn/")
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    // The following code: handle resp
  
* 访问文件


      f, err := os.Open("/home/golangshare/gopher.txt")
      if err != nil {
          return err
      }
      defer f.Close()
      // The following code: handle f

[存在问题] 实际上，这种写法是存在潜在问题的。

    defer x.Close() 会忽略它的返回值，但在执行 x.Close() 时，
    我们并不能保证 x 一定能正常关闭，万一它返回错误应该怎么办？这种写法，
    会让程序有可能出现非常难以排查的错误。

那么，Close() 方法会返回什么错误呢？在 POSIX 操作系统中，
例如 Linux 或者 maxOS，关闭文件的 Close() 函数最终是调用了系统方法 close()，
我们可以通过 man close 手册，查看 close() 可能会返回什么错误：
    
    1ERRORS
    2     The close() system call will fail if:
    3
    4     [EBADF]            fildes is not a valid, active file descriptor.
    5
    6     [EINTR]            Its execution was interrupted by a signal.
    7
    8     [EIO]              A previously-uncommitted write(2) encountered an
    9                        input/output error.

* 错误 EBADF 表示无效文件描述符 fd，与本文中的情况无关；
* EINTR 是指的 Unix 信号打断； 
* 那么本文中可能存在的错误是 EIO。EIO 的错误是指未提交读，这是什么错误呢？
  EIO 错误是指文件的 write() 的读还未提交时就调用了 close() 方法。

  计算机存储层次结构，是一个经典的计算机存储器层级结构，在这个层次结构中，
  从上至下，设备的访问速度越来越慢，容量越来越大。
  存储器层级结构的主要思想是上一层的存储器作为低一层存储器的高速缓存。

CPU 访问寄存器会非常之快，相比之下，访问 RAM 就会很慢，而访问磁盘或者网络，那意味着就是蹉跎光阴。
如果每个 write() 调用都将数据同步地提交到磁盘，那么系统的整体性能将会极度降低，而我们的计算机是不会这样工作的。
当我们调用 write() 时，数据并没有立即被写到目标载体上，计算机存储器每层载体都在缓存数据，在合适的时机下，
将数据刷到下一层载体，这将写入调用的同步、缓慢、阻塞的同步转为了快速、异步的过程。

这样看来，EIO 错误的确是我们需要提防的错误。
这意味着如果我们尝试将数据保存到磁盘，
在 defer x.Close() 执行时，操作系统还并未将数据刷到磁盘，
这时我们应该获取到该错误提示
（只要数据还未落盘，那数据就没有持久化成功，它就是有可能丢失的，
例如出现停电事故，这部分数据就永久消失了，且我们会毫不知情）。
但是按照上文的惯例写法，我们程序得到的是 nil 错误。

【解决方案】
我们针对关闭文件的情况，来探讨几种可行性改造方案、
* 第一种方案，那就是不使用 defer
  

      func solution01() error {
          f, err := os.Create("/home/golangshare/gopher.txt")
          if err != nil {
              return err
          }
      
          if _, err = io.WriteString(f, "hello gopher"); err != nil {
                // 这种写法就需要我们在 io.WriteString 执行失败时，明确调用 f.Close() 进行关闭。
              f.Close()
              return err
          }
      
          return f.Close()
      }
这种写法就需要我们在 io.WriteString 执行失败时，明确调用 f.Close() 进行关闭。
但是这种方案，需要在每个发生错误的地方都要加上关闭语句 f.Close()，
如果对 f 的写操作 case 较多，容易存在遗漏关闭文件的风险。

* 第二种方案是，通过命名返回值 err 和闭包来处理
  

      func solution02() (err error) {
          f, err := os.Create("/home/golangshare/gopher.txt")
          if err != nil {
              return
          }
        
          // 闭包处理
          defer func() {
              closeErr := f.Close()
              if err == nil {
                  err = closeErr
              }
          }()
      
          _, err = io.WriteString(f, "hello gopher")
          return
      }
这种方案解决了方案一中忘记关闭文件的风险，如果有更多 if err !=nil 的条件分支，这种模式可以有效降低代码行数。

* 第三种方案是，在函数最后 return 语句之前，显示调用一次 f.Close()

    
      func solution03() error {
          f, err := os.Create("/home/golangshare/gopher.txt")
          if err != nil {
              return err
          }
          // 这个close因为defer不会返回结果
          defer f.Close()
      
          if _, err := io.WriteString(f, "hello gopher"); err != nil {
              return err
          }
          // 这里的close会返回结果
          if err := f.Close(); err != nil {
              return err
          }
          return nil
      }
这种解决方案能在 io.WriteString 发生错误时，由于 defer f.Close() 的存在能得到 close 调用。
也能在 io.WriteString 未发生错误，但缓存未刷新到磁盘时，获得 err := f.Close() 的错误信息，
而且由于 defer f.Close() 并不会返回错误，所以并不担心两次 Close() 调用会将错误覆盖。

* 最后一种方案是，函数 return 时执行 f.Sync()
  

      func solution04() error {
         f, err := os.Create("/home/golangshare/gopher.txt")
          if err != nil {
              return err
          }
          defer f.Close()
      
          if _, err = io.WriteString(f, "hello world"); err != nil {
              return err
          }
          // 强制执行刷新
          return f.Sync()
      }

由于调用 close() 是最后一次获取操作系统返回错误的机会，但是在我们关闭文件时，
缓存不一定被会刷到磁盘上。那么，我们可以调用 f.Sync() （其内部调用系统函数 fsync ）
强制性让内核将缓存持久到磁盘上去。

    // Sync commits the current contents of the file to stable storage.
    // Typically, this means flushing the file system's in-memory copy
    // of recently written data to disk.
    func (f *File) Sync() error {
        if err := f.checkValid("sync"); err != nil {
            return err
        }
        if e := f.pfd.Fsync(); e != nil {
            return f.wrapErr("sync", e)
        }
        return nil
    }
由于 fsync 的调用，这种模式能很好地避免 close 出现的 EIO。可以预见的是，
由于强制性刷盘，这种方案虽然能很好地保证数据安全性，但是在执行效率上却会大打折扣。
        
    


    

