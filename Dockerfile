FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/api/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /server .
COPY --from=builder /app/.env .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./server"]