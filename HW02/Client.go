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

		if c.communiationSlice[0] == 0 {
			c.clientNr = rand.Intn(100) + 1 //avoiding 0
			c.communiationSlice[0] = c.clientNr
			fmt.Printf("%s %s", c.name, "Initiated Handshake \n")
			c.handshake <- c.communiationSlice

		} else if c.clientNr != 0 {
			c.communiationSlice = <-c.handshake

			c.serverNr = c.communiationSlice[1]

			if c.communicationCheck(c.communiationSlice[0]) {
				fmt.Printf("%s %s", c.name, "Finishing Handshake\n")

				c.communiationSlice[0] = c.clientNr
				c.communiationSlice[1] = c.serverNr
				fmt.Println(c.communiationSlice)
				c.handshake <- c.communiationSlice
				fmt.Println("c")
				c.infoChan <- "It works"
			} else {
				fmt.Printf("%s %s", c.name, "Initiated Handshake\n")
				c.handshake <- c.communiationSlice
				time.Sleep(time.Second)
			}

		}

	}
}

func (c *client) communicationCheck(input int) bool {
	if input == (c.clientNr + 1) {
		c.clientNr = input
		c.serverNr++
		return true
	} else {
		return false
	}
}
