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

########################################################################################################################

gen-goroutine-profile:
	go tool pprof --pdf ./profiling/pprof/goroutine.pprof > ./docs/reports/profiling/pprof/goroutine-profiling.pdf

gen-threadcreation-profile:
	go tool pprof --pdf ./profiling/pprof/threadcreation.pprof > ./docs/reports/profiling/pprof/threadcreation-profiling.pdf

profiling: gen-goroutine-profile gen-threadcreation-profile

run-main-test:
	go run -race cmd/test/handler-leak-tracing/manual-stopping/main.go

########################################################################################################################

TAG := $(shell cat VERSION)
tag:
	git tag $(TAG)

########################################################################################################################

build-pyroscope-profiling:
	docker image build --no-cache -t pietroski/pyroscope-profiling-test -f cmd/test/handler-leak-tracing/pyroscope-profiling/deploy-sample/docker/pyroscope-profiling.Dockerfile .

tagging-pyroscope-profiling-image:
	docker tag pietroski/pyroscope-profiling-test pietroski/pyroscope-profiling-test:v0.0.6

update-docker-pyroscope-profiling-image:
	docker push pietroski/pyroscope-profiling-test:v0.0.6

kube-pyroscope-profiling-deployment:
	kubectl-24 apply -f cmd/test/handler-leak-tracing/pyroscope-profiling/deploy-sample/k8s/k8s-deployment.yml

kube-pyroscope-patch: kube-pyroscope-profiling-deployment

kube-pyroscope-cleanup:
	kubectl delete deployments.apps parca-hello-world

########################################################################################################################
