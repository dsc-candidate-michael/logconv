FROM golang:1.10

ADD . /go/src/github.com/logconv

WORKDIR /go/src/github.com/logconv

RUN go build ./cmd...

