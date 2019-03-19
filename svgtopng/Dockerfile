FROM golang:alpine

WORKDIR /go/src/app
COPY . .

RUN apk add --no-cache git mercurial inkscape \
    && go get github.com/codegangsta/gin \
    && apk del git mercurial

RUN go get -d -v ./...
RUN go install -v ./...

CMD gin -p 8080 -a 8081 run main.go

EXPOSE 8080