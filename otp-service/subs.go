package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/pubsub"
)

func main() {
	pflag.String("subscriber", "verification", "name of subscriber")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	ctx := context.Background()
	proj := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if proj == "" {
		fmt.Fprintf(os.Stderr, "GOOGLE_CLOUD_PROJECT environment variable must be set.\n")
		os.Exit(1)
	}
	client, err := pubsub.NewClient(ctx, proj)
	if err != nil {
		log.Fatalf("Could not create pubsub Client: %v", err)
	}

	sub := viper.GetString("subscriber") // retrieve values from viper instead of pflag
	subExists := false

	fmt.Println("Listing all subscriptions from the project:")
	subs, err := list(client)
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range subs {
		fmt.Println("sub:", s, " s.string():", s.String())
		sa := strings.Split(s.String(), "/")
		salen := len(sa)
		if sa[salen-1] == sub {
			subExists = true
			break
		}
	}
	fmt.Println("saexists:", subExists)

	t := createTopicIfNotExists(client)

	// Create a new subscription.
	if !subExists {
		if err := create(client, sub, t); err != nil {
			log.Fatal(err)
		}
	}

	// Pull messages via the subscription.
	if err := pullMsgs(client, sub, t, false); err != nil {
		log.Fatal(err)
	}

	// Delete the subscription.
	if err := delete(client, sub); err != nil {
		log.Fatal(err)
	}

	// Print all the subscriptions in the project.

}

func list(client *pubsub.Client) ([]*pubsub.Subscription, error) {
	ctx := context.Background()
	// [START pubsub_list_subscriptions]
	var subs []*pubsub.Subscription
	it := client.Subscriptions(ctx)
	for {
		s, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	// [END pubsub_list_subscriptions]
	return subs, nil
}

func createTopicIfNotExists(c *pubsub.Client) *pubsub.Topic {
	ctx := context.Background()

	const topic = "verification"

	// Create a topic to subscribe to.
	t := c.Topic(topic)
	ok, err := t.Exists(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		return t
	}

	t, err = c.CreateTopic(ctx, topic)
	if err != nil {
		log.Fatalf("Failed to create the topic: %v", err)
	}
	return t
}

func create(client *pubsub.Client, name string, topic *pubsub.Topic) error {
	ctx := context.Background()
	// [START pubsub_create_pull_subscription]
	sub, err := client.CreateSubscription(ctx, name, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 20 * time.Second,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Created subscription: %v\n", sub)
	// [END pubsub_create_pull_subscription]
	return nil
}

func pullMsgs(client *pubsub.Client, name string, topic *pubsub.Topic, testPublish bool) error {
	ctx := context.Background()

	if testPublish {
		// Publish 10 messages on the topic.
		var results []*pubsub.PublishResult
		for i := 0; i < 10; i++ {
			res := topic.Publish(ctx, &pubsub.Message{
				Data: []byte(fmt.Sprintf("hello world #%d", i)),
			})
			results = append(results, res)
		}

		// Check that all messages were published.
		for _, r := range results {
			_, err := r.Get(ctx)
			if err != nil {
				return err
			}
		}
	}

	var mu sync.Mutex
	received := 0
	sub := client.Subscription(name)
	cctx, cancel := context.WithCancel(ctx)
	err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		fmt.Printf("Got message: %q\n", string(msg.Data))
		mu.Lock()
		defer mu.Unlock()
		received++
		if received == 10 {
			cancel()
		}
	})
	if err != nil {
		return err
	}

	return nil
}

func delete(client *pubsub.Client, name string) error {
	ctx := context.Background()
	// [START pubsub_delete_subscription]
	sub := client.Subscription(name)
	if err := sub.Delete(ctx); err != nil {
		return err
	}
	fmt.Println("Subscription deleted.")
	// [END pubsub_delete_subscription]
	return nil
}
