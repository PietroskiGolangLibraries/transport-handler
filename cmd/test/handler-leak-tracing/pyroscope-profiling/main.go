package main

import (
	"context"

	pprof_server "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/cmd/test/handler-leak-tracing/parca-profiling/pprof-server"
	pyroscope_server "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/cmd/test/handler-leak-tracing/pyroscope-profiling/pyroscope-server"
	mocked_transport_handlers "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/mocks/handlers"
	handlers_model "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/models/handlers"
	transporthandler "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/tools/handler"
	stack_tracer "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/tools/tracer/stack"
)

func main() {
	st := stack_tracer.NewGoroutineStackTracer()
	st.Trace()
	defer st.Trace()

	ctx, _ := context.WithCancel(context.Background())

	svr1 := mocked_transport_handlers.NewMockedNamedRunningServer("server-1")
	svr2 := mocked_transport_handlers.NewMockedNamedRunningServer("server-2")
	svr3 := mocked_transport_handlers.NewMockedNamedRunningServer("server-3")
	svr4 := mocked_transport_handlers.NewMockedNamedRunningServer("server-4")
	svr5 := mocked_transport_handlers.NewMockedNamedRunningServer("server-5")

	pprofSvr := pprof_server.NewPProfServer(ctx)
	pyroscopeSvr := pyroscope_server.NewPyroscopeServer()

	h := transporthandler.NewDefaultHandler()
	h.StartServers(
		map[string]handlers_model.Server{
			"server-1": svr1,
			"server-2": svr2,
			"server-3": svr3,
			"server-4": svr4,
			"server-5": svr5,

			"pprof":     pprofSvr,
			"pyroscope": pyroscopeSvr,
		},
	)
}
