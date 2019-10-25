FROM alpine
RUN apk add entr ca-certificates
EXPOSE 8080
RUN mkdir app
WORKDIR /app
COPY bin/ /app/
CMD ls binary | entr -r ./binary
