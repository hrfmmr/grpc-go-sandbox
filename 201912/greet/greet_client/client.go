package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/hrfmmr/grpc-go-sandbox/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello, I'm a client")

	certFile := "ssl/ca.crt"
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		log.Fatal(err)
	}
	opts := grpc.WithTransportCredentials(creds)
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect:%v", err)
	}
	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)
	doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDiStreaming(c)
	// doUnaryWithDeadline(c, "Alice", 5*time.Second) // should complete
	// doUnaryWithDeadline(c, "Bob", 1*time.Second)   // should timeout
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "Doe",
		},
	}
	rsp, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC:%v", err)
	}
	log.Printf("Response from Greet:%v", rsp.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "Doe",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Response:%v", msg.Result)
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	reqs := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
				LastName:  "Doe",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Alice",
				LastName:  "",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, req := range reqs {
		fmt.Printf("Sending req:%+v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	rsp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("LongGreet rsp %+v\n", rsp)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting BiDi streaming RPC...")
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	reqs := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
				LastName:  "Doe",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Taro",
				LastName:  "Yamada",
			},
		},
	}

	waitc := make(chan struct{})

	// send a bunch of messages
	go func() {
		for _, req := range reqs {
			fmt.Printf("Sending req = %+v\n", req)
			if err := stream.Send(req); err != nil {
				log.Fatal(err)
				close(waitc)
			}
			time.Sleep(500 * time.Millisecond)
		}
		if err := stream.CloseSend(); err != nil {
			log.Fatal(err)
			close(waitc)
		}
	}()

	// recv a bunch of messages
	go func() {
		for {
			rsp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
				break
			}
			fmt.Printf("Received rsp = %+v\n", rsp)
		}
		close(waitc)
	}()
	fmt.Println("👀 waiting stream finished...")
	<-waitc
	fmt.Println("✔Done")
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, firstName string, timeout time.Duration) {
	log.Printf("👉Starting unary with deadline RPC... name:%v\n", firstName)
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: firstName,
			LastName:  "Doe",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer func() {
		log.Printf("cancel name:%+v\n", firstName)
		cancel()
	}()
	rsp, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Printf("⏰Timeout!! name:%v\n", firstName)
			} else {
				log.Printf("❗Unexpected error:%v\n", statusErr)
			}
		} else {
			log.Fatal(err)
		}
		return
	}
	log.Printf("✅rsp = %+v\n", rsp)
}
