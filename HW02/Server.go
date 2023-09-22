package main

import (
	"fmt"
	"math/rand"
	"time"
)

type server struct {
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

func (s *server) run() {

	for true {
		// The server recieves the first handshake from the client it then increments that number,
		// makes it own sequence number and the send both back.
		s.communiationSlice = <-s.handshake
		if s.communiationSlice[1] == 0 {
			fmt.Printf("%s %s", s.name, "Recived handshake\n")
			s.acknowledgement(s.communiationSlice[0])
			fmt.Println(s.communiationSlice)
			s.handshake <- s.communiationSlice
			// here it gets the acknowledgement and date from the client it then prints the data it got and the handshake is done.
		} else if s.communicationCheck(s.communiationSlice[1]) {
			s.communiationSlice[0] = s.clientNr
			s.communiationSlice[1] = s.serverNr
			fmt.Println("a")
			fmt.Printf("%s %s", s.name, "Handshake done\n")
			fmt.Printf("%s", <-s.infoChan)
			// If we get the wrong information we try and start over.
		} else {
			fmt.Println("b")
			s.handshake <- s.communiationSlice
			time.Sleep(time.Second)
		}

	}
}

// here we make the servers sequence number and increment the clients sequence number.
func (s *server) acknowledgement(input int) {
	s.serverNr = rand.Intn(100) + 1 //avoiding 0
	s.communiationSlice[1] = s.serverNr
	s.communiationSlice[0] = input + 1
}

// this checks if the incoming acknowledgement is one number more than the sequence number it send itself.
// It then increments the sequence number it got in the same message
func (s *server) communicationCheck(input int) bool {
	if input == (s.serverNr + 1) {
		s.serverNr = input
		s.clientNr++
		return true
	} else {
		return false
	}
}
