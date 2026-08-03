package main

import (
	"context"
	"log"
	"os"
	"time"

	greetv1 "github.com/golang-study/15-grpc/api/greet/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	target := "localhost:50051"
	if v := os.Getenv("GRPC_TARGET"); v != "" {
		target = v
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	name := "Gopher"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	client := greetv1.NewGreeterClient(conn)
	resp, err := client.SayHello(ctx, &greetv1.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("SayHello: %v", err)
	}
	log.Printf("response: %s", resp.GetMessage())
}
