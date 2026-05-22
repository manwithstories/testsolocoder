FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o e-learning-platform ./cmd/main.go

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/e-learning-platform .
COPY config.yaml .

RUN mkdir -p uploads logs

EXPOSE 8080

CMD ["./e-learning-platform"]
