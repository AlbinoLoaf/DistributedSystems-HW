package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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
}

func main() {
	flag.Parse()
	fmt.Println("Starting server...")
	// log.SetOutput(os.Stdout)

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", *Address+":"+*Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a proto server object
	s := grpc.NewServer()

	// Attach the Node service to the server
	proto.RegisterNodeManagementServer(s, &Server{id: *idf})

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

	for {
		c.RequestReply(context.Background(), &proto.Request{Message: "go Ahead", Id: *idf, Timestamp: 1})
		i := rand.Intn(6)
		time.Sleep(time.Duration(i) + 3)
	}
}

func (s *Server) RequestReply(ctx context.Context, in *proto.Request) (*proto.Reply, error) {
	// If true this node accepts the requests replying go ahead
	if s.ResolveRequest(in.Id, in.Timestamp) {
		return &proto.Reply{Message: "Go ahead", Id: in.Id}, nil
	}
	return nil, nil

}
func (s *Server) ResolveRequest(ForeignTimestamp int64, foreignID int64) (Accaptance bool) {
	fmt.Printf("Server %d got contacted by server %d\n", s.id, foreignID)
	time.Sleep(time.Second)
	return true

}

// type Node struct {
// 	id                   int64
// 	timestamp            int64
// 	ipList               []int
// 	usingCriticalSection bool
// }

// var (
// 	nodePort = flag.Int("nPort", 0, "ode port number")
// 	idflag   = flag.Int("nId", 0, "Node id")
// 	//serverPort = flag.Int("sPort", 0, "server port number (should match the port used for the server)")
// )

// func (n *Node) nodemain() {
// 	flag.Parse()

// }

func GetPorts() []string {
	content, err := ioutil.ReadFile("ports.txt")
	if err != nil {
		log.Fatal(err)
	}
	portString := bytes.NewBuffer(content).String()
	portArr := strings.Split(portString, ",")

	return portArr
}

// func (n *Node) createNode() {

// }

// func (n *proto.NodeManagementClient) RequestReply(ctx context.Context, in *proto.Request) (*proto.Reply, error) {
// 	// If true this node accepts the requests replying go ahead
// 	if n.ResolveRequest(in.Id, in.Timestamp) {
// 		return &proto.Reply{Message: "Go ahead", Id: in.Id}, nil
// 	}
// 	return nil, nil

//}
