package main

import (
	"context"
	"fmt"
	"golang_learn/golang_learn/models/grpc/proto"
	"net"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) SayHello(ctx context.Context, r *proto.HelloRequest) (*proto.HelloResponse, error) {
	// 从context中获取metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println("get metadata error")
	}
	//获取md中的key-value
	if nameSlice, ok := md["name"]; ok {
		for i, e := range nameSlice {
			fmt.Sprintln(i, e)
		}
	}
	return &proto.HelloResponse{
		Message: "metadata have got",
	}, nil
}

func main() {
	s := grpc.NewServer()
	/*
		type Server struct {
			opts serverOptions

			mu  sync.Mutex // guards following
			lis map[net.Listener]bool
			// conns contains all active server transports. It is a map keyed on a
			// listener address with the value being the set of active transports
			// belonging to that listener.
			conns    map[string]map[transport.ServerTransport]bool
			serve    bool
			drain    bool
			cv       *sync.Cond              // signaled when connections close for GracefulStop
			services map[string]*serviceInfo // service 名字 -> service info
			events   trace.EventLog

			quit               *grpcsync.Event
			done               *grpcsync.Event
			channelzRemoveOnce sync.Once
			serveWG            sync.WaitGroup // counts active Serve goroutines for GracefulStop

			channelzID int64 // channelz unique identification number
			czData     *channelzData

			serverWorkerChannels []chan *serverWorkerData
		}
	*/
	// 将grpcserver注册到自行实现的server中
	/**
	1、首先判断server{}是否是实现了grpcserver
	2、然后将server添加到grpcserver的services中：
		a.构建serverinfo
				info := &serviceInfo{
					serviceImpl: ss,  这是自己实现的server实体
					methods:     make(map[string]*MethodDesc),
					streams:     make(map[string]*StreamDesc),
					mdata:       sd.Metadata,
				}
		b.将grpc.NewServer()创建的server中的method和stream加入到info中
		c.最后将 s.services[sd.ServiceName] = info
	*/
	proto.RegisterGreeterServer(s, &server{})
	listener, err := net.Listen("tcp", "0.0.0.0:5001")
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	// 交由grpc来处理listener
	err = s.Serve(listener)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}

}
