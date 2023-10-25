package main

import (
	proto "ChittyChat/grpc"
	"context"
	"flag"
	"log"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	id         int
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
	log.Printf("test %s", *clientName)
	//log.Print("test")
	CreateClient(*clientName)

	// Wait for the client (user) to ask for the time

	//go waitForTimeRequest(client)

	for {

	}
}
func CreateClient(Inputname string) {
	serverConnection, _ := connectToServer()
	client, err := serverConnection.ClientJoin(context.Background(), &proto.NewClient{
		Name: Inputname,
	})
	if err != nil {
		log.Fatalf("User creation fallied: %v", err)
	}
	log.Printf(`
	ClientName: %s 
	ClientID: %s`, client.GetName(), client.GetId())

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
