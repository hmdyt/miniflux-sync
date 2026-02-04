FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o miniflux-sync .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/miniflux-sync /usr/local/bin/

ENTRYPOINT ["miniflux-sync"]
