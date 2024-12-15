# Build
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compila il binario per Linux
WORKDIR /app/cmd/feature-flag-service-api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/api .

# Runtime
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/api .
COPY .env .

RUN chmod +x ./api

EXPOSE 3000

# Avvia l'applicazione
CMD ["./api"]
