package main

import (
	"fmt"
	"time"
)

type Philosopher struct {
	name                                string
	eating                              bool
	priorityFork, nonPriorityFork chan string
	timesEaten                          int
}

func (p *Philosopher) run() {
	var bucket string
	for p.timesEaten < 4 {
		fmt.Printf("%s is thinking\n", p.name)
		select {
		case p.priorityFork <- "I want to eat":
			p.nonPriorityFork <- "%s is eating\n"
			fmt.Printf("%s is eating\n", p.name)
			p.timesEaten++
			time.Sleep(time.Millisecond * 500)
			fmt.Printf("%s is done eating\n", p.name)
			bucket = <-p.priorityFork
			bucket = <-p.nonPriorityFork
			time.Sleep(time.Millisecond * 500)
		default:
			time.Sleep(time.Millisecond * 500)
		}

	}
	bucket = "%s is completly full\n"
	fmt.Printf(bucket, p.name)

}

