# -----------------------------------------------------------------
# Builder
# -----------------------------------------------------------------
FROM golang:1.20-alpine AS builder
WORKDIR /company
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o company .

# -----------------------------------------------------------------
# Runner
# -----------------------------------------------------------------
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /company .
COPY config.yaml .
EXPOSE 8080
CMD ["./company"]
