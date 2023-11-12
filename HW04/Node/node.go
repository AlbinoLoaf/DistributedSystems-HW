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
	Timestamp                               int64
	usingCriticalSection                    bool
	Nodes                                   []Node
	inquiries                               []int
}

func main() {
	flag.Parse()
	fmt.Println("Starting server...")
	// log.SetOutput(os.Stdout)
	server := &Server{
		id:                   *idf,
		Timestamp:            1,
		usingCriticalSection: false,
		Nodes:                make([]Node, 0),
	}
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", *Address+":"+*Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Address for server")
	fmt.Println(lis.Addr())
	// Create a proto server object
	s := grpc.NewServer()
	// Attach the Node service to the server
	proto.RegisterNodeManagementServer(s, server)

	// start sending requests to peer
	list := GetPorts()
	i := 0
	for i < len(list) {
		if !strings.Contains(list[i], *Port) {
			go server.runClient(list[i])
			fmt.Print("g")

		}
		i++
	}
	fmt.Print("h")
	randomSleep(5e3, 1e3)
	count := 0
	for count < 5 {
		for !server.usingCriticalSection {
			fmt.Print("E")
			server.usingCriticalSection = server.confirmation()
			fmt.Print("F")
			fmt.Print("tis the season")
			randomSleep(1e3, 9e3)
			i := 0
			for i < len(server.Nodes) {
				go server.Nodes[i].RequestCriticalSection()
				i++
			}
		}
		count++
	}
	// Serve proto server
	fmt.Println("Serving proto on port " + *Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	// code here is unreachable because protoServer.Serve occupies the current thread.
}

func (s *Server) runClient(connection string) {
	fmt.Println("starting client...")
	// Set up a connection to the server.
	fmt.Printf("\n%s\n", connection)
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	fmt.Print("A")
	var conn *grpc.ClientConn
	var err error

	for i := 0; i < 8; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		conn, err = grpc.DialContext(ctx, connection, opts...)
		if err != nil {
			// If there's an error, sleep for a while before retrying
			log.Printf("Attempt %d: did not connect: %v", i+1, err)
			time.Sleep(time.Second * 5)
		} else {
			// If connection is successful, break out of the loop
			break
		}
	}

	if conn == nil {
		log.Fatalf("Could not establish connection after %d attempts", 8)
	} else {
		defer conn.Close()
		fmt.Print("B")
	}
	c_ := proto.NewNodeManagementClient(conn)
	message := fmt.Sprintf("Saying hi from %d", *idf)
	fmt.Print("C")
	reply, err := c_.InitialContact(context.Background(), &proto.Request{Message: message, Id: s.id, Timestamp: 1})
	if err != nil {
		log.Fatalf("Reply cause fatal error: %v", err)
	}
	fmt.Print("D")
	fmt.Println(reply.Message)
	s.Nodes = append(s.Nodes, Node{
		c:           c_,
		Add:         connection,
		serverID:    reply.Id,
		QueryServer: s.id,
		permision:   false,
	})
}
func (n *Node) RequestCriticalSection() {
	message := fmt.Sprintf("Can %d have the critical section?", n.QueryServer)
	reply, err := n.c.RequestCriticalSection(context.Background(), &proto.Request{Message: message, Id: n.QueryServer, Timestamp: 1})
	if err != nil {
		log.Fatalf("Reply cause fatal error: %v", err)
	}
	n.permision = reply.Allowance
}

func (s *Server) InitialContact(ctx context.Context, in *proto.Request) (*proto.Reply, error) {
	// If true this node accepts the requests replying go ahead
	fmt.Print("Req Recieved")
	replyMSG := fmt.Sprintf("Saying hi from %d", s.id)
	return &proto.Reply{Message: replyMSG, Id: s.id}, nil
}
func (s *Server) RequestCriticalSection(ctx context.Context, in *proto.Request) (*proto.CriticalSectionAllowance, error) {
	for s.usingCriticalSection {
		fmt.Println(in.Message)
		randomSleep(1e2, 3e2)
	}
	fmt.Println(in.Message)
	fmt.Println("Approved")
	return &proto.CriticalSectionAllowance{Allowance: true, Id: s.id}, nil
}

type Node struct {
	c           proto.NodeManagementClient
	Add         string
	serverID    int64
	QueryServer int64
	permision   bool
}

func (s *Server) confirmation() bool {
	var Acknowldegments int
	var Acknowledged bool
	for !Acknowledged {
		Acknowldegments = 0
		i := 0
		for i < len(s.Nodes) {
			if s.Nodes[i].permision {
				Acknowldegments++
			}
			if Acknowldegments == len(s.Nodes) {
				Acknowledged = true
			}
		}
	}
	return Acknowledged
}

func (s *Server) working() (done bool) {
	randomSleep(5e3, 5e3)
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

func randomSleep(minimum int, rangeofNumbers int) {
	waitTimer := (rand.Intn(rangeofNumbers) + minimum)
	fmt.Println(waitTimer)
	time.Sleep(time.Duration(float64(waitTimer) * float64(time.Millisecond)))
}
