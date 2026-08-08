FROM golang:1.26-alpine AS base

WORKDIR /notezy-backend

RUN apk add --no-cache git

COPY contracts/go.mod contracts/go.sum contracts/
COPY shared/go.mod shared/go.sum shared/
COPY internal/gateway/go.mod internal/gateway/go.sum internal/gateway/
WORKDIR /notezy-backend/internal/gateway
RUN go mod download
WORKDIR /notezy-backend

COPY . .

FROM base AS dev

WORKDIR /notezy-backend/internal/gateway

RUN go install github.com/air-verse/air@latest

FROM base AS builder

WORKDIR /notezy-backend/internal/gateway

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/gateway ./commands

FROM alpine:latest AS production

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /out/gateway /app/gateway

EXPOSE 7777

CMD ["/app/gateway"]
