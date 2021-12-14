FROM golang:1.17-alpine as builder

RUN apk add git bash

ENV GO111MODULE=on

# Add our code
COPY ./ /src

# build
WORKDIR /src
RUN GOGC=off go build -mod=vendor -v -o /prometheus-browserless-exporter .

# multistage
FROM alpine:latest

RUN apk --update upgrade && \
    apk add curl ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

COPY --from=builder /prometheus-browserless-exporter /usr/bin/prometheus-browserless-exporter

# Run the image as a non-root user
RUN adduser -D prom
RUN chmod 0755 /usr/bin/prometheus-browserless-exporter

USER prom

CMD prometheus-browserless-exporter
