package main

import (
	"context"
	"os"

	mocked_transport_handlers "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/cmd/test/handler-leak-tracing/manual-stopping-ctx/handle-with-ctx-cancellation"
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

	svr1 := mocked_transport_handlers.NewMockedNamedRunningServer(ctx, "server-1")
	svr2 := mocked_transport_handlers.NewMockedNamedRunningServer(ctx, "server-2")
	svr3 := mocked_transport_handlers.NewMockedNamedRunningServer(ctx, "server-3")
	svr4 := mocked_transport_handlers.NewMockedNamedRunningServer(ctx, "server-4")
	svr5 := mocked_transport_handlers.NewMockedNamedRunningServer(ctx, "server-5")

	h := transporthandler.NewHandler(ctx, cancel, fakepprof.NewFakePProfProfiler(), os.Exit)
	h.StartServers(
		map[string]handlers_model.Server{
			"server-1": svr1,
			"server-2": svr2,
			"server-3": svr3,
			"server-4": svr4,
			"server-5": svr5,
		},
	)
}
