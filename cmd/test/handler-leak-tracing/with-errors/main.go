package main

import (
	"context"

	mocked_transport_handlers "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/mocks/handlers"
	"gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/mocks/os/exit/fake"
	fakepprof "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/mocks/profiling/pprof/fake"
	handlers_model "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/models/handlers"
	transporthandler "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/tools/handler"
	stack_tracer "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/tools/tracer/stack"
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
	sErr1 := &mocked_transport_handlers.MockedErrServer{}
	sErr2 := &mocked_transport_handlers.MockedErrServer{}
	sErr3 := &mocked_transport_handlers.MockedErrServer{}
	sErr4 := &mocked_transport_handlers.MockedErrServer{}
	sErr5 := &mocked_transport_handlers.MockedErrServer{}

	h := transporthandler.NewHandler(
		ctx,
		cancel,
		fakepprof.NewFakePProfProfiler(),
		fake.Exit,
	)
	h.StartServers(
		map[string]handlers_model.Server{
			"server-1":  svr1,
			"server-2":  svr2,
			"server-3":  svr3,
			"server-4":  svr4,
			"server-5":  svr5,
			"err-srv-1": sErr1,
			"err-srv-2": sErr2,
			"err-srv-3": sErr3,
			"err-srv-4": sErr4,
			"err-srv-5": sErr5,
		},
	)
}
