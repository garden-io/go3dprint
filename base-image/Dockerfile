FROM golang:1.11.9-alpine3.9

WORKDIR /go/src/app

RUN apk add --no-cache entr git gcc musl-dev

ENV GO111MODULE=on
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["./app"]
EXPOSE 8080

ONBUILD COPY go.mod hotreload.sh ./
ONBUILD RUN go mod vendor

ONBUILD COPY *.go ./
ONBUILD RUN go build .