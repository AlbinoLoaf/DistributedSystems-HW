package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	proto "hw05/grpc"

	"google.golang.org/grpc"
)

type Node struct {
	proto.UnimplementedAuctionServer
	id              int64
	Timestamp       int64
	hightestSeenBid int64
	RedundancyNodes map[int64]proto.AuctionClient
	ctx             context.Context
	mu              sync.Mutex
	knownClients    int64
	// clients         map[int64]proto.AuctionClient
}

func main() {
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int64(arg1) + 5000

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	n := &Node{
		id:              ownPort,
		Timestamp:       1,
		hightestSeenBid: 0,
		RedundancyNodes: make(map[int64]proto.AuctionClient),
		ctx:             ctx,
		knownClients:    0,
		// clients:         make(map[int64]proto.AuctionClient),
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterAuctionServer(s, n)
	fmt.Printf("server registered at port %d\n", ownPort)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to server %v\n", err)
		}
	}()
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}

	for i := 0; i < 2; i++ {
		port := int64(5000) + int64(i)
		fmt.Printf("my portn%d the other port %d\n", ownPort, port)
		if port == ownPort {
			continue
		} else {

			var conn *grpc.ClientConn
			fmt.Printf("Trying to dial: %v\n", port)
			defer cancel()
			conn, err = grpc.DialContext(ctx, fmt.Sprintf(":%v", port), opts...)
			if err != nil {
				log.Printf("Attempt %d: did not connect: %v\n", i+1, err)
				time.Sleep(time.Second * 5)
			} else {
				defer conn.Close()

				c := proto.NewAuctionClient(conn)
				n.RedundancyNodes[port] = c
				fmt.Printf("Map Mapping from %d to %v\n", port, n.RedundancyNodes[port])
			}
		}
	}
	for k, i := range n.RedundancyNodes {
		fmt.Printf("server node %d has map from %d to %d\n", ownPort, k, i)
	}
	go n.AuctionDuration()

	// // Create a channel to signal when checkTime() returns false
	// done := make(chan bool)

	// // Run checkTime() in a separate goroutine
	// go func() {
	// 	for n.checkTime() {
	// 		// Sleep for a while to avoid busy waiting
	// 		time.Sleep(time.Second)
	// 	}
	// 	// Send a signal when checkTime() returns false
	// 	done <- true
	// }()

	for {
		if n.checkTime() {
			fmt.Printf("The Auction is over!\nHighest bid was: %d", n.hightestSeenBid)
			os.Exit(0)
			// Sleep for a while to avoid busy waiting
		} else {

			continue

		}
	}
}
func (n *Node) checkTime() bool {
	//fmt.Printf("Checking time\nTime at %d\n%t\n", n.Timestamp, n.Timestamp > 10)
	if n.Timestamp < 30 {

		return false
	} else {
		return true
	}
}

func (n *Node) AuctionDuration() {
	var i int64
	for i = 0; i <= 30; i++ {
		//fmt.Printf("Time at %d\n",i)
		time.Sleep(time.Second)
		n.Timestamp = i
		//fmt.Printf("Time at %d\n", n.Timestamp)
	}

}
func (n *Node) SendBid(ctx context.Context, req *proto.Bid) (*proto.ServerReply, error) {
	//n.Timestamp += 1

	rep := &proto.ServerReply{Succes: true}
	if n.ResolveBid(req.Bid) {
		n.ServerRedundancy()
	}
	return rep, nil
}

func (n *Node) Redundancy(ctx context.Context, req *proto.Bid) (*proto.ServerReply, error) {
	n.Timestamp += 1
	_ = n.ResolveBid(req.Bid)
	rep := &proto.ServerReply{Succes: true}
	return rep, nil
}

func (n *Node) ServerRedundancy() {
	sendBid := &proto.Bid{Bid: n.hightestSeenBid, Id: n.id}

	for id, client := range n.RedundancyNodes {
		if client == nil {
			fmt.Printf("client %d is nil\n", id)
			continue
		}
		serverReply, err := client.Redundancy(n.ctx, sendBid)
		if err != nil {
			fmt.Println("something went wrong, removing client")
			delete(n.RedundancyNodes, id)
			continue
		}
		fmt.Printf("Got reply from id %v: %v\n", id, serverReply.Succes)
	}
}

func (n *Node) ResolveBid(bid int64) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	if bid > n.hightestSeenBid {
		n.hightestSeenBid = bid
		fmt.Printf("I server %d got a new bid %d\n", n.id, n.hightestSeenBid)
		return true
	} else {
		fmt.Printf("Bid %d was lower than %d \n", bid, n.hightestSeenBid)
		return false
	}
}

func (n *Node) RequestId(ctx context.Context, in *proto.RequestClientId) (*proto.ClientId, error) {
	n.knownClients += 1
	// Map client ids to client and send information to all other servers

	// n.clients[n.knownClients] = n.RedundancyNodes[n.knownClients]
	// for id, client := range n.RedundancyNodes {
	// 	reply, err := client.RequestId(n.ctx, &proto.RequestClientId{})
	// 	if err != nil {
	// 		fmt.Println("something went wrong")
	// 	}
	// 	fmt.Printf("Got reply from id %v: %v\n", id, reply.Id)
	// }
	return &proto.ClientId{Id: n.knownClients, Timestamp: n.Timestamp}, nil
}
