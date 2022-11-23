ARG PROJECT_DIR=/root/project/simple-game

FROM golang:1.19.3 AS builder
ARG PROJECT_DIR
COPY ./common $PROJECT_DIR/common
COPY ./web $PROJECT_DIR/web

WORKDIR $PROJECT_DIR/web/cmd/server
ENV GOPROXY=https://goproxy.io
RUN go build -o web

FROM alpine:latest
ARG PROJECT_DIR
WORKDIR $PROJECT_DIR/web/cmd/server

COPY --from=builder $PROJECT_DIR/web/cmd/server/web ./web
COPY ./web/config $PROJECT_DIR/web/config
RUN mkdir /lib64 && ln -s /lib/ld-musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

EXPOSE 8088
ENTRYPOINT ["./web"]
