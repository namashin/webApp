package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
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
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Println("killTimer", sig)
			fmt.Println(hc.GetCounterValue())
			os.Exit(1)
		}
	}()
}

func (hc *HealthCounter) Counting() {
	hc.lock.Lock()
	defer hc.lock.Unlock()
	hc.KillTimer()
	for {
		time.Sleep(1 * time.Second)
		hc.counter++
		if hc.counter == 10000 {
			hc.counter = 0
		}
	}
}
