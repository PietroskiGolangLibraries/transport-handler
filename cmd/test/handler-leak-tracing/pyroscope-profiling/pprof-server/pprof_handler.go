package pprof_server

import (
	"context"
	"log"
	"net/http"
	"net/http/pprof"
	"time"
)

type (
	PProfServer struct {
		server *http.Server

		name string
		ctx  context.Context
	}
)

func NewPProfServer(ctx context.Context) *PProfServer {
	s := &PProfServer{
		ctx: ctx,
	}

	s.Handle()

	return s
}

func (svr *PProfServer) Handle() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
	}

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/readiness", readinessHandler)

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/vars", http.DefaultServeMux)

	svr.server = server
}

func (svr *PProfServer) Start() error {
	return svr.server.ListenAndServe()
}

func (svr *PProfServer) Stop() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovering from panic")
		}
	}()

	if err := svr.server.Shutdown(svr.ctx); err != nil && err != http.ErrServerClosed {
		log.Printf("HTTP server ListenAndServe shutdown error: %v\n", err)

		return
	}

	log.Printf("HTTP server ListenAndServe shutdown ok")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
