FROM golang:alpine

WORKDIR /go/src/app
COPY . .

RUN apk --no-cache add inkscape
RUN go get -d -v ./...
RUN go install -v ./...

# RUN apk add --no-cache git mercurial \
#     && go get github.com/codegangsta/gin \
#     && apk del git mercurial

# CMD gin -p 8080 -a 8081 --path . run main.go

EXPOSE 8080

CMD ["app"]