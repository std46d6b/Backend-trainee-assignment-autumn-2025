FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o app ./cmd/app

FROM alpine:3.20
WORKDIR /srv
RUN apk add --no-cache ca-certificates && adduser -D -H -u 10001 appuser
COPY --from=builder /app/app /srv/app
EXPOSE 8080
USER appuser
ENTRYPOINT ["/srv/app"]