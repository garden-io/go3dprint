FROM golang:alpine

WORKDIR /go/src/app

RUN apk add --no-cache mercurial git \
    && go get github.com/codegangsta/gin \
    && apk del mercurial

# RUN go get -d -v ./...
ENV GO111MODULE=on
COPY go.mod ./
RUN go mod vendor

CMD ["gin", "-p", "8080", "-a", "8081", "run", "main.go"]
EXPOSE 8080

COPY main.go ./
RUN go install