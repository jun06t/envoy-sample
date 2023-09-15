package main

import (
	"context"
	"log"
	"time"

	pb "github.com/jun06t/grpc-sample/unary/proto"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Endpoint string `split_words:"true"`
}

var conf Config

func init() {
	envconfig.MustProcess("MYAPP", &conf)
}

func main() {
	conn, err := grpc.Dial(conf.Endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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
	resp, err := c.SayHello(context.TODO(), req)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Reply: ", resp.Message)
}
