package pubsub

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"

	twilio "flahmingo/sms_utils"
	messageUtils "flahmingo/utils/msg"
)

/* To do
1) Refactor pullMsgs to accept a function in args as handler for msgs received.
2) It is possible for Messages to be redelivered, even if Message.Ack has been called.
Client code must be robust to multiple deliveries of messages.
*/

func PullMsgs(projectID, subscriptionID string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	// Consume messages.
	var mu sync.Mutex
	//received := 0
	sub := client.Subscription(subscriptionID)
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		// This is the handler to the subscribed topic
		// ---------------------------------------------
		fmt.Printf("Got message: %q\n", string(msg.Data))
		phoneNumber := messageUtils.GetPhoneNumber(string(msg.Data))
		randomOtpCode := messageUtils.GetOtpCode(string(msg.Data)) 
		message := messageUtils.GetInstructionMessage(string(msg.Data))
		otpMsg := fmt.Sprintf("%s %s", message, randomOtpCode)
		twilio.SendOTP(phoneNumber, otpMsg)
		fmt.Printf("Sent sms %s: %s\n", phoneNumber, otpMsg)
		// ---------------------------------------------
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("Receive: %v", err)
	}
	return nil
}

