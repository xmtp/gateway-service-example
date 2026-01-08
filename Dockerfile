
# BUILD IMAGE --------------------------------------------------------
ARG GO_VERSION=1.25
FROM golang:${GO_VERSION}-alpine AS builder

# Get build tools and required header files
RUN apk add --no-cache build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the gateway binary from cmd/gateway
RUN go build -o /out/xmtpd-gateway ./cmd/gateway

# ACTUAL IMAGE -------------------------------------------------------

FROM alpine:3.21

LABEL maintainer="eng@ephemerahq.com"
LABEL source="https://github.com/xmtp/gateway-service-example"
LABEL description="XMTP Gateway Service"

# color, nocolor, json
ENV GOLOG_LOG_FMT=json

RUN apk add --no-cache curl

COPY --from=builder /out/xmtpd-gateway /usr/bin/xmtpd-gateway

ENTRYPOINT ["/usr/bin/xmtpd-gateway"]
