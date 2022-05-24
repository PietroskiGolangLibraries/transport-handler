package mocked_transport_handlers

import "log"

type MockedPanicServer struct{}

func (ms *MockedPanicServer) Handle() {}

func (ms *MockedPanicServer) Start() error {
	log.Panic("panic starting the server")
	return nil
}
