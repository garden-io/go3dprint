# ARG BASE_IMAGE_VERSION
# FROM base-image:${BASE_IMAGE_VERSION}

# ENTRYPOINT ["/bin/sh", "-c"]
# CMD ["find . | entr -r ./hotreload.sh"]

FROM golang:1.11.9-alpine3.9

WORKDIR /go/src/app

RUN apk add --no-cache entr git gcc musl-dev

ENV GO111MODULE=on
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["find . | entr -r ./hotreload.sh"]
EXPOSE 8080

COPY hotreload.sh ./
COPY go.mod ./
RUN go mod vendor
COPY *.go ./
RUN go build