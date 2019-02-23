FROM golang:alpine

WORKDIR /go/src/app

RUN apk add --no-cache mercurial git \
  && go get github.com/codegangsta/gin \
  && apk del mercurial

ENV GO111MODULE=on
CMD ["./app"]
EXPOSE 8080

ONBUILD COPY go.mod go.sum ./
ONBUILD RUN go mod vendor

ONBUILD COPY *.go ./
ONBUILD RUN go build