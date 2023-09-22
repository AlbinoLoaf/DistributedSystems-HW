package main

import (
	"fmt"
	"math/rand"
	"time"
)

type client struct {
	//Name of client
	name string
	//handshake channel
	handshake chan []int
	//Slice for acknowledgement and squence
	// index 0 is the client and index 1 is the server
	communiationSlice []int
	//Information channel
	infoChan chan string
	//Sequence Number
	clientNr int
	//Acknowledgement Number
	serverNr int
}

func (c *client) run() {
	for true {
		// the first part of the handshake where it sends a sequence number
		if c.communiationSlice[0] == 0 {
			c.clientNr = rand.Intn(100) + 1 //avoiding 0
			c.communiationSlice[0] = c.clientNr
			fmt.Printf("%s %s", c.name, "Initiated Handshake \n")
			c.handshake <- c.communiationSlice
			// gets the acknowledgement from the server gets it only if it itself has send a number before hand.
		} else if c.clientNr != 0 {
			c.communiationSlice = <-c.handshake

			c.serverNr = c.communiationSlice[1]
			//here we check if the aknowledgement is correct the we increment and the send the aknowledgement and data.
			if c.communicationCheck(c.communiationSlice[0]) {
				fmt.Printf("%s %s", c.name, "Finishing Handshake\n")

				c.communiationSlice[0] = c.clientNr
				c.communiationSlice[1] = c.serverNr
				fmt.Println(c.communiationSlice)
				c.handshake <- c.communiationSlice
				fmt.Println("c")
				c.infoChan <- "It works"
				// If we get the wrong information we try and start over.
			} else {
				fmt.Printf("%s %s", c.name, "Initiated Handshake\n")
				c.handshake <- c.communiationSlice
				time.Sleep(time.Second)
			}

		}

	}
}

// this checks if the incoming acknowledgement is one number more than the sequence number it send itself.
// It then increments the sequence number it got in the same message
func (c *client) communicationCheck(input int) bool {
	if input == (c.clientNr + 1) {
		c.clientNr = input
		c.serverNr++
		return true
	} else {
		return false
	}
}
