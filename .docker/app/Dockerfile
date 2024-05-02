FROM golang:1.22.2 AS builder
WORKDIR /build

COPY . .

RUN go build -o main cmd/app/main.go

FROM ubuntu AS runner
WORKDIR /app

COPY --from=builder /build/main .

ENV PORT=8080

EXPOSE ${PORT}

CMD ["./main"]