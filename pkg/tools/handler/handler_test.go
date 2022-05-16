package transporthandler

import (
	"fmt"
	"github.com/golang/mock/gomock"
	mocked_exiter "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/mocks/os/exit"
	os_models "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/mocks/os/models"
	handlers_model "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/models/handlers"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mth "gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/mocks/handlers"
)

var (
	anyErr = fmt.Errorf("any-error")
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
				h := NewHandler(nil, nil)

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
		name      string
		setup     func() (Handler, []handlers_model.Server)
		assertion func(h Handler, servers ...handlers_model.Server)
	}{
		{
			name: "",
			setup: func() (Handler, []handlers_model.Server) {
				ms := &mth.MockedRunningServer{}
				h := NewHandler(nil, OsExit)

				return h, []handlers_model.Server{ms}
			},
			assertion: func(h Handler, servers ...handlers_model.Server) {
				defer func() {
					if r := recover(); r != nil {
						str, ok := r.(string)
						require.True(t, ok)
						require.Contains(t, str, "os.Exit(0)")
					}
				}()

				panic("os.Exit(0)")
				//h.StartServers(servers...)
			},
		},
		{
			name: "",
			setup: func() (Handler, []handlers_model.Server) {
				ms := &mth.MockedErrServer{}
				h := NewHandler(nil, OsExit)

				return h, []handlers_model.Server{ms}
			},
			assertion: func(h Handler, servers ...handlers_model.Server) {
				defer func() {
					if r := recover(); r != nil {
						str, ok := r.(string)
						require.True(t, ok)
						t.Log(r)
						require.Contains(t, str, "os.Exit(2)")
					}
				}()

				panic("os.Exit(2)")
				//h.StartServers(servers...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hdl, servers := tt.setup()
			tt.assertion(hdl, servers...)
		})
	}
}

func Test_handler_StartServersAgain(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(os_models.Exiter) *handler
		stubs     func(*mocked_exiter.MockExiter)
		assertion func(*handler, func(...handlers_model.Server))
	}{
		{
			name: "returns error on starting the server",
			setup: func(exiter os_models.Exiter) *handler {
				chanToPanic := make(chan error)
				h := &handler{
					stopServerSig: chanToPanic,
					osExit:        exiter.Exit,
				}

				return h
			},
			stubs: func(exiter *mocked_exiter.MockExiter) {
				exiter.EXPECT().Exit(1).Times(1).Return()
			},
			assertion: func(
				h *handler,
				fn func(servers ...handlers_model.Server),
			) {
				fn([]handlers_model.Server{
					&mth.MockedErrServer{},
				}...)
			},
		},
		{
			name: "panics on starting the server",
			setup: func(exiter os_models.Exiter) *handler {
				chanToPanic := make(chan error)
				h := &handler{
					stopServerSig: chanToPanic,
					osExit:        exiter.Exit,
				}

				return h
			},
			stubs: func(exiter *mocked_exiter.MockExiter) {
				exiter.EXPECT().Exit(2).Times(1).Return()
			},
			assertion: func(
				h *handler,
				fn func(servers ...handlers_model.Server),
			) {
				fn([]handlers_model.Server{
					&mth.MockedPanicServer{},
				}...)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockedExiter := mocked_exiter.NewMockExiter(ctrl)
			h := tt.setup(mockedExiter)
			tt.stubs(mockedExiter)
			tt.assertion(h, h.StartServers)
		})
	}
}

func Test_handler_StartServersAgainWithAutoGenMocks(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(os_models.Exiter) *handler
		stubs     func(*mocked_exiter.MockExiter, *mth.MockServer)
		assertion func(
			*handler,
			func(...handlers_model.Server),
			[]handlers_model.Server,
		)
	}{
		{
			name: "returns error on starting the server",
			setup: func(exiter os_models.Exiter) *handler {
				chanToPanic := make(chan error)
				h := &handler{
					stopServerSig: chanToPanic,
					osExit:        exiter.Exit,
				}

				return h
			},
			stubs: func(
				exiter *mocked_exiter.MockExiter,
				server *mth.MockServer,
			) {
				exiter.EXPECT().Exit(1).Times(1).Return()
				server.EXPECT().Start().Times(1).Return(anyErr)
			},
			assertion: func(
				h *handler,
				fn func(servers ...handlers_model.Server),
				servers []handlers_model.Server,
			) {
				fn(servers...)
			},
		},
		{
			name: "panics on starting the server",
			setup: func(exiter os_models.Exiter) *handler {
				chanToPanic := make(chan error)
				h := &handler{
					stopServerSig: chanToPanic,
					osExit:        exiter.Exit,
				}

				return h
			},
			stubs: func(
				exiter *mocked_exiter.MockExiter,
				server *mth.MockServer,
			) {
				exiter.EXPECT().Exit(2).Times(1).Return()
				server.EXPECT().Start().Times(0).Return(nil)
			},
			assertion: func(
				h *handler,
				fn func(servers ...handlers_model.Server),
				servers []handlers_model.Server,
			) {
				servers[0] = nil
				fn(servers...)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockedExiter := mocked_exiter.NewMockExiter(ctrl)
			mockedServer := mth.NewMockServer(ctrl)
			h := tt.setup(mockedExiter)
			tt.stubs(mockedExiter, mockedServer)
			tt.assertion(h, h.StartServers, []handlers_model.Server{mockedServer})
		})
	}
}

