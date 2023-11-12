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
	usingCriticalSection                    bool
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
	fmt.Printf("Server type: %T\n", s)
	// Attach the Node service to the server
	proto.RegisterNodeManagementServer(s, &Server{id: *idf, usingCriticalSection: false})

	// start sending requests to peer
	list := GetPorts()
	i := 0
	for i < len(list) {
		if !strings.Contains(list[i], *Port) {
			go runClient(list[i], s)

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

func runClient(connection string, server *grpc.Server) {
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
	answers := 0

	for {
		reply, err := c.RequestReply(context.Background(), &proto.Request{Message: "go Ahead", Id: *idf, Timestamp: 1})
		if err != nil {
			log.Fatalf("Reply cause fatal error: %v", err)
		}
		//We know this is cursed but it works.
		//time.duration is in nanoseconds and
		//we desire the order of maginitude to be in seconds
		var magnitude int
		magnitude = 1e9
		waitTimer := (rand.Intn(4) + 1) * magnitude
		fmt.Println(waitTimer / magnitude)

		if "Go ahead" == reply.Message {
			answers++
			if answers == 2 {

				thinkTimer := (rand.Intn(10) + 5) * magnitude
				fmt.Println(thinkTimer / magnitude)
				//Take critical Setion.
				fmt.Print("I am the criticallest of sections\n")
				time.Sleep(time.Duration(thinkTimer))
				answers = 0
			}
		}
		time.Sleep(time.Duration(waitTimer))
		fmt.Printf("Server %d was told %s\n", *idf, reply.Message)
	}
}

func (s *Server) RequestReply(ctx context.Context, in *proto.Request) (*proto.Reply, error) {
	// If true this node accepts the requests replying go ahead
	if s.ResolveRequest(in.Timestamp, in.Id) {
		return &proto.Reply{Message: "Go ahead", Id: in.Id}, nil
	} else {
		return &proto.Reply{Message: "Don't go ahead", Id: in.Id}, nil
	}
}
func (s *Server) ResolveRequest(ForeignTimestamp int64, foreignID int64) (Accaptance bool) {
	fmt.Printf("Server %d got contacted by server %d\n", s.id, foreignID)
	for s.usingCriticalSection {
		time.Sleep(time.Millisecond * 500)
	}
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
