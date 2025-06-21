# Build stage
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
COPY go.work go.work.sum ./
COPY shared/ ./shared/
COPY api/ ./api/
COPY services/api-gateway/ ./services/api-gateway/

WORKDIR /app/services/api-gateway

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/services/api-gateway/main .

# EXPOSE 8000

CMD ["./main"]