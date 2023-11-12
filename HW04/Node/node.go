package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"time"

	proto "hw04/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var idf = flag.Int64("id", 0, "server Id")
var Port = flag.String("port", "5400", "Own Tcp server port")
var Address = flag.String("address", "localhost", "Tcp server address")

type Server struct {
	proto.UnimplementedNodeManagementServer // You need this line if you have a server
	id                                      int64
	Timestamp                               int64
	usingCriticalSection                    bool
}

func main() {
	flag.Parse()
	fmt.Println("Starting server...")
	// log.SetOutput(os.Stdout)
	server := &Server{
		id:                   *idf,
		Timestamp:            1,
		usingCriticalSection: false,
	}
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", *Address+":"+*Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a proto server object
	s := grpc.NewServer()
	fmt.Printf("Server type: %T\n", s)
	// Attach the Node service to the server
	proto.RegisterNodeManagementServer(s, server)

	// start sending requests to peer
	list := GetPorts()
	i := 0
	for i < len(list) {
		if !strings.Contains(list[i], *Port) {
			go runClient(list[i])

		}
		i++
	}
	// Serve proto server
	fmt.Println("Serving proto on port " + *Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	// code here is unreachable because protoServer.Serve occupies the current thread.
}

func (s *Server) RequestReply(ctx context.Context, in *proto.Request) (*proto.Reply, error) {
	// If true this node accepts the requests replying go ahead
	replyMSG := fmt.Sprintf("Saying hi from %d", s.id)
	return &proto.Reply{Message: replyMSG, Id: s.id}, nil
}
func (s *Server) ResolveRequest(ForeignTimestamp int64, foreignID int64) (Accaptance bool) {
	fmt.Printf("Server %d got contacted by server %d\n", s.id, foreignID)
	for s.usingCriticalSection {
		time.Sleep(time.Millisecond * 500)
	}
	return true

}

func GetPorts() []string {
	content, err := ioutil.ReadFile("ports.txt")
	if err != nil {
		log.Fatal(err)
	}
	portString := bytes.NewBuffer(content).String()
	portArr := strings.Split(portString, ",")

	return portArr
}

func runClient(connection string) {
	fmt.Println("starting client...")
	// Set up a connection to the server.
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(connection, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := proto.NewNodeManagementClient(conn)
	message := fmt.Sprintf("Saying hi from %d", *idf)
	reply, err := c.RequestReply(context.Background(), &proto.Request{Message: message, Id: *idf, Timestamp: 1})
	if err != nil {
		log.Fatalf("Reply cause fatal error: %v", err)
	}
	fmt.Println(reply.Message)

}
