FROM golang:1.24.3-alpine AS builder

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN go build -C cmd/gatherer/ -ldflags="-w -s" -o /go/bin/

FROM alpine

COPY --from=builder /go/bin/gatherer /go/bin/app

CMD ["/go/bin/app"]