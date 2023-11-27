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
	AuctionTime     int64
	server          proto.AuctionClient
	Timeer          bool
	conn            *grpc.ClientConn
}

func main() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	c := &BidClient{
		id:              1,
		Mybid:           0,
		AuctionDuration: 30,
		AuctionTime:     0,
		server:          nil,
		Timeer:          false,
		conn:            nil,
	}
	c.AttemptConnection(opts)
	c.getId()
	go c.AuctionTimer()
	scanner := bufio.NewScanner(os.Stdin)
	for {	
		if scanner.Scan() {
			text := scanner.Text()
			c.ResolveBid(text)
			}
		}
	}

func (c *BidClient) checkTime() bool {
	//fmt.Printf("Checking time\nTime at %d\n%t\n", n.Timestamp, n.Timestamp > 10)
	if c.Timeer {
		return false
	} else {
		return true
	}
}

func (c *BidClient) AuctionTimer() {
	var i int64
	for i = 0; i < c.AuctionDuration; i++ {
		time.Sleep(time.Second)
	}
	fmt.Println("Auction is over")
	os.Exit(0)
}

func (c *BidClient) getId() {
	reply, err := c.server.RequestId(context.Background(), &proto.RequestClientId{Message: "give me Id"})
	if err != nil {
		log.Fatalf("Getting id failed: %v", err)
	}
	fmt.Printf("my id before assignment %d\n", c.id)
	c.id = reply.Id
	c.AuctionTime = reply.Timestamp
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
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}
	if c.conn.GetState().String() != "READY" {
		// Close the existing connection before creating a new one
		c.conn.Close()
	}
	c.AttemptConnection(opts)

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

func (c *BidClient) AttemptConnection(opts []grpc.DialOption) *grpc.ClientConn {
	var err error
	for i := 0; i < 2; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		port := int64(5000) + int64(i)
		fmt.Printf("trying port %d\n", port)
		c.conn, err = grpc.DialContext(ctx, fmt.Sprintf(":%v", port), opts...)
		if err != nil {
			log.Printf("Attempt %d: did not connect: %v\n", i+1, err)
			time.Sleep(time.Second * 1)
		} else {
			c.server = proto.NewAuctionClient(c.conn)

			defer cancel()
			break
		}
	}
	return c.conn
}
