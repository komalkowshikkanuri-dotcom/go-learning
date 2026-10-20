FROM golang:1.26 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG APP=api
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/${APP}

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /out/app /app/app

ENTRYPOINT ["/app/app"]
