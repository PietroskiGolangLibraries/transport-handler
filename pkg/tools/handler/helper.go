package transporthandler

import (
	"context"
	handlers_model "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v3/pkg/models/handlers"
	"log"
	"os/signal"
	"syscall"
)

func (h *handler) stopServers() {
	for srvName, srv := range h.serverMapping {
		if srv == nil {
			continue
		}

		h.goPool.wg.Add(1)
		go func(srvName string, srv handlers_model.Server) {
			defer h.handleStopPanic()
			defer h.goPool.wg.Done()

			h.stopCall(srvName, srv)
		}(srvName, srv)
	}
}

func (h *handler) stopCall(srvName string, srv handlers_model.Server) {
	log.Printf("stopping server: %v\n", srvName)
	srv.Stop()
	log.Printf("stopped server: %v\n", srvName)
}

func (h *handler) makeSrvChan(srvLen int) {
	h.srvChan = srvChan{
		stopSig:  makeStopSrvSig(),
		errSig:   makeErrSrvSig(srvLen),
		panicSig: makePanicSrvSig(srvLen),
	}
}

func (h *handler) closeSrvSigChan() {
	defer h.handleCloseChanPanic()

	log.Println("closing srv sig channel...")
	if !isChanClosed(h.srvChan.stopSig) {
		signal.Stop(h.srvChan.stopSig)
		close(h.srvChan.stopSig)
	}
	log.Println("srv sig channel successfully closed.")
}

func (h *handler) closeErrChan() {
	defer h.handleCloseChanPanic()

	log.Println("closing error channel...")
	if !isChanClosed(h.srvChan.errSig) {
		close(h.srvChan.errSig)
	}
	log.Println("error channel successfully closed.")
}

func (h *handler) closePanicChan() {
	defer h.handleCloseChanPanic()

	log.Println("closing panic channel...")
	if !isChanClosed(h.srvChan.panicSig) {
		close(h.srvChan.panicSig)
	}
	log.Println("panic channel successfully closed.")
}

func (h *handler) handleCloseChanPanic() {
	if r := recover(); r != nil {
		log.Printf("recovering from close channel panic: %v", r)
		h.goPool.gst.Trace()
		h.osExit(2)
		return
	}
}

func (h *handler) handleWaiting() {
	log.Println("waiting for goroutines to stop")
	h.goPool.wg.Wait()
	log.Println("goroutines stopped")
}

func (h *handler) handleErr(err error) {
	log.Printf("error from server: %v", err)
	if !isChanClosed(h.srvChan.errSig) {
		h.srvChan.errSig <- err
		log.Println("post-err")
	}
}

func (h *handler) handlePanic() {
	if r := recover(); r != nil {
		log.Printf("recovering from runtime panic: %v", r)
		if !isChanClosed(h.srvChan.panicSig) {
			h.srvChan.panicSig <- true
			log.Println("post-panic")
		}
	}
}

func (h *handler) handleStopPanic() {
	if r := recover(); r != nil {
		log.Printf("recovering from stopping runtime panic: %v", r)
	}
}

func (h *handler) sigKill() {
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGQUIT)
}

func handleCtxGen(ctx context.Context, cancelFn context.CancelFunc) (context.Context, context.CancelFunc) {
	if ctx == nil && cancelFn == nil {
		ctx, cancelFn = context.WithCancel(context.Background())
	} else if ctx != nil && cancelFn == nil {
		ctx, cancelFn = context.WithCancel(ctx)
	}

	return ctx, cancelFn
}

func isChanClosed[T any](ch chan T) bool {
	log.Println("analysing channel")
	select {
	case _, ok := <-ch:
		log.Println("channel is:", ok)
		return !ok
	default:
		return false
	}
}
