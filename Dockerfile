FROM golang:1.23-alpine
WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off


COPY . .
RUN go mod tidy
RUN go build -o server ./cmd/server
EXPOSE 8080
CMD ["./server"]