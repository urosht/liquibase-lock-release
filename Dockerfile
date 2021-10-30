FROM golang:alpine AS builder
WORKDIR $GOPATH/src/app/

COPY main.go .
COPY go.mod .

RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s -extldflags "-static"' -a -o /go/bin/app .

FROM scratch
COPY --from=builder /go/bin/app /go/bin/app

ENTRYPOINT ["/go/bin/app"]