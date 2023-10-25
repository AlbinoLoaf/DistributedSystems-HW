// Inspired by the following link:
// https://medium.com/outreach-prague/lamport-clocks-determining-the-order-of-events-in-distributed-systems-41a9a8489177
package main

import (
	"sync"
)

type LamportClock struct {
	latestTime int
	mutex      sync.Mutex
}

func tick(l *LamportClock) {
	l.mutex.Lock()
	l.latestTime++
	l.mutex.Unlock()
}

func compareAndUpdate(recievedTimestamp int, l *LamportClock) {
	l.mutex.Lock()
	if l.latestTime < recievedTimestamp {
		l.latestTime = recievedTimestamp
	}
	l.latestTime++
	l.mutex.Unlock()
}

func getTime(l *LamportClock) (latestTime int) {
	return l.latestTime
}
