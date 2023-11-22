package main

import (
	"bufio"
	"context"
	"fmt"
	proto "hw05/grpc"
	"log"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BidClient struct {
	id              int64
	Mybid           int64
	AuctionDuration int64
	server          proto.AuctionClient
	ctx             context.Context
}

func main() {

	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	var conn *grpc.ClientConn
	for i := 0; i < 4; i++ {
		port := int64(5000) + int64(i/2)
		fmt.Printf("trying port %d\n", port)
		defer cancel()
		conn, err := grpc.DialContext(ctx, fmt.Sprintf(":%v", port), opts...)
		if err != nil {
			log.Printf("Attempt %d: did not connect: %v\n", i+1, err)
			time.Sleep(time.Second * 5)
		} else {
			defer conn.Close()
			break
		}
	}

	// if conn == nil {
	// 	log.Fatalf("Could not establish connection after %d attempts", 2)
	// } else {
	// 	defer conn.Close()
	// }
	c_ := proto.NewAuctionClient(conn)
	c := &BidClient{
		id:              1,
		Mybid:           0,
		AuctionDuration: 3,
		server:          c_,
		ctx:             ctx,
	}
	for true {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			c.ResolveBid(text)

		}

	}
}

func (c *BidClient) ResolveBid(bid string) {
	if num, err := strconv.ParseInt(bid, 10, 64); err == nil {
		fmt.Printf("You entered the integer: %d\n", num)
		if num < c.Mybid {
			fmt.Printf("The integer %d is not a valid bid \nPlease enter lager bid than last bet of %d", num, c.Mybid)
		}
		c.BidAuction(num)

	} else {
		fmt.Println("You did not enter an integer.")
	}
}

func (c *BidClient) BidAuction(mybid int64) {
	reply, err := c.server.SendBid(context.Background(), &proto.Bid{
		Bid: mybid,
		Id:  c.id})
	if err != nil {
		log.Fatalf("Bidding failed: %v", err)
	}
	fmt.Printf("The bid was %v", reply.Succes)

}
