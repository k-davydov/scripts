# Start from the official Go image for building the binary
FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

COPY . ./

RUN go build -o tools

# Final minimal image
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/tools .
EXPOSE 8080
ENTRYPOINT ["./tools"]
