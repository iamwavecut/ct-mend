FROM golang:1.19-alpine AS build
ENV GOPATH /go
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && \
    apk update && \
    apk add --no-cache upx gcc g++ && \
    go get github.com/golang-migrate/migrate/v4/cmd/migrate && \
    go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate
COPY . .
RUN migrate -database sqlite3://db.sqlite -path resources/migrations up && \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -ldflags='-w -s -extldflags "-static"' -o server cmd/server/main.go && chmod +x server && upx ./server

RUN wget --header="accept-encoding: gzip" https://github.com/robxu9/bash-static/releases/download/5.1.016-1.2.3/bash-linux-x86_64 -O bash && \
    chmod +x bash && \
    upx ./bash

FROM gcr.io/distroless/static

ARG TLS_ADDR
ARG STORAGE_TYPE
ARG STORAGE_ADDR
ARG LOG_LEVEL
ARG GRACEFUL_TIMEOUT

ENV TLS_ADDR=${TLS_ADDR} \
    STORAGE_TYPE=${STORAGE_TYPE} \
    STORAGE_ADDR=${STORAGE_ADDR} \
    LOG_LEVEL=${LOG_LEVEL} \
    GRACEFUL_TIMEOUT=${GRACEFUL_TIMEOUT}

COPY --from=build /build/bash /usr/bin/
COPY --from=build /build/server /app/
COPY --from=build /build/db.sqlite /app/

WORKDIR /app
ENTRYPOINT ["./server"]
