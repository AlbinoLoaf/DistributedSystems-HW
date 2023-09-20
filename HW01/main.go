package main

import (
	"fmt"
	"time"
)
/* HOW DOES OUR SOLUTION WORK (ensuring no deadlocks):
Our solution is logically meant to work by saying that every philosopher has a priority and non-priority fork, think of this as being a
right or left where each philosopher EXCEPT ONE prefers left. Because every philosopher except one will ask for the right hand fork
a deadlock will be prevented, the main reasson it gets prevented is due to our select statement which ensures that no channel gets blocked
before both forks are available. To begin with only the lefthanded philosopher can eat, but this will create a domino effect affecting the entire
table ensuring everyone gets a turn eating with only the forks that are available without any deadlock. */
//func main() {
	forkOneCommunication := make(chan string, 1)   // Fork between philosofer 1 and 2
	forkTwoCommunication := make(chan string, 2)   // Fork between philosofer 2 and 3
	forkThreeCommunication := make(chan string, 1) // Fork between philosofer 3 and 4
	forkFourCommunication := make(chan string, 1)  // Fork between philosofer 4 and 5
	forkFiveCommunication := make(chan string, 1)  // Fork between philosofer 5 and 1
	forkOneReq := make(chan string, 1)
	forkTwoReq := make(chan string, 1)
	forkThreeReq := make(chan string, 1)
	forkFourReq := make(chan string, 1)
	forkFiveReq := make(chan string, 1)
	philOne := Philosopher{"Philosopher One", false, forkOneCommunication, forkFiveCommunication,
		forkOneReq, forkFiveReq, 0}
	philTwo := Philosopher{"Philosopher Two", false, forkTwoCommunication, forkOneCommunication,
		forkTwoReq, forkOneReq, 0}
	philThree := Philosopher{"Philosopher Three", false, forkThreeCommunication, forkTwoCommunication,
		forkThreeReq, forkTwoReq, 0}
	philFour := Philosopher{"Philosopher Four", false, forkFourCommunication, forkThreeCommunication,
		forkFourReq, forkThreeReq, 0}
	PhilFive := Philosopher{"Philosopher Five", false, forkFourCommunication, forkFiveCommunication,
		forkFourReq, forkFiveReq, 0} // This is the "lefthanded" philosopher
    
	forkOne := Fork{"Fork one", forkOneReq}
	forkTwo := Fork{"Fork Two", forkTwoReq}
	forkThree := Fork{"Fork Three", forkThreeReq}
	forkFour := Fork{"Fork Four", forkFourReq}
	forkFive := Fork{"Fork Five", forkFiveReq}
	go forkOne.run()
	go forkTwo.run()
	go forkThree.run()
	go forkFour.run()
	go forkFive.run()

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
// Our initial pseudo-code before writing the program, could be relevant to see our thought process.
/*
main.
go fork 1 - 5
go phi 1 -5
while
 phi 1 - 5 task
 phi times eaten > 3 end
 await 2 second
*/
