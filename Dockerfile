FROM golang:1.13-alpine AS build-dist
ENV GOPROXY='https://goproxy.cn,direct'

WORKDIR /go/cache

ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /go/release

ADD . .

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -tags netgo -installsuffix cgo -o /bin/app main.go

FROM alpine as prod

COPY --from=build-dist /bin/app /bin/app

RUN apk add --no-cache -U  tzdata && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    chmod +x /bin/app

ENV REDIS_HOST='redis'
ENV REDIS_PORT=6379
ENV REDIS_CHANNEL=0
ENV ENCRYPT_KEY="xxxxxxxxxxxxxxxx"

CMD ["/bin/app"]
EXPOSE 80
EXPOSE 9000