FROM golang:1.22.2 AS go
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 

COPY . .

RUN go build -o main cmd/app/main.go

FROM ubuntu:latest
WORKDIR /app
COPY --from=go /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go /app/main .
USER 1001

CMD ["/app/main"]