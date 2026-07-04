# Stage 1: Build
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# Download go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o pendekin-backend ./cmd/main.go

# Stage 2: Run
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/pendekin-backend .
EXPOSE 8080

CMD ["./pendekin-backend"]
