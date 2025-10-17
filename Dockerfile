
# BUILD IMAGE --------------------------------------------------------
ARG GO_VERSION=1.25
FROM golang:${GO_VERSION}-alpine AS builder

# Get build tools and required header files
RUN apk add --no-cache build-base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG VERSION=0.1
RUN go build -ldflags="-X 'main.Version=$VERSION'" -o bin/xmtpd-gateway src/main.go

# ACTUAL IMAGE -------------------------------------------------------

FROM alpine:3.21

LABEL maintainer="eng@ephemerahq.com"
LABEL source="https://github.com/xmtp/gateway-service-example"
LABEL description="XMTP Gateway Service"

# color, nocolor, json
ENV GOLOG_LOG_FMT=json

RUN apk add --no-cache curl

COPY --from=builder /app/bin/xmtpd-gateway /usr/bin/

ENTRYPOINT ["/usr/bin/xmtpd-gateway"]
