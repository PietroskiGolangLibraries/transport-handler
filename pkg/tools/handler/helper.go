package transporthandler

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"unsafe"
)

func (h *handler) makeSrvChan(srvLen int) {
	h.srvChan = srvChan{
		stopSig:  privateStopSrvSig,
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

func (h *handler) sigKill() {
	err := syscall.Kill(syscall.Getpid(), syscall.SIGQUIT)
	if err != nil {
		log.Printf("failed to send a quit seignal into the system: %v", err)
	}
}

func handleCtxGen(ctx context.Context, cancelFn context.CancelFunc) (context.Context, context.CancelFunc) {
	if ctx == nil || cancelFn == nil {
		ctx, cancelFn = context.WithCancel(context.Background())
	} else if ctx != nil || cancelFn == nil {
		ctx, cancelFn = context.WithCancel(ctx)
	}

	return ctx, cancelFn
}

func isChanClosed(ch interface{}) bool {
	//if reflect.TypeOf(ch).Kind() != reflect.Chan {
	//	panic("only channels!")
	//}

	// get interface value pointer, from cgo_export
	// typedef struct { void *t; void *v; } GoInterface;
	// then get channel real pointer
	cptr := *(*uintptr)(unsafe.Pointer(
		unsafe.Pointer(uintptr(unsafe.Pointer(&ch)) + unsafe.Sizeof(uint(0))),
	))

	// this function will return true if chan.closed > 0
	// see hchan on https://github.com/golang/go/blob/master/src/runtime/chan.go
	// type hchan struct {
	// qcount   uint           // total data in the queue
	// dataqsiz uint           // size of the circular queue
	// buf      unsafe.Pointer // points to an array of dataqsiz elements
	// elemsize uint16
	// closed   uint32
	// **

	cptr += unsafe.Sizeof(uint(0)) * 2
	cptr += unsafe.Sizeof(unsafe.Pointer(uintptr(0)))
	cptr += unsafe.Sizeof(uint16(0))
	return *(*uint32)(unsafe.Pointer(cptr)) > 0
}
