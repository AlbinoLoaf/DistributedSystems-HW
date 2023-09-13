package main

import (
	"fmt"
	"time"
)

type Philosopher struct {
	name                          string // Name is to differ between philosopher
	eating                        bool
	priorityFork, nonPriorityFork chan string // Determines preffered fork
	prioReq, nonPrioReq           chan string // Telling the fork to announce it is being used
	timesEaten                    int // Counter of times eaten
}

func (p *Philosopher) run() {
	var bucket string// Bucket used to empty channel contents
	for p.timesEaten < 4 {// For loop ensuring every philosopher eats 3 times
		fmt.Printf("%s is thinking\n", p.name)
		select {
		case p.priorityFork <- "I want to eat":// Case: a philosopher attempts to eat if forks are available
			p.prioReq <- p.name
			p.nonPrioReq <- p.name
			p.nonPriorityFork <- "%s is eating\n"
			fmt.Printf("%s is eating\n", p.name)
			p.timesEaten++
			time.Sleep(time.Millisecond * 500)
			fmt.Printf("%s is done eating\n", p.name)
			bucket = <-p.priorityFork
			bucket = <-p.nonPriorityFork
			time.Sleep(time.Millisecond * 500)
		default:// Default: if case is unreachable a philosopher will think and try again after half a second
			time.Sleep(time.Millisecond * 500)
		}

	}
	bucket = "%s is completly full\n"
	fmt.Printf(bucket, p.name)

}
