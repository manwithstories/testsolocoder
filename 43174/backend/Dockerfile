FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o campus-trade-platform ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/campus-trade-platform .
COPY --from=builder /app/.env .

RUN mkdir -p uploads exports
RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

CMD ["./campus-trade-platform"]
