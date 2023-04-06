package main

import (
	"context"
	"log"
	"time"

	pb "github.com/jun06t/grpc-sample/unary/proto"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Config struct {
	KeepaliveTime    time.Duration `split_words:"true" default:"10s"`
	KeepaliveTimeout time.Duration `split_words:"true" default:"5s"`
	RpcTimeout       time.Duration `split_words:"true" default:"60s"`
	Endpoint         string        `split_words:"true"`
}

var conf Config

func init() {
	envconfig.MustProcess("MYAPP", &conf)
}

func main() {
	kacp := keepalive.ClientParameters{
		Time:                conf.KeepaliveTime,
		Timeout:             conf.KeepaliveTime,
		PermitWithoutStream: true,
	}
	conn, err := grpc.Dial(conf.Endpoint,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(kacp),
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
	ctx, cancel := context.WithTimeout(context.TODO(), conf.RpcTimeout)
	defer cancel()
	//resp, err := c.SayHello(ctx, req, grpc.WaitForReady(true))
	resp, err := c.SayHello(ctx, req)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Reply: ", resp.Message)
}
