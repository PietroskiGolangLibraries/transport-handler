package transporthandler

import (
	"context"
	"log"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	fakepprof "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/mocks/profiling/pprof/fake"

	"github.com/golang/mock/gomock"
	mocked_exiter "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/mocks/os/exit"
	mocked_profiler "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/mocks/profiling/pprof"
	stack_tracer "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/tools/tracer/stack"
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

func Test_handler_handleStopPanic(t *testing.T) {
	t.Run(
		"it panics but recovers",
		func(t *testing.T) {
			h := &handler{}
			defer h.handleStopPanic()
			panic("any-panic")
		},
	)

	t.Run(
		"does not panic neither recovers",
		func(t *testing.T) {
			h := &handler{}
			defer h.handleStopPanic()
		},
	)
}

func Test_handleCtxGen(t *testing.T) {
	fakepprof := fakepprof.NewFakePProfProfiler()
	var exiter = func(i int) {
		log.Printf("exiting with code: %v", i)
	}
	tests := []struct {
		name      string
		setup     func() (context.Context, context.CancelFunc)
		assertion func(context.Context, context.CancelFunc)
	}{
		{
			name: "ctx and cancelFn as not nil",
			setup: func() (context.Context, context.CancelFunc) {
				ctx, cancelFn := context.WithCancel(context.Background())
				return ctx, cancelFn
			},
			assertion: func(ctx context.Context, cancelFunc context.CancelFunc) {
				require.NotNil(t, ctx)
				require.NotNil(t, cancelFunc)
				h := NewHandler(ctx, cancelFunc, fakepprof, exiter)
				require.NotNil(t, h)
			},
		},
		{
			name: "ctx and cancelFn as nil",
			setup: func() (context.Context, context.CancelFunc) {
				return nil, nil
			},
			assertion: func(ctx context.Context, cancelFunc context.CancelFunc) {
				require.Nil(t, ctx)
				require.Nil(t, cancelFunc)
				h := NewHandler(ctx, cancelFunc, fakepprof, exiter)
				require.NotNil(t, h)
			},
		},
		{
			name: "ctx not will but cancelFn as nil",
			setup: func() (context.Context, context.CancelFunc) {
				ctx := context.Background()
				return ctx, nil
			},
			assertion: func(ctx context.Context, cancelFunc context.CancelFunc) {
				require.NotNil(t, ctx)
				require.Nil(t, cancelFunc)
				h := NewHandler(ctx, cancelFunc, fakepprof, exiter)
				require.NotNil(t, h)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(tt.setup())
		})
	}
}
