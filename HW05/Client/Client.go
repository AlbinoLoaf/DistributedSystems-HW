package main

import (
	"bufio"
	"context"
	"fmt"
	proto "hw05/grpc"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BidClient struct {
	id              int64
	Mybid           int64
	AuctionDuration int64
	server          proto.AuctionClient
}

func main() {
	c := &BidClient{
		id:              1,
		Mybid:           2,
		AuctionDuration: 3,
		server:          proto.AuctionClient,
	}

	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	for i := 0; i < 4; i++ {
		port := int64(5000) + int64(i/2)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		var conn *grpc.ClientConn
		conn, err := grpc.DialContext(ctx, fmt.Sprintf(":%v", port), opts...)

		if err != nil {
			log.Printf("Attempt %d: did not connect: %v", i+1, err)
			time.Sleep(time.Second * 5)
		} else {
			defer conn.Close()
			break
		}
		if conn == nil {
			log.Fatalf("Could not establish connection after %d attempts", 8)
		} else {
			defer conn.Close()
			fmt.Print("B")
		}
		// c_ := proto.NewAuctionClient(conn)
	}
	for true {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Printf("%v", scanner.Text())
		}
	}
}
