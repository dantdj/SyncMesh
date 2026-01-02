# Build stage
FROM golang:1.25 AS builder

WORKDIR /app

# Copy dependency files first to leverage Docker layer caching
COPY go.mod ./
RUN go mod download

# Copy source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o local-client .

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/local-client .

EXPOSE 8080
CMD ["./local-client"]

