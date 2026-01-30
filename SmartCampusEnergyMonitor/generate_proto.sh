#!/bin/bash
set -e

echo "Generating Go stubs..."
docker run --rm -v $(pwd):/workspace -w /workspace golang:latest sh -c "
    apt-get update -qq && apt-get install -y -qq protobuf-compiler && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    export PATH=\$PATH:\$(go env GOPATH)/bin && \
    mkdir -p aggregator-service/pb && \
    protoc -Iproto --go_out=aggregator-service/pb --go_opt=paths=source_relative --go-grpc_out=aggregator-service/pb --go-grpc_opt=paths=source_relative proto/sensor.proto
"

echo "Generating Python stubs..."
docker run --rm -v $(pwd):/workspace -w /workspace python:3.9 sh -c "
    pip install -q grpcio-tools && \
    python -m grpc_tools.protoc -Iproto --python_out=sensor-simulator --grpc_python_out=sensor-simulator proto/sensor.proto
"
echo "Done."
