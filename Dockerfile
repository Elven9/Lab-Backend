FROM golang:1.14.0-alpine3.11 AS build

# Copy File
COPY . /root/source/

WORKDIR /root/source/

RUN [ "go", "build", "-o", "server", "." ]

FROM alpine:3.11.3

WORKDIR /root

# Install curl and git and kubectl
RUN apk update \
    && apk add curl git \
    && curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl \
    && chmod +x ./kubectl \
    && mv ./kubectl /bin/kubectl

COPY --from=build /root/source/server /root/server

EXPOSE 8080/tcp

ENV GIN_MODE=release

ENTRYPOINT ["./server"]