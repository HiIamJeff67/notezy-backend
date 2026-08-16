FROM golang:1.26-alpine AS base

WORKDIR /notegic-backend

RUN apk add --no-cache git

COPY contracts/go.mod contracts/go.sum contracts/
COPY shared/go.mod shared/go.sum shared/
COPY internal/realtimegateway/go.mod internal/realtimegateway/go.sum internal/realtimegateway/
WORKDIR /notegic-backend/internal/realtimegateway
RUN go mod download
WORKDIR /notegic-backend

COPY . .

FROM base AS dev

WORKDIR /notegic-backend/internal/realtimegateway

RUN go install github.com/air-verse/air@latest

FROM base AS builder

WORKDIR /notegic-backend/internal/realtimegateway

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/realtimegateway ./commands

FROM alpine:latest AS production

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /out/realtimegateway /app/realtimegateway

EXPOSE 7779

CMD ["/app/realtimegateway"]
