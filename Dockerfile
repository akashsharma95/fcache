## Build Stage
FROM golang:1.17-alpine AS builder
RUN apk add --no-cache git make
WORKDIR /opt/app/
COPY go.mod go.sum ./
RUN go mod download && \
    go mod verify
COPY . .
RUN make build

## Final image
FROM scratch
COPY --from=builder /go/bin/inmemcache /inmemcache
ENTRYPOINT ["/inmemcache"]
LABEL Name=inmemcache Version=0.0.1
