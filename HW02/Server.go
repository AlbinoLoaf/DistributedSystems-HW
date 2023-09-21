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
		s.communiationSlice = <-s.handshake
		if s.communiationSlice[1] == 0 {
			fmt.Printf("%s %s", s.name, "Recived handshake\n")
			s.acknowledgement(s.communiationSlice[0])
			fmt.Println(s.communiationSlice)
			s.handshake <- s.communiationSlice
		} else if s.communicationCheck(s.communiationSlice[1]) {
			s.communiationSlice[0] = s.clientNr
			s.communiationSlice[1] = s.serverNr
			fmt.Println("a")
			fmt.Printf("%s %s", s.name, "Handshake done\n")
			fmt.Printf("%s", <-s.infoChan)
		} else {
			fmt.Println("b")
			s.handshake <- s.communiationSlice
			time.Sleep(time.Second)
		}

	}
}
func (s *server) acknowledgement(input int) {
	s.serverNr = rand.Intn(100) + 1 //avoiding 0
	s.communiationSlice[1] = s.serverNr
	s.communiationSlice[0] = input + 1
}

func (s *server) communicationCheck(input int) bool {
	if input == (s.serverNr + 1) {
		s.serverNr = input
		s.clientNr++
		return true
	} else {
		return false
	}
}
