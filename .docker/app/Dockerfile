FROM golang:1.22.2 AS go
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download 

COPY . .

RUN go build -o main cmd/app/main.go

FROM gcr.io/distroless/static-debian10
WORKDIR /app
COPY --from=go /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go /build/main .
USER 1001

CMD ["/app/main"]