package mock_generator

import _ "github.com/golang/mock/mockgen/model"

//go:generate mockgen -package mocked_exiter -destination ../../../pkg/mocks/os/exit/exit.go gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/mocks/os/models Exiter
//go:generate mockgen -package mocked_transport_handlers -destination ../../../pkg/mocks/handlers/handler.go gitlab.com/pietroski-software-company/load-test/gotest/pkg/transport-handler/pkg/models/handlers Server
