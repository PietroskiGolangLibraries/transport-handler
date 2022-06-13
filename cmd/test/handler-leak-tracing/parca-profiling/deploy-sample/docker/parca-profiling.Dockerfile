FROM golang:alpine as builder
WORKDIR /app
COPY . .
#COPY ../../../../../../ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o parca-profiling cmd/test/handler-leak-tracing/parca-profiling/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/parca-profiling /usr/bin/
ENTRYPOINT ["parca-profiling"]
