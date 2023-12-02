package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aungkoko/livescore-admin-client/cmd"
	"github.com/aungkoko/livescore-admin-client/pb"
)

type ResponseData struct {
	Message string `json:"message"`
}

type Event struct {
	EventID     string `json:"event_id"`
	EventType   string `json:"event_type"`
	Description string `json:"description"`
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not allowed method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var data Event
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	event := &pb.MatchEvent{EventId: data.EventID, EventType: data.EventType, Description: data.Description}

	response := ResponseData{
		Message: "Your request is accepted",
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	go func() {
		cmd.Dispatch(event)
		if err != nil {
			log.Printf("Error making gRPC call: %v", err)
		}
	}()
}

func main() {
	http.HandleFunc("/events", sendMessage)

	fmt.Println("Server listening on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("failed to start http server")
	}
	fmt.Println("Server listening on :8080")
}
