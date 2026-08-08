FROM golang:1.26-alpine AS base

WORKDIR /notezy-backend

RUN apk add --no-cache git

COPY go.mod go.sum ./
COPY internal/gateway/go.mod internal/gateway/go.sum internal/gateway/
COPY internal/durablejob/go.mod internal/durablejob/go.sum internal/durablejob/
COPY internal/email/go.mod internal/email/go.sum internal/email/
COPY internal/realtimegateway/go.mod internal/realtimegateway/go.sum internal/realtimegateway/
RUN go mod download

COPY . .

FROM base AS dev

WORKDIR /notezy-backend/internal/gateway

RUN go install github.com/air-verse/air@latest

FROM base AS builder

WORKDIR /notezy-backend/internal/gateway

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/gateway ./commands

FROM alpine:latest AS production

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /out/gateway /app/gateway

EXPOSE 7777

CMD ["/app/gateway"]
