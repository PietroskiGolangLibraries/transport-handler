# Makefile command list

mock-generate:
	go get -d github.com/golang/mock/mockgen
	go mod vendor
	go generate ./...
	go mod vendor

test-unit:
	go test -race -v `go list ./... | grep -v ./pkg/mocks`

test-unit-cover:
	go test -race -v -coverprofile ./docs/reports/tests/unit/cover.out `go list ./... | grep -v ./pkg/mocks`

test-unit-cover-silent:
	go test -race -coverprofile ./docs/reports/tests/unit/cover.out `go list ./... | grep -v ./pkg/mocks`

test-unit-cover-all:
	go test -race -v -coverprofile ./docs/reports/tests/unit/cover-all.out ./...

test-unit-cover-all-silent:
	go test -race -coverprofile ./docs/reports/tests/unit/cover-all.out ./...

test-unit-cover-report:
	go tool cover -html=./docs/reports/tests/unit/cover.out

test-unit-cover-all-report:
	go tool cover -html=./docs/reports/tests/unit/cover-all.out

gen-goroutine-profile:
	go tool pprof --pdf ./profiling/pprof/goroutine.pprof > ./docs/reports/profiling/pprof/goroutine-profiling.pdf

gen-threadcreation-profile:
	go tool pprof --pdf ./profiling/pprof/threadcreation.pprof > ./docs/reports/profiling/pprof/threadcreation-profiling.pdf

profiling: gen-goroutine-profile gen-threadcreation-profile

run-main-test:
	go run -race cmd/test/handler-leak-tracing/manual-stopping/main.go
