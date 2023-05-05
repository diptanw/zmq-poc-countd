FROM golang:1.20-alpine3.17 AS builder

ARG package

WORKDIR /src/
COPY . .

RUN apk add --update --no-cache \
    zeromq-dev pkgconfig alpine-sdk libsodium-dev libzmq-static libsodium-static

RUN CGO_LDFLAGS="$CGO_LDFLAGS -lstdc++ -lm -lsodium" \
    CGO_ENABLED=1 \
    go build -v -a -ldflags '-s -w -extldflags "-static"' -o exec ./cmd/${package}

FROM alpine:3.17

RUN apk upgrade --update --no-cache \
	&& addgroup -S 65011 \
	&& adduser -D -S -G 65011 65011

USER 65011:65011

COPY --from=builder /src/exec /usr/local/bin/

ENTRYPOINT [ "exec" ]

