package main

import (
	"fmt"
	"log"
	"net"

	pb "go.quinn.io/dataq/rpc"
	"google.golang.org/grpc"
)

func main() {
	plugin := New()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	s := grpc.NewServer()
	pb.RegisterDataQPluginServer(s, NewServer(plugin))
	
	fmt.Printf("Gmail plugin server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
