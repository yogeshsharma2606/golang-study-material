package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	greetv1 "github.com/golang-study/15-grpc/api/greet/v1"
	"google.golang.org/grpc"
)

type greeterServer struct {
	greetv1.UnimplementedGreeterServer
}

func (greeterServer) SayHello(_ context.Context, req *greetv1.HelloRequest) (*greetv1.HelloReply, error) {
	name := req.GetName()
	if name == "" {
		name = "world"
	}
	return &greetv1.HelloReply{Message: fmt.Sprintf("Hello, %s!", name)}, nil
}

func main() {
	addr := ":50051"
	if v := os.Getenv("GRPC_ADDR"); v != "" {
		addr = v
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	srv := grpc.NewServer()
	greetv1.RegisterGreeterServer(srv, greeterServer{})
	go func() {
		log.Printf("gRPC server on %s", addr)
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down gRPC server")
	srv.GracefulStop()
}
