package main

import (
	"fmt"
	"time"
)

func main() {
	forkOneCommunication := make(chan string)   // Fork between philosofer 1 and 2
	forkTwoCommunication := make(chan string)   // Fork between philosofer 2 and 3
	forkThreeCommunication := make(chan string) // Fork between philosofer 3 and 4
	forkFourCommunication := make(chan string)  // Fork between philosofer 4 and 5
	forkFiveCommunication := make(chan string)  // Fork between philosofer 5 and 1

	ForkOneDone := make(chan string)   // Fork between philosofer 1 and 2
	ForkTwoDone := make(chan string)   // Fork between philosofer 2 and 3
	ForkThreeDone := make(chan string) // Fork between philosofer 3 and 4
	ForkFourDone := make(chan string)  // Fork between philosofer 4 and 5
	ForkFiveDone := make(chan string)  // Fork between philosofer 5 and 1

	forkOne := Fork{forkOneCommunication, ForkOneDone}
	forkTwo := Fork{forkTwoCommunication, ForkTwoDone}
	forkThree := Fork{forkThreeCommunication, ForkThreeDone}
	forkFour := Fork{forkFourCommunication, ForkFourDone}
	forkFive := Fork{forkFiveCommunication, ForkFiveDone}

	go forkmain(forkOne)
	go forkmain(forkTwo)
	go forkmain(forkThree)
	go forkmain(forkFour)
	go forkmain(forkFive)

	philOne := Philosopher{"One", false, forkOneCommunication, forkFiveCommunication, ForkOneDone, 0}
	philTwo := Philosopher{"Two", false, forkTwoCommunication, forkOneCommunication, ForkTwoDone, 0}
	philThree := Philosopher{"Three", false, forkThreeCommunication, forkTwoCommunication, ForkThreeDone, 0}
	philFour := Philosopher{"Four", false, forkFourCommunication, forkThreeCommunication, ForkFourDone, 0}
	PhilFive := Philosopher{"Five", false, forkFourCommunication, forkFiveCommunication, ForkFiveDone, 0} // Loking the other way

	go PhilFive.run()
	time.Sleep(time.Second * 3)
	go philOne.run()
	time.Sleep(time.Second)
	go philTwo.run()
	time.Sleep(time.Second)
	go philThree.run()
	time.Sleep(time.Second)
	go philFour.run()

	i := 1
	for i < 60 {
		time.Sleep(time.Second * 1)
		fmt.Print("time running ")
		fmt.Println(i)
		//fmt.Println(philOne.timesEaten, PhilFive.timesEaten)
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

/*
Fork One long switch
Used bool

Channel

*/
// type Fork struct{
// 	inUse bool
// 	forkChannel chan string
// 	inquiry string
// }
// func New(ForkChannel chan string ) Fork{
// 	f := Fork{false, ForkChannel,""}
// }
// func forkmain(givenChannel chan string) *Fork{
// 	fork = givenChannel
// 	var inquiry string

// for (true){
// 		inquiry = <-fork
// 		ForkLogic(inquirPrintln(a)
// }
// }

// func ForkLogic(fork Fork,inquiry string) {

// 	inquiry =<-fork.forkChannel
// 	switch inquiry {
// 	case "Request fork":
// 		fork.inUse = true
// 	case "Put down fork":
// 		fork.inUse = false
// 	case "":
// 		fmt.Println("Fork have done nothing")
// 	}
// }
