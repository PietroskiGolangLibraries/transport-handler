package stack_tracer

import (
	"fmt"
	tracer_models "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/models/tracer"
	"runtime"
)

type goroutineStackTracer struct{}

func NewGoroutineStackTracer() tracer_models.Tracer {
	tracer := &goroutineStackTracer{}
	return tracer
}

func NewGST() tracer_models.Tracer {
	tracer := &goroutineStackTracer{}
	tracer.Trace()
	return tracer
}

func (goroutineStackTracer) Trace() {
	fmt.Printf("running goroutine number: %v\n", runtime.NumGoroutine())
}
