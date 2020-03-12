FROM golang:1.14.0-alpine3.11 AS build

# Copy File
COPY router /root/source/router
COPY go.mod go.sum main.go /root/source/

WORKDIR /root/source/

RUN [ "go", "build", "-o", "server", "." ]

FROM alpine:3.11.3

COPY --from=build /root/source/server /root/server

WORKDIR /root

EXPOSE 8080/tcp

CMD ["./server"]