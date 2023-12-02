package cmd

import (
	"context"
	"log"

	"github.com/aungkoko/livescore-admin-client/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Dispatch(event *pb.MatchEvent) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewLiveScoreClient(conn)

	stream, err := client.StreamMatchEvents(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	done := make(chan struct{})

	go sendMatchEvent(stream, client, event)

	<-done
}

func sendMatchEvent(stream pb.LiveScore_StreamMatchEventsClient, client pb.LiveScoreClient, event *pb.MatchEvent) {

	if err := stream.Send(event); err != nil {
		log.Fatalf("Error sending data: %v", err)
	}

	if err := stream.CloseSend(); err != nil {
		log.Fatalf("Error closing stream: %v", err)
	}
}
