package pubsub

import (
	"context"
	"fmt"
	"grpc-example/src/proto"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	protoLib "google.golang.org/protobuf/proto"
)

type OTPSubscriber struct {
}

func (otp *OTPSubscriber) Listen(channel chan *proto.OTPValidationRequest) {
	ctx := context.Background()
	initialize(ctx)

	log.Print(os.Getenv("PUBSUB_SUBSCRIPTION"))
	sub := client.Subscription(os.Getenv("PUBSUB_SUBSCRIPTION"))

	cctx, cancel := context.WithCancel(ctx)
	count := 0
	err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		log.Printf("xx")
		msg.Ack()
		fmt.Printf("Got message: %q\n", string(msg.Data))
		message := &proto.OTPValidationRequest{}
		protoLib.Unmarshal(msg.Data, message)
		channel <- message
		count += 1
		if count > 100 {
			cancel()
			otp.Listen(channel)
		}
	})
	if err != nil {
		log.Fatal("Error subscribing: ", err)
	}
}
