package main

import (
	"log"
	"sync"
	"syscall"
	"testing"
	"time"
)

var (
	wg sync.WaitGroup
)

func TestMain(t *testing.M) {
	wg.Add(1)
	go main()
	wg.Done()
	time.Sleep(time.Second * 2)
	log.Println("sending quit signal")
	syscall.Kill(syscall.Getpid(), syscall.SIGQUIT)
	wg.Wait()
}
