FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY main.go ./
RUN go build -o proxy-service main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/proxy-service .

EXPOSE 80
CMD ["./proxy-service"]
