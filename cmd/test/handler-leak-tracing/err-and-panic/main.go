package main

import (
	"context"
	mocked_transport_handlers "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/mocks/handlers"
	transporthandler "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/tools/handler"
	stack_tracer "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/tools/tracer/stack"
)

func main() {
	st := stack_tracer.NewGoroutineStackTracer()
	st.Trace()
	defer st.Trace()

	ctx, cancel := context.WithCancel(context.Background())

	svr1 := mocked_transport_handlers.NewMockedNamedRunningServer("server-1")
	svr2 := mocked_transport_handlers.NewMockedNamedRunningServer("server-2")
	svr3 := mocked_transport_handlers.NewMockedNamedRunningServer("server-3")
	svr4 := mocked_transport_handlers.NewMockedNamedRunningServer("server-4")
	svr5 := mocked_transport_handlers.NewMockedNamedRunningServer("server-5")
	sPanic1 := &mocked_transport_handlers.MockedPanicServer{}
	sPanic2 := &mocked_transport_handlers.MockedPanicServer{}
	sPanic3 := &mocked_transport_handlers.MockedPanicServer{}
	sPanic4 := &mocked_transport_handlers.MockedPanicServer{}
	sPanic5 := &mocked_transport_handlers.MockedPanicServer{}
	sErr1 := &mocked_transport_handlers.MockedErrServer{}
	sErr2 := &mocked_transport_handlers.MockedErrServer{}
	sErr3 := &mocked_transport_handlers.MockedErrServer{}
	sErr4 := &mocked_transport_handlers.MockedErrServer{}
	sErr5 := &mocked_transport_handlers.MockedErrServer{}

	h := transporthandler.NewDefaultHandler(ctx, cancel)
	h.StartServers(
		svr1, svr2, svr3, svr4, svr5,
		sErr1, sErr2, sErr3, sErr4, sErr5,
		sPanic1, sPanic2, sPanic3, sPanic4, sPanic5,
		nil, nil, nil, nil, nil,
	)
}
