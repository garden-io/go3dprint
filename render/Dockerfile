# ARG BASE_IMAGE_VERSION
# FROM base-image:${BASE_IMAGE_VERSION}

# CMD ["./render"]

FROM golang:1.11.9-alpine3.9

WORKDIR /go/src/app

RUN apk add --no-cache mercurial git gcc musl-dev \
  && go get github.com/codegangsta/gin \
  && apk del mercurial

ENV GO111MODULE=on
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["./render"]

EXPOSE 8080

COPY go.mod ./
RUN go mod vendor

COPY *.go ./
RUN go build .