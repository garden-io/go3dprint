FROM golang:alpine

WORKDIR /go/src/app
COPY . .

RUN apk add --no-cache git
ENV GO111MODULE=on
RUN go mod vendor

RUN apk add --no-cache mercurial \
    && go get github.com/codegangsta/gin \
    && apk del git mercurial

RUN go get -d -v ./...
RUN go install -v ./...

CMD gin -p 8080 -a 8081 run main.go
EXPOSE 8080

# CMD ["app"]