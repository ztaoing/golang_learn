在客户端代码中通过调用Conn.Close()方法发起了关闭TCP连接的请求，这是一种默认的关闭连接的方式：
    默认关闭需要四次回收的确认过程，这是一种商量的方式，而tcp为我们提供了另外一种强制的关闭模式
    如何强制关闭呢？具体在go代码中如何实现呢？

    以客户端主动关闭连接为例，当它调用close后，会向服务端发送FIN报文，如果服务端socket接收缓冲区里已经没有数据了，
    那么服务端的read将会得到一个EOF错误。
    发起关闭会精力：FIN_WAIT_1-->FIN_WAIT_2->TIME_WAIT->CLOSE的状态变化，这些状态需要得到被关闭方的反馈而发生更新。
    
    强制关闭：
    默认的关闭方式，不管是客户端还是服务端主动发起关闭，都要经过对方的应答，才能最终实现真正的关闭连接。

    在发起关闭时，能不能不经过对方的同意，就结束掉连接呢？可以
    解：tcp协议提供了一个rst的标志位，当连接的一方认为该连接异常时，可以通过发送rst包来立即关闭连接，而不用等待
        被关闭方的ack确认。

    在go中可以通过net.TCPConn.SetLinger()方法来实现:函数的注释已经非常清晰了，但是需要读者要有socket缓冲区的概念。
    
    // SetLinger sets the behavior of Close on a connection which still
    // has data waiting to be sent or to be acknowledged.
    //
    // sec < 0 (默认的关闭方式), 操作系统会将缓冲区中未完全处理完的数据都处理完，然后关闭连接
    //
    //
    // sec == 0（强制关闭）, 操作系统会直接丢弃缓冲区中的数据
    //
    //sec > 0,操作系统会以默认关闭方式运行，但是超过定义的时间sec后，如果还没处理完缓冲区中的数据， the data is sent in the background as with sec < 0. On
    // 在某些操作系统下，缓冲区中的数据可能会被丢弃

    func (c *TCPConn) SetLinger(sec int) error {}

    socket缓冲区：
    当应用层通过socket进行读与写操作师，实质上经过了一层socket缓冲区，它分为发送缓冲区和接收缓冲区。

    缓冲区可以通过netstat -nt命令查看
    $ netstat -nt
    Active Internet connections
    Proto Recv-Q Send-Q  Local Address          Foreign Address        (state)
    tcp4       0      0  127.0.0.1.57721        127.0.0.1.49448        ESTABLISHED

    Recv-Q：接收缓冲区
    Send-Q：发送缓冲区
    
    

    