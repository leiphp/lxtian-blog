FROM golang:1.22-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

# 设置 GOFLAGS 环境变量
ENV GOFLAGS="-mod=mod"

COPY . .
COPY gateway/etc /app/etc
RUN go mod download
RUN go build -ldflags="-s -w" -o /app/gateway ./gateway/gateway.go

FROM scratch

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/gateway /app/gateway
COPY --from=builder /app/etc /app/etc

CMD ["./gateway", "-f", "etc/gateway-api.yaml"]
