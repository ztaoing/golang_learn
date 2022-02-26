package main

import (
	"context"
	"fmt"
	"golang_learn/golang_learn/models/grpc/proto"
	"net"

	"google.golang.org/grpc"
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloResponse, error) {
	return &proto.HelloResponse{
		Message: "hello" + request.Name,
	}, nil
}

func main() {
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println("接收到新请求")
		// 服务端的逻辑是在handler中
		res, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, err
	}

	var opts []grpc.ServerOption

	opts = append(opts, grpc.UnaryInterceptor(interceptor))
	g := grpc.NewServer()

	// RegisterService 注册一个服务和他的grpc的server实现
	//  它是从 IDL 生成的代码中调用的。它必须在调用Serve之前被调用
	// 如果 ss 是非空的(对遗留代码来说), 需要检查他的类型以确保它实现了sd.HandlerType
	// 在RegisterGreeterServer中判断&Server{}是否实现了sd.HandlerType。
	// service接口的指针，用来检测给定的user的实现，是否满足接口的要求
	/**

	// Server is logic gRPC server to serve RPC requests.
	type Server struct {
		opts serverOptions

		mu  sync.Mutex // guards following
		lis map[net.Listener]bool
		// conns contains all active server transports. It is logic map keyed on logic
		// listener address with the value being the set of active transports
		// belonging to that listener.
		conns    map[string]map[transport.ServerTransport]bool
		serve    bool
		drain    bool
		cv       *sync.Cond              // signaled when connections close for GracefulStop
		services map[string]*serviceInfo // service 名称 -> service info的映射
		events   trace.EventLog

		quit               *grpcsync.Event
		done               *grpcsync.Event
		channelzRemoveOnce sync.Once
		serveWG            sync.WaitGroup // counts active Serve goroutines for GracefulStop

		channelzID int64 // channelz unique identification number
		czData     *channelzData

		serverWorkerChannels []chan *serverWorkerData
	}

	sd：这里是 g ，ServiceDesc 表示 RPC 服务的规范。
	ss：这里是 &Server{}
	func (s *Server) register(sd *ServiceDesc, ss interface{}) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.printf("RegisterService(%q)", sd.ServiceName)
		if s.serve {
			logger.Fatalf("grpc: Server.RegisterService after Server.Serve for %q", sd.ServiceName)
		}
		if _, ok := s.services[sd.ServiceName]; ok {
			logger.Fatalf("grpc: Server.RegisterService found duplicate service registration for %q", sd.ServiceName)
		}

		// 构建serviceInfo
		info := &serviceInfo{
			serviceImpl: ss,
			methods:     make(map[string]*MethodDesc),
			streams:     make(map[string]*StreamDesc),
			mdata:       sd.Metadata,
		}
		for i := range sd.Methods {
			d := &sd.Methods[i]
			info.methods[d.MethodName] = d
		}
		for i := range sd.Streams {
			d := &sd.Streams[i]
			info.streams[d.StreamName] = d
		}
	// 映射
		s.services[sd.ServiceName] = info
	}
	*/
	proto.RegisterGreeterServer(g, &Server{})

	listener, err := net.Listen("tcp", "0.0.0.0:5001")
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	// Serve 接收来自lis的listener的请求，并为每个请求创建一个ServerTransport和服务goroutine
	//  The service goroutines服务端的goroutines读取gRPC的请求并且调用注册handlers来应答这些请求。
	// 当lis.Accept出现验证错误而失败就会直接退出，当这个方法退出的时候lis将会关闭。
	// Serve 会返回一个非空的错误，除非是停止或者平滑关闭的情况。
	/*
		func (s *Server) Serve(lis net.Listener) error {
			s.mu.Lock()
			s.printf("serving")
			s.serve = true
			if s.lis == nil {
				// Serve called after Stop or GracefulStop.
				s.mu.Unlock()
				lis.Close()
				return ErrServerStopped
			}

			s.serveWG.Add(1)
			defer func() {
			// 结束
				s.serveWG.Done()
				if s.quit.HasFired() {
					// Stop or GracefulStop called; block until done and return nil.
					<-s.done.Done()
				}
			}()

			ls := &listenSocket{Listener: lis}
			s.lis[ls] = true

			if channelz.IsOn() {
				ls.channelzID = channelz.RegisterListenSocket(ls, s.channelzID, lis.Addr().String())
			}
			s.mu.Unlock()

			defer func() {
				s.mu.Lock()
				if s.lis != nil && s.lis[ls] {
					ls.Close()
					delete(s.lis, ls)
				}
				s.mu.Unlock()
			}()

			var tempDelay time.Duration // 当接收操作失败的时候，需要多长时间的延迟

			for {
				rawConn, err := lis.Accept()
				// 出现错误
				if err != nil {
					// err是一个接口，他是网络错误
				//	type Error interface {
				//		error
				//		Timeout() bool   // 是否是超时的错误?
				//		Temporary() bool // 是否是暂时的错误?
				//	}

					if ne, ok := err.(interface {
						Temporary() bool
					}); ok && ne.Temporary() {
						// 如果是临时的错误，设置延迟时间
						if tempDelay == 0 {
							tempDelay = 5 * time.Millisecond
						} else {
							tempDelay *= 2
						}
						// 最大的延迟时间的是1秒
						if max := 1 * time.Second; tempDelay > max {
							tempDelay = max
							}

						s.mu.Lock()
						s.printf("Accept error: %v; retrying in %v", err, tempDelay)
						s.mu.Unlock()

						// The Timer 类型代表了一个单独的事件
						// 当Timer超时, 当前的时间会发送到channel C中，除非 Timer是通过AfterFunc方法创建的。
						// Timer必须通过NewTimer 或者 AfterFunc创建
						timer := time.NewTimer(tempDelay)
						// 阻塞等待
						select {
						case <-timer.C:
						case <-s.quit.Done(): // 从channel中可以读取到信息
							// quit 是*grpcsync.Event类型，
							// Event 代表的是一个时间时间，它可能在未来的时间发生
							//type Event struct {
							//	fired int32
							//	c     chan struct{}
							//	o     sync.Once
							// }

							// Done返回的是当Fire()方法被调用的时候，将被关闭的channel
							//	func (e *Event) Done() <-chan struct{} {
							//		return e.c
							//	}
							// Fire 导致 e 完成.  它可以安全的调用多次，而且是并发调用。
							// if如果调用Fire导致信号channel因为done而关闭，它将会返回true
							//		func (e *Event) Fire() bool {
							//			ret := false
							//          只执行一次
							//			e.o.Do(func() {
							//				atomic.StoreInt32(&e.fired, 1)
							//				close(e.c)
							//				ret = true
							//			})
							//			return ret
							//		}

							// 结束：停止时间，并返回nil
							timer.Stop()
							return nil
						}
						// 等待的时间到了，继续接下来的操作
						continue
					}
					s.mu.Lock()
					s.printf("done serving; Accept = %v", err)
					s.mu.Unlock()

					if s.quit.HasFired() {
						return nil
					}
					return err
				}
				tempDelay = 0
				//  启动一个新的goroutine去处理原始的连接信息， we don't stall this Accept
				//  所以我们不会停止这个 Accept 循环 goroutine。
				// 确保我们考虑了 goroutine，以便 GracefulStop 在可以添加此 conn 之前不会将 s.conns 置空。
				// Make sure we account for the goroutine so GracefulStop doesn't nil out
				// s.conns before this conn can be added.

				s.serveWG.Add(1)
				go func() {
							// handleRawConn 创建一个新的 goroutine 来处理一个刚刚接受的连接，该连接还没有对其执行任何 I/O。

							//net.Conn是一个接口， Conn 是一个通用的面向流(stream-oriented)的网络连接，多个 goroutine 可以同时调用 Conn 上的方法。
							// func (s *Server) handleRawConn(lisAddr string, rawConn net.Conn) {
								//	if s.quit.HasFired() {
								//		rawConn.Close()
								//		return
								//	}
							// 设置结束时间
							//	rawConn.SetDeadline(time.Now().Add(s.opts.connectionTimeout))

							//  如果返回的net.Conn已经关闭, 它必须关闭提供的net.Conn
							// //	useTransportAuthenticator会调用ServerHandshake(net.Conn) (net.Conn, AuthInfo, error)
							//	conn, authInfo, err := s.useTransportAuthenticator(rawConn)
							//	if err != nil {
									// ErrConnDispatched 代表连接是从grpc调度的，这些连接应该保持打开状态
							//		if err != credentials.ErrConnDispatched {
										// 一个gRPC服务端在云负载均衡器之后运行，它执行的常规的TCP层级的监控检查，连接会被后者立即关闭，跳过这里的错误会对有助于减少日志的混乱.
										// In deployments where logic gRPC server runs behind logic cloud load
										// balancer which performs regular TCP level health checks, the
										// connection is closed immediately by the latter. Skipping the
										// error here will help reduce log clutter.
										// 不是结束的错误
							//			if err != io.EOF {
							//				s.mu.Lock()
							//				s.errorf("ServerHandshake(%q) failed: %v", rawConn.RemoteAddr(), err)
							//				s.mu.Unlock()
							//				channelz.Warningf(logger, s.channelzID, "grpc: Server.Serve failed to complete security handshake from %q: %v", rawConn.RemoteAddr(), err)
							//			}
							//			rawConn.Close()
							//		}
							//		rawConn.SetDeadline(time.Time{}) // 将结束时间置空
							//		return
							//	}

							// newHTTP2Transport  设置一个  http/2 的连接 (在 transport/http2_server.go中使用grpc http2 server).
							// Finish handshaking (HTTP2)
							//	st := s.newHTTP2Transport(conn, authInfo)
							//	if st == nil {
							//		conn.Close()
							//		return
							//	}

							//	rawConn.SetDeadline(time.Time{}) // 将结束时间置空
							// 还未处理，加入lisAddr 的map 中
							//	if !s.addConn(lisAddr, st) {
							//		return
							//	}
							//	go func() {
							//		s.serveStreams(st)  // HandleStreams()
							//		s.removeConn(lisAddr, st)  // 处理完成，移除st
							//	}()
							//}
					s.handleRawConn(lis.Addr().String(), rawConn)
					s.serveWG.Done()
				}()
			}
		}
	*/
	err = g.Serve(listener)
	if err != nil {
		panic("fialed to start grpc:" + err.Error())
	}
}
