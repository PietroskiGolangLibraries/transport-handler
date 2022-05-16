package mocked_transport_handlers

import (
	"log"
)

type MockedRunningNoExitServer struct{}

func (ms *MockedRunningNoExitServer) Handle() {}

func (ms *MockedRunningNoExitServer) Start() error {
	log.Println("server is up and running")
	return nil
}
