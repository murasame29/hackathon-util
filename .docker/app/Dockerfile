FROM golang:1.22.2 AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download 

COPY . .

RUN go build -o /bin/app cmd/app/main.go

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=builder /bin/app /app/app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/app"]