FROM golang:1.22.2 AS builder
WORKDIR /build
COPY . .
RUN go build -o main ./cmd/app/main.go

FROM alpine:3.17
WORKDIR /
COPY --from=builder /build/main .

EXPOSE 8080

CMD [ "./main"]