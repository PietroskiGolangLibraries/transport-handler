package transporthandler

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
)

func (h *handler) closeSrvSigChan() {
	defer h.handleCloseChanPanic()

	log.Println("closing srv sig channel...")
	if !isChanClosed(h.stopSrvSig) {
		signal.Stop(h.stopSrvSig)
		close(h.stopSrvSig)
	}
	log.Println("srv sig channel successfully closed.")
}

func (h *handler) closeErrChan() {
	defer h.handleCloseChanPanic()

	log.Println("closing error channel...")
	if !isChanClosed(h.errSrvSig) {
		close(h.errSrvSig)
	}
	log.Println("error channel successfully closed.")
}

func (h *handler) closePanicChan() {
	defer h.handleCloseChanPanic()

	log.Println("closing panic channel...")
	if !isChanClosed(h.panicSrvSig) {
		close(h.panicSrvSig)
	}
	log.Println("panic channel successfully closed.")
}

func (h *handler) handleCloseChanPanic() {
	if r := recover(); r != nil {
		log.Printf("recovering from close channel panic: %v", r)
		h.goPool.gst.Trace()
		h.osExit(2)
	}
}

func (h *handler) handleWaiting() {
	log.Println("waiting for goroutines to stop")
	h.goPool.wg.Wait()
	log.Println("goroutines stopped")
}

func (h *handler) handleErr(err error) {
	log.Printf("error from server: %v", err)
	if !isChanClosed(h.errSrvSig) {
		//h.errChanMonitor.CountElem()
		h.errSrvSig <- err
	}
}

func (h *handler) handlePanic() {
	if r := recover(); r != nil {
		err := fmt.Errorf("recovering from runtime panic: %v", r)
		log.Println(err)
		//log.Printf("recovering from runtime panic: %v", r)
		if !isChanClosed(h.panicSrvSig) {
			//h.panicChanMonitor.CountElem()
			h.panicSrvSig <- true
			log.Println("post-panic")
		}
	}
}

func (h *handler) sigKill() {
	err := syscall.Kill(syscall.Getpid(), syscall.SIGQUIT)
	if err != nil {
		log.Printf("failed to send a quit seignal into the system: %v", err)
	}
	//h.sigChanMonitor.CountElem()
}

func isChanClosed[T any](ch chan T) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func handleCtxGen(ctx context.Context, cancelFn context.CancelFunc) (context.Context, context.CancelFunc) {
	if ctx == nil || cancelFn == nil {
		ctx, cancelFn = context.WithCancel(context.Background())
	} else if ctx != nil || cancelFn == nil {
		ctx, cancelFn = context.WithCancel(ctx)
	}

	return ctx, cancelFn
}
