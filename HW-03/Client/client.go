package main

import (
	proto "ChittyChat/grpc"
	"time"

	"context"
	"flag"
	"log"

	"bufio"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	name      string
	id        int64
	Timestamp int64
}

var (
	//clientPort = flag.Int("cPort", 0, "client port number")
	serverPort = flag.Int("sPort", 0, "server port number (should match the port used for the server)")
	clientName = flag.String("n", " ", "Client name")
)

func main() {
	// Parse the flags to get the port for the client
	flag.Parse()
	c := &Client{}
	c.CreateClient(*clientName)

	// // Wait for the client (user) to ask for the time
	go c.AwaitChange()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if strings.ToLower(scanner.Text()) == "leave" {
			err := c.leaveServer()
			if err != nil {
				log.Printf("Client leaving failed: %v", err)
			} else {
				log.Printf("Client %s with id %d left", c.name, c.id)
				os.Exit(0)
			}
		} else {

			log.Printf("Publish: %s", scanner.Text())
		}
	}
	//go waitForTimeRequest(client)

	for {

	}
}

func (c *Client) compareAndUpdate(serverConnection proto.UsermanagementClient) {
	Server, err := serverConnection.RequestChange(context.Background(), &proto.Timestamp{
		Time: c.Timestamp,
	})
	if err != nil {
		log.Fatalf("Compare And update sad %v", err)
	}
	if Server.Time > c.Timestamp {
		i := c.Timestamp + 1
		for i <= Server.Time {
			c.RequestEvent(i)
			i++
		}
		c.Timestamp = Server.Time
	}
}

func (c *Client) AwaitChange() {
	serverConnection, err := connectToServer()
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)

	}
	for true {
		time.Sleep(time.Millisecond * 1000)
		c.compareAndUpdate(serverConnection)
	}

}

func (c *Client) Publish(message string) {
	serverConnection, err := connectToServer()
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)

	}
	c.compareAndUpdate(serverConnection)
	Server, err := serverConnection.SendMessage(context.Background(), &proto.PublishMessage{
		Name:    c.name,
		Message: message,
		Id:      c.id,
	})
	c.Timestamp = Server.Time

}

func (c *Client) RequestEvent(wantedEvent int64) {
	serverConnection, err := connectToServer()
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)

	}
	Event, err := serverConnection.RequestBroadcast(context.Background(), &proto.Timestamp{
		Time: wantedEvent,
	})
	log.Print(Event)
}
func (c *Client) leaveServer() error {
	serverConnection, err := connectToServer()
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)

	}
	state, err := serverConnection.LeaveClient(context.Background(), &proto.Client{
		Name: c.name,
		Id:   c.id,
	})
	if err != nil {
		return err
	}
	if state.Accept {
		return nil
	}
	return nil
}
func (c *Client) CreateClient(Inputname string) {
	serverConnection, err := connectToServer()
	if err != nil {
		log.Fatalf("User creation fallied: %v", err)

	}
	client, err := serverConnection.ClientJoin(context.Background(), &proto.NewClient{
		Name: Inputname,
	})
	if err != nil {
		log.Fatalf("User creation fallied: %v", err)
	}
	log.Printf(`
	ClientName: %s
	ClientID: %d`, client.GetName(), client.GetId())

	c.name = client.GetName()
	c.id = client.GetId()
	c.Timestamp = client.GetTimestamp()
}

// func waitForTimeRequest(client *Client) {
// 	// Connect to the server
// 	serverConnection, _ := connectToServer()

// 	// Wait for input in the client terminal
// 	scanner := bufio.NewScanner(os.Stdin)
// 	for scanner.Scan() {
// 		input := scanner.Text()
// 		log.Printf("Client asked for time with input: %s\n", input)

// 		// Ask the server for the time
// 		timeReturnMessage, err := serverConnection.NewClient(context.Background(), &proto.AskForTimeMessage{
// 			ClientId: int64(client.id),
// 		})

// 		if err != nil {
// 			log.Printf(err.Error())
// 		} else {
// 			log.Printf("Server %s says the time is %s\n", timeReturnMessage.ServerName, timeReturnMessage.Time)
// 		}
// 	}
// }

func connectToServer() (proto.UsermanagementClient, error) {
	// Dial the server at the specified port.
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", *serverPort)
	} else {
		log.Printf("Connected to the server at port %d\n", *serverPort)
	}
	return proto.NewUsermanagementClient(conn), nil
}
