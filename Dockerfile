FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o babel cmd/web/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/babel .
COPY --from=builder /app/web ./web
EXPOSE 8080

CMD ["./babel"]
