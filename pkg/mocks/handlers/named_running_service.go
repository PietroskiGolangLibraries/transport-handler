package mocked_transport_handlers

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type MockedNamedRunningServer struct {
	name string
	//ctx      context.Context
	//cancelFn context.CancelFunc
	conn chan int
	sig  chan os.Signal
	wg   sync.WaitGroup
}

func NewMockedNamedRunningServer(serverName string) *MockedNamedRunningServer {
	svr := &MockedNamedRunningServer{
		name: serverName,
		//ctx:      ctx,
		//cancelFn: cancelFn,
		conn: make(chan int),
		sig:  make(chan os.Signal),
		wg:   sync.WaitGroup{},
	}

	svr.Handle()

	return svr
}

func (svr *MockedNamedRunningServer) Handle() {
	signal.Notify(svr.sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	fmt.Println("handling creation for service of name:" + svr.name)
}

func (svr *MockedNamedRunningServer) Start() error {
	log.Printf("server of name '%v' is up and running\n", svr.name)
	//time.Sleep(time.Hour)
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
	<-svr.sig
	signal.Stop(svr.sig)
	log.Printf("server of name '%v' waiting for subprocess to stop\n", svr.name)
	svr.conn <- 0
	svr.wg.Wait()
	log.Printf("server of name '%v' subprocess stopped\n", svr.name)
	log.Printf("server of name '%v' closing sig channel\n", svr.name)
	close(svr.sig)
	log.Printf("server of name '%v' closed sig channel\n", svr.name)
}
