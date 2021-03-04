FROM golang:1.16-alpine as build
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk add --no-cache make git

WORKDIR /go/src/github.com/websublime/courier

# Pulling dependencies
COPY ./Makefile ./go.* ./
RUN make deps

# Building stuff
COPY . /go/src/github.com/websublime/courier
RUN make build

FROM alpine:3.7
RUN adduser -D -u 1000 websublime

RUN apk add --no-cache ca-certificates
COPY --from=build /go/src/github.com/websublime/courier/courier /usr/local/bin/courier

ENV COURIER_PRODUCTION true
ENV COURIER_PORT 8883
ENV COURIER_HOST localhost
ENV COURIER_WS_URL ws://localhost:8883/ws
ENV COURIER_DATABASE_URL postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
ENV COURIER_DATABASE_NAMESPACE courier
ENV COURIER_JWT_SECRET 3EK6FD+o0+c7tzBNVfjpMkNDi2yARAAKzQlk8O2IKoxQu4nF7EdAh8s3TwpHwrdWT6R
ENV COURIER_KEY_SECRET kNDKzQlk8ONVfjpMKo2I75tg67ujmki8

USER websublime
CMD ["courier"]