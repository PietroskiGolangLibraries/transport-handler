//go:build (darwin && cgo) || linux

package pprof

import (
	"github.com/pkg/profile"
	profiler_models "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/models/profiler"
)

type (
	pprof profiler_models.Profiler
)

func NewPProfProfiler(options ...func(*profile.Profile)) profiler_models.Profiler {
	return profile.Start(options...)
}

func NewDefaultPProfProfiler() profiler_models.Profiler {
	return profile.Start(
		profile.GoroutineProfile,
		profile.ProfilePath("./profiling/pprof"),
	)
}
