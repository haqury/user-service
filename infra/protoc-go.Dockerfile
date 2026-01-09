FROM golang:1.22-alpine AS builder

RUN apk add --no-cache \
    protobuf \
    protobuf-dev \
    git \
    make \
    curl

# Устанавливаем protoc-gen-go, protoc-gen-go-grpc и protoc-gen-grpc-gateway
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.0 && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.0

# Для поддержки импортов google/api/annotations.proto
RUN mkdir -p /include && \
    wget -q -O /tmp/googleapis.zip https://github.com/googleapis/googleapis/archive/master.zip && \
    unzip -q /tmp/googleapis.zip -d /tmp && \
    mv /tmp/googleapis-master/google /include/ && \
    rm -rf /tmp/googleapis.zip /tmp/googleapis-master

# Улучшенный entrypoint скрипт
COPY infra/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

WORKDIR /workspace
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]