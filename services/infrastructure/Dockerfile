# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
COPY go.work go.work.sum ./
COPY shared/ ./shared/
COPY api/proto/ ./api/proto/
COPY services/infrastructure/ ./services/infrastructure/
WORKDIR /app/services/infrastructure

RUN go mod download
RUN go build -o infrastructure main.go

# Runtime stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app
COPY --from=builder /app/services/infrastructure/infrastructure .


EXPOSE 50052

CMD ["./infrastructure"]