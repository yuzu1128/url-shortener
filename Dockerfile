# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN go build -o url-shortener .

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/url-shortener .
COPY index.html .
EXPOSE 8080
CMD ["./url-shortener"]
