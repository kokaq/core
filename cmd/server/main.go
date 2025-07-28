package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	// pb "github.com/kokaq/core/proto"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	// pb.RegisterNamespaceServiceServer(grpcServer, NewNamespaceServer())
	log.Println("gRPC server started on :50051")
	grpcServer.Serve(lis)
}
