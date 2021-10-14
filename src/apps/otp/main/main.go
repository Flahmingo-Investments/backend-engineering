package main

import (
	"context"
	"grpc-example/src/apps/otp/services"
	"grpc-example/src/commons/services/pubsub"
	"grpc-example/src/proto"
	"log"
)

func main() {
	subscriber := &pubsub.OTPSubscriber{}
	otpService := &services.OTPService{}

	channel := make(chan *proto.OTPValidationRequest)
	go subscriber.Listen(channel)

	log.Print("Listening pubsub...")
	ctx := context.Background()
	for {
		message := <-channel
		response, _ := otpService.SendOTP(ctx, message)
		if response.Status == proto.Status_SUCCESS {
			log.Print("Sent from Twilio " + message.Phone + ": " + message.Otp)
		} else {
			log.Print("Message not sent")
		}
	}
}
