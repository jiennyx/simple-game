ARG PROJECT_DIR=/root/project/simple-game

FROM golang:1.19.3 AS builder
ARG PROJECT_DIR
COPY ./common $PROJECT_DIR/common
COPY ./userservice $PROJECT_DIR/userservice

WORKDIR $PROJECT_DIR/userservice/cmd/rpc
ENV GOPROXY=https://goproxy.io
RUN go build -o userservice

FROM alpine:latest
ARG PROJECT_DIR
WORKDIR $PROJECT_DIR/userservice/cmd/rpc

COPY --from=builder $PROJECT_DIR/userservice/cmd/rpc/userservice ./userservice
COPY ./userservice/config $PROJECT_DIR/userservice/config
RUN mkdir /lib64 && ln -s /lib/ld-musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

EXPOSE 80
ENTRYPOINT ["./userservice"]
