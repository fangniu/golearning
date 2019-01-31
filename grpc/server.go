package main

import (
	"log"
	"net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/golearning/grpc/hello"
	)

const (
	port = ":50051"
)

type server struct {}

func (s *server) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloResponse, error) {
	return &hello.HelloResponse{Reply: "Hello " + in.Greeting, Number: []int32{1,2,3}}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}
	s := grpc.NewServer()
	hello.RegisterHelloServiceServer(s, &server{})
	s.Serve(lis)
}

