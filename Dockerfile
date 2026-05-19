FROM golang:1.26.3-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux go build -C cmd/gatherer/ -o /out/gatherer

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata postgresql17-client

WORKDIR /app

COPY --from=builder /out/gatherer /app/gatherer

EXPOSE 8080

ENTRYPOINT ["/app/gatherer"]