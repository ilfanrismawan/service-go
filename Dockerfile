FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Update go.mod and go.sum to ensure consistency
RUN go mod tidy

# Generate swagger docs (continue build even if swagger fails)
RUN swag init -g cmd/app/main.go || echo "Warning: Swagger generation failed, continuing build..."

# Build the application and migration tool
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app && \
  CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate ./cmd/migrate

# Final stage
FROM alpine:latest

# Install ca-certificates and tzdata for timezone support
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binaries from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .

EXPOSE 8080

# Run migrations and start the application
CMD ["sh", "-c", "./migrate && ./main"]

