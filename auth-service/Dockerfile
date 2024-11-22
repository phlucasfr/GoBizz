FROM golang:1.23.3-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o main ./cmd

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache libc6-compat

COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main"]
