package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/jun06t/grpc-sample/unary/proto"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
)

type Config struct {
	Addr      string        `split_words:"true" default:":8001"`
	SlowStart time.Duration `split_words:"true" default:"10s"`
}

var conf Config

func init() {
	envconfig.MustProcess("MYAPP", &conf)
}

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Println(in.String())
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	time.Sleep(conf.SlowStart)
	lis, err := net.Listen("tcp", conf.Addr)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	err = s.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
