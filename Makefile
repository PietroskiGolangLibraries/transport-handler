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

test-unit-cover-report:
	go tool cover -html=./docs/reports/tests/unit/cover.out
