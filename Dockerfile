FROM golang:1.24

WORKDIR /notezy-backend

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o notezy-backend .

EXPOSE 7777
