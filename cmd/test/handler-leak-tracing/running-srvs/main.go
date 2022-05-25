package main

import (
	"context"
	mocked_transport_handlers "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/mocks/handlers"
	"gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/mocks/profiling/pprof/fake"
	transporthandler "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/tools/handler"
	stack_tracer "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/tools/tracer/stack"
	"log"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"
)

var (
	wg = &sync.WaitGroup{}
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

	h := transporthandler.NewHandler(
		ctx,
		cancel,
		fakepprof.NewFakePProfProfiler(),
		func(i int) {
			log.Printf("exiting with code: %v", i)
		},
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		h.StartServers(svr1, svr2, svr3, svr4, svr5)
	}()

	time.Sleep(time.Second * 1)
	_ = syscall.Kill(os.Getpid(), syscall.SIGQUIT)
	wg.Wait()

	log.Printf("goroutine count is: %v", runtime.NumGoroutine())
}
