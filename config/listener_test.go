package config

import (
	"log"
	"sync"
	"testing"
)

func TestListener(t *testing.T) {
	wd := sync.WaitGroup{}
	wd.Add(1)
	f := func() {
		log.Println("receive reload signal!")
		wd.Done()
	}
	AddReloadListener(f)
	publishReloadSignal()
	wd.Wait()
}
