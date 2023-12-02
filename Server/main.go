package main

import (
	"io"
	"log"
	"net"
	"sync"

	pb "github.com/aungkoko/livescore-server/pb"
	"github.com/google/uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type liveScoreServer struct {
	mu      sync.Mutex
	clients map[string]pb.LiveScore_StreamMatchEventsServer
	events  []*pb.MatchEvent
	pb.UnimplementedLiveScoreServer
}

func (s *liveScoreServer) StreamMatchEvents(stream pb.LiveScore_StreamMatchEventsServer) error {
	s.mu.Lock()
	clientID := uuid.New().String()
	s.clients[clientID] = stream
	s.mu.Unlock()

	for _, event := range s.events {
		if err := stream.Send(event); err != nil {
			log.Printf("Error sending event to client: %v", err)
			return err
		}
	}

	for {
		event, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error receiving event from client: %v", err)
			break
		}

		s.mu.Lock()
		s.events = append(s.events, event)
		s.mu.Unlock()

		s.broadcast(event, clientID)
	}

	s.mu.Lock()
	delete(s.clients, clientID)
	s.mu.Unlock()

	log.Printf("Client %s disconnected", clientID)

	return nil
}

func (s *liveScoreServer) broadcast(event *pb.MatchEvent, senderClientID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for clientID, client := range s.clients {
		if clientID != senderClientID {
			if err := client.Send(event); err != nil {
				log.Printf("Error broadcasting event to client: %v", err)
			}
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterLiveScoreServer(server, &liveScoreServer{clients: make(map[string]pb.LiveScore_StreamMatchEventsServer)})
	reflection.Register(server)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
