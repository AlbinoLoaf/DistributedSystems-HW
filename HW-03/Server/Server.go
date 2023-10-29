package main

import (
	proto "ChittyChat/grpc"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedUsermanagementServer // Necessary
	name                                    string
	port                                    int
	capacity                                []bool
	log                                     []string
	Timestamp                               int64
}

var port = flag.Int("port", 0, "server port number")

func main() {
	// Get the port from the command line when the server is run

	flag.Parse()

	// Create a server struct
	server := &Server{
		name: "serverName",
		port: *port,
		// declare server capacity
		capacity:  make([]bool, 5),
		log:       make([]string, 1),
		Timestamp: 1,
	}
	server.log[0] = "start of the server"
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

func GenerateId(s *Server) (id int) {
	var num int
	for i, v := range s.capacity {
		if v == false {
			//log.Printf("Value = %v Index = %d", v, i)
			num = i
			s.capacity[i] = true
			//log.Printf("Value = %v Index = %d", v, i)
			break
		}

	}
	return num + 1
}

func deleteID(s *Server, id int64) {
	if s.capacity[id-1] {
		s.capacity[id-1] = false
	} else {
		log.Fatalf("The wrong ID recieved")
	}
}

func (c *Server) LeaveClient(ctx context.Context, in *proto.Client) (*proto.Confirmation, error) {
	deleteID(c, in.Id)
	if c.capacity[in.Id-1] {
		log.Print("Couldn't delete Client")
		return nil, nil
	} else {
		event := fmt.Sprintf("%s left the server what a looser", in.Name)
		c.LogEvent(event)
		return &proto.Confirmation{Accept: true}, nil
	}

}
func (c *Server) SendMessage(ctx context.Context, in *proto.PublishMessage) (*proto.Timestamp, error) {
	Message := fmt.Sprintf("%s was set at %d by %s", in.Message, c.Timestamp, in.Name)
	c.LogEvent(Message)
	return &proto.Timestamp{Time: c.Timestamp}, nil
}

func (c *Server) RequestChange(ctx context.Context, in *proto.Timestamp) (*proto.Timestamp, error) {
	return &proto.Timestamp{Time: c.Timestamp}, nil
}

func (c *Server) LogEvent(Event string) {
	log.Printf("Lamport time: %d", c.Timestamp)
	c.log = append(c.log, Event)
	c.Timestamp++
}

func (c *Server) ClientJoin(ctx context.Context, in *proto.NewClient) (*proto.Client, error) {
	var _id int64 = int64(GenerateId(c))
	event := fmt.Sprintf("A new person have joined say hi to %s who joined at lamport time %d", in.Name, (c.Timestamp)) // Without ID
	log.Printf("Client %s joined and got assigned ID %d", in.Name, _id)
	c.LogEvent(event)
	return &proto.Client{Name: in.Name, Id: _id, Timestamp: c.Timestamp - 1}, nil
}

func (c *Server) RequestBroadcast(ctx context.Context, in *proto.Timestamp) (*proto.Broadcast, error) {
	//return &proto.Broadcast{message: broadcastString}, nil
	//log.Printf("Server time %d ---- incomming time: %d \n log length %d", c.Timestamp, in.Time, len(c.log))
	return &proto.Broadcast{Message: c.log[in.Time]}, nil
}
