源码分析:

用的go版本是1.12.7。
从发起一个网络请求开始:这里很多代码，其实只是为了展示这部分代码是怎么跟踪下来的，方便大家去看源码的时候去跟一下。

    1res, err := client.Do(req)
    2func (c *Client) Do(req *Request) (*Response, error) {
    3    return c.do(req)
    4}
    5
    6func (c *Client) do(req *Request) {
    7    // ...
    8    if resp, didTimeout, err = c.send(req, deadline); err != nil {
    9    // ...
    10  }
    11    // ...  
    12}  
    13func send(ireq *Request, rt RoundTripper, deadline time.Time) {
    14    // ...    
    15    resp, err = rt.RoundTrip(req)
    16     // ...  
    17}
    18
    19// 从这里进入 RoundTrip 逻辑
    20/src/net/http/roundtrip.go: 16
    21func (t *Transport) RoundTrip(req *Request) (*Response, error) {
    22    return t.roundTrip(req)
    23}
    24
    25func (t *Transport) roundTrip(req *Request) (*Response, error) {
    26    // 尝试去获取一个空闲连接，用于发起 http 连接
    27  pconn, err := t.getConn(treq, cm)
    28  // ...
    29}
    30
    31// 重点关注这个函数，返回是一个长连接
    32func (t *Transport) getConn(treq *transportRequest, cm connectMethod) (*persistConn, error) {
    33  // 省略了大量逻辑，只关注下面两点
    34    // 有空闲连接就返回
    35    pc := <-t.getIdleConnCh(cm)
    36
    37  // 没有创建连接
    38  pc, err := t.dialConn(ctx, cm)
    39
    40}
    
    最后一个上面的代码里有个 getConn 方法。在发起网络请求的时候，会先取一个网络连接，取连接有两个来源。
    如果有空闲连接，就拿空闲连接:
        /src/net/http/tansport.go:810
    2func (t *Transport) getIdleConnCh(cm connectMethod) chan *persistConn {
    3 // 返回放空闲连接的chan
    4 ch, ok := t.idleConnCh[key]
    5   // ...
    6 return ch
    7}

    没有空闲连接，就创建长连接:
    1/src/net/http/tansport.go:1357
        2func (t *Transport) dialConn() {
        3  //...
        4  conn, err := t.dial(ctx, "tcp", cm.addr())
        5  // ...
        6  go pconn.readLoop()
        7  go pconn.writeLoop()
        8  // ...
        9}

    当第一次发起一个http请求时，这时候肯定没有空闲连接，会建立一个新连接。
    同时会创建一个读goroutine和一个写goroutine。
    
    注意上面代码里的t.dial(ctx, "tcp", cm.addr())，如果像文章开头那样设置了 http.Transport的:
        1  Dial: func(netw, addr string) (net.Conn, error) {
        2   conn, err := net.DialTimeout(netw, addr, time.Second*2) //设置建立连接超时
        3   if err != nil {
        4      return nil, err
        5   }
        6   err = conn.SetDeadline(time.Now().Add(time.Second * 3)) //设置发送接受数据超时
        7   if err != nil {
        8      return nil, err
        9   }
        10   return conn, nil
        11},
    
    那么这里就会在下面的dial里被执行到:
    1   func (t *Transport) dial(ctx context.Context, network, addr string) (net.Conn, error) {
    2   // ...
    3  c, err := t.Dial(network, addr)
    4  // ...
    5}
    
    这里面调用的设置超时，会执行到:
        /src/net/net.go
    2func (c *conn) SetDeadline(t time.Time) error {
    3    //...
    4    c.fd.SetDeadline(t)
    5    //...
    6}
    7
    8//...
    9
    10func setDeadlineImpl(fd *FD, t time.Time, mode int) error {
    11    // ...
    12    runtime_pollSetDeadline(fd.pd.runtimeCtx, d, mode)
    13    return nil
    14}
    15
    16
    17//go:linkname poll_runtime_pollSetDeadline internal/poll.runtime_pollSetDeadline
    18func poll_runtime_pollSetDeadline(pd *pollDesc, d int64, mode int) {
    19    // ...
    20  // 设置一个定时器事件
    21  rtf = netpollDeadline
    22    // 并将事件注册到定时器里
    23  modtimer(&pd.rt, pd.rd, 0, rtf, pd, pd.rseq)
    24}  
    
    上面的源码，简单来说就是，当第一次调用请求的，会建立个连接，这时候还会注册一个定时器事件，假设时间设了3s，
    那么这个事件会在3s后发生，然后执行注册事件的逻辑。而这个注册事件就是netpollDeadline。
    注意这个netpollDeadline，待会会提到。

    设置了超时事件，且超时事件是3s后之后，发生。在此期间正常收发数据。一切如常。
    直到3s过后，这时候看读goroutine，会等待网络数据返回。
            1/src/net/http/tansport.go:1642
        2func (pc *persistConn) readLoop() {
        3    //...
        4    for alive {
        5        _, err := pc.br.Peek(1)  // 阻塞读取服务端返回的数据
        6    //...
        7}
    
    然后就是一直跟代码。
    src/bufio/bufio.go: 129
    2func (b *Reader) Peek(n int) ([]byte, error) {
    3   // ...
    4   b.fill()
    5   // ...   
    6}
    7
    8func (b *Reader) fill() {
    9    // ...
    10    n, err := b.rd.Read(b.buf[b.w:])
    11    // ...
    12}
    13
    14/src/net/http/transport.go: 1517
    15func (pc *persistConn) Read(p []byte) (n int, err error) {
    16    // ...
    17    n, err = pc.conn.Read(p)
    18    // ...
    19}
    20
    21// /src/net/net.go: 173
    22func (c *conn) Read(b []byte) (int, error) {
    23    // ...
    24    n, err := c.fd.Read(b)
    25    // ...
    26}
    27
    28func (fd *netFD) Read(p []byte) (n int, err error) {
    29    n, err = fd.pfd.Read(p)
    30    // ...
    31}
    32
    33/src/internal/poll/fd_unix.go:
    34func (fd *FD) Read(p []byte) (int, error) {
    35    //...
    36  if err = fd.pd.waitRead(fd.isFile); err == nil {
    37    continue
    38  }
    39    // ...
    40}
    41
    42func (pd *pollDesc) waitRead(isFile bool) error {
    43    return pd.wait('r', isFile)
    44}
    45
    46func (pd *pollDesc) wait(mode int, isFile bool) error {
    47    // ...
    48  res := runtime_pollWait(pd.runtimeCtx, mode)
    49    return convertErr(res, isFile)
    50}

    直到跟到 runtime_pollWait，这个可以简单认为是等待服务端数据返回：
        1//go:linkname poll_runtime_pollWait internal/poll.runtime_pollWait
        2func poll_runtime_pollWait(pd *pollDesc, mode int) int {
        3
        4    // 1.如果网络正常返回数据就跳出
        5  for !netpollblock(pd, int32(mode), false) {
        6    // 2.如果有出错情况也跳出
        7        err = netpollcheckerr(pd, int32(mode))
        8        if err != 0 {
        9            return err
        10        }
        11    }
        12    return 0
        13}

    整条链路跟下来，就是会一直等待数据，等待的结果只有两个：

    1、有可以读的数据
    2、出现报错

    这里面的报错，又有那么两种：

    1、连接关闭
    2、超时

    1func netpollcheckerr(pd *pollDesc, mode int32) int {
    2    if pd.closing {
    3        return 1 // errClosing
    4    }
    5    if (mode == 'r' && pd.rd < 0) || (mode == 'w' && pd.wd < 0) {
    6        return 2 // errTimeout
    7    }
    8    return 0
    9}

    其中提到的超时，就是指这里面返回的数字2，会通过下面的函数，转化为 ErrTimeout，
    而 ErrTimeout.Error() 其实就是 i/o timeout。
        1func convertErr(res int, isFile bool) error {
        2    switch res {
        3    case 0:
        4        return nil
        5    case 1:
        6        return errClosing(isFile)
        7    case 2:
        8        return ErrTimeout // ErrTimeout.Error() 就是 "i/o timeout"
        9    }
        10    println("unreachable: ", res)
        11    panic("unreachable")
        12}

    那么问题来了。上面返回的超时错误，也就是返回2的时候的条件是怎么满足的？
    1    if (mode == 'r' && pd.rd < 0) || (mode == 'w' && pd.wd < 0) {
    2        return 2 // errTimeout
    3    }

    还记得刚刚提到的 netpollDeadline吗？这里面放了定时器3s到点时执行的逻辑。
        1func timerproc(tb *timersBucket) {
        2    // 计时器到设定时间点了，触发之前注册函数
        3    f(arg, seq) // 之前注册的是 netpollDeadline
        4}
        5
        6func netpollDeadline(arg interface{}, seq uintptr) {
        7    netpolldeadlineimpl(arg.(*pollDesc), seq, true, true)
        8}
        9
        10/src/runtime/netpoll.go: 428
        11func netpolldeadlineimpl(pd *pollDesc, seq uintptr, read, write bool) {
        12    //...
        13    if read {
        14        pd.rd = -1
        15        rg = netpollunblock(pd, 'r', false)
        16    }
        17    //...
        18}

    这里会设置pd.rd=-1，是指 poller descriptor.read deadline ，含义网络轮询器文件描述符的读超时时间， 我们知道在linux里万物皆文件，这里的文件其实是指这次网络通讯中使用到的socket。

    这时候再回去看发生超时的条件就是if (mode == 'r' && pd.rd < 0)。

    至此。我们的代码里就收到了 io timeout 的报错。

    总结：
    1、不要在 http.Transport中设置超时，那是连接的超时，不是请求的超时。否则可能会出现莫名 io timeout报错。

    2、请求的超时在创建client里设置。
