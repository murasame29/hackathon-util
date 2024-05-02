FROM golang:1.22.2 AS builder
WORKDIR /build

COPY . .

RUN go build -o app ./cmd/app/main.go

FROM alpine:latest AS app
WORKDIR /app

COPY --from=builder ./build/app /app

ENV PORT=8080

EXPOSE ${PORT}

CMD ["/app"]