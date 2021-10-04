// Package main implements a client for OTPService.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "flahmingo/grpc_service"

	"google.golang.org/grpc"
)

const (
	address            = "localhost:50051"
	defaultPhoneNumber = "+1DDDDDDDDDD"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewOTPServiceClient(conn)

	fmt.Println(os.Args)
	// Contact the server and print out its response.
	phoneNumber := defaultPhoneNumber
	if len(os.Args) > 1 {
		phoneNumber = string(os.Args[1])
		fmt.Printf("Received phoneNumber %s\n", string(phoneNumber))
	} else {
		panic("Missing phoneNumber in cmd line arg 1")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SendOTP(ctx, &pb.OTPRequest{PhoneNumber: phoneNumber})
	if err != nil {
		log.Fatalf("Failed to connect to the server : %v", err)
	}
	log.Printf("Response: %s", r.GetMessage())
}
