package main

import (
	"bufio"
	"context"
	"fmt"
	proto "hw05/grpc"
	"log"
	"os"
	"reflect"
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
	defer cancel()
	var conn *grpc.ClientConn
	var err error

	for i := 0; i < 2; i++ {
		port := int64(5000) + int64(i)
		fmt.Printf("trying port %d\n", port)
		conn, err = grpc.DialContext(ctx, fmt.Sprintf(":%v", port), opts...)
		if err != nil {
			log.Printf("Attempt %d: did not connect: %v\n", i+1, err)
			time.Sleep(time.Second * 1)
		} else {
			break
		}
	}

	if conn == nil {
		log.Fatalf("Could not establish connection after %d attempts", 2)
	} else {
		defer conn.Close()
	}
	c_ := proto.NewAuctionClient(conn)
	c := &BidClient{
		id:              1,
		Mybid:           0,
		AuctionDuration: 3,
		server:          c_,
		ctx:             ctx,
	}
	c.getId()
	for true {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			c.ResolveBid(text)

		}

	}
}

func (c *BidClient) getId() {
	reply, err := c.server.RequestId(context.Background(), &proto.RequestClientId{Message: "give me Id"})
	if err != nil {
		log.Fatalf("Getting id failed: %v", err)
	}
	fmt.Printf("my id before assignment %d\n", c.id)
	c.id = reply.Id
	fmt.Printf("recieved id %d and my assigned id %d\n", reply.Id, c.id)

}

func (c *BidClient) ResolveBid(bid string) {
	if num, err := strconv.ParseInt(bid, 10, 64); err == nil {
		fmt.Printf("You entered the integer: %d\n", num)
		//Checks that the bid is larger than your last bid,
		//initial bid is zero thu no negative bids
		if num < c.Mybid {
			fmt.Printf("The integer %d is not a valid bid \nPlease enter lager bid than last bet of %d\n", num, c.Mybid)
		} else {
			c.BidAuction(num)
		}
	} else {
		fmt.Println("You did not enter an integer.")
	}
}

func (c *BidClient) BidAuction(mybid int64) {
	c.Mybid = mybid
	fmt.Printf("Bidding %d\nwith Id: %d\n", mybid, c.id)
	fmt.Println(reflect.TypeOf(mybid))
	reply, err := c.server.SendBid(context.Background(), &proto.Bid{
		Bid: mybid,
		Id:  c.id})
	if err != nil {
		log.Fatalf("Bidding failed: %v", err)
	}
	fmt.Printf("The bid was %v\n", reply.Succes)

}
