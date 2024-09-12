FROM golang:1.23.1-alpine as builder
RUN apk update && apk add build-base cmake gcc git
WORKDIR /go/src/github.com/hovanhoa/go-url-shortener
ADD .. .
WORKDIR /go/src/github.com/hovanhoa/go-url-shortener/cmd/api
RUN go install
WORKDIR /go/bin

FROM alpine:3.18
RUN apk add ca-certificates curl
# install migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate.linux-amd64 /usr/bin/migrate
ENV PATH="${PATH}:/go/bin"
WORKDIR /go/bin
COPY --from=builder /go/src/github.com/hovanhoa/go-url-shortener/*.json /go/bin/
COPY --from=builder /go/bin/* /go/bin/
