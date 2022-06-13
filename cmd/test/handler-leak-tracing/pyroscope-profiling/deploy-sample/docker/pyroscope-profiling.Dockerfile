FROM golang:alpine as builder
WORKDIR /app
COPY . .
# ../../../../../../ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o pyroscope-profiling cmd/test/handler-leak-tracing/pyroscope-profiling/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/pyroscope-profiling /usr/bin/
ENTRYPOINT ["pyroscope-profiling"]
