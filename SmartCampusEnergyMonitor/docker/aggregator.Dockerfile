FROM golang:1.24 as builder

WORKDIR /app

# Copy go.mod and go.sum (if exists)
COPY aggregator-service/go.mod ./
# Copy all source code including generated pb
COPY aggregator-service/ .

# Tidy dependencies (ensure everything matches source)
RUN go mod tidy

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o aggregator main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/aggregator .
EXPOSE 50051 2112
CMD ["./aggregator"]
