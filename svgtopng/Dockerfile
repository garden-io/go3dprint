FROM golang:alpine

WORKDIR /go/src/app
COPY . .

RUN apk --no-cache add inkscape
RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080

CMD ["app"]