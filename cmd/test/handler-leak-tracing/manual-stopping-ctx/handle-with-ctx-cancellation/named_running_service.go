package mocked_transport_handlers

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
)

type MockedNamedRunningServer struct {
	name string
	ctx  context.Context
	conn chan int
	sig  chan os.Signal
	wg   sync.WaitGroup
}

func NewMockedNamedRunningServer(ctx context.Context, serverName string) *MockedNamedRunningServer {
	svr := &MockedNamedRunningServer{
		name: serverName,
		ctx:  ctx,
		conn: make(chan int),
		sig:  make(chan os.Signal),
		wg:   sync.WaitGroup{},
	}

	svr.Handle()

	return svr
}

func (svr *MockedNamedRunningServer) Handle() {
	fmt.Println("handling creation for service of name:" + svr.name)
}

func (svr *MockedNamedRunningServer) Start() error {
	log.Printf("server of name '%v' is up and running\n", svr.name)
	svr.wg.Add(1)
	go func() {
		defer svr.wg.Done()
		log.Printf("server of name '%v' is up and running in a subprocess\n", svr.name)
		<-svr.conn
		log.Printf("server of name '%v' is stopping a subprocess\n", svr.name)
		return
	}()

	svr.Stop()
	return nil
}

func (svr *MockedNamedRunningServer) Stop() {
	log.Println("CTX HERE!!", &svr.ctx)
	<-svr.ctx.Done()
	log.Println("CANCELLING IT HERE!!")
	log.Printf("server of name '%v' waiting for subprocess to stop\n", svr.name)
	svr.conn <- 0
	svr.wg.Wait()
	close(svr.conn)
	log.Printf("server of name '%v' subprocess stopped\n", svr.name)
	log.Printf("server of name '%v' closing sig channel\n", svr.name)
	log.Printf("server of name '%v' closed sig channel\n", svr.name)
}
