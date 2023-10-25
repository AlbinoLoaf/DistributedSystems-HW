package main

import (
	proto "ChittyChat/grpc"

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
	name       string
	id         int64
	portNumber int
}

var (
	clientPort = flag.Int("cPort", 0, "client port number")
	serverPort = flag.Int("sPort", 0, "server port number (should match the port used for the server)")
	clientName = flag.String("n", " ", "Client name")
)

func main() {
	// Parse the flags to get the port for the client
	flag.Parse()
	c := &Client{}
	c.CreateClient(*clientName, *clientPort)

	// // Wait for the client (user) to ask for the time
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
func (c *Client) CreateClient(Inputname string, clientport int) {
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
	c.portNumber = clientport
	//return returnClient, nil
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
