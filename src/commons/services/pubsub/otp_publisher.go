package pubsub

import (
	"context"
	"grpc-example/src/proto"
	"log"

	"cloud.google.com/go/pubsub"
	protoLib "google.golang.org/protobuf/proto"
)

type OTPPublisher struct {
}

func (otp *OTPPublisher) SendOTP(message *proto.OTPValidationRequest) bool {
	ctx := context.Background()
	initialize(ctx)
	data, _ := protoLib.Marshal(message)
	msg := &pubsub.Message{
		Data: data,
	}
	serverID, err := topic.Publish(ctx, msg).Get(ctx)
	log.Print(serverID)
	if err != nil {
		log.Fatal(err)
	}
	return err == nil
}
