FROM golang:1.13.0-alpine3.10 AS builder

WORKDIR /app
RUN apk add --no-cache --virtual .go-deps git gcc musl-dev openssl pkgconf
COPY go.mod .
COPY go.sum .
RUN go mod download \
  && go get github.com/rakyll/statik

COPY public /app/public
COPY main.go .

# Build static binary
RUN go generate
RUN GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o speedtest main.go

ENTRYPOINT ["/app/speedtest"]

# Build the smallest image possible
FROM alpine:latest AS runner
COPY --from=builder /app/speedtest /bin/speedtest
ENTRYPOINT ["/bin/speedtest"]
