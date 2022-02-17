package server

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type HealthCounter struct {
	lock    sync.Mutex
	counter int
}

func (hc *HealthCounter) GetCounterValue() int {
	return hc.counter
}

func Start() {
	hc := &HealthCounter{counter: 0}
	hc.Counting()
}

func (hc *HealthCounter) KillTimer() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("KillTimer", c)
		os.Exit(2)
	}()
}

func (hc *HealthCounter) Counting() {
	hc.KillTimer()
	for {
		time.Sleep(1 * time.Second)
		hc.lock.Lock()
		hc.counter++
		hc.lock.Unlock()

		if hc.counter == 10000 {
			hc.counter = 0
		}
	}
}
