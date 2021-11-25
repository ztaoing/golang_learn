package main

import (
	"context"
	"fmt"
	"golang_learn/golang_learn/models/grpc/proto"

	"google.golang.org/grpc"
)

type customCredential struct{}

func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appid":  "101010",
		"appkey": "i am key",
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires
// transport security.
func (c customCredential) RequireTransportSecurity() bool {
	return false
}

func main() {
	// 在拦截器中使用token
	/*interceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		md := metadata.New(map[string]string{
			"appid":  "10001",
			"appkey": "this is a appkey",
		})
		ctx = metadata.NewOutgoingContext(context.Background(), md)
		err := invoker(ctx, method, req, reply, cc, opts...)
		fmt.Printf("耗时：%s\n", time.Since(start))
		return err
	}*/
	// 这里边有一个拦截器
	// WithPerRPCCredentials returns a DialOption which sets credentials and places
	// auth state on each outbound RPC.
	// PerRPCCredentials defines the common interface for the credentials which need to
	// attach security information to every RPC (e.g., oauth2).
	/*
		type PerRPCCredentials interface {
			// GetRequestMetadata gets the current request metadata, refreshing
			// tokens if required. This should be called by the transport layer on
			// each request, and the data should be populated in headers or other
			// context. If a status code is returned, it will be used as the status
			// for the RPC. uri is the URI of the entry point for the request.
			// When supported by the underlying implementation, ctx can be used for
			// timeout and cancellation. Additionally, RequestInfo data will be
			// available via ctx to this call.
			// TODO(zhaoq): Define the set of the qualified keys instead of leaving
			// it as an arbitrary string.
			GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
			// RequireTransportSecurity indicates whether the credentials requires
			// transport security.
			RequireTransportSecurity() bool
		}
	*/
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
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
