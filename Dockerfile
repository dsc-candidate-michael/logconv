FROM golang:1.10

ADD . /go/src/batch-log-converter

WORKDIR /go/src/batch-log-converter 

RUN go build ./cmd...

