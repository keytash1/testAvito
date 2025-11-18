FROM golang:1.23-alpine

WORKDIR /app

ENV GOPROXY=https://goproxy.io,direct
ENV GOSUMDB=off

RUN apk add --no-cache git

COPY . .
RUN go build -o server ./cmd/server

EXPOSE 8080
CMD ["./server"]