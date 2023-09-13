package main

import (
	"fmt"
	"time"
)

type Philosopher struct {
	name                                string
	eating                              bool
	priorityFork, nonPriorityFork, done chan string
	timesEaten                          int
}

func (p *Philosopher) run() {
	var bucket string
	for p.timesEaten < 4 {
		fmt.Printf("%s is thinking\n", p.name)
		fmt.Printf("%s Asking \n", p.name)
		p.priorityFork <- "I want to eat"
		p.nonPriorityFork <- "I want to eat"
		fmt.Printf("%s is eating\n", p.name)
		p.timesEaten++
		time.Sleep(time.Second * 2)
		p.done <- "done"
		bucket = <-p.priorityFork
		bucket = <-p.nonPriorityFork
	}
	fmt.Println(bucket)

}

// func (p *Philosopher) request(fork chan string) bool {
// 	var answer string
// 	p.priorityFork <- "Request fork"
// 	answer = <-p.priorityFork
// 	fmt.Print("Philosophers Answer from the fork: ")
// 	fmt.Println(answer)
// 	if strings.Compare(answer, "Request Accepted") == 0 {
// 		return true
// 	} else {
// 		return false
// 	}

// }

// func (p *Philosopher) run() {
// 	for p.timesEaten < 4 {
// 		if p.request(p.priorityFork) {
// 			if p.request(p.nonPriorityFork) {
// 				p.eating = !p.eating
// 				fmt.Printf("%s Eating\n", p.name)
// 				time.Sleep(time.Second * 2)
// 				p.priorityFork <- "Put down fork"
// 				p.nonPriorityFork <- "Put down fork"
// 				// wait
// 				p.timesEaten++
// 				p.eating = !p.eating
// 				fmt.Printf("%s Done eating\n", p.name)
// 			} else {
// 				p.priorityFork <- "Put down fork"
// 				time.Sleep(time.Second)
// 			}
// 		}
// 	}
// }

/*
Philo
Eating bool
priofork string
timesEaten int
Function request bool

"main"
while (truPrint(a)
	if requist(prio)
		if request(nonPrio)
			eat flip
			wait time
			eat flip
	wait


task
	print timesEaten
	if eating
		print phil {ID} eating
	else pring phil {id} thinking

*/
