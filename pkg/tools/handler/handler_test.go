package transporthandler

import (
	"context"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mocked_transport_handlers "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/mocks/handlers"
	mocked_exiter "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/mocks/os/exit"
	"gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/mocks/os/exit/fake"
	mocked_profiler "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/mocks/profiling/pprof"
	fakepprof "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/mocks/profiling/pprof/fake"
	handlers_model "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/v4/pkg/models/handlers"
)

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name  string
		setup func() Handler
		want  func(Handler)
	}{
		{
			name: "new handler initialisation",
			setup: func() Handler {
				ctx, cancel := context.WithCancel(context.Background())
				h := NewHandler(
					ctx,
					cancel,
					fakepprof.NewFakePProfProfiler(),
					fake.Exit,
				)

				return h
			},
			want: func(h Handler) {
				require.NotNil(t, h)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want(tt.setup())
		})
	}
}

func TestNewDefaultHandler(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() Handler
		assertion func(*testing.T, Handler)
	}{
		{
			name: "returns a default handler",
			setup: func() Handler {
				h := NewDefaultHandler()
				return h
			},
			assertion: func(t *testing.T, h Handler) {
				require.NotNil(t, h)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.setup())
		})
	}
}

func Test_handler_StartServers(t *testing.T) {
	tests := []struct {
		name  string
		setup func(
			mockProfiler *mocked_profiler.MockProfiler,
			exiter *mocked_exiter.MockExiter,
		) Handler
		stubs func(
			mockProfiler *mocked_profiler.MockProfiler,
			exiter *mocked_exiter.MockExiter,
		)
		assertion func(Handler)
	}{
		{
			name: "quit sig sent - exits with code 0 - sigKill",
			setup: func(
				profiler *mocked_profiler.MockProfiler,
				exiter *mocked_exiter.MockExiter,
			) Handler {
				ctx, cancel := context.WithCancel(context.Background())
				h := NewHandler(ctx, cancel, profiler, exiter.Exit)
				return h
			},
			stubs: func(
				mockedProfiler *mocked_profiler.MockProfiler,
				mockedExiter *mocked_exiter.MockExiter,
			) {
				mockedProfiler.EXPECT().Stop().Times(1).Return()
				mockedExiter.EXPECT().Exit(0).Times(1).Return()
			},
			assertion: func(h Handler) {
				wg := &sync.WaitGroup{}
				svr1 := mocked_transport_handlers.NewMockedNamedRunningServer("server-1")
				svr2 := mocked_transport_handlers.NewMockedNamedRunningServer("server-2")

				wg.Add(1)
				go func() {
					defer wg.Done()
					h.StartServers(
						map[string]handlers_model.Server{
							"server-1": svr1,
							"server-2": svr2,
						},
					)
				}()

				time.Sleep(time.Millisecond * 500)
				_ = syscall.Kill(os.Getpid(), syscall.SIGQUIT)
				//h.Cancel()
				wg.Wait()
			},
		},
		{
			name: "err sig sent - exits with code 1",
			setup: func(
				profiler *mocked_profiler.MockProfiler,
				exiter *mocked_exiter.MockExiter,
			) Handler {
				ctx, cancel := context.WithCancel(context.Background())
				h := NewHandler(ctx, cancel, profiler, exiter.Exit)
				return h
			},
			stubs: func(mockedProfiler *mocked_profiler.MockProfiler, mockedExiter *mocked_exiter.MockExiter) {
				mockedProfiler.EXPECT().Stop().Times(1).Return()
				mockedExiter.EXPECT().Exit(1).Times(1).Return()
			},
			assertion: func(h Handler) {
				svr1 := mocked_transport_handlers.NewMockedNamedRunningServer("server-1")
				svr2 := mocked_transport_handlers.NewMockedNamedRunningServer("server-2")
				sErr1 := &mocked_transport_handlers.MockedErrServer{}
				sErr2 := &mocked_transport_handlers.MockedErrServer{}

				h.StartServers(
					map[string]handlers_model.Server{
						"server-1":  svr1,
						"server-2":  svr2,
						"err-srv-1": sErr1,
						"err-srv-2": sErr2,
					},
				)
			},
		},
		{
			name: "panic sig sent - exits with code 2",
			setup: func(
				profiler *mocked_profiler.MockProfiler,
				exiter *mocked_exiter.MockExiter,
			) Handler {
				ctx, cancel := context.WithCancel(context.Background())
				h := NewHandler(ctx, cancel, profiler, exiter.Exit)
				return h
			},
			stubs: func(mockedProfiler *mocked_profiler.MockProfiler, mockedExiter *mocked_exiter.MockExiter) {
				mockedProfiler.EXPECT().Stop().Times(1).Return()
				mockedExiter.EXPECT().Exit(2).Times(1).Return()
			},
			assertion: func(h Handler) {
				svr1 := mocked_transport_handlers.NewMockedNamedRunningServer("server-1")
				svr2 := mocked_transport_handlers.NewMockedNamedRunningServer("server-2")
				sPanic1 := &mocked_transport_handlers.MockedPanicServer{}
				sPanic2 := &mocked_transport_handlers.MockedPanicServer{}

				h.StartServers(
					map[string]handlers_model.Server{
						"server-1":    svr1,
						"server-2":    svr2,
						"panic-srv-1": sPanic1,
						"panic-srv-2": sPanic2,
					},
				)
			},
		},
		{
			name: "quit sig sent - exits with code 0 - context",
			setup: func(
				profiler *mocked_profiler.MockProfiler,
				exiter *mocked_exiter.MockExiter,
			) Handler {
				ctx, cancel := context.WithCancel(context.Background())
				h := NewHandler(ctx, cancel, profiler, exiter.Exit)
				return h
			},
			stubs: func(
				mockedProfiler *mocked_profiler.MockProfiler,
				mockedExiter *mocked_exiter.MockExiter,
			) {
				mockedProfiler.EXPECT().Stop().Times(1).Return()
				mockedExiter.EXPECT().Exit(1).Times(1).Return()
			},
			assertion: func(h Handler) {
				wg := &sync.WaitGroup{}
				svr1 := mocked_transport_handlers.NewMockedNamedRunningServer("server-1")
				svr2 := mocked_transport_handlers.NewMockedNamedRunningServer("server-2")

				wg.Add(1)
				go func() {
					defer wg.Done()
					h.StartServers(
						map[string]handlers_model.Server{
							"server-1": svr1,
							"server-2": svr2,
						},
					)
				}()

				time.Sleep(time.Millisecond * 500)
				//_ = syscall.Kill(os.Getpid(), syscall.SIGQUIT)
				h.Cancel()
				wg.Wait()
			},
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
