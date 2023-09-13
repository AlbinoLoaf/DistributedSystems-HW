package main

import (
	"fmt"
	"time"
)

func main() {
	forkOneCommunication := make(chan string, 1)   // Fork between philosofer 1 and 2
	forkTwoCommunication := make(chan string, 2)   // Fork between philosofer 2 and 3
	forkThreeCommunication := make(chan string, 1) // Fork between philosofer 3 and 4
	forkFourCommunication := make(chan string, 1)  // Fork between philosofer 4 and 5
	forkFiveCommunication := make(chan string, 1)  // Fork between philosofer 5 and 1

	philOne := Philosopher{"One", false, forkOneCommunication, forkFiveCommunication, 0}
	philTwo := Philosopher{"Two", false, forkTwoCommunication, forkOneCommunication, 0}
	philThree := Philosopher{"Three", false, forkThreeCommunication, forkTwoCommunication, 0}
	philFour := Philosopher{"Four", false, forkFourCommunication, forkThreeCommunication, 0}
	PhilFive := Philosopher{"Five", false, forkFourCommunication, forkFiveCommunication, 0} // Loking the other way
	
	go PhilFive.run()
	go philOne.run()
	go philTwo.run()
	go philThree.run()
	go philFour.run()

	i := 1
	for i < 60 {
		time.Sleep(time.Second * 1)
		fmt.Print("time running ")
		fmt.Println(i)
		i++
	}
}

/*
main.
go fork 1 - 5
go phi 1 -5
while
 phi 1 - 5 task
 phi times eaten > 3 end
 await 2 second
*/
