FROM golang:1.17 as builder

WORKDIR /usr/src/app

RUN go env -w GO111MODULE=auto \
  && go env -w CGO_ENABLED=0 \
  && go env -w GOPROXY=https://goproxy.cn,direct

COPY . .

RUN go mod tidy

RUN set -ex \
    && cd /usr/src/app \
    && go build -ldflags "-s -w -extldflags '-static'" -o zb

FROM alpine:latest

COPY --from=builder /usr/src/app/zb /usr/bin/zb
COPY --from=builder /usr/src/app/hua.ttf /data/hua.ttf
COPY --from=builder /usr/src/app/images /data/images
RUN chmod +x /usr/bin/zb

WORKDIR /data

CMD [ "/usr/bin/zb", "bot"]