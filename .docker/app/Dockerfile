FROM golang:1.22.2 AS builder
WORKDIR /build
COPY . .
RUN go build -o main ./cmd/app/main.go

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /build/main /app/main

EXPOSE 8080

CMD [ "/app/main"]