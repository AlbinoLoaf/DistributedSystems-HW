package main

import (
	"fmt"
	"time"
)

func main() {
	var handshake = make(chan []int, 2)
	var info = make(chan string)
	server := server{"server", handshake, make([]int, 2), info, 0, 0}
	client := client{"client", handshake, make([]int, 2), info, 0, 0}
	go client.run()
	go server.run()
	i := 1
	for i < 60 {
		time.Sleep(time.Second)
		fmt.Println(i)
		i++
	}
	// create channel
	// go run server
	// go run client
}
