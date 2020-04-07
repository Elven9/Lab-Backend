FROM golang:1.14.0-alpine3.11 AS build

# Copy File
COPY router /root/source/router
COPY utils /root/source/utils
COPY go.mod go.sum main.go /root/source/

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

COPY k8s-config /root/.kube/config

COPY --from=build /root/source/server /root/server
COPY hardwareInfo.json /root/hardwareInfo.json

EXPOSE 8080/tcp

ENV GIN_MODE=release

CMD ["./server"]