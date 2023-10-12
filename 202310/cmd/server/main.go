package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	hellopb "mygrpc/pkg/grpc"
)

type myGreetingServer struct {
	hellopb.UnimplementedGreetingServiceServer
}

func (s *myGreetingServer) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func NewMyGreetingServer() *myGreetingServer {
	return &myGreetingServer{}
}

func main() {
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	hellopb.RegisterGreetingServiceServer(s, NewMyGreetingServer())
	reflection.Register(s)

	go func() {
		log.Printf("💨 Start gRPC server port:%+v\n", port)
		s.Serve(listener)
	}()

	q := make(chan os.Signal, 1)
	signal.Notify(q, os.Interrupt)
	<-q
	log.Println("👋 Stopping gRPC server")
	s.GracefulStop()
}
