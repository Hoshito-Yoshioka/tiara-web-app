# ==============================
# Go Backend - Multi-stage build
# ==============================
FROM golang:1.25-bookworm AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# ==============================
# Runtime
# ==============================
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates postgresql-client \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/server .
COPY migrations/ ./migrations/
COPY scripts/ ./scripts/

RUN mkdir -p uploads/staff

EXPOSE 1323
CMD ["./server"]
