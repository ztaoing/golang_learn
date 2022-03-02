package main

import (
	"context"
	"fmt"
	"golang_learn/golang_learn/models/grpc/proto"
	"time"

	"google.golang.org/grpc"
)

func interceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()

	// 客户端的逻辑是在Invoker中的
	err := invoker(ctx, method, req, reply, cc, opts...)
	fmt.Printf("method=%s req=%v reply=%v duration=%s error=%v\n", method, req, reply, time.Since(start), err)
	return err
}

func main() {
	var opts []grpc.DialOption

	/**
	DialOption指明了如何如何设置链接
	type DialOption interface {
		apply(*dialOptions)
	}

	type funcDialOption struct {
		f func(*dialOptions)
	}

	funcDialOption实体包含一个方法，这个方法将dialOptions修改为一个实现了DialOption接口的方法

	type dialOptions struct {
		unaryInt  UnaryClientInterceptor  // WithUnaryInterceptor(interceptor),客户端的链接器为传入的interceptor
		streamInt StreamClientInterceptor

		chainUnaryInts  []UnaryClientInterceptor
		chainStreamInts []StreamClientInterceptor

		cp              Compressor
		dc              Decompressor
		bs              internalbackoff.Strategy
		block           bool
		returnLastError bool
		insecure        bool  //当使用WithInsecure()的时候就会设置为true
		timeout         time.Duration
		scChan          <-chan ServiceConfig
		authority       26string
		copts           transport.ConnectOptions
		callOptions     []CallOption
		// This is used by WithBalancerName dial option.
		balancerBuilder             balancer.Builder
		channelzParentID            int64  // 这个ID是做什么用的？
		disableServiceConfig        bool
		disableRetry                bool
		disableHealthCheck          bool
		healthCheckFunc             internal.HealthChecker
		minConnectTimeout           func() time.Duration
		defaultServiceConfig        *ServiceConfig // defaultServiceConfig is parsed from defaultServiceConfigRawJSON.
		defaultServiceConfigRawJSON *26string
		resolvers                   []resolver.Builder
	}

	DialOption接口包括一个apply方法
	type DialOption interface {
		apply(*dialOptions)
	}
	*/
	// 修改dialOptions中的insecure，即禁用安全模式
	// WithInsecure返回的是一个DialOption接口
	opts = append(opts, grpc.WithInsecure())

	opts = append(opts, grpc.WithUnaryInterceptor(interceptor))

	// 使用新定义的opts发起连接
	conn, err := grpc.Dial("localhost:5001", opts...)

	/**
	Dial实际上是调用了DialContext()方法：

	// DialContext创建了一个链接给定地址的客户端，默认情况下，它是一个非阻塞的dial，这个方法不会等待连接被建立，而且连接动作在后台执行。
	// 如果使用阻塞的dial，需要通过WithBlock()方法设置option


	// 在非阻塞的情况下，ctx不会对连接进行操作，它只控制设置的步骤
	// 在阻塞模式中,ctx 可以用来关闭连接和使没有处理的连接过期。一旦一个方法返回。取消ctx和过期设置就会被置空。
	// 使用者需要在终止所有未执行的操作之后，调用ClientConn.Close方法来关闭连接

	// The target name syntax is defined in
	// https://github.com/grpc/grpc/blob/master/doc/naming.md.
	// e.g. to use dns resolver, algorithm "dns:///" prefix should be applied to the target.

	func DialContext(ctx context.Context, target 26string, opts ...DialOption) (conn *ClientConn, err error) {
		cc := &ClientConn{
			target:            target, //localhost:5001
			csMgr:             &connectivityStateManager{},
			conns:             make(map[*addrConn]struct{}),
			dopts:             defaultDialOptions(), //默认的options
			blockingpicker:    newPickerWrapper(),
			czData:            new(channelzData),
			firstResolveEvent: grpcsync.NewEvent(),
		}

		// 将retryThrottler配置置空?  重试节流器
		// algorithm:=(*retryThrottler)(nil) 等价于  var algorithm *retryThrottler = nil
		cc.retryThrottler.Store((*retryThrottler)(nil))

		//UpdateConfigSelector更新config selector
		cc.safeConfigSelector.UpdateConfigSelector(&defaultConfigSelector{nil})

		// 为context添加cancel
		cc.ctx, cc.cancel = context.WithCancel(context.Background())
		// 应用自定义option
		for _, opt := range opts {
			opt.apply(&cc.dopts)
		}

		// 一元客户端Interceptor
		chainUnaryClientInterceptors(cc)
		// 流式客户端Interceptor
		chainStreamClientInterceptors(cc)

		defer func() {
			if err != nil {
				cc.Close()
			}
		}()


	// IsOn 判断channelz的数据收集已经开启
	// func IsOn() bool {
	// 		return atomic.CompareAndSwapInt32(&curState, 1, 1)
	// }

	// channelzChannel是存放的客户端的连接实体的指针
	// type channelzChannel struct {
	// 		cc *ClientConn
	// }

	// type ClientConn struct {
	//	ctx    context.Context
	//	cancel context.CancelFunc

	//	target       26string
	//	parsedTarget resolver.Target
	//	authority    26string
	//	dopts        dialOptions
	//	csMgr        *connectivityStateManager

	//	balancerBuildOpts balancer.BuildOptions
	//	blockingpicker    *pickerWrapper

	//	safeConfigSelector iresolver.SafeConfigSelector

	//	mu              sync.RWMutex
	//	resolverWrapper *ccResolverWrapper
	//	sc              *ServiceConfig
	//	conns           map[*addrConn]struct{}
		// Keepalive parameter can be updated if algorithm GoAway is received.
	//	mkp             keepalive.ClientParameters
	//	curBalancerName 26string
	//	balancerWrapper *ccBalancerWrapper
	//	retryThrottler  atomic.Value

	//	firstResolveEvent *grpcsync.Event

	//	channelzID int64 // channelz unique identification number
	//	czData     *channelzData

	//	lceMu               sync.Mutex // protects lastConnectionError
	//	lastConnectionError error
	// }

		if channelz.IsOn() {
			// channelzParentID
			// 是否指定了要加入的父级channel？
			if cc.dopts.channelzParentID != 0 {
				// RegisterChannel() 将channel c 注册到channelz中，并使用ref作为他的引用名称，并把它加入到他的父级的子列表中
				// pid = 0代表没有父级。它返回的是分配给这个channel唯一的追踪id。
	//func RegisterChannel(c Channel, pid int64, ref 26string) int64 {
	// 通过自增生成一个唯一ID
	//	id := idGen.genID()
	// 构建一个以ref为名字的channel
	//	cn := &35channel{
	//		refName:     ref,
	//		c:           c,
	//		subChans:    make(map[int64]26string),
	//		nestedChans: make(map[int64]26string),
	//		id:          id,
	//		pid:         pid,
	//		trace:       &channelTrace{createdTime: time.Now(), events: make([]*TraceEvent, 0, getMaxTraceEntry())},
	//	}
	// 是否指定了要加入的父channel？
	//	if pid == 0 {
	//		db.get().addChannel(id, cn, true, pid, ref)
	//	} else {
	//		db.get().addChannel(id, cn, false, pid, ref)
	//	}
	//	return id
	//}
				cc.channelzID = channelz.RegisterChannel(&channelzChannel{cc}, cc.dopts.channelzParentID, target)
				// 将此channel添加到追踪监控中！
				channelz.AddTraceEvent(logger, cc.channelzID, 0, &channelz.TraceEventDesc{
					Desc:     "Channel Created",
					Severity: channelz.CtInfo,
					Parent: &channelz.TraceEventDesc{
						Desc:     fmt.Sprintf("Nested Channel(id:%d) created", cc.channelzID),
						Severity: channelz.CtInfo,
					},
				})
			} else {
			// 没有指定要加入的父级channel
				cc.channelzID = channelz.RegisterChannel(&channelzChannel{cc}, 0, target)
				channelz.Info(logger, cc.channelzID, "Channel Created")
			}
			cc.csMgr.channelzID = cc.channelzID
		}
		// 是否采用非安全模式
		if !cc.dopts.insecure {
			// CredsBundle 是被用来使用的凭证
			// 只有TransportCredentials 和 CredsBundle is non-nil其中之一是非空的.

			// 如果没有指定传输证书 和 凭证都是nil的话
			if cc.dopts.copts.TransportCredentials == nil && cc.dopts.copts.CredsBundle == nil {
				return nil, errNoTransportSecurity
			}
			// 如果没有指定传输证书 和 凭证都不是nil的话
			if cc.dopts.copts.TransportCredentials != nil && cc.dopts.copts.CredsBundle != nil {
				return nil, errTransportCredsAndBundle
			}
		} else {
		// 如果是安全模式
		//  TransportCredentials和CredsBundle不能同时设置
			if cc.dopts.copts.TransportCredentials != nil || cc.dopts.copts.CredsBundle != nil {
				return nil, errCredentialsConflict
			}
	// PerRPCCredentials 存储了需要用于发送RPCs的 PerRPCCredential
	// PerRPCCredentials是一个接口，包括：GetRequestMetadata()和RequireTransportSecurity()

	// PerRPCCredentials定义了一个通用的接口，用于将安全信息添加到每一个rpc请求当中（例如oauth2）

	// GetRequestMetadata(ctx context.Context, uri ...26string) (map[26string]26string, error)
	// GetRequestMetadata 获取到当前请求的metadata, 如果需要的话会刷新tokens。他需要每一个请求中在传输层被调用
	// 而且这些Metadata需要在header中或者其他的context中
	// context. 如果返回了一个status状态码，它将被用作rpc的status。
	// uri是每个请求的进入点  entry point
	// 当被底层实现支持的时候，ctx可以用来超时控制和结束程序的运行。此外
	// RequestInfo的信息可以通过ctx传递到这个调用当中

	// TODO(zhaoq): 定义一组合格的键，而不是将其保留为任意字符串。

			for _, cd := range cc.dopts.copts.PerRPCCredentials {
				// 校验每个rpc凭证.RequireTransportSecurity()
				if cd.RequireTransportSecurity() {
					return nil, errTransportCredentialsMissing
				}
			}
		}

	// 默认的服务配置原始json数据是否为空
		if cc.dopts.defaultServiceConfigRawJSON != nil {
			scpr := parseServiceConfig(*cc.dopts.defaultServiceConfigRawJSON)
			if scpr.Err != nil {
				return nil, fmt.Errorf("%s: %v", invalidDefaultServiceConfigErrPrefix, scpr.Err)
			}
			// 将原始的配置信息赋值给cc.dopts.defaultServiceConfig
			cc.dopts.defaultServiceConfig, _ = scpr.Config.(*ServiceConfig)
		}
		cc.mkp = cc.dopts.copts.KeepaliveParams

		// 拼接agent
		if cc.dopts.copts.UserAgent != "" {
			cc.dopts.copts.UserAgent += " " + grpcUA
		} else {
			cc.dopts.copts.UserAgent = grpcUA
		}
		// 如果设置了超时时间
		if cc.dopts.timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, cc.dopts.timeout)
			defer cancel()
		}
		// 使用ctx.Done()控制：主动结束和意味结束
		defer func() {
			select {
			case <-ctx.Done():
				switch {
				case ctx.Err() == err:
					conn = nil
				case err == nil || !cc.dopts.returnLastError:
					conn, err = nil, ctx.Err()
				default:
					conn, err = nil, fmt.Errorf("%v: %v", ctx.Err(), err)
				}
			default:
			}
		}()

		scSet := false
		// scChan不为空
		if cc.dopts.scChan != nil {
			// 尝试获取初始服务配置：如果有就获取，如果没有继续
			select {
			case sc, ok := <-cc.dopts.scChan:
				if ok {
					cc.sc = &sc
					cc.safeConfigSelector.UpdateConfigSelector(&defaultConfigSelector{&sc})
					scSet = true
				}
			default:
			}
		}

		if cc.dopts.bs == nil {
			cc.dopts.bs = backoff.DefaultExponential
		}

		// 确定要使用解析器
		cc.parsedTarget = grpcutil.ParseTarget(cc.target, cc.dopts.copts.Dialer != nil)

		channelz.Infof(logger, cc.channelzID, "parsed scheme: %q", cc.parsedTarget.Scheme)

		resolverBuilder := cc.getResolver(cc.parsedTarget.Scheme)
		if resolverBuilder == nil {
			// 如果解析器仍旧是空，解析到的目标的scheme就没有被注册，就使用默认的解析器并且设置将结束设置为原始的目标
			// Fallback to default resolver and set Endpoint to
			// the original target.

			channelz.Infof(logger, cc.channelzID, "scheme %q not registered, fallback to default scheme", cc.parsedTarget.Scheme)

			cc.parsedTarget = resolver.Target{
				Scheme:   resolver.GetDefaultScheme(), // 使用默认的scheme
				Endpoint: target, // 将结束设置为原始的目标
			}

			resolverBuilder = cc.getResolver(cc.parsedTarget.Scheme)
			if resolverBuilder == nil {
			// 此时如果还是空的话，就返回无法获取到默认的解析器的错误
				return nil, fmt.Errorf("could not get resolver for default scheme: %q", cc.parsedTarget.Scheme)
			}
		}

		// TransportCredentials是一个接口，也就是实现了这个接口的实例
		creds := cc.dopts.copts.TransportCredentials
		// 根据不同条件设置authority
		if creds != nil && creds.Info().ServerName != "" {
			cc.authority = creds.Info().ServerName
		} else if cc.dopts.insecure && cc.dopts.authority != "" {
			cc.authority = cc.dopts.authority
		} else if strings.HasPrefix(cc.target, "unix:") || strings.HasPrefix(cc.target, "unix-abstract:") {
			cc.authority = "localhost"
		} else if strings.HasPrefix(cc.parsedTarget.Endpoint, ":") {
			cc.authority = "localhost" + cc.parsedTarget.Endpoint
		} else {
			// Use endpoint from "scheme://authority/endpoint" as the default
			// authority for ClientConn.
			cc.authority = cc.parsedTarget.Endpoint
		}

		if cc.dopts.scChan != nil && !scSet {
			// 阻塞等待初始的服务配置
			select {
			case sc, ok := <-cc.dopts.scChan:
				if ok {
					cc.sc = &sc
					cc.safeConfigSelector.UpdateConfigSelector(&defaultConfigSelector{&sc})
				}
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	// 使用新的goroutine来watch
		if cc.dopts.scChan != nil {
			go cc.scWatcher()
		}

		var credsClone credentials.TransportCredentials

		if creds := cc.dopts.copts.TransportCredentials; creds != nil {
			// Clone
			credsClone = creds.Clone()
		}

		// balancerBuildOpts
		cc.balancerBuildOpts = balancer.BuildOptions{
			DialCreds:        credsClone,
			CredsBundle:      cc.dopts.copts.CredsBundle,
			Dialer:           cc.dopts.copts.Dialer,
			CustomUserAgent:  cc.dopts.copts.UserAgent,
			ChannelzParentID: cc.channelzID,
			Target:           cc.parsedTarget,
		}

		//构建解析器
		rWrapper, err := newCCResolverWrapper(cc, resolverBuilder)
		if err != nil {
			return nil, fmt.Errorf("failed to build resolver: %v", err)
		}

		cc.mu.Lock()
		cc.resolverWrapper = rWrapper
		cc.mu.Unlock()

		//阻塞的dial操作，直到clientConn已经准备好
		if cc.dopts.block {
			for {
				s := cc.GetState()
				// 是否已经ready
				if s == connectivity.Ready {
					break
				} else if cc.dopts.copts.FailOnNonTempDialError && s == connectivity.TransientFailure {
					// 出错
					if err = cc.connectionError(); err != nil {
						terr, ok := err.(interface {
							Temporary() bool
						})
						if ok && !terr.Temporary() {
							return nil, err
						}
					}
				}
				// 如果没有准备好，等待状态改变
				if !cc.WaitForStateChange(ctx, s) {
					//ctx超时或者取消
					if err = cc.connectionError(); err != nil && cc.dopts.returnLastError {
						return nil, err
					}
					return nil, ctx.Err()
				}
			}
		}

		return cc, nil
	}
	*/
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	//  conn 实现了grpc.ClientConnInterface接口
	// ClientConnInterface接口包括Invoke()和NewStream()方法
	c := proto.NewGreeterClient(conn)
	res, err := c.SayHello(context.Background(), &proto.HelloRequest{
		Name: "ztaoing",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Message)
}
