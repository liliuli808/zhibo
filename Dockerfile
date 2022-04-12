FROM golang:1.17 as builder

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy

RUN go build -ldflags "-s -w" -o zb main.go

COPY zb /usr/local/bin/

CMD ["zb"]