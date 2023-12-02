package main

import (
	"context"
	"log"

	"github.com/aungkoko/livescore-receiver-client/pb" // Update with your actual path

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func receiveStreamData(client pb.LiveScoreClient) error {
	stream, err := client.StreamMatchEvents(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	// Receive and print the stramed messages from the server
	for {
		data, err := stream.Recv()
		if err == nil {
			log.Printf("Received data from server: %s",
				"EventID "+data.GetEventId()+" Type "+
					data.EventType+" Description "+data.GetDescription())
		} else {
			break
		}
	}

	return nil
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewLiveScoreClient(conn)

	if err := receiveStreamData(client); err != nil {
		log.Fatalf("Error receiving data: %v", err)
	}
}
