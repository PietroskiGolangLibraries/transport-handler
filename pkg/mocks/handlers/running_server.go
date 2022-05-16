package mocked_transport_handlers

import (
	"log"
	"os"
)

type MockedRunningServer struct{}

func (ms *MockedRunningServer) Handle() {}

func (ms *MockedRunningServer) Start() error {
	log.Println("server is up and running")
	os.Exit(0)
	return nil
}
