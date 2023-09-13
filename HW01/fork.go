package main

import "fmt"

/*
Fork One long switch
Used bool

Channel
*/
type Fork struct {
	inUse       bool
	forkChannel chan string
}

func forkmain(fork Fork) {
	var inquiry string
	for true {
		inquiry = <-fork.forkChannel
		ForkLogic(fork, inquiry)
	}
}

func ForkLogic(fork Fork, inquiry string) {
	fmt.Print("Fork Being used: ")
	fmt.Println(fork.inUse)
	fmt.Println(inquiry)
	switch inquiry {
	case "Request fork":
		// Cheking boolian
		if fork.inUse {
			fork.forkChannel <- "Request Denied"
		} else {
			fork.inUse = true
			fork.forkChannel <- "Request Accepted"
		}
	case "Put down fork":
		fork.inUse = false

	}
}
