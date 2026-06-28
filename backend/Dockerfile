FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o crm-api ./cmd/api

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/crm-api .

EXPOSE 8080

CMD ["./crm-api"]