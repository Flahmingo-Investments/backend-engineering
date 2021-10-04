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

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	constants "flahmingo/constants"
	pb "flahmingo/grpc_service"
	pubsub "flahmingo/pubsub/topics"
	cryptUtils "flahmingo/utils/crypt"
	db "flahmingo/utils/db"
	message "flahmingo/utils/msg"

	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

// *******************************************
// GRPC method Implementations for AuthService
// *******************************************

// server is used to implement flahmingo.grpc.service.AuthService
type server struct {
	pb.UnimplementedAuthServiceServer
}

// SignupWithPhoneNumber implements flahmingo.grpc.service.AuthService.SignupWithPhoneNumber
func (s *server) SignupWithPhoneNumber(ctx context.Context, in *pb.SignupWithPhoneNumberRequest) (*pb.SignupWithPhoneNumberResponse, error) {
	name := in.GetName()
	phoneNumber := in.GetPhoneNumber()
	log.Printf("Received call for SignupWithPhoneNumber : %s %s", name, phoneNumber)
	// Actual Invocation of logic
	response := signupWithPhoneNumber(phoneNumber, name)
	fmt.Println(response)
	// To do check -> In case of error blank message being passed back
	return &pb.SignupWithPhoneNumberResponse{Message: response}, nil
}

// VerifyPhoneNumber implements flahmingo.grpc.service.AuthService.
func (s *server) VerifyPhoneNumber(ctx context.Context, in *pb.VerifyPhoneNumberRequest) (*pb.VerifyPhoneNumberResponse, error) {
	phoneNumber := in.GetPhoneNumber()
	log.Printf("Received call for VerifyPhoneNumber : %s", phoneNumber)
	// Actual Invocation of logic
	response := strconv.FormatBool(verifyPhoneNumber(phoneNumber))
	return &pb.VerifyPhoneNumberResponse{Message: response}, nil
}

// LoginWithPhoneNumber implements flahmingo.grpc.service.AuthService.LoginWithPhoneNumber
func (s *server) LoginWithPhoneNumber(ctx context.Context, in *pb.LoginWithPhoneNumberRequest) (*pb.LoginWithPhoneNumberResponse, error) {
	phoneNumber := in.GetPhoneNumber()
	log.Printf("Received call for LoginWithPhoneNumber : %s", phoneNumber)
	// Actual Invocation of logic

	response := strconv.FormatBool(loginWithPhoneNumber(phoneNumber))
	fmt.Println(response)
	return &pb.LoginWithPhoneNumberResponse{Message: response}, nil
}

// ValidatePhoneNumberLogin implements flahmingo.grpc.service.AuthService.ValidatePhoneNumberLogin
func (s *server) ValidatePhoneNumberLogin(ctx context.Context, in *pb.ValidatePhoneNumberLoginRequest) (*pb.ValidatePhoneNumberLoginResponse, error) {
	phoneNumber := in.GetPhoneNumber()
	otpCode := in.GetOtpCode()
	log.Printf("Received call for ValidatePhoneNumberLogin : %s %s", phoneNumber, otpCode)
	// Actual Invocation of logic
	response := strconv.FormatBool(validatePhoneNumberLogin(phoneNumber, otpCode))
	fmt.Println(response)
	return &pb.ValidatePhoneNumberLoginResponse{Message: response}, nil
}

// GetProfile implements flahmingo.grpc.service.AuthService.GetProfile
func (s *server) GetProfile(ctx context.Context, in *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	var responseMsg string
	phoneNumber := in.GetPhoneNumber()
	log.Printf("Received call for GetProfile : %s", phoneNumber)
	// Actual Invocation of logic
	var user db.User
	user = getProfile(phoneNumber)
	db.PrintAll()
	if user.PhoneNumber == "" {
		responseMsg = "User Not Registered"
	} else {
		responseMsg = fmt.Sprintf("User details : {Name : %s Phone : %s}", user.Name, user.PhoneNumber)
	}
	fmt.Println(responseMsg)
	return &pb.GetProfileResponse{Message: responseMsg}, nil
}

// Logout implements flahmingo.grpc.service.AuthService.Logout
func (s *server) Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	phoneNumber := in.GetPhoneNumber()
	log.Printf("Received call for Logout : %s", phoneNumber)
	// Actual Invocation of logic
	response := strconv.FormatBool(logout(phoneNumber))
	fmt.Println(response)
	return &pb.LogoutResponse{Message: response}, nil
}

// *******************************************
// *******************************************

func signupWithPhoneNumber(phoneNumber, name string) string {
	response := "Phone already registered, use another !"
	defer func() {
		if err := recover(); err != nil {
			response = "Phone already registered, use another !"
			log.Println("panic occurred:", err)
		}
	}()
	db.RegisterUser(name, phoneNumber)
	response = "OK"
	return response
}
func verifyPhoneNumber(phoneNumber string) bool {
	// Use msg format to indicate "OTP for Phone Verification"
	randomOtpCode := cryptUtils.GenerateOTPMsg(constants.OTP_MIN_INT, constants.OTP_MAX_INT)
	pubsub.Publish(constants.PUBSUB_PROJECT_ID, constants.PUBSUB_TOPIC, message.CreateMsg(phoneNumber, randomOtpCode, constants.VERIFICATION_INSTR_MSG_STRING))
	return true
}
func loginWithPhoneNumber(phoneNumber string) bool {
	// Use msg format to indicate "OTP for Login"
	randomOtpCode := cryptUtils.GenerateOTPMsg(constants.OTP_MIN_INT, constants.OTP_MAX_INT)
	db.UpdateLoginOTP(phoneNumber, randomOtpCode)
	pubsub.Publish(constants.PUBSUB_PROJECT_ID, constants.PUBSUB_TOPIC, message.CreateMsg(phoneNumber, randomOtpCode, constants.LOGIN_INSTR_MSG_STRING))
	return true
}
func validatePhoneNumberLogin(phoneNumber, otpCode string) bool {
	return db.LoginUsingOTP(phoneNumber, otpCode)
}

func getProfile(phoneNumber string) db.User {
	return db.GetUserDetails(phoneNumber)
}

func logout(phoneNumber string) bool {
	return db.Logout(phoneNumber)
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &server{})
	log.Printf("auth service server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
