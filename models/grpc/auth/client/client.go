package main

import (
	"context"
	"fmt"
	"golang_learn/golang_learn/models/grpc/proto"

	"google.golang.org/grpc"
)

// 实现了PerRPCCredentials接口
type customCredential struct{}

func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appid":  "101010",
		"appkey": "i am key",
	}, nil
}

// RequireTransportSecurity 指明credentials是否需要传输安全
func (c customCredential) RequireTransportSecurity() bool {
	return false
}

func main() {
	// 在拦截器中使用token
	/*interceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		md := metadata.New(map[string]string{
			"appid":  "10001",
			"appkey": "this is logic appkey",
		})
		ctx = metadata.NewOutgoingContext(context.Background(), md)
		err := invoker(ctx, method, req, reply, cc, opts...)
		fmt.Printf("耗时：%s\n", time.Since(start))
		return err
	}*/
	// 这里边有一个拦截器
	// WithPerRPCCredentials 返回的是一个DialOption，这个DialOption设置了credentials 和
	// 权限状态 作用在每一个发出的RPC请求.
	// PerRPCCredentials 定义了一个通用接口，这个接口用于将需要的安全信息添加到每一个rpc请求中(例如 oauth2).
	/*
		type PerRPCCredentials interface {
			// GetRequestMetadata 获取当前请求的metadata,如果需要的话要更新token，
			// tokens if required. 这个方法需要在每个请求的传输层被调用，数据需要保存在headers或者其他的context中。
			// context. 如果返回了一个状态码，它将会被用作rpc请求的状态码返回
			// 当底层实现支持的时候，ctx能够用于超时控制和取消控制。
			// 此外, RequestInfo的信息 将会通过ctx传递给这个调用请求
			// TODO(zhaoq): 定义一组有质量的key，而不是把它定义为无异议的字符串

			GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)

			// RequireTransportSecurity 指明credentials是否需要传输安全
			RequireTransportSecurity() bool
		}
	*/
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	/*
		// 传入了一个匿名方法
		// func(o *dialOptions) {
		//	o.copts.PerRPCCredentials = append(o.copts.PerRPCCredentials, creds)
		//}

		func WithPerRPCCredentials(creds credentials.PerRPCCredentials) DialOption {
			return newFuncDialOption(func(o *dialOptions) {
			o.copts.PerRPCCredentials = append(o.copts.PerRPCCredentials, creds)
		})
		}
		// 将传入的方法保存在funcDialOption中
		func newFuncDialOption(f func(*dialOptions)) *funcDialOption {
			return &funcDialOption{
			f: f,
		}
		}
		// 应用dialOptions的时候就是执行传入的匿名方法
		func (fdo *funcDialOption) apply(do *dialOptions) {
			fdo.f(do)
		}
	*/
	opts = append(opts, grpc.WithPerRPCCredentials(customCredential{}))
	conn, err := grpc.Dial("127.0.0.1:50001", opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewGreeterClient(conn)
	resp, err := c.SayHello(context.Background(), &proto.HelloRequest{
		Name: "tao",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Message)
}
