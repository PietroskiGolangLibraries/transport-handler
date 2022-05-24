//go:build (darwin && cgo) || linux

package transporthandler

import (
	"context"
	"fmt"
	"gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/models/handlers"
	tracer_models "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/tools/tracer/models"
	stack_tracer "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/tools/tracer/stack"
	"gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/tools/worker"
	"log"
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

	handler struct {
		ctx      context.Context
		cancelFn context.CancelFunc

		stopSrvSig  chan os.Signal
		errSrvSig   chan error
		panicSrvSig chan bool
		osExit      func(code int)

		sigChanMonitor   *worker.ChanMonitor
		errChanMonitor   *worker.ChanMonitor
		panicChanMonitor *worker.ChanMonitor

		goPool goPool
	}

	goPool struct {
		control int
		goCount int64
		wg      *sync.WaitGroup
		mtx     *sync.Mutex
		gst     tracer_models.Tracer
	}
)

var (
	privateStopSrvSig  = make(chan os.Signal)
	privateErrSrvSig   = make(chan error)
	privatePanicSrvSig = make(chan bool)
	OsExit             = os.Exit
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

		stopSrvSig:  privateStopSrvSig,
		errSrvSig:   privateErrSrvSig,
		panicSrvSig: privatePanicSrvSig,
		osExit:      OsExit,

		sigChanMonitor:   worker.NewChanMonitor(),
		errChanMonitor:   worker.NewChanMonitor(),
		panicChanMonitor: worker.NewChanMonitor(),

		goPool: goPool{
			goCount: 0,
			wg:      &sync.WaitGroup{},
			mtx:     &sync.Mutex{},
			gst:     stack_tracer.NewGST(),
		},
	}
}

// StartServers starts all the variadic given servers and blocks the main thread.
func (h *handler) StartServers(servers ...handlers_model.Server) {
	signal.Notify(h.stopSrvSig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	log.Println("Goroutine initial count ->", h.goPool.goCount)
	for _, s := range servers {
		h.goPool.wg.Add(1)
		go func(s handlers_model.Server) {
			defer h.handlePanic()
			defer h.goPool.wg.Done()

			if s == nil {
				return
			}
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
		case <-h.stopSrvSig:
			fmt.Println("\nstop server sig!!")
			h.handleShutdown()
			h.osExit(0)
		case <-h.errSrvSig:
			fmt.Println("\nerr server sig!!")
			h.handleShutdown()
			h.osExit(1)
		case <-h.panicSrvSig:
			fmt.Println("\npanic server sig!!")
			h.handleShutdown()
			h.osExit(2)
		}
	}
}

func (h *handler) handleShutdown() {
	h.sigKill()
	h.handleWaiting()
	h.closeSrvSigChan()
	time.Sleep(time.Millisecond * 500)
	h.closeErrChan()
	time.Sleep(time.Millisecond * 500)
	h.closePanicChan()
	h.goPool.gst.Trace()
}
