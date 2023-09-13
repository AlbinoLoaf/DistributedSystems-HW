package main

import (
	"fmt"
)

type Fork struct {
	name    string
	request chan string
}

func (f *Fork) run() {
	// var x string;
	for true {
		x := <-f.request
		fmt.Printf("%s is being used by %s \n", f.name, x)

	}
}
