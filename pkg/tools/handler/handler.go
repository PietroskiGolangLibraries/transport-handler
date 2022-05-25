//go:build (darwin && cgo) || linux

package transporthandler

import (
	"context"
	"github.com/pkg/profile"
	"gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/models/handlers"
	tracer_models "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/models/tracer"
	stack_tracer "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/tools/tracer/stack"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type (
	Handler interface {
		StartServers(servers map[string]handlers_model.Server)
		Cancel()
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
	makeStopSrvSig = func() chan os.Signal {
		return make(chan os.Signal)
	}

	makeErrSrvSig = func(n int) chan error {
		return make(chan error, n)
	}
	makePanicSrvSig = func(n int) chan bool {
		return make(chan bool, n)
	}

	OsExit = os.Exit
)

func NewHandler(
	ctx context.Context,
	cancelFn context.CancelFunc,
	prof Profiler,
	exiter func(int),
) Handler {
	ctx, cancelFn = handleCtxGen(ctx, cancelFn)

	return &handler{
		ctx:      ctx,
		cancelFn: cancelFn,
		osExit:   exiter,

		profiler: profiler{
			pprof: prof,
		},

		goPool: goPool{
			wg:  &sync.WaitGroup{},
			gst: stack_tracer.NewGST(),
		},
	}
}

func NewDefaultHandler() Handler {
	ctx, cancelFn := handleCtxGen(nil, nil)

	return &handler{
		ctx:      ctx,
		cancelFn: cancelFn,
		osExit:   OsExit,

		profiler: profiler{
			pprof: profile.Start(
				profile.GoroutineProfile,
				profile.ProfilePath("./profiling/pprof"),
			),
		},

		goPool: goPool{
			wg:  &sync.WaitGroup{},
			gst: stack_tracer.NewGST(),
		},
	}
}

// StartServers starts all the variadic given servers and blocks the main thread.
func (h *handler) StartServers(servers map[string]handlers_model.Server) {
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
			log.Println("\nstop server sig!!")
			h.handleShutdown()
			h.osExit(0)
			return
		case <-h.srvChan.errSig:
			log.Println("\nerr server sig!!")
			h.handleShutdown()
			h.osExit(1)
			return
		case <-h.srvChan.panicSig:
			log.Println("\npanic server sig!!")
			h.handleShutdown()
			h.osExit(2)
			return
		case <-h.ctx.Done():
			log.Println("\nctx server sig!!")
			h.handleShutdown()
			h.osExit(0)
			return
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

func (h *handler) Cancel() {
	h.cancelFn()
}
