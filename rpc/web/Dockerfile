#golang:1.22-alpine
FROM ccr.ccs.tencentyun.com/lxtian/golang:1.22-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata ca-certificates

WORKDIR /build

# 设置 GOFLAGS 环境变量
ENV GOFLAGS="-mod=mod"

COPY . .
COPY rpc/web/etc /app/etc
RUN go mod download
RUN go build -ldflags="-s -w" -o /app/web ./rpc/web/web.go


FROM scratch

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

# 从构建阶段复制根证书
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /app
COPY --from=builder /app/web /app/web
COPY --from=builder /app/etc /app/etc

CMD ["./web", "-f", "etc/web.yaml"]
