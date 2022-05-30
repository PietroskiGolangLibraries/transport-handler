package mocked_transport_handlers

import (
	"fmt"
	"log"
)

type MockedErrServer struct{}

func (ms *MockedErrServer) Handle() {}

func (ms *MockedErrServer) Start() error {
	return fmt.Errorf("\nerror to start the server\n")
}

func (ms *MockedErrServer) Stop() {
	log.Println("stop called")
}
