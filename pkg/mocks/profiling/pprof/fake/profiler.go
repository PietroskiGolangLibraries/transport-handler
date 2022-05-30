//go:build (darwin && cgo) || linux

package fakepprof

import (
	"log"

	prof "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v3/pkg/models/profiler"
)

type (
	fakepprof struct{}
)

func NewFakePProfProfiler() prof.Profiler {
	return &fakepprof{}
}

func (ffpp *fakepprof) Stop() {
	log.Printf("stoping fake profiling")
}
