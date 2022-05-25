package transporthandler

import (
	"context"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	mocked_exiter "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/mocks/os/exit"
	mocked_profiler "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/mocks/profiling/pprof"
	stack_tracer "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/tools/tracer/stack"
)

func Test_handler_handleCloseChanPanic(t *testing.T) {
	tests := []struct {
		name  string
		setup func(
			mockProfiler *mocked_profiler.MockProfiler,
			exiter *mocked_exiter.MockExiter,
		) *handler
		stubs func(
			mockProfiler *mocked_profiler.MockProfiler,
			exiter *mocked_exiter.MockExiter,
		)
		assertion func(*handler)
	}{
		{
			name: "recovers from panic",
			setup: func(
				prof *mocked_profiler.MockProfiler,
				exiter *mocked_exiter.MockExiter,
			) *handler {
				ctx, cancel := context.WithCancel(context.Background())
				h := &handler{
					ctx:      ctx,
					cancelFn: cancel,
					osExit:   exiter.Exit,
					goPool: goPool{
						wg:  &sync.WaitGroup{},
						gst: stack_tracer.NewGST(),
					},
					srvChan:  srvChan{},
					profiler: profiler{pprof: prof},
				}
				h.makeSrvChan(0)
				return h
			},
			stubs: func(
				mockedProfiler *mocked_profiler.MockProfiler,
				mockedExiter *mocked_exiter.MockExiter,
			) {
				mockedProfiler.EXPECT().Stop().Times(0).Return()
				mockedExiter.EXPECT().Exit(2).Times(1).Return()
			},
			assertion: func(h *handler) {
				defer h.handleCloseChanPanic()
				close(h.srvChan.panicSig)
				close(h.srvChan.panicSig)
			},
		},
		{
			name: "does not panic",
			setup: func(
				mockProfiler *mocked_profiler.MockProfiler,
				exiter *mocked_exiter.MockExiter,
			) *handler {
				return &handler{}
			},
			stubs: func(
				mockProfiler *mocked_profiler.MockProfiler,
				exiter *mocked_exiter.MockExiter,
			) {
				return
			},
			assertion: func(h *handler) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockedProfiler := mocked_profiler.NewMockProfiler(ctrl)
			mockedExiter := mocked_exiter.NewMockExiter(ctrl)
			tt.stubs(mockedProfiler, mockedExiter)
			h := tt.setup(mockedProfiler, mockedExiter)
			tt.assertion(h)
		})
	}
}
