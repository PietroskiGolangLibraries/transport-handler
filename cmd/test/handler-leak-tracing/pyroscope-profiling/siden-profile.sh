#!/usr/bin/env bash

PORT=$(kubectl describe service "$1" | grep http | grep -o "[0-9][0-9][0-9][0-9]")
LOCAL_PORT_FORWARDING=${2:-"8001"}
OPEN_AT_LOCAL_PORT=${3:-"8081"}
TYPE=${4:-"heap"}
kubectl port-forward service/"$1" "$LOCAL_PORT_FORWARDING:$PORT" | $(sleep 5 && go tool pprof -http=":$OPEN_AT_LOCAL_PORT" "http://localhost:$LOCAL_PORT_FORWARDING/debug/pprof/$TYPE")

siden_custom_profile() {
    PORT=$(kubectl describe service "$1" | grep http | grep -o "[0-9][0-9][0-9][0-9]")
    LOCAL_PORT_FORWARDING=${2:-"8001"}
    OPEN_AT_LOCAL_PORT=${3:-"8081"}
    TYPE=${4:-"heap"}
    kubectl port-forward service/"$1" "$LOCAL_PORT_FORWARDING:$PORT" | $(sleep 5 && go tool pprof -http=":$OPEN_AT_LOCAL_PORT" "http://localhost:$LOCAL_PORT_FORWARDING/debug/pprof/$TYPE")
}
