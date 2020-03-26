FROM golang@sha256:291f4c64d9e3c397aaf0bd94ccb3a573b918be2417dbf5907768f29e8de73929 as build

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR $GOPATH/src/github.com/mxmxcz/qshare/
COPY . .

RUN go test ./...
RUN go build -o /go/bin/qshare ./cmd/qshare
COPY static /go/bin/static

FROM scratch
WORKDIR /
COPY --from=build /go/bin/qshare /qshare
COPY --from=build /go/bin/static /static
EXPOSE 8080
ENTRYPOINT ["./qshare"]