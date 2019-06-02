FROM golang:1.12.5-alpine

WORKDIR /go/src/github.com/kind84/hlx
COPY . .

RUN apk update && apk add git gcc libc-dev
RUN go get -d -v ./...