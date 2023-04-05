package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/jun06t/envoy-sample/keepalive/proto"
	"google.golang.org/grpc"
)

func main() {
	addr := os.Getenv("ENDPOINT")
	conn, err := grpc.Dial(addr,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(
			grpc.WaitForReady(true),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	for {
		hello(c)
		time.Sleep(1 * time.Second)
	}
}

func hello(c pb.GreeterClient) {
	req := &pb.HelloRequest{
		Name: "alice",
		Age:  10,
		Man:  true,
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	//resp, err := c.SayHello(ctx, req, grpc.WaitForReady(true))
	resp, err := c.SayHello(ctx, req)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Reply: ", resp.Message)
}
