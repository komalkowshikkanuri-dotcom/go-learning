# Build stage
FROM golang:1.26.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /url-shortener .

# Runtime stage
FROM scratch

COPY --from=builder /url-shortener /url-shortener

EXPOSE 8080

ENTRYPOINT ["/url-shortener"]
