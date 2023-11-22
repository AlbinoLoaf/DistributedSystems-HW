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

	proto "hw05/grpc"

	"google.golang.org/grpc"
)

type Node struct {
	proto.UnimplementedAuctionServer
	id              int64
	Timestamp       int64
	hightestSeenBid int64
	RedundancyNodes map[int64]proto.NodeServer
	ctx             context.Context
	mu              sync.Mutex
}

func main() {
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int64(arg1) + 5000

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	n := &Node{
		id:              ownPort,
		Timestamp:       1,
		hightestSeenBid: 1,
		RedundancyNodes: make(map[int64]proto.AuctionClient),
		ctx:             ctx,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterAuctionServer(s, n)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()

	for i := 0; i < 3; i++ {
		port := int64(5000) + int64(i)

		if port == ownPort {
			continue
		}

		var conn *grpc.ClientConn
		fmt.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		c := proto.NewAuctionClient(conn)
		n.RedundancyNodes[port] = c
	}
	for true {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			n.sendPingToAll()
		}
	}
}
func (n *Node) sendBid(ctx context.Context, req *proto.Request) (*proto.serverReply, error) {
	n.Timestamp += 1

	rep := &proto.serverReply{Succes: true}
	return rep, nil
}

func (n *Node) sendPingToAll() {
	sendBid := &proto.sendBid{bid: n.hightestSeenBid, id: n.id}
	for id, client := range n.RedundancyNodes {
		serverReply, err := client.sendBid(n.ctx, sendBid)
		if err != nil {
			fmt.Println("something went wrong")
		}
		fmt.Printf("Got reply from id %v: %v\n", id, serverReply.Succes)
	}

}
