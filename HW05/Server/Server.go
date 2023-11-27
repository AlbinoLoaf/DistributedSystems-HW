package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	proto "hw05/grpc"

	"github.com/bradhe/stopwatch"
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
	clients         map[int64]proto.AuctionClient
}

func main() {
	watch := stopwatch.Start()
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int64(arg1) + 5000

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	n := &Node{
		id:              ownPort,
		Timestamp:       1,
		hightestSeenBid: 1,
		RedundancyNodes: make(map[int64]proto.AuctionClient),
		ctx:             ctx,
		knownClients:    2,
		//clients:         make(map[int64]proto.AuctionClient),
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
			fmt.Printf("connection: %v\n", conn)
			if err != nil {
				log.Printf("Attempt %d: did not connect: %v\n", i+1, err)
				time.Sleep(time.Second * 5)
			} else {
				defer conn.Close()

				fmt.Printf("connection assigned: %v\n", conn)
				c := proto.NewAuctionClient(conn)
				n.RedundancyNodes[port] = c
				fmt.Printf("Map Mapping from %d to %v\n", port, n.RedundancyNodes[port])
			}
		}
	}
	for k, i := range n.RedundancyNodes {
		fmt.Printf("server node %d has map from %d to %d\n", ownPort, k, i)
	}
	for true {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			sendBid := &proto.Bid{Bid: n.hightestSeenBid, Id: n.id}

			for id, client := range n.RedundancyNodes {
				serverReply, err := client.SendBid(n.ctx, sendBid)
				if err != nil {
					watch.Stop()
					fmt.Printf("Milliseconds elapsed: %v\n", watch.Milliseconds())
					fmt.Println("something went wrong\n")
				}
				fmt.Printf("Got reply from id %v: %v\n", id, serverReply.Succes)
			}
		}
	}
}

func (n *Node) SendBid(ctx context.Context, req *proto.Bid) (*proto.ServerReply, error) {
	n.Timestamp += 1

	rep := &proto.ServerReply{Succes: true}
	return rep, nil
}

func (n *Node) sendPingToAll() {
	sendBid := &proto.Bid{Bid: n.hightestSeenBid, Id: n.id}

	for id, client := range n.RedundancyNodes {
		serverReply, err := client.SendBid(n.ctx, sendBid)
		if err != nil {
			fmt.Println("something went wrong\n")
		}
		fmt.Printf("Got reply from id %v: %v\n", id, serverReply.Succes)
	}
}

func (n *Node) RequestId(ctx context.Context, in *proto.RequestClientId) (*proto.ClientId, error) {
	n.knownClients += 1
	//n.clients[n.knownClients] =
	return &proto.ClientId{Id: n.knownClients}, nil
}
