# Build Stage
FROM golang:1.26.6-alpine3.24 AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

# Run Stage
FROM alpine:3.24
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
RUN apk add --no-cache sqlite
CMD ["/app/main"]
