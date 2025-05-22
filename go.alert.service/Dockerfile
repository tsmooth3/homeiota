FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go.alert.service main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/go.alert.service .

# run the binary every 5 minutes
CMD ["sh", "-c", "while true; do ./go.alert.service; sleep 300; done"]
