package main

import (
	"log"
	"net"

	"github.com/ViniciusSevero/fc2-grpc/pb"
	"github.com/ViniciusSevero/fc2-grpc/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Couldn't connect %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &services.UserServiceImpl{})
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Couldn't serve %v", err)
	}
}
