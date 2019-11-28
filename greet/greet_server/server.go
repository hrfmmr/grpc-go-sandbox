package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/hrfmmr/grpc-go-sandbox/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	rsp := &greetpb.GreetResponse{
		Result: result,
	}
	return rsp, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoke with req:%+v\n", req)
	firstName := req.GetGreeting().FirstName
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number:" + strconv.Itoa(i)
		rsp := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(rsp)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})
	fmt.Println("Listening greeting request...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve:%v", err)
	}
}
