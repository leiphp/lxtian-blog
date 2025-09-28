FROM golang:1.21-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build/zero

RUN apk update --no-cache && apk add --no-cache tzdata

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
COPY ./etc/payment.yaml /app/etc/
RUN go build -ldflags="-s -w" -o /app/payment ./payment.go


FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/payment /app/payment
COPY --from=builder /app/etc /app/etc

CMD ["./payment", "-f", "etc/payment.yaml"]
