package main

import (
	"context"
	"golang_learn/golang_learn/models/grpc/proto"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloResponse, error) {
	// 将状态码装换为error，返回给客户端
	return nil, status.Errorf(codes.NotFound, "记录未找到：%s", request.Name)
	//return &proto.HelloReply{
	//	Message: "hello "+request.Name,
	//}, nil
}

func main() {
	g := grpc.NewServer()
	proto.RegisterGreeterServer(g, &Server{})
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	err = g.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
