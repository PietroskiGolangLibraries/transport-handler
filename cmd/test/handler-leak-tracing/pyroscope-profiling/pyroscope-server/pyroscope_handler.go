package pyroscope_server

import (
	"context"
	"log"
	"net/http"

	"github.com/pyroscope-io/client/pyroscope"
)

const (
	ctxErr = "context-error"
)

type (
	PyroscopeServer struct {
		ctx      context.Context
		profiler *pyroscope.Profiler
	}
)

func NewPyroscopeServer() *PyroscopeServer {
	s := &PyroscopeServer{
		ctx: context.Background(),
	}

	s.Handle()

	return s
}

func (svr *PyroscopeServer) Handle() {
	pyroscopeProfiler, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: "pietroski.pyroscope.profiling.demo",

		// replace this with the address of pyroscope server
		//ServerAddress: "http://pyroscope-server:4040",
		//ServerAddress: "http://localhost:4040",
		ServerAddress: "http://pyroscope.pyroscope-demo.svc.cluster.local:4040",

		Tags: map[string]string{
			"pietroski": "pietroski.pyroscope.profiling.test.tag.value",
		},

		// you can disable logging by setting this to nil
		Logger: pyroscope.StandardLogger,

		// optionally, if authentication is enabled, specify the API key:
		// AuthToken: os.Getenv("PYROSCOPE_AUTH_TOKEN"),

		// by default all profilers are enabled,
		// but you can select the ones you want to use:
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
		},
	})
	if err != nil {
		svr.ctx = context.WithValue(svr.ctx, ctxErr, err)
	}

	svr.profiler = pyroscopeProfiler
}

func (svr *PyroscopeServer) Start() error {
	if err, ok := svr.ctx.Value(ctxErr).(error); ok && err != nil {
		return err
	}

	return nil
}

func (svr *PyroscopeServer) Stop() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovering from panic")
		}
	}()

	if err := svr.profiler.Stop(); err != nil && err != http.ErrServerClosed {
		log.Printf("HTTP server ListenAndServe shutdown error: %v\n", err)

		return
	}

	log.Printf("HTTP server ListenAndServe shutdown ok")
}
