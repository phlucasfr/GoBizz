FROM golang:1.23.3 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy

RUN go build -o main .

FROM gcr.io/distroless/base-debian11

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main"]
