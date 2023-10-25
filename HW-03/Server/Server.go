package main

import (
	proto "ChittyChat/grpc"
	"context"
	"flag"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedUsermanagementServer // Necessary
	name                                    string
	port                                    int
}

var port = flag.Int("port", 0, "server port number")

func main() {
	// Get the port from the command line when the server is run
	flag.Parse()

	// Create a server struct
	server := &Server{
		name: "serverName",
		port: *port,
	}

	// Start the server
	go startServer(server)

	// Keep the server running until it is manually quit
	for {

	}
}
func startServer(server *Server) {

	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.port))

	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	log.Printf("Started server at port: %d\n", server.port)

	// Register the grpc server and serve its listener
	proto.RegisterUsermanagementServer(grpcServer, server)
	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

// func (c *Server) AskForTime(ctx context.Context, in *proto.AskForTimeMessage) (*proto.TimeMessage, error) {
// 	log.Printf("Client asking with ID: %d is asing for time", in.ClientId)
// 	return &proto.TimeMessage{Time: time.Now().String()}, nil
// }

func (c *Server) ClientJoin(ctx context.Context, in *proto.NewClient) (*proto.Client, error) {
	var _id int64 = 1
	log.Printf("Client %s joined and got assigned ID %d", in.Name, _id)
	return &proto.Client{Name: in.Name, Id: _id}, nil
}

func (c *Server) BroadcastMessage(ctx context.Context, incomming *proto.PublishMessage) (*proto.Broadcast, error) {
	var broadcastString string = "User " + strconv.Itoa(int(incomming.ClientId)) + " said " + incomming.Message + " at time Lamport time: 4"
	log.Print(broadcastString)
	//return &proto.Broadcast{message: broadcastString}, nil
	return nil, nil
}
