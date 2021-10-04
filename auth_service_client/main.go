/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	pb "flahmingo/grpc_service"

	"google.golang.org/grpc"
)

const (
	address       = "localhost:50052"
	defaultOption = "printHelp"
)

func readArg(position int) string {
	return os.Args[position]
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to the server : %v", err)
	}
	defer conn.Close()
	c := pb.NewAuthServiceClient(conn)

	option := readArg(1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if strings.EqualFold(option, "SignupWithPhoneNumber") {
		phoneNumber := readArg(2)
		name := readArg(3)
		r, err := c.SignupWithPhoneNumber(ctx, &pb.SignupWithPhoneNumberRequest{PhoneNumber: phoneNumber, Name: name})
		if err != nil {
			log.Fatalf("could not signup: %v", err)
		}
		log.Printf("Response: %s", r.GetMessage())
	} else if strings.EqualFold(option, "VerifyPhoneNumber") {
		phoneNumber := readArg(2)
		r, err := c.VerifyPhoneNumber(ctx, &pb.VerifyPhoneNumberRequest{PhoneNumber: phoneNumber})
		if err != nil {
			log.Fatalf("could not verify: %v", err)
		}
		log.Printf("Response: %s", r.GetMessage())
	} else if strings.EqualFold(option, "LoginWithPhoneNumber") {
		phoneNumber := readArg(2)
		r, err := c.LoginWithPhoneNumber(ctx, &pb.LoginWithPhoneNumberRequest{PhoneNumber: phoneNumber})
		if err != nil {
			log.Fatalf("could not login: %v", err)
		}
		log.Printf("Response: %s", r.GetMessage())
	} else if strings.EqualFold(option, "ValidatePhoneNumberLogin") {
		phoneNumber := readArg(2)
		otpCode := readArg(3)
		r, err := c.ValidatePhoneNumberLogin(ctx, &pb.ValidatePhoneNumberLoginRequest{PhoneNumber: phoneNumber, OtpCode: otpCode})
		if err != nil {
			log.Fatalf("could not validate: %v", err)
		}
		log.Printf("Response: %s", r.GetMessage())
	} else if strings.EqualFold(option, "GetProfile") {
		phoneNumber := readArg(2)
		r, err := c.GetProfile(ctx, &pb.GetProfileRequest{PhoneNumber: phoneNumber})
		if err != nil {
			log.Fatalf("could not getprofile: %v", err)
		}
		log.Printf("Response: %s", r.GetMessage())
	} else if strings.EqualFold(option, "Logout") {
		phoneNumber := readArg(2)
		r, err := c.Logout(ctx, &pb.LogoutRequest{PhoneNumber: phoneNumber})
		if err != nil {
			log.Fatalf("could not logout: %v", err)
		}
		log.Printf("Response: %s", r.GetMessage())
	} else {
		log.Printf("Use one of the following method for 1st arg: [SignupWithPhoneNumber,VerifyPhoneNumber,LoginWithPhoneNumber,ValidatePhoneNumberLogin,GetProfile,Logout] ")
		log.Printf("Mention phoneNumber for 2nd arg ")
		log.Printf("Mention name/otpcode for 3rd arg ")
	}
}
