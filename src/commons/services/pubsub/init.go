package pubsub

import (
	"context"
	"grpc-example/src/proto"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

var (
	topic  *pubsub.Topic
	client *pubsub.Client
)

type OTPPubSub interface {
	SendOTP(*proto.OTPValidationRequest) bool
}

func initialize(ctx context.Context) {
	if topic == nil {
		var err error
		options := option.WithCredentialsJSON([]byte(os.Getenv("GOOGLE_CLOUD_CREDENTIALS")))
		client, err = pubsub.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"), options)
		if err != nil {
			log.Fatal(err)
		}
		topic = client.Topic(os.Getenv("PUBSUB_TOPIC"))
		exists, err := topic.Exists(ctx)
		if err != nil {
			log.Fatal(err)
		}
		if !exists {
			log.Fatal(err)
		}
	}
}
