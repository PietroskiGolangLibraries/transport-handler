package transporthandler

import (
	"log"
	"os"
	"strings"

	"gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/models/handlers"
)

type (
	Handler interface {
		StartServers(servers ...handlers_model.Server)

		handleError()
		handlePanic()
		closeChan()
		handleCloseChanPanic()
		verifyCodeZero(r interface{})
	}

	handler struct {
		stopServerSig chan error
		osExit        func(code int)
	}
)

var (
	OsExit               = os.Exit
	privateStopServerSig = make(chan error)
)

func NewHandler(stopServerSig chan error, exiter func(int)) Handler {
	if stopServerSig == nil {
		stopServerSig = privateStopServerSig
	}
	if exiter == nil {
		exiter = OsExit
	}

	return &handler{
		stopServerSig: stopServerSig,
		osExit:        exiter,
	}
}

// StartServers starts all the variadic given servers and blocks the main thread.
func (h *handler) StartServers(servers ...handlers_model.Server) {
	for _, s := range servers {
		go func(stopServerSig chan error, s handlers_model.Server) {
			defer h.handlePanic()
			if err := s.Start(); err != nil {
				stopServerSig <- err
			}
		}(h.stopServerSig, s)
	}

	h.handleError()
}

func (h *handler) handleError() {
	for {
		select {
		case err := <-h.stopServerSig:
			if err != nil {
				h.closeChan()
				log.Println(err)
				h.osExit(1)
				return
			}

			h.closeChan()
			return
		}
	}
}

func (h *handler) handlePanic() {
	if r := recover(); r != nil {
		h.verifyCodeZero(r)
		log.Printf("recovering from panic: %v", r)
		//h.closeChan()
		h.osExit(2)
		h.stopServerSig <- nil
	}
}

func (h *handler) closeChan() {
	defer h.handleCloseChanPanic()

	log.Println("closing channel...")
	close(h.stopServerSig)
	log.Println("channel successfully closed.")
}

func (h *handler) handleCloseChanPanic() {
	if r := recover(); r != nil {
		h.verifyCodeZero(r)
		log.Printf("recovering from close channel panic: %v", r)
		h.osExit(2)
	}
}

func (h *handler) verifyCodeZero(r interface{}) {
	str, ok := r.(string)
	if ok && strings.Contains(str, "os.Exit(0)") {
		h.osExit(0)
	}
}
