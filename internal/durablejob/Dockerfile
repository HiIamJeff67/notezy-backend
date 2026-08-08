FROM golang:1.26-alpine AS base

WORKDIR /notezy-backend

RUN apk add --no-cache git

COPY contracts/go.mod contracts/go.sum contracts/
COPY shared/go.mod shared/go.sum shared/
COPY internal/durablejob/go.mod internal/durablejob/go.sum internal/durablejob/
WORKDIR /notezy-backend/internal/durablejob
RUN go mod download
WORKDIR /notezy-backend

COPY . .

FROM base AS dev

WORKDIR /notezy-backend/internal/durablejob

RUN go install github.com/air-verse/air@latest

FROM base AS builder

WORKDIR /notezy-backend/internal/durablejob

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/durablejob ./commands

FROM alpine:latest AS production

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /out/durablejob /app/durablejob

EXPOSE 8082

CMD ["/app/durablejob"]
