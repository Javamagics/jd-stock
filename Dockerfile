FROM golang:1.21.11 as builder
#ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
ENV GO111MODULE=on
ENV GOCACHE=/go/pkg/.cache/go-build

ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /app/jd-stock /app/runner.go

FROM alpine:3.6 as alpine
RUN apk update && \
    apk add -U --no-cache ca-certificates tzdata

FROM alpine:3.6
MAINTAINER zhuweitung
LABEL maintainer="zhuweitung" \
    email="zhuweitung@foxmail.com"

ENV TZ="Asia/Shanghai"

COPY --from=alpine /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/jd-stock /app/jd-stock
COPY --from=builder /app/data /app/data
COPY --from=builder /app/static /app/static

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo $TZ > /etc/timezone

WORKDIR /app
# 暴露端口 8080
EXPOSE 8080
CMD ["./jd-stock"]

