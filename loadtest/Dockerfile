FROM alpine:latest

RUN apk add --no-cache curl && \
    curl -L https://github.com/grafana/k6/releases/download/v0.54.0/k6-v0.54.0-linux-amd64.tar.gz | tar xvz -C /tmp && \
    mv /tmp/k6-v0.54.0-linux-amd64/k6 /usr/local/bin/ && \
    rm -rf /tmp/k6-v0.54.0-linux-amd64

CMD ["sleep", "infinity"]