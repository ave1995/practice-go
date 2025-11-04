FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o grpc-client .

# --- runtime stage ---
FROM alpine:3.22.2
WORKDIR /app
COPY --from=builder /app/grpc-client .

EXPOSE 8080
CMD ["./grpc-client"]