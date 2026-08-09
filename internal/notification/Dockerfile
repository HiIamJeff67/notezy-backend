FROM golang:1.26-alpine AS base

WORKDIR /notezy-backend

RUN apk add --no-cache git

COPY contracts/go.mod contracts/go.sum contracts/
COPY shared/go.mod shared/go.sum shared/
COPY internal/notification/go.mod internal/notification/go.sum internal/notification/
WORKDIR /notezy-backend/internal/notification
RUN go mod download
WORKDIR /notezy-backend

COPY . .

FROM base AS dev

WORKDIR /notezy-backend/internal/notification

RUN go install github.com/air-verse/air@latest

FROM base AS builder

WORKDIR /notezy-backend/internal/notification

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/notification ./commands

FROM alpine:latest AS production

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /out/notification /app/notification

EXPOSE 7781

CMD ["/app/notification"]
