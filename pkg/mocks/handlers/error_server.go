package mocked_transport_handlers

import (
	"fmt"
)

type MockedErrServer struct{}

func (ms *MockedErrServer) Handle() {}

func (ms *MockedErrServer) Start() error {
	return fmt.Errorf("\nerror to start the server\n")
}
