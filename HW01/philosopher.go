package main

import(
    "strings"
)

type philosopher struct {
	eating bool;
	priorityFork, nonPriorityFork chan string;
	timesEaten int;
}

func (p *philosopher) request(fork chan string) bool{
	var answer string
	p.priorityFork <- "Request fork"
	answer =<- p.priorityFork
	if (strings.Compare(answer,"Request accepted" )==0){
		return true
	} else { 
		return false;
	}

}

func (p *philosopher) run(){
	for p.timesEaten<3{
		if p.request(p.priorityFork){
			if p.request(p.nonPriorityFork){
				p.eating = !p.eating
				// wait
				p.timesEaten++
				p.eating = !p.eating
			}
		}
	}
}

/* 
Philo
Eating bool
priofork string
timesEaten int
Function request bool

"main"
while (true)
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