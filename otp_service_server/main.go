// Package main implements a server for OTPService.
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"flahmingo/constants"
	pb "flahmingo/grpc_service"
	pubsub "flahmingo/pubsub/subscriptions"
	twilio "flahmingo/sms_utils"
	cryptUtils "flahmingo/utils/crypt"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement flahmingo.grpc.service.OTPService.
type server struct {
	pb.UnimplementedOTPServiceServer
}

// SendOTP implements flahmingo.grpc.service.OTPService.SendOTP
func (s *server) SendOTP(ctx context.Context, in *pb.OTPRequest) (*pb.OTPResponse, error) {
	phoneNumber := in.GetPhoneNumber()
	log.Printf("Received call for SendOTP : %v", phoneNumber)
	// Actual Invocation of logic
	randomOtpCode := cryptUtils.GenerateOTPMsg(constants.OTP_MIN_INT, constants.OTP_MAX_INT)
	message := constants.VERIFICATION_INSTR_MSG_STRING
	otpMsg := fmt.Sprintf("%s %s", message, randomOtpCode)
	twilio.SendOTP(phoneNumber, otpMsg)
	fmt.Printf("Sent sms to phoneNumber: %s with message: %s\n", phoneNumber, otpMsg)
	return &pb.OTPResponse{Message: "Hello " + in.GetPhoneNumber()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOTPServiceServer(s, &server{})
	log.Printf("otp service server listening at %v", lis.Addr())
	log.Println("Starting google pubsub topic subscription listener")
	// Start subscription handler
	pubsub.PullMsgs(constants.PUBSUB_PROJECT_ID, constants.PUBSUB_SUBSCRIPTION_ID)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
