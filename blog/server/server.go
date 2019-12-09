package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/hrfmmr/grpc-go-sandbox/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var collection *mongo.Collection

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorId string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

type server struct{}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.Blog
	data := &blogItem{
		AuthorId: blog.AuthorId,
		Title:    blog.Title,
		Content:  blog.Content,
	}
	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		log.Fatal(err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error:%v", err))
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot convert to ObjectID"))
	}
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.AuthorId,
			Title:    blog.Title,
			Content:  blog.Content,
		},
	}, nil
}

func main() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// connect to MongoDB
	fmt.Println("🔗Connecting to MongoDB")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		fmt.Println("Closing MongoDB connection")
		client.Disconnect(context.TODO())
	}()
	collection = client.Database("grpctest").Collection("blog")

	fmt.Println("🚀Blog Service Started")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	defer func() {
		fmt.Println("Closing the listener")
		lis.Close()
	}()
	if err != nil {
		log.Fatal(err)
		return
	}
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	defer func() {
		fmt.Println("Stop the server")
		s.Stop()
	}()
	blogpb.RegisterBlogServiceServer(s, &server{})
	go func() {
		fmt.Println("Starting the server")
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
			return
		}
	}()
	// Wait for Ctl-C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	// Block until a signal is received
	<-ch
}
