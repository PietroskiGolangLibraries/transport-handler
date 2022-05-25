//go:build (darwin && cgo) || linux

package transporthandler

import (
	"context"
	"fmt"
	"github.com/pkg/profile"
	"gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/models/handlers"
	tracer_models "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/tools/tracer/models"
	stack_tracer "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/tools/tracer/stack"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type (
	Handler interface {
		StartServers(servers ...handlers_model.Server)
	}

	Profiler interface {
		Stop()
	}

	profiler struct {
		pprof Profiler
	}

	goPool struct {
		wg  *sync.WaitGroup
		gst tracer_models.Tracer
	}

	srvChan struct {
		stopSig  chan os.Signal
		errSig   chan error
		panicSig chan bool
	}

	handler struct {
		ctx      context.Context
		cancelFn context.CancelFunc
		osExit   func(code int)

		goPool   goPool
		srvChan  srvChan
		profiler profiler
	}
)

var (
	privateStopSrvSig  = make(chan os.Signal)
	privateErrSrvSig   = make(chan error)
	privatePanicSrvSig = make(chan bool)

	makeErrSrvSig = func(n int) chan error {
		return make(chan error, n)
	}
	makePanicSrvSig = func(n int) chan bool {
		return make(chan bool, n)
	}

	OsExit = os.Exit
)

//func NewHandler(
//	ctx context.Context,
//	cancelFn context.CancelFunc,
//	stopServerSig chan os.Signal,
//	stopServerErrSig chan error,
//	exiter func(int),
//) Handler {
//	ctx, cancelFn = handleCtxGen(ctx, cancelFn)
//	stopServerSig, stopServerErrSig = handleStopChanGen(stopServerSig, stopServerErrSig)
//
//	if exiter == nil {
//		exiter = OsExit
//	}
//
//	return &handler{
//		ctx:              ctx,
//		cancelFn:         cancelFn,
//		stopServerSig:    stopServerSig,
//		stopServerErrSig: stopServerErrSig,
//		osExit:           exiter,
//		sysNotifier:      nil,
//		chanMonitor:      worker.NewChanMonitor(),
//		errChanMonitor:   worker.NewChanMonitor(),
//		goPool: goPool{
//			goCount: 0,
//			wg:      &sync.WaitGroup{},
//			mtx:     &sync.Mutex{},
//			gst:     stack_tracer.NewGST(),
//		},
//	}
//}

func NewDefaultHandler(
	ctx context.Context,
	cancelFn context.CancelFunc,
) Handler {
	ctx, cancelFn = handleCtxGen(ctx, cancelFn)

	return &handler{
		ctx:      ctx,
		cancelFn: cancelFn,
		osExit:   OsExit,

		profiler: profiler{
			pprof: profile.Start(profile.GoroutineProfile, profile.ProfilePath("./pprof")),
		},

		goPool: goPool{
			wg:  &sync.WaitGroup{},
			gst: stack_tracer.NewGST(),
		},
	}
}

// StartServers starts all the variadic given servers and blocks the main thread.
func (h *handler) StartServers(servers ...handlers_model.Server) {
	h.makeSrvChan(len(servers))
	signal.Notify(h.srvChan.stopSig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	for _, s := range servers {
		h.goPool.wg.Add(1)
		go func(s handlers_model.Server) {
			defer h.handlePanic()
			defer h.goPool.wg.Done()

			if err := s.Start(); err != nil {
				h.handleErr(err)
			}
		}(s)
		time.Sleep(time.Millisecond * 125)
	}

	h.handleServer()
}

func (h *handler) handleServer() {
	for {
		select {
		case <-h.srvChan.stopSig:
			fmt.Println("\nstop server sig!!")
			h.handleShutdown()
			h.osExit(0)
		case <-h.srvChan.errSig:
			fmt.Println("\nerr server sig!!")
			h.handleShutdown()
			h.osExit(1)
		case <-h.srvChan.panicSig:
			fmt.Println("\npanic server sig!!")
			h.handleShutdown()
			h.osExit(2)
		}
	}
}

func (h *handler) handleShutdown() {
	h.sigKill()
	h.handleWaiting()
	time.Sleep(time.Millisecond * 250)
	h.closeSrvSigChan()
	time.Sleep(time.Millisecond * 250)
	h.closeErrChan()
	time.Sleep(time.Millisecond * 250)
	h.closePanicChan()
	h.goPool.gst.Trace()
	h.profiler.pprof.Stop()
}
