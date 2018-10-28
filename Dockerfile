FROM golang:1.10

ADD . /go/src/logconv

WORKDIR /go/src/logconv

RUN go build ./cmd...

