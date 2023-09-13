package main

import (
	"fmt"
)

func main() {
	fork1 := make(chan string) // Fork between philosofer 1 and 2
	fork2 := make(chan string) // Fork between philosofer 2 and 3
	fork3 := make(chan string) // Fork between philosofer 3 and 4
	fork4 := make(chan string) // Fork between philosofer 4 and 5 
	fork5 := make(chan string) // Fork between philosofer 5 and 1 
	fmt.Println(fork1, fork2, fork3, fork4, fork5)
	//go person{fork1, fork2}
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
