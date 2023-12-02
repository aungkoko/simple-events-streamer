# simple-events-streamer

After exploring about gRPC and HTTP2, I would like to create my own project and want to get some exposure with gRPC and go.
So, I created this project using gRPC bi-directional streaming.

# project's idea
Imagine a livescore app, when a match event triggers, the users want to know the match stats simultaneously. eg. Player A scores at 35th mins, Player C fouls player B.
So, for thise case, we need a sender who sends the match event to server.
Then server needs to stream the match events to all clients who connect to server except the one who sends the match event.


# architecture overview
I design 3 components, 
* Server -> receive stream events from the sender and send to the clients
* SenderClient -> accepts the match event via REST api and call gRPC bidirectional stream method from server and stream
* ReceiverClient -> receive stream data from gRPC server and displays(currently logging the event data)

# prerequisites
* need to install go

# how to test
* go to sever and run ```go run main.go```, it will serve the gRPC server to accept and send steram data
* go to sender client and run ```go run main.go```, it will serve the simple http server on port 8080 and accepts this POST request
```
url http://127.0.0.1:8080/events
payload {
	"event_id": "Hola",
	"event_type": "Type",
	"description": "ak kicks the ball"
}```

