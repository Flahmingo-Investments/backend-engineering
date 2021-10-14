package main

import (
	"flag"
	authService "grpc-example/src/apps/auth/services"
	commons "grpc-example/src/commons/utils"
	"grpc-example/src/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	go commons.DatabaseConnect()

	flag.Parse()
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("Failed to listen on port 5000 %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	proto.RegisterAuthServiceServer(grpcServer, authService.Initialize())

	log.Print("Listening on 5000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve on port 5000 %v", err)
	}
}
