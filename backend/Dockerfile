FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o crm-api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o crm-worker ./cmd/worker

FROM alpine:3.20

RUN adduser -D -u 10001 appuser

WORKDIR /home/appuser

COPY --from=builder /app/crm-api .
COPY --from=builder /app/crm-worker .

USER appuser

EXPOSE 8080

# Default command runs the API. The worker service overrides this with
# ./crm-worker in docker-compose / Kubernetes.
CMD ["./crm-api"]
