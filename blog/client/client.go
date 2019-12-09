package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hrfmmr/grpc-go-sandbox/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()
	c := blogpb.NewBlogServiceClient(cc)
	// create blog
	blog := &blogpb.Blog{
		AuthorId: "JohnDoe",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}
	fmt.Println("Create the Blog")
	rsp, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("✔Created blog:%+v\n", rsp.Blog)
}
