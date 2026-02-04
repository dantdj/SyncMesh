# Build stage
FROM golang:1.25 AS builder

WORKDIR /app

# Copy dependency files first to leverage Docker layer caching
COPY signalling-server/go.mod signalling-server/go.sum ./signalling-server/
COPY api/go.mod ./api/go.mod
RUN cd signalling-server && go mod download

# Copy source code and build
COPY signalling-server/ ./signalling-server/
COPY api ./api
RUN cd signalling-server && CGO_ENABLED=0 GOOS=linux go build -o signalling-server .

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/signalling-server .

EXPOSE 8080
CMD ["./signalling-server"]
