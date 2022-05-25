package main

import (
	"context"

	mocked_transport_handlers "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/mocks/handlers"
	"gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/mocks/os/exit/fake"
	fakepprof "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/mocks/profiling/pprof/fake"
	handlers_model "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/models/handlers"
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

	h := transporthandler.NewHandler(
		ctx,
		cancel,
		fakepprof.NewFakePProfProfiler(),
		fake.Exit,
	)
	h.StartServers(
		map[string]handlers_model.Server{
			"server-1":    svr1,
			"server-2":    svr2,
			"server-3":    svr3,
			"server-4":    svr4,
			"server-5":    svr5,
			"panic-srv-1": sPanic1,
			"panic-srv-2": sPanic2,
			"panic-srv-3": sPanic3,
			"panic-srv-4": sPanic4,
			"panic-srv-5": sPanic5,
		},
	)
}
