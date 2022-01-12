package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
)

func main() {
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

	const topic = "verification"
	// Create a new topic called my-topic.
	// if err := create(client, topic); err != nil {
	// 	log.Fatalf("Failed to create a topic: %v", err)
	// }

	// List all the topics from the project.
	fmt.Println("Listing all topics from the project:")
	topics, err := list(client)
	if err != nil {
		log.Fatalf("Failed to list topics: %v", err)
	}
	for _, t := range topics {
		fmt.Println(t)
	}

	// Publish a text message on the created topic.
	if err := publish(client, topic, "sendOTP"); err != nil {
		log.Fatalf("Failed to publish: %v", err)
	}

	// // Delete the topic.
	// if err := delete(client, topic); err != nil {
	// 	log.Fatalf("Failed to delete the topic: %v", err)
	// }
}
func list(client *pubsub.Client) ([]*pubsub.Topic, error) {
	ctx := context.Background()

	// [START pubsub_list_topics]
	var topics []*pubsub.Topic

	it := client.Topics(ctx)
	for {
		topic, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}

	return topics, nil
	// [END pubsub_list_topics]
}
func create(client *pubsub.Client, topic string) error {
	ctx := context.Background()
	// [START pubsub_create_topic]
	t, err := client.CreateTopic(ctx, topic)
	if err != nil {
		return err
	}
	fmt.Printf("Topic created: %v\n", t)
	// [END pubsub_create_topic]
	return nil
}

func publish(client *pubsub.Client, topic, msg string) error {
	ctx := context.Background()
	// [START pubsub_publish]
	// [START pubsub_quickstart_publisher]
	t := client.Topic(topic)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Published a message; msg ID: %v\n", id)
	// [END pubsub_publish]
	// [END pubsub_quickstart_publisher]
	return nil
}

func delete(client *pubsub.Client, topic string) error {
	ctx := context.Background()
	// [START pubsub_delete_topic]
	t := client.Topic(topic)
	if err := t.Delete(ctx); err != nil {
		return err
	}
	fmt.Printf("Deleted topic: %v\n", t)
	// [END pubsub_delete_topic]
	return nil
}
