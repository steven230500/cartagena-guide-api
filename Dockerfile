# syntax=docker/dockerfile:1.7

FROM golang:1.27-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api
RUN CGO_ENABLED=0 go build -o /out/seed ./cmd/seed

FROM alpine:3.20 AS runner
WORKDIR /app
RUN adduser -D -u 10001 appuser
COPY --from=builder /out/api ./api
COPY --from=builder /out/seed ./seed
USER appuser

EXPOSE 8080
CMD ["./api"]