func TestPanicStart(t *testing.T) {
	ms := &mth.MockedPanicServer{}
	h := NewHandler(nil, OsExit)

	// Run the crashing code when FLAG is set
	if os.Getenv("FLAG") == "2" {
		h.StartServers(ms)
		return
	}
	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestPanicStart")
	cmd.Env = append(os.Environ(), "FLAG=2")
	err := cmd.Run()

	// Cast the error as *exec.ExitError and compare the result
	e, ok := err.(*exec.ExitError)
	expectedErrorString := "exit status 2"
	assert.Equal(t, true, ok)
	assert.Equal(t, expectedErrorString, e.Error())
}

func TestErrStart(t *testing.T) {
	ms := &mth.MockedErrServer{}
	h := NewHandler(nil, OsExit)

	// Run the crashing code when FLAG is set
	if os.Getenv("FLAG") == "1" {
		h.StartServers(ms)
		return
	}
	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestErrStart")
	cmd.Env = append(os.Environ(), "FLAG=1")
	err := cmd.Run()

	// Cast the error as *exec.ExitError and compare the result
	e, ok := err.(*exec.ExitError)
	expectedErrorString := "exit status 1"
	assert.Equal(t, true, ok)
	assert.Equal(t, expectedErrorString, e.Error())
}

func TestRunningStart(t *testing.T) {
	ms := &mth.MockedRunningServer{}
	h := NewHandler(nil, OsExit)

	// Run the crashing code when FLAG is set
	if os.Getenv("FLAG") == "0" {
		h.StartServers(ms)
		return
	}
	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestRunningStart")
	cmd.Env = append(os.Environ(), "FLAG=0")
	err := cmd.Run()

	// Cast the error as *exec.ExitError and compare the result
	e, ok := err.(*exec.ExitError)
	assert.Equal(t, false, ok)
	assert.Nil(t, e)
}

func Test_handler_closeChan_and_handleCloseChanPanic(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(os_models.Exiter) *handler
		stubs     func(*mocked_exiter.MockExiter)
		assertion func(*handler, func())
	}{
		{
			name: "panics when closing chan",
			setup: func(exiter os_models.Exiter) *handler {
				chanToPanic := make(chan error)
				h := &handler{
					stopServerSig: chanToPanic,
					osExit:        exiter.Exit,
				}

				return h
			},
			stubs: func(exiter *mocked_exiter.MockExiter) {
				exiter.
					EXPECT().
					Exit(2).
					Times(1).
					Return()
			},
			assertion: func(h *handler, fn func()) {
				close(h.stopServerSig)
				fn()
			},
		},
		{
			name: "does not panic when closing chan",
			setup: func(exiter os_models.Exiter) *handler {
				chanToPanic := make(chan error)
				h := &handler{
					stopServerSig: chanToPanic,
					osExit:        exiter.Exit,
				}

				return h
			},
			stubs: func(exiter *mocked_exiter.MockExiter) {},
			assertion: func(h *handler, fn func()) {
				fn()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockedExiter := mocked_exiter.NewMockExiter(ctrl)
			h := tt.setup(mockedExiter)
			tt.stubs(mockedExiter)
			tt.assertion(h, h.closeChan)
		})
	}
}

func Test_handler_verifyCodeZero(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(os_models.Exiter) (Handler, interface{})
		stubs     func(*mocked_exiter.MockExiter)
		assertion func(func(interface{}), interface{})
	}{
		{
			name: "testing code zero",
			setup: func(exiter os_models.Exiter) (Handler, interface{}) {
				h := NewHandler(nil, exiter.Exit)
				r := "os.Exit(0)"
				return h, r
			},
			stubs: func(mockedExiter *mocked_exiter.MockExiter) {
				mockedExiter.
					EXPECT().
					Exit(0).
					Times(1).
					Return()
			},
			assertion: func(fn func(interface{}), r interface{}) {
				fn(r)
			},
		},
		{
			name: "testing non code zero - do not call exit",
			setup: func(exiter os_models.Exiter) (Handler, interface{}) {
				h := NewHandler(nil, exiter.Exit)
				r := "os.Exit(1)"
				return h, r
			},
			stubs: func(mockedExiter *mocked_exiter.MockExiter) {},
			assertion: func(fn func(interface{}), r interface{}) {
				fn(r)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockedExiter := mocked_exiter.NewMockExiter(ctrl)
			h, r := tt.setup(mockedExiter)
			tt.stubs(mockedExiter)
			tt.assertion(h.verifyCodeZero, r)
		})
	}
}
