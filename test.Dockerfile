FROM golang:1.14.1-alpine as build

ENV CGO_ENABLED=0

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR $GOPATH/src/github.com/mxmxcz/qshare/
COPY . .

RUN go test ./...
CMD go build -o /go/bin/qshare ./cmd/qshare