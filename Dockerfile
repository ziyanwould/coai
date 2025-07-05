# Author: ProgramZmh
# License: Apache-2.0
# Description: Dockerfile for chatnio

FROM --platform=$BUILDPLATFORM golang:1.20-alpine AS backend

WORKDIR /backend
COPY . .

# Set go proxy to https://goproxy.cn (open for vps in China Mainland)
# RUN go env -w GOPROXY=https://goproxy.cn,direct
ARG TARGETARCH
ARG TARGETOS
ENV GOOS=$TARGETOS GOARCH=$TARGETARCH GO111MODULE=on CGO_ENABLED=1

# Install build dependencies
RUN apk update && \
    apk add --no-cache \
    gcc \
    musl-dev \
    g++ \
    make \
    linux-headers \
    wget \
    tar \
    git \
    bash

# Install cross-compilation toolchain for ARM64
RUN if [ "$TARGETARCH" = "arm64" ]; then \
    git clone https://github.com/richfelker/musl-cross-make.git /tmp/musl-cross-make && \
    cd /tmp/musl-cross-make && \
    echo "TARGETS = aarch64-linux-musl" > config.mak && \
    echo "OUTPUT = /usr/local" >> config.mak && \
    make -j$(nproc) && \
    make install && \
    cd / && \
    rm -rf /tmp/musl-cross-make; \
    fi

# Build backend
RUN if [ "$TARGETARCH" = "arm64" ]; then \
    CC=/usr/local/bin/aarch64-linux-musl-gcc \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=arm64 \
    go build -o chat .; \
    else \
    go install && \
    go build .; \
    fi

FROM node:18 AS frontend

WORKDIR /app
COPY ./app .

RUN npm install -g pnpm && \
    pnpm install && \
    pnpm run build && \
    rm -rf node_modules src


FROM alpine

# Install dependencies
RUN apk upgrade --no-cache && \
    apk add --no-cache wget ca-certificates tzdata && \
    update-ca-certificates 2>/dev/null || true

# Set timezone
RUN echo "Asia/Shanghai" > /etc/timezone && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /

# Copy dist
COPY --from=backend /backend/chat /chat
COPY --from=backend /backend/config.example.yaml /config.example.yaml
COPY --from=backend /backend/utils/templates /utils/templates
COPY --from=backend /backend/addition/article/template.docx /addition/article/template.docx
COPY --from=frontend /app/dist /app/dist

# Volumes
VOLUME ["/config", "/logs", "/storage"]

# Expose port
EXPOSE 8094

# Run application
CMD ["./chat"]
